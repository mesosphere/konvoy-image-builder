package app

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
)

type ArtifactsCmdFlags struct {
	OSPackagesBundleFile     string
	ContainerdBundleFile     string
	PIPPackagesBundleFile    string
	ContainerImagesBundleDir string
	Inventory                string
	RootFlags
	Overrides []string
	ExtraVars []string
	WorkDir   string
}

type ArtifactUploader struct {
	workDir string
}

func NewArtifactUploader(buildName string) (*ArtifactUploader, error) {
	if buildName == "" {
		name := "artifact-upload"
		buildName = fmt.Sprintf("%s-%d", name, time.Now().Unix())
	}
	workDir, err := createRunDirectory(buildName, OutputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize artifact uploader %w", err)
	}
	return &ArtifactUploader{
		workDir: workDir,
	}, nil
}

func (a *ArtifactUploader) UploadArtifacts(artifactFlags ArtifactsCmdFlags) error {
	playbookOptions, err := a.playbookOptionsFromFlag(artifactFlags)
	if err != nil {
		return err
	}
	log.Printf("writing new configuration to %s", a.workDir)
	playbook := appansible.NewPlaybook("upload-artifacts", artifactFlags.Inventory, playbookOptions)
	return playbook.Run(NewRunOptions(artifactFlags.RootFlags))
}

func (a *ArtifactUploader) playbookOptionsFromFlag(artifactFlags ArtifactsCmdFlags) (*ansible.PlaybookOptions, error) {
	osPackagesBundleFile, err := filepath.Abs(artifactFlags.OSPackagesBundleFile)
	if err != nil {
		return nil, fmt.Errorf("failed to find absolute path for --os-packages-bundle %w", err)
	}
	containerdBundleFile, err := filepath.Abs(artifactFlags.ContainerdBundleFile)
	if err != nil {
		return nil, fmt.Errorf("failed to find absolute path for --containerd-bundle %w", err)
	}
	pipPackagesBundleFile, err := filepath.Abs(artifactFlags.PIPPackagesBundleFile)
	if err != nil {
		return nil, fmt.Errorf("failed to find absolute path for --pip-packages-bundle %w", err)
	}
	containerImagesDir, err := filepath.Abs(artifactFlags.ContainerImagesBundleDir)
	if err != nil {
		return nil, fmt.Errorf("failed to find absolute path for --container-images-dir %w", err)
	}
	args := make(map[string]interface{})
	if err = mergeUserOverridesToMap(artifactFlags.Overrides, args); err != nil {
		return nil, fmt.Errorf("error merging overrides: %w", err)
	}
	if err = addExtraVarsToMap(artifactFlags.ExtraVars, args); err != nil {
		//nolint:golint // error has context needed
		return nil, err
	}
	extraVarsPath, varsErr := filepath.Abs(filepath.Join(a.workDir, ansibleVarsFilename))
	if varsErr != nil {
		return nil, fmt.Errorf("failed to create vars file %w", varsErr)
	}

	passedUserArgs := map[string]interface{}{
		"os_packages_local_bundle_file":  osPackagesBundleFile,
		"containerd_local_bundle_file":   containerdBundleFile,
		"pip_packages_local_bundle_file": pipPackagesBundleFile,
		"images_local_bundle_dir":        containerImagesDir,
	}
	// passedUserArgs take highest precedence
	if err = MergeMapsOverwrite(args, passedUserArgs); err != nil {
		return nil, fmt.Errorf("failed to merge user override %w", err)
	}
	if err = initAnsibleConfig(extraVarsPath, args); err != nil {
		//nolint:golint // error has context needed
		return nil, err
	}
	playbookOptions := &ansible.PlaybookOptions{
		ExtraVars: []string{
			fmt.Sprintf("@%s", extraVarsPath),
		},
	}
	return playbookOptions, nil
}
