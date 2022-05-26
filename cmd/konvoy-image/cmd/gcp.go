package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	gcpExample = "gcp ... images/gcp/centos-7.yaml"
	gcpUse     = "gcp <image.yaml>"
)

var gcpBuildCmd = &cobra.Command{
	Use:     gcpUse,
	Short:   "build and provision gcp images",
	Example: gcpExample,
	Args:    cobra.ExactArgs(1),
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		runBuild(args[0])
	},
}

func initBuildGCP() {
	fs := gcpBuildCmd.Flags()

	initGenerateFlags(fs, &buildFlags.generateCLIFlags)
	initGCPFlags(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, &buildFlags)
}

func initGCPFlags(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.Azure = &app.AzureArgs{}
	addGCPArgs(fs, gFlags.userArgs.Azure)
}

func addGCPArgs(fs *flag.FlagSet, azure *app.AzureArgs) {

}
