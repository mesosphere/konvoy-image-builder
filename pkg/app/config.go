package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
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

// NOTE(jkoelker) `strval` and `strslice` are taken from https://github.com/Masterminds/sprig.
func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func strslice(v interface{}) ([]string, error) {
	switch v := v.(type) {
	case []string:
		return v, nil
	case []interface{}:
		b := make([]string, 0, len(v))
		for _, s := range v {
			if s != nil {
				b = append(b, strval(s))
			}
		}
		return b, nil
	default:
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			l := val.Len()
			b := make([]string, 0, l)
			for i := 0; i < l; i++ {
				value := val.Index(i).Interface()
				if value != nil {
					b = append(b, strval(value))
				}
			}
			return b, nil
		default:
			if v == nil {
				return nil, nil
			}

			return nil, ErrPathNotSlice
		}
	}
}

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

func (config Config) AddBuildNameExtra(base string) string {
	return fmt.Sprintf("%s%s", base, config.Get(BuildNameExtraKey))
}

func (config Config) GetSliceWithError(path string) ([]string, error) {
	value, err := config.get(path)
	if err != nil {
		return nil, err
	}

	slice, err := strslice(value)
	if err != nil {
		return nil, fmt.Errorf("error %s is not a string slice (%T): %w", path, value, err)
	}

	return slice, nil
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

func (config Config) GetWithEnvironment(path string, environmentVariable string) (string, error) {
	value, ok := os.LookupEnv(environmentVariable)
	if ok {
		return value, nil
	}

	value, err := config.GetWithError(path)
	if err != nil {
		return "", fmt.Errorf("error getting %s from config: %w", path, err)
	}

	return value, nil
}

func (config Config) Set(path string, value interface{}) error {
	if value == "" {
		return nil
	}

	return config.SetAllowEmpty(path, value)
}

func (config Config) SetAllowEmpty(path string, value interface{}) error {
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

	if buildName == "" {
		buildName = DefaultBuildName
	}

	return config.AddBuildNameExtra(buildName)
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

func genPackerGPUVars(config Config) error {
	gpuEnabled := false
	gpuTypes := ""
	gpuNvidiaVersion := ""
	// NOTE(jkoelker) Copy GPU configuration from top level into packer.
	//                If it is not found, that is OK, GPU is just not enabled.
	types, err := config.GetSliceWithError(GPUTypesKey)
	if err != nil && !errors.Is(err, pointerstructure.ErrNotFound) {
		return fmt.Errorf("failed to get gpu types: %w", err)
	}

	if len(types) > 0 {
		gpuEnabled = true
		gpuTypes = strings.Join(types, ",")
		gpuNvidiaVersion = config.Get(GPUNvidiaVersion)
	}

	if err := config.Set(PackerGPUPath, gpuEnabled); err != nil {
		return fmt.Errorf("failed to set packer gpu enabled: %w", err)
	}

	if err := config.SetAllowEmpty(PackerGPUTypes, gpuTypes); err != nil {
		return fmt.Errorf("failed to set packer gpu enabled: %w", err)
	}

	if err := config.SetAllowEmpty(PackerGPUNvidiaVersion, gpuNvidiaVersion); err != nil {
		return fmt.Errorf("failed to set packer gpu enabled: %w", err)
	}

	return nil
}

func GenPackerVars(config Config, extraVarsPath string) ([]byte, error) {
	if err := genPackerGPUVars(config); err != nil {
		return nil, fmt.Errorf("error copying gpu config to packer config: %w", err)
	}

	i, found := config["packer"]
	p := make(map[string]interface{})
	if found {
		for k, v := range i.(map[interface{}]interface{}) {
			key := k.(string)

			switch v := v.(type) {
			case string:
				p[key] = v
			case []byte:
				p[key] = string(v)
			case fmt.Stringer:
				p[key] = v.String()
			case nil:
				p[key] = ""
			case []string:
				p[key] = strings.Join(v, ",")
			default:
				p[key] = fmt.Sprintf("%v", v)
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
	if err := config.Set(KubernetesVersionKey, userArgs.KubernetesVersion); err != nil {
		return fmt.Errorf("failed to set %s: %w", KubernetesVersionKey, err)
	}

	if err := config.Set(ContainerdVersionKey, userArgs.ContainerdVersion); err != nil {
		return fmt.Errorf("failed to set %s: %w", ContainerdVersionKey, err)
	}

	if userArgs.Amazon != nil {
		if err := MergeAmazonUserArgs(config, userArgs.Amazon); err != nil {
			return fmt.Errorf("failed to set amazon args: %w", err)
		}
	}

	if userArgs.Azure != nil {
		if err := MergeAzureUserArgs(config, userArgs.Azure); err != nil {
			return fmt.Errorf("failed to set azure args: %w", err)
		}
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

func MergeAmazonUserArgs(config Config, amazonArgs *AmazonArgs) error {
	if err := config.Set(PackerBuilderRegionPath, amazonArgs.AWSBuilderRegion); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerBuilderRegionPath, err)
	}

	if amazonArgs.SourceAMI != "" {
		if err := setSourceAMI(config, amazonArgs.SourceAMI); err != nil {
			return fmt.Errorf("failed to source ami: %w", err)
		}
	} else {
		if err := config.Set(PackerFilterNamePath, amazonArgs.AMIFilterName); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerFilterNamePath, err)
		}

		if err := config.Set(PackerFilterOwnerPath, amazonArgs.AMIFilterOwner); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerFilterOwnerPath, err)
		}
	}

	if err := config.Set(PackerInstanceTypePath, amazonArgs.AWSInstanceType); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerInstanceTypePath, err)
	}

	if len(amazonArgs.AMIRegions) > 0 {
		value := strings.Join(amazonArgs.AMIRegions, ",")
		if err := config.Set(PackerAMIRegionsPath, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerAMIRegionsPath, err)
		}
	}

	if len(amazonArgs.AMIUsers) > 0 {
		value := strings.Join(amazonArgs.AMIUsers, ",")
		if err := config.Set(PackerAMIUsersPath, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerAMIUsersPath, err)
		}
	}

	if len(amazonArgs.AMIGroups) > 0 {
		value := strings.Join(amazonArgs.AMIGroups, ",")
		if err := config.Set(PackerAMIGroupsPath, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", PackerAMIGroupsPath, err)
		}
	}

	return nil
}

