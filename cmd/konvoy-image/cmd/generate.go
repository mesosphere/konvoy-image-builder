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
	Use:   "generate <image.yaml>",
	Short: "generate files relating to building images",
	Long: "Generate files relating to building images. Specifying AWS arguments is deprecated " +
		"and will be removed in a future version. Use the `aws` subcommand instead.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runGenerate(args[0])
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

func runGenerate(image string) {
	builder := newBuilder()
	_, err := builder.InitConfig(newInitOptions(image, generateFlags))
	if err != nil {
		bail("error rendering builder configuration", err, 2)
	}
}

func init() {
	initGenerateAws()
	initGenerateAzure()

	fs := generateCmd.Flags()

	initGenerateFlags(fs, &generateFlags)
	initAmazonFlags(fs, &generateFlags)

	generateCmd.AddCommand(awsGenerateCmd)
	generateCmd.AddCommand(azureGenerateCmd)
}

func initGenerateFlags(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	addOverridesArg(fs, &gFlags.overrides)
	addClusterArgs(
		fs,
		&gFlags.userArgs.ClusterArgs.KubernetesVersion,
		&gFlags.userArgs.ClusterArgs.ContainerdVersion,
	)

	addExtraVarsArg(fs, &gFlags.userArgs.ExtraVars)
}
