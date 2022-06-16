//

package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var (
	azureExample = "azure --location westus2 --subscription-id <sub_id> images/azure/centos-79.yaml"
	azureUse     = "azure <image.yaml>"
)

var azureBuildCmd = &cobra.Command{
	Use:     azureUse,
	Short:   "build and provision azure images",
	Example: azureExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runBuild(args[0])
	},
}

var azureGenerateCmd = &cobra.Command{
	Use:     azureUse,
	Short:   "generate files relating to building azure images",
	Example: azureExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runGenerate(args[0])
	},
}

func initBuildAzure() {
	fs := azureBuildCmd.Flags()

	initGenerateFlags(fs, &buildFlags.generateCLIFlags)
	initAzureFlags(fs, &buildFlags.generateCLIFlags)

	addBuildArgs(fs, &buildFlags)
}

func initGenerateAzure() {
	fs := azureGenerateCmd.Flags()

	initGenerateFlags(fs, &generateFlags)
	initAzureFlags(fs, &generateFlags)
}

func initAzureFlags(fs *flag.FlagSet, gFlags *generateCLIFlags) {
	gFlags.userArgs.Azure = &app.AzureArgs{}
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
		[]string{},
		"a list of locatins to publish the image (default same as `location`)",
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
}
