//

package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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
	ID       string
	Secret   string
	TenantID string
}

func NewCredentials(clientID string, clientSecret string, tenantID string) (*Credentials, error) {
	return &Credentials{
		ID:       clientID,
		Secret:   clientSecret,
		TenantID: tenantID,
	}, nil
}

func createGalleryImage(
	ctx context.Context,
	cred azcore.TokenCredential,
	description *ImageDescription,
	location string,
	subscriptionID string,
) (*armcompute.GalleryImage, error) {
	galleryImageClient := armcompute.NewGalleryImagesClient(subscriptionID, cred, nil)

	pollerResp, err := galleryImageClient.BeginCreateOrUpdate(
		ctx,
		description.ResourceGroupName,
		description.GalleryName,
		description.GalleryImageName,
		armcompute.GalleryImage{
			Location: to.StringPtr(location),
			Properties: &armcompute.GalleryImageProperties{
				OSType:           armcompute.OperatingSystemTypesLinux.ToPtr(),
				OSState:          armcompute.OperatingSystemStateTypesGeneralized.ToPtr(),
				HyperVGeneration: armcompute.HyperVGenerationV2.ToPtr(),
				Identifier: &armcompute.GalleryImageIdentifier{
					Offer:     to.StringPtr(description.Offer),
					Publisher: to.StringPtr(description.Publisher),
					SKU:       to.StringPtr(description.SKU),
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

	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
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
) (*armcompute.Gallery, error) {
	galleriesClient := armcompute.NewGalleriesClient(subscriptionID, cred, nil)

	pollerResp, err := galleriesClient.BeginCreateOrUpdate(
		ctx,
		description.ResourceGroupName,
		description.GalleryName,
		armcompute.Gallery{
			Location: to.StringPtr(location),
			Properties: &armcompute.GalleryProperties{
				Description: to.StringPtr("DKP Gallery."),
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

	resp, err := pollerResp.PollUntilDone(ctx, 10*time.Second)
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
) (*armresources.ResourceGroup, error) {
	resourceGroupClient := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)

	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		description.ResourceGroupName,
		armresources.ResourceGroup{
			Location: to.StringPtr(location),
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
		_, err = createResourceGroup(ctx, cred, description, location, subscriptionID)
		if err != nil {
			return fmt.Errorf(
				"failed to create resource group %s in %s: %w",
				description.ResourceGroupName,
				location,
				err,
			)
		}

		_, err = createGallery(ctx, cred, description, location, subscriptionID)
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
