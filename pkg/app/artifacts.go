package app

import (
	"fmt"
	"path/filepath"

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
}

func UploadArtifacts(artifactFlags ArtifactsCmdFlags) error {
	playbookOptions, err := playbookOptionsFromFlag(artifactFlags)
	if err != nil {
		return err
	}
	playbook := appansible.NewPlaybook("upload-artifacts", artifactFlags.Inventory, playbookOptions)
	return playbook.Run(NewRunOptions(artifactFlags.RootFlags))
}

func playbookOptionsFromFlag(artifactFlags ArtifactsCmdFlags) (*ansible.PlaybookOptions, error) {
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
	passedUserArgs := map[string]interface{}{
		"os_packages_local_bundle_file":  osPackagesBundleFile,
		"containerd_local_bundle_file":   containerdBundleFile,
		"pip_packages_local_bundle_file": pipPackagesBundleFile,
		"images_local_bundle_dir":        containerImagesDir,
	}

	if err := mergeUserOverridesToMap(artifactFlags.Overrides, passedUserArgs); err != nil {
		return nil, fmt.Errorf("error merging overrides: %w", err)
	}

	if err = addExtraVarsToMap(artifactFlags.ExtraVars, passedUserArgs); err != nil {
		//nolint:golint // error has context needed
		return nil, err
	}
	playbookOptions := &ansible.PlaybookOptions{
		ExtraVarsMap: passedUserArgs,
	}
	return playbookOptions, nil
}
