package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

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
	CloudEndpoint  *AzureCloudFlag
}

type AzureCloudFlag struct {
	Endpoint AzureCloudEndpoint
}

type AzureCloudEndpoint string

var (
	AzureCloudEndpointChina      AzureCloudEndpoint = "China"
	AzureCloudEndpointPublic     AzureCloudEndpoint = "Public"
	AzureCloudEndpointGovernment AzureCloudEndpoint = "USGovernment"
)

func (a *AzureCloudFlag) String() string {
	return string(a.Endpoint)
}

func (a *AzureCloudFlag) Set(s string) error {
	switch s {
	case "China":
		a.Endpoint = AzureCloudEndpointChina
		return nil
	case "USGovernment":
		a.Endpoint = AzureCloudEndpointGovernment
		return nil
	case "Public":
		a.Endpoint = AzureCloudEndpointPublic
		return nil
	default:
		return fmt.Errorf("flag must be set to one of %v", ListAzureEndpoints())
	}
}

func (a *AzureCloudFlag) Type() string {
	return "string"
}

func ListAzureEndpoints() []AzureCloudEndpoint {
	return []AzureCloudEndpoint{
		AzureCloudEndpointPublic,
		AzureCloudEndpointGovernment,
		AzureCloudEndpointChina,
	}
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

	endpoint, err := config.GetWithError(PackerAzureCloudEndpointPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud endpoint: %w", err)
	}
	var cloudConfig cloud.Configuration
	switch endpoint {
	case string(AzureCloudEndpointChina):
		cloudConfig = cloud.AzureChina
	case string(AzureCloudEndpointGovernment):
		cloudConfig = cloud.AzureGovernment

	default:
		cloudConfig = cloud.AzurePublic
	}
	credentials, err := azure.NewCredentials(clientID, clientSecret, tenantID, cloudConfig)
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
	if err != nil {
		return nil, fmt.Errorf("failed to get azure gallery image sku: %w", err)
	}

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

	location, err := config.GetWithError(PackerAzureLocation)
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
		location,
		subscriptionID,
	); err != nil {
		return fmt.Errorf("failed to ensure azure image descriptions: %w", err)
	}

	return nil
}
