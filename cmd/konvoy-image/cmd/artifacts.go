package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var artifactsFlags app.ArtifactsCmdFlags

var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "upload artifacts to hosts defined in inventory-file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		uploader, err := app.NewArtifactUploader(artifactsFlags.WorkDir)
		err = uploader.UploadArtifacts(artifactsFlags)
		if err != nil {
			return fmt.Errorf("failed to upload artifacts %w", err)
		}
		return nil
	},
}

func init() {
	fs := artifactsCmd.Flags()
	fs.StringVar(&artifactsFlags.Inventory, "inventory-file", "inventory.yaml", "an ansible inventory defining your infrastructure")
	fs.StringVar(&artifactsFlags.OSPackagesBundleFile, "os-packages-bundle", "", "path to os-packages tar file for install on remote hosts.")
	fs.StringVar(&artifactsFlags.ContainerdBundleFile, "containerd-bundle", "", "path to Containerd tar file for install on remote hosts.")
	fs.StringVar(&artifactsFlags.PIPPackagesBundleFile, "pip-packages-bundle", "", "path to pip-packages tar file"+
		" for install on remote hosts.")
	fs.StringVar(&artifactsFlags.ContainerImagesBundleDir, "container-images-dir", "", "path to container images for install on remote hosts.")
	addOverridesArg(fs, &artifactsFlags.Overrides)
	addWorkDirArg(fs, &artifactsFlags.WorkDir)
	addExtraVarsArg(fs, &artifactsFlags.ExtraVars)
}
