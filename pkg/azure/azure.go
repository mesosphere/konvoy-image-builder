//

package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type ImageDescription struct {
	GalleryName       string
	GalleryImageName  string
	Offer             string
	Publisher         string
	ResourceGroupName string
	SKU               string
}

func NewImageDescription(
	galleryName string,
	galleryImageName string,
	offer string,
	publisher string,
	resourceGroupName string,
	sku string,
) (*ImageDescription, error) {
	return &ImageDescription{
		GalleryName:       galleryName,
		GalleryImageName:  galleryImageName,
		Offer:             offer,
		Publisher:         publisher,
		ResourceGroupName: resourceGroupName,
		SKU:               sku,
	}, nil
}

type Credentials struct {
	ID          string
	Secret      string
	TenantID    string
	CloudConfig cloud.Configuration
}

func NewCredentials(clientID, clientSecret, tenantID string, CloudConfig cloud.Configuration) (*Credentials, error) {
	return &Credentials{
		ID:          clientID,
		Secret:      clientSecret,
		TenantID:    tenantID,
		CloudConfig: CloudConfig,
	}, nil
}

func createGalleryImage(
	ctx context.Context,
	cred azcore.TokenCredential,
	description *ImageDescription,
	location string,
	subscriptionID string,
	options *arm.ClientOptions,
) (*armcompute.GalleryImage, error) {
	galleryImageClient, err := armcompute.NewGalleryImagesClient(subscriptionID, cred, options)
	if err != nil {
		return nil, fmt.Errorf("failed to azure gallery image client: %w", err)
	}

	pollerResp, err := galleryImageClient.BeginCreateOrUpdate(
		ctx,
		description.ResourceGroupName,
		description.GalleryName,
		description.GalleryImageName,
		armcompute.GalleryImage{
			Location: to.Ptr(location),
			Properties: &armcompute.GalleryImageProperties{
				OSType:           to.Ptr(armcompute.OperatingSystemTypesLinux),
				OSState:          to.Ptr(armcompute.OperatingSystemStateTypesGeneralized),
				HyperVGeneration: to.Ptr(armcompute.HyperVGenerationV2),
				Identifier: &armcompute.GalleryImageIdentifier{
					Offer:     to.Ptr(description.Offer),
					Publisher: to.Ptr(description.Publisher),
					SKU:       to.Ptr(description.SKU),
				},
			},
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create gallery image %s in %s: %w",
			description.GalleryImageName,
			description.GalleryName,
			err,
		)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating gallery image %s in %s: %w",
			description.GalleryImageName,
			description.GalleryName,
			err,
		)
	}

	return &resp.GalleryImage, nil
}

func createGallery(
	ctx context.Context,
	cred azcore.TokenCredential,
	description *ImageDescription,
	location string,
	subscriptionID string,
	options *arm.ClientOptions,
) (*armcompute.Gallery, error) {
	galleriesClient, err := armcompute.NewGalleriesClient(subscriptionID, cred, options)
	if err != nil {
		return nil, fmt.Errorf("failed to azure galleries client: %w", err)
	}

	pollerResp, err := galleriesClient.BeginCreateOrUpdate(
		ctx,
		description.ResourceGroupName,
		description.GalleryName,
		armcompute.Gallery{
			Location: to.Ptr(location),
			Properties: &armcompute.GalleryProperties{
				Description: to.Ptr("DKP Gallery."),
			},
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create gallery %s: %w",
			description.GalleryName,
			err,
		)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating gallery %s: %w",
			description.GalleryName,
			err,
		)
	}

	return &resp.Gallery, nil
}

func createResourceGroup(
	ctx context.Context,
	cred azcore.TokenCredential,
	description *ImageDescription,
	location string,
	subscriptionID string,
	options *arm.ClientOptions,
) (*armresources.ResourceGroup, error) {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, options)
	if err != nil {
		return nil, fmt.Errorf("failed to azure resource groups client: %w", err)
	}

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		description.ResourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(location),
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating resource group %s: %w",
			description.ResourceGroupName,
			err,
		)
	}

	return &resourceGroupResp.ResourceGroup, nil
}

func EnsureImageDescriptions(
	ctx context.Context,
	credentials *Credentials,
	description *ImageDescription,
	locations []string,
	subscriptionID string,
) error {
	cloudConfig := credentials.CloudConfig
	options := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloudConfig,
		},
	}
	cred, err := azidentity.NewClientSecretCredential(
		credentials.TenantID,
		credentials.ID,
		credentials.Secret,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}

	for _, location := range locations {
		_, err = createResourceGroup(ctx, cred, description, location, subscriptionID, options)
		if err != nil {
			return fmt.Errorf(
				"failed to create resource group %s in %s: %w",
				description.ResourceGroupName,
				location,
				err,
			)
		}

		_, err = createGallery(ctx, cred, description, location, subscriptionID, options)
		if err != nil {
			return fmt.Errorf(
				"failed to create image gallery %s in %s: %w",
				description.GalleryName,
				location,
				err,
			)
		}

		_, err = createGalleryImage(
			ctx,
			cred,
			description,
			location,
			subscriptionID,
			options,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to create image gallery %s in %s: %w",
				description.GalleryName,
				location,
				err,
			)
		}
	}

	return nil
}
