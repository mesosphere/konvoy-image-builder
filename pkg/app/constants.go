package app

const (
	AnsibleExtraVarsKey = "ansible_extra_vars"

	BuildNameKey      = "build_name"
	BuildNameExtraKey = "build_name_extra"

	CommonConfigDefaultPath = "./images/common.yaml"
	ContainerdVersionKey    = "containerd_version"

	DefaultBuildName = "provision_build"

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

	OutputDir = "./work"
)

var AnsibleRunsDirectory = "ansible-runs"
