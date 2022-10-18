//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	azureExample = "azure --location westus2 --subscription-id <sub_id> images/azure/centos-79.yaml"
	azureUse     = "azure <image.yaml>"
)

func NewAzureBuildCmd() *cobra.Command {
	flags := &buildCLIFlags{}
	cmd := &cobra.Command{
		Use:     azureUse,
		Short:   "build and provision azure images",
		Example: azureExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runBuild(args[0], flags)
		},
	}

	initBuildAzureFlags(cmd.Flags(), flags)
	return cmd
}

func NewAzureGenerateCmd() *cobra.Command {
	flags := &generateCLIFlags{}
	cmd := &cobra.Command{
		Use:     azureUse,
		Short:   "generate files relating to building azure images",
		Example: azureExample,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runGenerate(args[0], flags)
		},
	}

	initGenerateAzureFlags(cmd.Flags(), flags)
	return cmd
}

func initBuildAzureFlags(fs *flag.FlagSet, buildFlags *buildCLIFlags) {
	initGenerateArgs(fs, &buildFlags.generateCLIFlags)
	initAzurergs(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, buildFlags)
}

func initGenerateAzureFlags(fs *flag.FlagSet, generateFlags *generateCLIFlags) {
	initGenerateArgs(fs, generateFlags)
	initAzurergs(fs, generateFlags)
}

func initAzurergs(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.Azure = &app.AzureArgs{
		CloudEndpoint: &app.AzureCloudFlag{
			Endpoint: app.AzureCloudEndpointPublic,
		},
	}
	addAzureArgs(fs, gFlags.userArgs.Azure)
}

func addAzureArgs(fs *flag.FlagSet, azure *app.AzureArgs) {
	fs.StringVar(
		&azure.ClientID,
		"client-id",
		"",
		"the client id to use for the build",
	)

	fs.StringVar(
		&azure.GalleryName,
		"gallery-name",
		"dkp",
		"the gallery name to publish the image in",
	)

	fs.StringVar(
		&azure.GalleryImageName,
		"gallery-image-name",
		"",
		"the gallery image name to publish the image to",
	)

	fs.StringVar(
		&azure.GalleryImageOffer,
		"gallery-image-offer",
		"dkp",
		"the gallery image offer to set",
	)

	fs.StringVar(
		&azure.GalleryImagePublisher,
		"gallery-image-publisher",
		"dkp",
		"the gallery image publisher to set",
	)

	fs.StringVar(
		&azure.GalleryImageSKU,
		"gallery-image-sku",
		"",
		"the gallery image sku to set",
	)

	fs.StringVar(
		&azure.Location,
		"location",
		"westus2",
		"the location in which to build the image",
	)

	fs.StringArrayVar(
		&azure.GalleryImageLocations,
		"gallery-image-locations",
		[]string{"westus"},
		"a list of locations to publish the image",
	)

	fs.StringVar(
		&azure.ResourceGroupName,
		"resource-group",
		"dkp",
		"the resource group to create the image in",
	)

	fs.StringVar(
		&azure.SubscriptionID,
		"subscription-id",
		"",
		"the subscription id to use for the build",
	)

	fs.StringVar(
		&azure.TenantID,
		"tenant-id",
		"",
		"the tenant id to use for the build",
	)

	fs.StringVar(
		&azure.InstanceType,
		"instance-type",
		"Standard_D2s_v3",
		"the Instance Type to use for the build",
	)
	fs.Var(
		azure.CloudEndpoint,
		"cloud-endpoint",
		fmt.Sprintf("Azure cloud endpoint. Which can be one of %v", app.ListAzureEndpoints()),
	)
}
