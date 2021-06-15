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
	fs.StringArrayVar(&gFlags.overrides, "overrides", []string{}, "a list of override YAML files")
	fs.StringVar(&gFlags.userArgs.KubernetesVersion, "kubernetes-version", "", "the version of kubernetes to install")
	fs.StringVar(&gFlags.userArgs.ContainerdVersion, "containerd-version", "", "the version of containerd to install")
	fs.StringVar(&gFlags.userArgs.AWSBuilderRegion, "region", "", "the aws region to run the builder")
	fs.StringArrayVar(&gFlags.userArgs.AMIRegions, "ami-regions", []string{}, "a list of regions to publish amis")
	fs.StringVar(&gFlags.userArgs.SourceAMI, "source-ami", "", "a specific ami available in the builder region to source from")
	fs.StringVar(&gFlags.userArgs.AMIFilterName, "source-ami-filter-name", "", "a ami name filter on for selecting the source image")
	fs.StringVar(&gFlags.userArgs.AMIFilterOwner, "source-ami-filter-owner", "", "only search AMIs belonging to this owner id")
}
