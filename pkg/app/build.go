package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/packer"
	"github.com/mesosphere/konvoy-image-builder/pkg/stringutil"
)

const (
	ansibleVarsFilename      = "ansible_vars.yaml"
	manifestFileName         = "packer.json"
	runDirectorySuffixLength = 5
)

type Builder struct{}

type InitOptions struct {
	CommonConfigPath string
	Image            string
	Overrides        []string
	UserArgs         UserArgs

	// ExtraVarsOnly is true when only ansible variables should only ansible variables
	// should be generated. Omitting packer variables. This is useful for working with
	// preprovisioned infrastructure
	ExtraVarsOnly bool
}

type BuildOptions struct {
	PackerPath         string
	PackerBuildFlags   packer.BuildFlags
	CustomManifestPath string
	DryRun             bool
}

type AmazonArgs struct {
	// AMI options
	SourceAMI        string   `json:"source_ami"`
	AMIFilterName    string   `json:"ami_filter_name"`
	AMIFilterOwner   string   `json:"ami_filter_owner"`
	AWSBuilderRegion string   `json:"aws_region"`
	AMIRegions       []string `json:"ami_regions"`
	AWSInstanceType  string   `json:"aws_instance_type"`

	AMIUsers  []string `json:"ami_users"`
	AMIGroups []string `json:"ami_groups"`
}

type ClusterArgs struct {
	KubernetesVersion string `json:"kubernetes_version" yaml:"kubernetes_version"`
	ContainerdVersion string `json:"containerd_version" yaml:"containerd_version"`
}

type UserArgs struct {
	ClusterArgs

	Amazon *AmazonArgs
	Azure  *AzureArgs

	// ExtraVars provided to ansible
	ExtraVars []string
}

//nolint:gocyclo // this will be refactored
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

	if err = MergeMapsOverwrite(config, imageConfig); err != nil {
		return "", fmt.Errorf("error merging image config: %w", err)
	}

	overrides, err := getOverrides(initOptions.Overrides)
	if err != nil {
		return "", fmt.Errorf("error getting overrides: %w", err)
	}

	if err = MergeMapsOverwrite(config, overrides...); err != nil {
		return "", fmt.Errorf("error merging overrides: %w", err)
	}

	if err = EnrichKubernetesFullVersion(config, initOptions.UserArgs.KubernetesVersion); err != nil {
		//nolint:golint // error has context needed
		return "", err
	}

	if err = MergeUserArgs(config, initOptions.UserArgs); err != nil {
		return "", fmt.Errorf("error merging user args: %w", err)
	}

	buildName := BuildName(config)
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

	// merge extraVars passed through args into config -- which will
	// show up in extraVarsPath
	extraVarSet := make(map[string]interface{})
	for _, extraVars := range initOptions.UserArgs.ExtraVars {
		set := strings.Split(extraVars, "=")
		//nolint:gomnd // the code is splitting on the equal
		if len(set) == 2 {
			k := set[0]
			v := set[1]
			extraVarSet[k] = v
		}
	}

	if err = MergeMapsOverwrite(config, extraVarSet); err != nil {
		return "", fmt.Errorf("error merging overrides: %w", err)
	}

	if err = initAnsibleConfig(extraVarsPath, config); err != nil {
		return "", err
	}

	if !initOptions.ExtraVarsOnly {
		if err = initPackerConfig(workDir, extraVarsPath, config); err != nil {
			return "", err
		}
	}
	return workDir, err
}

func (b *Builder) Run(workDir string, buildOptions BuildOptions) error {
	config, err := configFromWorkDir(workDir, ansibleVarsFilename)
	if err != nil {
		return err
	}

	builderType := config.Get(PackerBuilderTypePath)
	if builderType == "" {
		return BuildError(
			fmt.Sprintf("%s is not defined in image manifest", PackerBuilderTypePath),
		)
	}

	var manifestPath string
	if buildOptions.CustomManifestPath == "" {
		// copy internal manifest to working directory
		opts := packer.RenderOptions{
			SourceAMIDefined: isSourceAMIProvided(config),
			DryRun:           buildOptions.DryRun,
		}

		var data []byte
		data, err = packer.GetManifest(builderType, &opts)
		if err != nil {
			return fmt.Errorf("error getting internal manifest: %w", err)
		}
		manifestPath = filepath.Join(workDir, manifestFileName)
		if err = ioutil.WriteFile(manifestPath, data, 0o600); err != nil {
			return fmt.Errorf("error writing packer manifest: %w", err)
		}
	} else {
		manifestPath = buildOptions.CustomManifestPath
	}

	switch builderType {
	case BuildTypeAzure:
		if err = ensureAzure(config); err != nil {
			return fmt.Errorf("error ensuring azure config: %w", err)
		}
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

// Provision will run ansible playbook directly on an existing set of hosts.
func (b *Builder) Provision(workDir string, flags ProvisionFlags) error {
	extraVarsPath, err := filepath.Abs(filepath.Join(workDir, ansibleVarsFilename))
	if err != nil {
		return InitConfigError("failed to get ansible variables path", err)
	}
	playbook := appansible.NewPlaybook(
		"provision", flags.Inventory, &ansible.PlaybookOptions{
			ExtraVars: []string{
				fmt.Sprintf("@%s", extraVarsPath),
			},
			ExtraVarsMap: map[string]interface{}{
				"sysprep":             false,
				"packer_builder_type": flags.Provider,
			},
		})

	if err := playbook.Run(NewRunOptions(flags.RootFlags)); err != nil {
		return fmt.Errorf("error running playbook: %w", err)
	}

	return nil
}

// recursively merges maps into orig, orig is modified.
func MergeMapsOverwrite(orig map[string]interface{}, maps ...map[string]interface{}) error {
	for _, m := range maps {
		if err := mergo.Merge(&orig, m, mergo.WithOverride); err != nil {
			return fmt.Errorf("error merging: %w", err)
		}
	}

	return nil
}

func createRunDirectory(buildName, dir string) (string, error) {
	f := fmt.Sprintf("%s-%d-%s", buildName, time.Now().Unix(), stringutil.RandString(runDirectorySuffixLength))
	s := filepath.Join(dir, f)
	if err := os.MkdirAll(s, 0o755); err != nil {
		return "", fmt.Errorf("error creating work directory: %w", err)
	}
	return s, nil
}
