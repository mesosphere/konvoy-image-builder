package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/mesosphere/konvoy-image-builder/pkg/packer"
	"github.com/mesosphere/konvoy-image-builder/pkg/stringutil"
)

const (
	CommonConfigDefaultPath = "./images/common.yaml"
	OutputDir               = "./work"

	kubernetesVersionKey       = "kubernetes_version"
	kubernetesFullVersionKey   = "kubernetes_full_version"
	kubernetesBuildMetadataKey = "kubernetes_build_metadata"
	containerdVersionKey       = "containerd_version"
	buildNameKey               = "build_name"
	buildNameExtraKey          = "build_name_extra"
	ansibleExtraVarsKey        = "ansible_extra_vars"
	manifestFileName           = "packer.json"
	packerBuilderTypeKey       = "packer_builder_type"
	packerSourceAMIKey         = "source_ami"
	packerFilterNameKey        = "ami_filter_name"
	packerFilterOwnerKey       = "ami_filter_owners"
	packerBuilderRegionKey     = "aws_region"
	packerAMIRegionsKey        = "ami_regions"

	ansibleVarsFilename = "ansible_vars.yaml"
)

type Builder struct{}

type InitOptions struct {
	CommonConfigPath string
	Image            string
	Overrides        []string
	UserArgs         UserArgs
}

type BuildOptions struct {
	PackerPath         string
	PackerBuildFlags   packer.BuildFlags
	CustomManifestPath string
}

type UserArgs struct {
	KubernetesVersion string `json:"kubernetes_version" yaml:"kubernetes_version"`
	ContainerdVersion string `json:"containerd_version" yaml:"containerd_version"`

	// AMI options
	SourceAMI        string   `json:"source_ami"`
	AMIFilterName    string   `json:"ami_filter_name"`
	AMIFilterOwner   string   `json:"ami_filter_owner"`
	AWSBuilderRegion string   `json:"aws_region"`
	AMIRegions       []string `json:"ami_regions"`
}

func (b *Builder) InitConfig(initOptions InitOptions) (string, error) {
	config, err := loadYAML(initOptions.CommonConfigPath)
	if err != nil {
		return "", fmt.Errorf("error parsing common vars [%s]: %w", initOptions.CommonConfigPath, err)
	}

	imageConfig, err := loadYAML(initOptions.Image)
	if err != nil {
		op := fmt.Sprintf("error parsing image config [%s]", initOptions.Image)
		return "", InitConfigError(op, err)
	}

	if err = mergeMapsOverwrite(config, imageConfig); err != nil {
		return "", fmt.Errorf("error merging image config: %w", err)
	}

	// TODO: move to function
	overrides := make([]map[string]interface{}, 0, len(initOptions.Overrides))
	for _, o := range initOptions.Overrides {
		data, errIO := loadYAML(o)
		if errIO != nil {
			return "", fmt.Errorf("error loading override: %w", errIO)
		}
		overrides = append(overrides, data)
	}

	if err = mergeMapsOverwrite(config, overrides...); err != nil {
		return "", fmt.Errorf("error merging overrides: %w", err)
	}

	enrichKubernetesFullVersion(config)
	mergeUserArgs(config, initOptions)

	buildName := buildName(config)
	if buildName == "" {
		return "", InitConfigError("build name is not defined", nil)
	}

	workDir, err := createRunDirectory(buildName, OutputDir)
	if err != nil {
		return "", InitConfigError("error creating work directory", err)
	}

	extraVarsPath, err := filepath.Abs(filepath.Join(workDir, ansibleVarsFilename))
	if err != nil {
		return "", InitConfigError("failed to get ansible variables path", err)
	}

	packerData, err := genPackerVars(config, extraVarsPath)
	if err != nil {
		return "", fmt.Errorf("error rendering packer vars: %w", err)
	}

	log.Printf("writing new configuration to %s", workDir)
	if err = ioutil.WriteFile(
		filepath.Join(workDir, "packer_vars.json"), packerData, 0600); err != nil {
		return "", fmt.Errorf("error writing packer variables: %w", err)
	}

	ansibleData, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("error marshelling ansible data: %w", err)
	}
	if err = ioutil.WriteFile(extraVarsPath, ansibleData, 0600); err != nil {
		return "", fmt.Errorf("error writing ansible vars: %w", err)
	}
	return workDir, nil
}

func (b *Builder) Run(workDir string, buildOptions BuildOptions) error {
	config, err := configFromWorkDir(workDir)
	if err != nil {
		return err
	}

	var manifestPath string
	if buildOptions.CustomManifestPath == "" {
		// copy internal manifest to working directory
		builderType := getString(config, packerBuilderTypeKey)
		if builderType == "" {
			return BuildError(fmt.Sprintf("%s is not defined in image manifest", packerBuilderTypeKey))
		}
		opts := packer.RenderOptions{SourceAMIDefined: isSourceAMIProvided(config)}

		var data []byte
		data, err = packer.GetManifest(builderType, &opts)
		if err != nil {
			return fmt.Errorf("error getting internal manifest: %w", err)
		}
		manifestPath = filepath.Join(workDir, manifestFileName)
		if err = ioutil.WriteFile(manifestPath, data, 0600); err != nil {
			return fmt.Errorf("error writing packer manifest: %w", err)
		}
	} else {
		manifestPath = buildOptions.CustomManifestPath
	}

	// TODO: consider supporting these externally and doing a deepcopy instead of manipulating the options
	packerBuildFlags := buildOptions.PackerBuildFlags
	packerBuildFlags.Force = false
	packerBuildFlags.Vars = map[string]string{}

	packerBuildFlags.VarFiles = []string{filepath.Join(workDir, "packer_vars.json")}

	packerCLI := packer.CLIRunner{
		Path: buildOptions.PackerPath,

		// TODO: use multi writer for log output
		Out:    os.Stdout,
		OutErr: os.Stderr,
	}

	log.Print("starting packer build")
	_, err = packerCLI.Build(manifestPath, packerBuildFlags)
	if err != nil {
		return fmt.Errorf("error running packer build: %w", err)
	}

	return nil
}

