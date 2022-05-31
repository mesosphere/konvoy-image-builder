package app

import (
	"context"
	"fmt"
	"os"

	"github.com/mesosphere/konvoy-image-builder/pkg/azure"
)

type AzureArgs struct {
	ClientID string

	InstanceType string

	GalleryImageLocations []string
	GalleryImageName      string
	GalleryImageOffer     string
	GalleryImagePublisher string
	GalleryImageSKU       string
	GalleryName           string

	Location          string
	ResourceGroupName string

	SubscriptionID string
	TenantID       string
}

func azureCredentials(config Config) (*azure.Credentials, error) {
	clientID, err := config.GetWithEnvironment(PackerAzureClientIDPath, AzureClientIDEnvVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to get client id: %w", err)
	}

	clientSecret, ok := os.LookupEnv(AzureClientSecretEnvVariable)
	if !ok {
		return nil, fmt.Errorf(
			"failed to get client secret (AZURE_CLIENT_SECRET): %w",
			ErrConfigRequired,
		)
	}

	tenantID, err := config.GetWithEnvironment(
		PackerAzureTenantIDPath,
		AzureTenantIDEnvVariable,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant id: %w", err)
	}

	credentials, err := azure.NewCredentials(clientID, clientSecret, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed create credentials: %w", err)
	}

	return credentials, nil
}

func azureImageDescription(config Config) (*azure.ImageDescription, error) {
	galleryName, err := config.GetWithError(PackerAzureGalleryNamePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery name: %w", err)
	}

	galleryImageName, err := config.GetWithError(PackerAzureGalleryImageNamePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery image name: %w", err)
	}

	offer, err := config.GetWithError(PackerAzureGalleryImageOfferPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery image offer: %w", err)
	}

	publisher, err := config.GetWithError(PackerAzureGalleryImagePublisherPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery image publisher: %w", err)
	}

	resourceGroupName, err := config.GetWithError(PackerAzureResourceGroupNamePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource group name: %w", err)
	}

	sku, err := config.GetWithError(PackerAzureGalleryImageSKU)
	if err != nil || sku == "" {
		// NOTE(jkoelker) fall back to mirroring the source
		sku, err = config.GetWithError(PackerAzureDistributionVersionPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get gallery image sku: %w", err)
		}

		if sku == "" {
			return nil, fmt.Errorf("failed to get gallery image sku: %w", ErrConfigRequired)
		}
	}

	// NOTE(jkoelker) Append the build extra to the sku to prevent conflicts between
	//                base images and their special flavors (e.g. centos7 and centos7-nvidia)
	sku = config.AddBuildNameExtra(sku)

	description, err := azure.NewImageDescription(
		galleryName,
		galleryImageName,
		offer,
		publisher,
		resourceGroupName,
		sku,
	)
	if err != nil {
		return nil, fmt.Errorf("failed create image description: %w", err)
	}

	return description, nil
}

func ensureAzure(config Config) error {
	credentials, err := azureCredentials(config)
	if err != nil {
		return fmt.Errorf("failed get credentials: %w", err)
	}

	description, err := azureImageDescription(config)
	if err != nil {
		return fmt.Errorf("failed get image description: %w", err)
	}

	locations, err := config.GetSliceWithError(PackerAzureGalleryLocations)
	if err != nil {
		return fmt.Errorf("failed to get location: %w", err)
	}

	subscriptionID, err := config.GetWithEnvironment(
		PackerAzureSubscriptionIDPath,
		AzureSubscriptionIDEnvVariable,
	)
	if err != nil {
		return fmt.Errorf("failed to get subscription id: %w", err)
	}

	ctx := context.Background()

	if err = azure.EnsureImageDescriptions(
		ctx,
		credentials,
		description,
		locations,
		subscriptionID,
	); err != nil {
		return fmt.Errorf("failed to ensure azure image descriptions: %w", err)
	}

	return nil
}
