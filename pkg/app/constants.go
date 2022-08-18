package app

const (
	AnsibleExtraVarsKey = "ansible_extra_vars"

	BuildNameKey      = "build_name"
	BuildNameExtraKey = "build_name_extra"

	BuildTypeAmazon = "amazon"
	BuildTypeAzure  = "azure"
	BuildTypeGCP    = "gcp"

	CommonConfigDefaultPath = "./images/common.yaml"
	ContainerdVersionKey    = "containerd_version"

	DefaultBuildName = "provision_build"

	GPUTypesKey      = "/gpu/types"
	GPUNvidiaVersion = "/nvidia_cuda_version"

	KubernetesVersionKey       = "kubernetes_version"
	KubernetesFullVersionKey   = "kubernetes_full_version"
	KubernetesBuildMetadataKey = "kubernetes_build_metadata"

	PackerAMIGroupsPath     = "/packer/ami_groups"
	PackerAMIRegionsPath    = "/packer/ami_regions"
	PackerAMIUsersPath      = "/packer/ami_users"
	PackerBuilderRegionPath = "/packer/aws_region"
	PackerBuilderTypePath   = "/packer_builder_type"
	PackerFilterNamePath    = "/packer/ami_filter_name"
	PackerFilterOwnerPath   = "/packer/ami_filter_owners"
	PackerInstanceTypePath  = "/packer/aws_instance_type"
	PackerSourceAMIPath     = "/packer/source_ami"

	PackerAzureClientIDPath              = "/packer/client_id"
	PackerAzureDistributionVersionPath   = "/packer/distribution_version"
	PackerAzureInstanceType              = "/packer/vm_size"
	PackerAzureGalleryLocations          = "/packer/gallery_image_locations"
	PackerAzureGalleryImageNamePath      = "/packer/gallery_image_name"
	PackerAzureGalleryImageOfferPath     = "/packer/gallery_image_offer"
	PackerAzureGalleryImagePublisherPath = "/packer/gallery_image_publisher"
	PackerAzureGalleryImageSKU           = "/packer/gallery_image_sku"
	PackerAzureLocation                  = "/packer/location"
	PackerAzureGalleryNamePath           = "/packer/gallery_name"
	PackerAzureResourceGroupNamePath     = "/packer/resource_group_name"
	PackerAzureSubscriptionIDPath        = "/packer/subscription_id"
	PackerAzureTenantIDPath              = "/packer/tenant_id"
	PackerAzureCloudEndpointPath         = "/packer/cloud_environment_name"

	PackerGCPProjectIDPath = "/packer/project_id"
	PackerGCPNetworkPath   = "/packer/network"
	PackerGCPRegionPath    = "/packer/region"

	PackerGPUPath          = "/packer/gpu"
	PackerGPUTypes         = "/packer/gpu_types"
	PackerGPUNvidiaVersion = "/packer/gpu_nvidia_version"

	AzureClientIDEnvVariable = "AZURE_CLIENT_ID"
	//nolint:gosec // environment var set by user
	AzureClientSecretEnvVariable   = "AZURE_CLIENT_SECRET"
	AzureSubscriptionIDEnvVariable = "AZURE_SUBSCRIPTION_ID"
	AzureTenantIDEnvVariable       = "AZURE_TENANT_ID"
	AzureCloudEndpointForClient    = "/packer/cloud_environment_client"

	//nolint:gosec // environment var set by user
	GCPCredentialEnvVariable = "GOOGLE_APPLICATION_CREDENTIALS"

	OutputDir = "./work"
)

var AnsibleRunsDirectory = "ansible-runs"
