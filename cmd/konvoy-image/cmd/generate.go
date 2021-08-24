package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

type generateCLIFlags struct {
	overrides []string

	userArgs app.UserArgs
}

var generateFlags generateCLIFlags

var generateCmd = &cobra.Command{
	Use:     "generate <image.yaml>",
	Short:   "generate files relating to building images",
	Example: "generate --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-7.yaml",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		builder := newBuilder()
		_, err := builder.InitConfig(newInitOptions(args[0], generateFlags))
		if err != nil {
			bail("error rendering builder configuration", err, 2)
		}
	},
}

func newInitOptions(image string, flags generateCLIFlags) app.InitOptions {
	return app.InitOptions{
		CommonConfigPath: app.CommonConfigDefaultPath,
		Image:            image,
		Overrides:        flags.overrides,
		UserArgs:         flags.userArgs,
	}
}

func init() {
	initGenerateFlags(generateCmd.Flags(), &generateFlags)
}

func initGenerateFlags(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	addOverridesArg(fs, &gFlags.overrides)
	addClusterArgs(
		fs,
		&gFlags.userArgs.ClusterArgs.KubernetesVersion,
		&gFlags.userArgs.ClusterArgs.ContainerdVersion,
	)
	addAWSUserArgs(fs, &gFlags.userArgs)
}