func configFromWorkDir(workDir string) (map[string]interface{}, error) {
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

func buildName(config map[string]interface{}) string {
	buildName := getString(config, buildNameKey)

	buildNameExtra := getString(config, buildNameExtraKey)
	if buildNameExtra != "" {
		return fmt.Sprintf("%s%s", buildName, buildNameExtra)
	}
	return buildName
}

func genPackerVars(config map[string]interface{}, extraVarsPath string) ([]byte, error) {
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
	p[kubernetesFullVersionKey] = getString(config, kubernetesFullVersionKey)
	p[containerdVersionKey] = getString(config, containerdVersionKey)
	p[buildNameKey] = buildName(config)
	p[buildNameExtraKey] = getString(config, buildNameExtraKey)
	p[ansibleExtraVarsKey] = fmt.Sprintf("@%s", extraVarsPath)

	data, err := json.Marshal(p)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling packer vars")
	}

	return data, nil
}

// enrichKubernetesFullVersion will enrich the kubernetes semver with build metadata added in via
// overrides. A common example is the fips override, which adds +fips.buildnumber to the version.
func enrichKubernetesFullVersion(config map[string]interface{}) {
	metadata, ok := config[kubernetesBuildMetadataKey]
	if !ok {
		config[kubernetesFullVersionKey] = config[kubernetesVersionKey]
	} else {
		config[kubernetesFullVersionKey] = fmt.Sprintf(
			"%s+%s", config[kubernetesVersionKey], metadata)
	}
}

func getString(config map[string]interface{}, key string) string {
	i, ok := config[key]
	if !ok {
		return ""
	}
	s, ok := i.(string)
	if !ok {
		return ""
	}
	return s
}

// recursively merges maps into orig, orig is modified.
func mergeMapsOverwrite(orig map[string]interface{}, maps ...map[string]interface{}) error {
	for _, m := range maps {
		if err := mergo.Merge(&orig, m, mergo.WithOverride); err != nil {
			return fmt.Errorf("error merging: %w", err)
		}
	}

	return nil
}

func loadYAML(path string) (map[string]interface{}, error) {
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

const runDirectorySuffixLength = 5

func createRunDirectory(buildName, dir string) (string, error) {
	f := fmt.Sprintf("%s-%d-%s", buildName, time.Now().Unix(), stringutil.RandString(runDirectorySuffixLength))
	s := filepath.Join(dir, f)
	if err := os.MkdirAll(s, 0755); err != nil {
		return "", fmt.Errorf("error creating work directory: %w", err)
	}
	return s, nil
}

func mergeUserArgs(config map[string]interface{}, initOptions InitOptions) {
	// TODO: more platforms will cause this to grow, perhaps write a generator for it
	if initOptions.UserArgs.KubernetesVersion != "" {
		config[kubernetesVersionKey] = initOptions.UserArgs.KubernetesVersion
	}

	if initOptions.UserArgs.ContainerdVersion != "" {
		config[containerdVersionKey] = initOptions.UserArgs.ContainerdVersion
	}

	p, ok := config["packer"]
	if !ok {
		return
	}

	packerMap, ok := p.(map[interface{}]interface{})
	if !ok {
		return
	}

	if initOptions.UserArgs.AWSBuilderRegion != "" {
		packerMap[packerBuilderRegionKey] = initOptions.UserArgs.AWSBuilderRegion
	}

	if initOptions.UserArgs.AMIFilterName != "" {
		packerMap[packerFilterNameKey] = initOptions.UserArgs.AMIFilterName
	}

	if initOptions.UserArgs.AMIFilterOwner != "" {
		packerMap[packerFilterOwnerKey] = initOptions.UserArgs.AMIFilterOwner
	}

	if len(initOptions.UserArgs.AMIRegions) > 0 {
		packerMap[packerAMIRegionsKey] = initOptions.UserArgs.AMIRegions
	}

	if initOptions.UserArgs.SourceAMI != "" {
		packerMap[packerSourceAMIKey] = initOptions.UserArgs.SourceAMI
		// using a specific AMI, clear filters
		packerMap[packerFilterNameKey] = ""
		packerMap[packerFilterOwnerKey] = ""
	}
}

func isSourceAMIProvided(config map[string]interface{}) bool {
	p, ok := config["packer"]
	if !ok {
		return false
	}

	packerMap, ok := p.(map[interface{}]interface{})
	if !ok {
		return false
	}
	d, ok := packerMap[packerSourceAMIKey]
	if !ok {
		return false
	}

	return d.(string) != ""
}
