package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/pointerstructure"
	"gopkg.in/yaml.v2"

	"github.com/mesosphere/konvoy-image-builder/pkg/version"
)

const (
	httpProxyKey                  = "http_proxy"
	httpsProxyKey                 = "https_proxy"
	noProxyKey                    = "no_proxy"
	packerKIBVersionKey           = "konvoy_image_builder_version"
	packerSSHBastionHostKey       = "ssh_bastion_host"
	packerSSHBastionUsernameKey   = "ssh_bastion_username"
	packerSSHBastionPasswordKey   = "ssh_bastion_password"         //nolint:gosec // just a key
	packerSSHBastionPrivateKeyKey = "ssh_bastion_private_key_file" //nolint:gosec // just a key

)

type Config map[string]interface{}

func (config Config) get(path string) (interface{}, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	value, err := pointerstructure.Get(config, path)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", path, err)
	}

	return value, nil
}

func (config Config) GetSliceWithError(path string) ([]string, error) {
	value, err := config.get(path)
	if err != nil {
		return nil, err
	}

	str, ok := value.([]string)
	if !ok {
		return nil, fmt.Errorf("error %s is not a slice: %w", path, ErrPathNotSlice)
	}

	return str, nil
}

func (config Config) GetWithError(path string) (string, error) {
	value, err := config.get(path)
	if err != nil {
		return "", err
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("error %s is not a string: %w", path, ErrPathNotString)
	}

	return str, nil
}

func (config Config) Get(path string) string {
	str, err := config.GetWithError(path)
	if err != nil {
		return ""
	}

	return str
}

func (config Config) Set(path string, value interface{}) error {
	if value == "" {
		return nil
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	_, err := pointerstructure.Set(config, path, value)
	if err != nil {
		return fmt.Errorf("failed to set %s to %s: %w", path, value, err)
	}

	return nil
}

func (config Config) Delete(path string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	pointer, err := pointerstructure.Parse(path)
	if err != nil {
		return fmt.Errorf("failed to parse path %s: %w", path, err)
	}

	_, err = pointer.Delete(config)
	if err != nil {
		return fmt.Errorf("failed to delete path %s: %w", path, err)
	}

	return nil
}

func BuildName(config Config) string {
	buildName := config.Get(BuildNameKey)

	buildNameExtra := config.Get(BuildNameExtraKey)
	if buildName == "" {
		buildName = DefaultBuildName
	}

	if buildNameExtra != "" {
		return fmt.Sprintf("%s%s", buildName, buildNameExtra)
	}

	return buildName
}

func configFromWorkDir(workDir string, ansibleVarsFilename string) (Config, error) {
	bytes, err := ioutil.ReadFile(path.Join(workDir, ansibleVarsFilename))
	if err != nil {
		return nil, err
	}
	var config map[string]interface{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// enrichKubernetesFullVersion will enrich the kubernetes semver with build metadata added in via
// overrides. A common example is the fips override, which adds +fips.buildnumber to the version.
func EnrichKubernetesFullVersion(config Config, userDefinedKubernetesVersion string) error {
	k8sVersion := config.Get(KubernetesVersionKey)

	// If we have something from the user, use that
	if userDefinedKubernetesVersion != "" {
		k8sVersion = userDefinedKubernetesVersion
	}

	// if we couldn't find it in the config or the user didn't define it
	if len(k8sVersion) == 0 {
		return ErrKubernetesVersionMissing
	}

	metadata := config.Get(KubernetesBuildMetadataKey)
	if metadata != "" {
		k8sVersion = fmt.Sprintf("%s+%s", k8sVersion, metadata)
	}

	err := config.Set(KubernetesFullVersionKey, k8sVersion)
	if err != nil {
		return fmt.Errorf("failed to set enriched kubernetes full version: %w", err)
	}

	return nil
}

func GenPackerVars(config Config, extraVarsPath string) ([]byte, error) {
	i, found := config["packer"]
	p := make(map[string]string)
	if found {
		for k, v := range i.(map[interface{}]interface{}) {
			switch v := v.(type) {
			case string:
				p[k.(string)] = v
			case []byte:
				p[k.(string)] = string(v)
			case fmt.Stringer:
				p[k.(string)] = v.String()
			case nil:
				p[k.(string)] = ""
			default:
				p[k.(string)] = fmt.Sprintf("%v", v)
			}
		}
	}
	// Common vars
	// TODO: make this a map
	p[KubernetesFullVersionKey] = config.Get(KubernetesFullVersionKey)
	p[ContainerdVersionKey] = config.Get(ContainerdVersionKey)
	p[BuildNameKey] = BuildName(config)
	p[BuildNameExtraKey] = config.Get(BuildNameExtraKey)
	p[AnsibleExtraVarsKey] = fmt.Sprintf("@%s", extraVarsPath)
	p[httpProxyKey] = config.Get(httpProxyKey)
	p[httpsProxyKey] = config.Get(httpsProxyKey)
	p[noProxyKey] = config.Get(noProxyKey)
	p[packerKIBVersionKey] = version.Version()

	// if there's no Bastion default to empty values
	// to satisfy packer
	if _, ok := p[packerSSHBastionHostKey]; !ok {
		p[packerSSHBastionUsernameKey] = ""
		p[packerSSHBastionHostKey] = ""
		p[packerSSHBastionPasswordKey] = ""
		p[packerSSHBastionPrivateKeyKey] = ""
	}

	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error marshaling packer vars: %w", err)
	}

	return data, nil
}

func getOverrides(paths []string) ([]map[string]interface{}, error) {
	overrides := make([]map[string]interface{}, 0, len(paths))

	for _, path := range paths {
		data, err := loadYAML(path)
		if err != nil {
			return nil, fmt.Errorf("error loading override: %w", err)
		}

		overrides = append(overrides, data)
	}

	return overrides, nil
}

func initAnsibleConfig(path string, config Config) error {
	ansibleData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshelling ansible data: %w", err)
	}
	if err = ioutil.WriteFile(path, ansibleData, 0o600); err != nil {
		return fmt.Errorf("error writing ansible vars: %w", err)
	}
	return nil
}

func initPackerConfig(workDir, extraVarsPath string, config Config) error {
	packerData, err := GenPackerVars(config, extraVarsPath)
	if err != nil {
		return fmt.Errorf("error rendering packer vars: %w", err)
	}

	log.Printf("writing new packer configuration to %s", workDir)
	if err = ioutil.WriteFile(
		filepath.Join(workDir, "packer_vars.json"), packerData, 0o600); err != nil {
		return fmt.Errorf("error writing packer variables: %w", err)
	}
	return nil
}

func isSourceAMIProvided(config Config) bool {
	return config.Get(PackerSourceAMIPath) != ""
}

func loadYAML(path string) (Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", path, err)
	}

	m := make(map[string]interface{})
	if err := yaml.Unmarshal(data, m); err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", path, err)
	}

	return m, nil
}

