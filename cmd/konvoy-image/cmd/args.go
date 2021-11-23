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
	fs.StringVar(
		&userArgs.AWSBuilderRegion,
		"region",
		"",
		"the region in which to build the AMI",
	)
	fs.StringArrayVar(
		&userArgs.AMIRegions,
		"ami-regions",
		[]string{},
		"a list of regions to publish amis",
	)
	fs.StringVar(
		&userArgs.SourceAMI,
		"source-ami",
		"",
		"the ID of the AMI to use as the source; must be present in the region in which the AMI is built",
	)
	fs.StringVar(
		&userArgs.AMIFilterName,
		"source-ami-filter-name",
		"",
		"restricts the set of source AMIs to ones whose Name matches filter",
	)
	fs.StringVar(
		&userArgs.AMIFilterOwner,
		"source-ami-filter-owner",
		"",
		"restricts the source AMI to ones with this owner ID"  
	)
	fs.StringVar(
		&userArgs.AWSInstanceType,
		"aws-instance-type",
		"",
		"instance type used to build the AMI; the type must be present in the region in which the AMI is built",
	)
	fs.StringArrayVar(
		&userArgs.AMIUsers,
		"ami-users",
		[]string{},
		"a list AWS user accounts which are allowed use the image",
	)
	fs.StringArrayVar(
		&userArgs.AMIGroups,
		"ami-groups",
		[]string{},
		"a list of AWS groups which are allowed use the image, using 'all' result in a public image",
	)
}
