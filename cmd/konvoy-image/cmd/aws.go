//

package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	awsExample = "aws --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-79.yaml"
	awsUse     = "aws <image.yaml>"
)

func NewAWSBuildCmd() *cobra.Command {
	flags := &buildCLIFlags{}
	cmd := &cobra.Command{
		Use:     awsUse,
		Short:   "build and provision aws images",
		Example: awsExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runBuild(args[0], flags)
		},
	}

	initBuildAWSFlags(cmd.Flags(), flags)
	return cmd
}

func NewAWSGenerateCmd() *cobra.Command {
	flags := &generateCLIFlags{}
	cmd := &cobra.Command{
		Use:     awsUse,
		Short:   "generate files relating to building aws images",
		Example: awsExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runGenerate(args[0], flags)
		},
	}

	initGenerateAWSFlags(cmd.Flags(), flags)
	return cmd
}

func initBuildAWSFlags(fs *flag.FlagSet, buildFlags *buildCLIFlags) {
	initGenerateArgs(fs, &buildFlags.generateCLIFlags)
	initAWSArgs(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, buildFlags)
}

func initGenerateAWSFlags(fs *flag.FlagSet, generateFlags *generateCLIFlags) {
	initGenerateArgs(fs, generateFlags)
	initAWSArgs(fs, generateFlags)
}

func initAWSArgs(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.Amazon = &app.AmazonArgs{}
	addAWSArgs(fs, gFlags.userArgs.Amazon)
}

func addAWSArgs(fs *flag.FlagSet, amazonArgs *app.AmazonArgs) {
	fs.StringVar(
		&amazonArgs.AWSBuilderRegion,
		"region",
		"",
		"the region in which to build the AMI",
	)
	fs.StringArrayVar(
		&amazonArgs.AMIRegions,
		"ami-regions",
		[]string{},
		"a list of regions to publish amis",
	)
	fs.StringVar(
		&amazonArgs.SourceAMI,
		"source-ami",
		"",
		"the ID of the AMI to use as the source; must be present in the region in which the AMI is built",
	)
	fs.StringVar(
		&amazonArgs.AMIFilterName,
		"source-ami-filter-name",
		"",
		"restricts the set of source AMIs to ones whose Name matches filter",
	)
	fs.StringVar(
		&amazonArgs.AMIFilterOwner,
		"source-ami-filter-owner",
		"",
		"restricts the source AMI to ones with this owner ID",
	)
	fs.StringVar(
		&amazonArgs.AWSInstanceType,
		"aws-instance-type",
		"",
		"instance type used to build the AMI; the type must be present in the region in which the AMI is built",
	)
	_ = fs.MarkDeprecated("aws-instance-type", "please use `--instance-type`.")
	fs.StringVar(
		&amazonArgs.AWSInstanceType,
		"instance-type",
		"",
		"instance type used to build the AMI; the type must be present in the region in which the AMI is built",
	)
	fs.StringArrayVar(
		&amazonArgs.AMIUsers,
		"ami-users",
		[]string{},
		"a list AWS user accounts which are allowed use the image",
	)
	fs.StringArrayVar(
		&amazonArgs.AMIGroups,
		"ami-groups",
		[]string{},
		"a list of AWS groups which are allowed use the image, using 'all' result in a public image",
	)
}
