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

func NewGenereateCmd() *cobra.Command {
	flags := &generateCLIFlags{}

	cmd := &cobra.Command{
		Use:   "generate <image.yaml>",
		Short: "generate files relating to building images",
		Long: "Generate files relating to building images. Specifying AWS arguments is deprecated " +
			"and will be removed in a future version. Use the `aws` subcommand instead.",
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			runGenerate(args[0], flags)
		},
	}

	cmd.AddCommand(NewAWSGenerateCmd())
	cmd.AddCommand(NewAzureGenerateCmd())

	initGenerateFlags(cmd.Flags(), flags)

	return cmd
}

func newInitOptions(image string, flags *generateCLIFlags) app.InitOptions {
	return app.InitOptions{
		CommonConfigPath: app.CommonConfigDefaultPath,
		Image:            image,
		Overrides:        flags.overrides,
		UserArgs:         flags.userArgs,
	}
}

func runGenerate(image string, generateFlags *generateCLIFlags) {
	builder := newBuilder()
	_, err := builder.InitConfig(newInitOptions(image, generateFlags))
	if err != nil {
		bail("error rendering builder configuration", err, 2)
	}
}

func initGenerateFlags(fs *flag.FlagSet, generateFlags *generateCLIFlags) {
	initGenerateArgs(fs, generateFlags)
	initAWSArgs(fs, generateFlags)
}

func initGenerateArgs(fs *flag.FlagSet, generateFlags *generateCLIFlags) {
	addOverridesArg(fs, &generateFlags.overrides)
	addClusterArgs(
		fs,
		&generateFlags.userArgs.ClusterArgs.KubernetesVersion,
		&generateFlags.userArgs.ClusterArgs.ContainerdVersion,
	)

	addExtraVarsArg(fs, &generateFlags.userArgs.ExtraVars)
}
