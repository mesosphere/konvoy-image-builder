package cmd

import (
	"github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

func addOverridesArg(fs *pflag.FlagSet, overrides *[]string) {
	fs.StringArrayVar(overrides, "overrides", []string{}, "a list of override YAML files")
}

func addClusterArgs(fs *pflag.FlagSet, kubernetesVersion, containerdVersion *string) {
	fs.StringVar(kubernetesVersion, "kubernetes-version", "", "the version of kubernetes to install")
	fs.StringVar(containerdVersion, "containerd-version", "", "the version of containerd to install")
}

func addAWSUserArgs(fs *pflag.FlagSet, userArgs *app.UserArgs) {
	fs.StringVar(&userArgs.AWSBuilderRegion, "region", "", "the aws region to run the builder")
	fs.StringArrayVar(&userArgs.AMIRegions, "ami-regions", []string{}, "a list of regions to publish amis")
	fs.StringVar(&userArgs.SourceAMI, "source-ami", "", "a specific ami available in the builder region to source from")
	fs.StringVar(&userArgs.AMIFilterName, "source-ami-filter-name", "", "a ami name filter on for selecting the source image")
	fs.StringVar(&userArgs.AMIFilterOwner, "source-ami-filter-owner", "", "only search AMIs belonging to this owner id")
}