func MergeAzureUserArgs(config Config, azureArgs *AzureArgs) error {
	if err := config.Set(PackerAzureClientIDPath, azureArgs.ClientID); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureTenantIDPath, err)
	}

	if err := config.Set(PackerAzureInstanceType, azureArgs.InstanceType); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureInstanceType, err)
	}

	galleryImageLocations := azureArgs.GalleryImageLocations
	if len(galleryImageLocations) == 0 {
		galleryImageLocations = []string{azureArgs.Location}
	}

	if err := config.Set(PackerAzureGalleryLocations, galleryImageLocations); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureGalleryLocations, err)
	}

	galleryImageName := azureArgs.GalleryImageName
	if galleryImageName == "" {
		galleryImageName = fmt.Sprintf("dkp-%s", BuildName(config))
	}

	if err := config.Set(PackerAzureGalleryImageNamePath, galleryImageName); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureGalleryImageNamePath, err)
	}

	if err := config.Set(
		PackerAzureGalleryImageOfferPath,
		azureArgs.GalleryImageOffer,
	); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureGalleryImageOfferPath, err)
	}

	if err := config.Set(
		PackerAzureGalleryImagePublisherPath,
		azureArgs.GalleryImagePublisher,
	); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureGalleryImagePublisherPath, err)
	}

	if err := config.Set(PackerAzureGalleryImageSKU, azureArgs.GalleryImageSKU); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureGalleryImageSKU, err)
	}

	if err := config.Set(PackerAzureLocation, azureArgs.Location); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureLocation, err)
	}

	if err := config.Set(PackerAzureGalleryNamePath, azureArgs.GalleryName); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureGalleryNamePath, err)
	}

	if err := config.Set(
		PackerAzureResourceGroupNamePath,
		azureArgs.ResourceGroupName,
	); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureResourceGroupNamePath, err)
	}

	if err := config.Set(PackerAzureSubscriptionIDPath, azureArgs.SubscriptionID); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureSubscriptionIDPath, err)
	}

	if err := config.Set(PackerAzureTenantIDPath, azureArgs.TenantID); err != nil {
		return fmt.Errorf("failed to set %s: %w", PackerAzureTenantIDPath, err)
	}

	return nil
}
