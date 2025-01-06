package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	packageCmd        = "create-package-bundle"
	packageCmdExample = "create-package-bundle --os redhat-8.4 --output-directory=artifacts"
	validOS           = []string{
		"centos-7.9",
		"redhat-7.9",
		"redhat-8.4",
		"redhat-8.6",
		"redhat-8.8",
		"rocky-9.1",
		"ubuntu-20.04",
		"oracle-9.4",
	}
)

type packageBundleCmdFlags struct {
	targetOS          string
	kubernetesVersion string
	fips              bool
	outputDirectory   string
	containerImage    string
}

var packageBundleFlags packageBundleCmdFlags

// NOTE Usage information always comes from the konvoy-image command, not the konvoy-image-wrapper command.
// We define this command only to show usage; it does nothing else.
var createPackageBundleCmd = &cobra.Command{
	Use:     packageCmd,
	Short:   "build os package bundles for airgapped installs",
	Example: packageCmdExample,
	Args:    cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}

func init() {
	fs := createPackageBundleCmd.Flags()
	fs.StringVar(
		&packageBundleFlags.targetOS,
		"os",
		"",
		fmt.Sprintf("The target OS you wish to create a package bundle for. Must be one of %v", validOS),
	)
	fs.StringVar(
		&packageBundleFlags.kubernetesVersion,
		"kubernetes-version",
		"",
		"The version of kubernetes to download packages for.",
	)
	fs.BoolVar(
		&packageBundleFlags.fips,
		"fips",
		false,
		"If the package bundle should include fips packages.")
	fs.StringVar(&packageBundleFlags.outputDirectory,
		"output-directory",
		"artifacts",
		"The directory to place the bundle in.")

	fs.StringVar(&packageBundleFlags.containerImage,
		"container-image",
		"",
		"A container image to use for building the package bundles")
}