func MergeUserArgs(config Config, userArgs UserArgs) error {
	err := config.Set(KubernetesVersionKey, userArgs.KubernetesVersion)
	if err != nil {
		return fmt.Errorf("failed to set %s: %w", KubernetesVersionKey, err)
	}

	err = config.Set(ContainerdVersionKey, userArgs.ContainerdVersion)
	if err != nil {
		return fmt.Errorf("failed to set %s: %w", ContainerdVersionKey, err)
	}

	err = MergeAmazonUserArgs(config, userArgs)
	if err != nil {
		return fmt.Errorf("failed to set amazon args: %w", err)
	}

	return nil
}

func setSourceAMI(config Config, sourceAMI string) error {
	if err := config.Set(PackerSourceAMIPath, sourceAMI); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerSourceAMIPath, err)
	}

	if err := config.Delete(PackerFilterNamePath); err != nil {
		return fmt.Errorf("failed to delete %s: %w", PackerFilterNamePath, err)
	}

	if err := config.Delete(PackerFilterOwnerPath); err != nil {
		return fmt.Errorf("failed to delete %s: %w", PackerFilterOwnerPath, err)
	}

	return nil
}

func MergeAmazonUserArgs(config Config, userArgs UserArgs) error {
	if err := config.Set(PackerBuilderRegionPath, userArgs.AWSBuilderRegion); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerBuilderRegionPath, err)
	}

	if userArgs.SourceAMI != "" {
		if err := setSourceAMI(config, userArgs.SourceAMI); err != nil {
			return fmt.Errorf("failed to source ami: %w", err)
		}
	} else {
		if err := config.Set(PackerFilterNamePath, userArgs.AMIFilterName); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerFilterNamePath, err)
		}

		if err := config.Set(PackerFilterOwnerPath, userArgs.AMIFilterOwner); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerFilterOwnerPath, err)
		}
	}

	if err := config.Set(PackerInstanceTypePath, userArgs.AWSInstanceType); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerInstanceTypePath, err)
	}

	if len(userArgs.AMIRegions) > 0 {
		value := strings.Join(userArgs.AMIRegions, ",")
		if err := config.Set(PackerAMIRegionsPath, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerAMIRegionsPath, err)
		}
	}

	if len(userArgs.AMIUsers) > 0 {
		value := strings.Join(userArgs.AMIUsers, ",")
		if err := config.Set(PackerAMIUsersPath, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerAMIUsersPath, err)
		}
	}

	if len(userArgs.AMIGroups) > 0 {
		value := strings.Join(userArgs.AMIGroups, ",")
		if err := config.Set(PackerAMIGroupsPath, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerAMIGroupsPath, err)
		}
	}

	return nil
}
