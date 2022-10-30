package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	gcpExample = "gcp ... images/gcp/centos-79.yaml"
	gcpUse     = "gcp <image.yaml>"
)

const (
	//nolint:gosec // environment var set by user
	envGCPApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
)

func NewGCPBuildCmd() *cobra.Command {
	flags := &buildCLIFlags{}
	cmd := &cobra.Command{
		Use:     gcpUse,
		Short:   "build and provision gcp images",
		Example: gcpExample,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			filePath, found := os.LookupEnv(envGCPApplicationCredentials)
			if !found {
				return fmt.Errorf("envornment variable %s is not set", envGCPApplicationCredentials)
			}
			_, err := os.Stat(filePath)
			if err != nil {
				return fmt.Errorf(
					"could not determine if file %q assigned to %s environment variable exists: %w",
					filePath,
					envGCPApplicationCredentials,
					err)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runBuild(args[0], flags)
		},
	}

	initBuildGCPFlags(cmd.Flags(), flags)
	return cmd
}

func initBuildGCPFlags(fs *flag.FlagSet, buildFlags *buildCLIFlags) {
	initGenerateArgs(fs, &buildFlags.generateCLIFlags)
	initGCPArgs(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, buildFlags)
}

func initGCPArgs(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.GCP = &app.GCPArgs{}
	addGCPArgs(fs, gFlags.userArgs.GCP)
}

func addGCPArgs(fs *flag.FlagSet, gcp *app.GCPArgs) {
	fs.StringVar(
		&gcp.ProjectID,
		"project-id",
		"",
		"the project id to use when storing created image",
	)

	fs.StringVar(
		&gcp.Network,
		"network",
		"",
		"the network to use when creating an image",
	)

	fs.StringVar(
		&gcp.Region,
		"region",
		"us-west1",
		"the region in which to launch the instance",
	)
}
