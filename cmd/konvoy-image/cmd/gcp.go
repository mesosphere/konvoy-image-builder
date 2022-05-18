package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	gcpExample = "gcp ... images/gcp/centos-79.yaml"
	gcpUse     = "gcp <image.yaml>"
)

var gcpBuildCmd = &cobra.Command{
	Use:     gcpUse,
	Short:   "build and provision gcp images",
	Example: gcpExample,
	Args:    cobra.ExactArgs(1),
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
	gFlags.userArgs.GCP = &app.GCPArgs{}
	addGCPArgs(fs, gFlags.userArgs.GCP)
}

func addGCPArgs(fs *flag.FlagSet, gcp *app.GCPArgs) {
	fs.StringVar(
		&gcp.ProjectID,
		"project-id",
		"default",
		"the project id to use when storing created image",
	)

	fs.StringVar(
		&gcp.Network,
		"network",
		"default",
		"the network to use when creating an image",
	)

	fs.StringVar(
		&gcp.Zone,
		"zone",
		"default",
		"the zone to use when storing a created image",
	)
}
