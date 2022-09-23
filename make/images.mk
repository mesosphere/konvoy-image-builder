#
# Image Building targets
#

# BUILD_DRY_RUN determines the value of the --dry-run flag of the build command. Should be 'true' or 'false'.
BUILD_DRY_RUN ?= true
ifeq ($(BUILD_DRY_RUN),true)
$(warning Warning: BUILD_DRY_RUN is true)
endif
VERBOSITY ?= 0
COMMA := ,
NULL :=
SPACE := $(NULL) $(NULL)

AIRGAPPED_BUNDLE_URL_PREFIX ?= downloads.d2iq.com/dkp
ARTIFACTS_DIR ?= artifacts/
CONTAINERD_URL ?= https://packages.d2iq.com/dkp/containerd
NVIDIA_URL ?= https://download.nvidia.com/XFree86/Linux-x86_64

NVIDIA_DRIVER_VERSION ?= $(shell \
	grep -E -e "nvidia_driver_version:" ansible/group_vars/all/defaults.yaml | \
	cut -d\" -f2 \
)


DEFAULT_KUBERNETES_VERSION_SEMVER ?= $(shell \
	grep -E -e "kubernetes_version:" ansible/group_vars/all/defaults.yaml | \
	cut -d\" -f2 \
)

DEFAULT_CONTAINERD_VERSION ?= $(shell \
	grep -E -e "containerd_version:" ansible/group_vars/all/defaults.yaml | \
	cut -d\" -f2 \
)

K8S_VER_MINOR := $(shell echo $(DEFAULT_KUBERNETES_VERSION_SEMVER) | cut -f2 -d.)
K8S_MINOR_LT_24 := $(shell [ $(K8S_VER_MINOR) -lt 24 ] && echo true)

# NOTE(jkoelker) Extract the provider as the first part (same as `cut -d- -f1`)
provider = $(firstword $(subst -,$(SPACE),$(1)))

#NOTE(jkoelker) Extract the distro as the second part (same as `cut -d- -f2`)
distro = $(wordlist 2,2,$(subst -,$(SPACE),$(1)))

#NOTE(jkoelker) Extract the version as the third part (same as `cut -d- -f3`)
version = $(wordlist 3,3,$(subst -,$(SPACE),$(1)))

#NOTE(jkoelker) Extract the major as the first part (same as `cut -d. -f1`)
major_version = $(firstword $(subst .,$(SPACE),$(1)))


#NOTE(jkoelker) Convert the distro to the package bundle distro name
os_distro = $(subst rhel,redhat,$(1))

os_distro_os_release = $(subst oracle,ol,$(subst redhat,rhel,$(1)))

# NOTE(jkoelker) Convert the provider to an image subdir
image_dir = $(subst aws,ami,$(call provider, $(1)))

# NOTE(jkoelker) Extract the file from the last part (same as `cut -d- -f2-`)
#                and squashes major and minor, e.g 7.9 -> 79, 8.2 -> 82
image_file = $(subst .,,$(subst $(SPACE),-,$(wordlist 2, 3, $(subst -,$(SPACE),$(1)))))

azure_vm_size = --instance-type Standard_B2ms
# NOTE(jkoelker) Set the VM Size argument for the provider if not already
#                in the ADDITIONAL_ARGS.
vm_size = $(if $(findstring --instance-type,$(2)),,$($(1)_vm_size))

$(ARTIFACTS_DIR):
	mkdir -p $(ARTIFACTS_DIR)

$(ARTIFACTS_DIR)/images:
	mkdir -p $(ARTIFACTS_DIR)/images

# TODO(jkoelker) UnPHONYify these targets
.PHONY: download-images-bundle
download-images-bundle: $(ARTIFACTS_DIR)/images
ifeq ($(K8S_MINOR_LT_24),true)
	curl -o $(ARTIFACTS_DIR)/images/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_images$(bundle_suffix).tar.gz -fsSL https://$(AIRGAPPED_BUNDLE_URL_PREFIX)/airgapped/kubernetes-images/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_images$(bundle_suffix).tar.gz
else
	curl -o $(ARTIFACTS_DIR)/images/kubernetes-images-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-d2iq.1$(bundle_suffix).tar -fsSL https://$(AIRGAPPED_BUNDLE_URL_PREFIX)/airgapped/kubernetes-images/kubernetes-images-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-d2iq.1$(bundle_suffix).tar
endif
.PHONY: download-os-packages-bundle
download-os-packages-bundle: $(ARTIFACTS_DIR)
	curl -o $(ARTIFACTS_DIR)/containerd-$(DEFAULT_CONTAINERD_VERSION)-d2iq.1-$(os_distribution_os_release)-$(os_distribution_major_minor_version)-$(os_distribution_arch)$(bundle_suffix).tar.gz -fsSL $(CONTAINERD_URL)/containerd-$(DEFAULT_CONTAINERD_VERSION)-d2iq.1-$(os_distribution_os_release)-$(os_distribution_major_minor_version)-$(os_distribution_arch)$(bundle_suffix).tar.gz
	curl -o $(ARTIFACTS_DIR)/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_$(os_distribution)_$(os_distribution_major_version)_$(os_distribution_arch)$(bundle_suffix).tar.gz -fsSL https://$(AIRGAPPED_BUNDLE_URL_PREFIX)/airgapped/os-packages/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_$(os_distribution)_$(os_distribution_major_version)_$(os_distribution_arch)$(bundle_suffix).tar.gz

.PHONY: download-nvidia-runfile
download-nvidia-runfile: $(ARTIFACTS_DIR)
	curl -o $(ARTIFACTS_DIR)/NVIDIA-Linux-x86_64-$(NVIDIA_DRIVER_VERSION).run -fsSL $(NVIDIA_URL)/$(NVIDIA_DRIVER_VERSION)/NVIDIA-Linux-x86_64-$(NVIDIA_DRIVER_VERSION).run

# NOTE(jkoelker) set no-op cleanup targets for providers that support `DryRun`.
.PHONY: aws-build-image-cleanup
aws-build-image-cleanup: ;

.PHONY: ova-build-image-cleanup
ova-build-image-cleanup: ;

.PHONY: gcp-build-image-cleanup
gcp-build-image-cleanup: ;

.PHONY: azure-build-image-cleanup
azure-build-image-cleanup: ;

# NOTE(jkoelker) The common build target every other target ends up calling.
.PHONY: build-image
build-image: build
build-image: $(IMAGE)
build-image: ## Build an image on a provider
	./bin/konvoy-image build $(subst ova,,$(PROVIDER)) $(IMAGE) \
	--dry-run=$(BUILD_DRY_RUN) \
	-v ${VERBOSITY} \
	$(if $(ADDITIONAL_OVERRIDES),--overrides=$(ADDITIONAL_OVERRIDES)) \
	$(call vm_size,$(PROVIDER),$(ADDITIONAL_ARGS)) \
	$(ADDITIONAL_ARGS)
	$(MAKE) $(PROVIDER)-build-image-cleanup

# NOTE(jkoelker) Parses the `PROVIDER` and `IMAGE` from the target name. E.g
#                `build-aws-centos-8.4` will set `PROVIDER=aws` and
#                `IMAGE=images/ami/centos-84.yaml1.
.PHONY: build-%
build-%:
	$(MAKE) build-image \
		PROVIDER=$(call provider,$*) \
		ADDITIONAL_OVERRIDES=$(ADDITIONAL_OVERRIDES) \
		ADDITIONAL_ARGS="$(ADDITIONAL_ARGS)" \
		IMAGE=images/$(call image_dir,$*)/$(call image_file,$*).yaml \
		VERBOSITY=$(VERBOSITY) \
		BUILD_DRY_RUN=$(BUILD_DRY_RUN)

.PHONY: %_fips
%_fips:
	$(MAKE) build-$* \
		ADDITIONAL_ARGS="$(ADDITIONAL_ARGS)" \
		ADDITIONAL_OVERRIDES=overrides/fips.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)${ADDITIONAL_OVERRIDES}) \
		VERBOSITY=$(VERBOSITY) \
		BUILD_DRY_RUN=$(BUILD_DRY_RUN)

.PHONY: %_offline
%_offline:
	# NOTE(jkoelker) Fail fast if offline is not supported for provider
	$(MAKE) devkit.run WHAT="make packer-$(call provider,$*)-offline-override.yaml"
	$(MAKE) os_distribution=$(call os_distro,$(call distro,$*)) \
		os_distribution_major_version=$(call major_version,$(call version,$*)) \
		os_distribution_os_release=$(call os_distro_os_release,$(call distro,$*)) \
		os_distribution_major_minor_version=$(call version,$*) \
		os_distribution_arch=x86_64 \
		bundle_suffix= \
		download-os-packages-bundle
	$(MAKE) pip-packages-artifacts
	$(MAKE) bundle_suffix= download-images-bundle
	$(MAKE) devkit.run WHAT="make build-$* \
		BUILD_DRY_RUN=$(BUILD_DRY_RUN) \
		VERBOSITY=$(VERBOSITY) \
		ADDITIONAL_ARGS=\"$(ADDITIONAL_ARGS)\" \
		ADDITIONAL_OVERRIDES=overrides/offline.yaml,packer-$(call provider, $*)-offline-override.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)$(ADDITIONAL_OVERRIDES))"

.PHONY: %_offline-fips
%_offline-fips:
	$(MAKE) devkit.run WHAT="make packer-$(call provider,$*)-offline-override.yaml"
	$(MAKE) os_distribution=$(call os_distro,$(call distro,$*)) \
		os_distribution_os_release=$(call os_distro_os_release,$(call distro,$*)) \
		os_distribution_major_minor_version=$(call version,$*) \
		os_distribution_major_version=$(call major_version,$(call version,$*)) \
		os_distribution_arch=x86_64 \
		bundle_suffix=_fips \
		download-os-packages-bundle
	$(MAKE) pip-packages-artifacts
	$(MAKE) download-images-bundle bundle_suffix=$$( if [ $$(echo "$(DEFAULT_KUBERNETES_VERSION_SEMVER)" | cut -f2 -d.) -lt 24 ];then echo "_fips"; else echo "-fips";fi )
	$(MAKE) devkit.run WHAT="make $*_fips \
		BUILD_DRY_RUN=${BUILD_DRY_RUN} \
		VERBOSITY=$(VERBOSITY) \
		ADDITIONAL_ARGS=\"$(ADDITIONAL_ARGS)\" \
		ADDITIONAL_OVERRIDES=overrides/offline-fips.yaml,packer-$(call provider,$*)-offline-override.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)${ADDITIONAL_OVERRIDES})"

.PHONY: %_offline-nvidia
%_offline-nvidia:
	$(MAKE) devkit.run WHAT="make packer-$(call provider,$*)-offline-override.yaml"
	$(MAKE) os_distribution=$(call os_distro,$(call distro,$*)) \
		os_distribution_os_release=$(call os_distro_os_release,$(call distro,$*)) \
		os_distribution_major_minor_version=$(call version,$*) \
		os_distribution_major_version=$(call major_version,$(call version,$*)) \
		os_distribution_arch=x86_64 \
		download-os-packages-bundle
	$(MAKE) download-nvidia-runfile
	$(MAKE) pip-packages-artifacts
	$(MAKE) download-images-bundle
	$(MAKE) devkit.run WHAT="make build-$* \
		BUILD_DRY_RUN=${BUILD_DRY_RUN} \
		VERBOSITY=$(VERBOSITY) \
		ADDITIONAL_ARGS="--instance-type=g4dn.2xlarge$(if $(ADDITIONAL_ARGS),$(SPACE)$(ADDITIONAL_ARGS))" \
		ADDITIONAL_OVERRIDES=overrides/offline.yaml,overrides/offline-nvidia.yaml,packer-$(call provider,$*)-offline-override.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)${ADDITIONAL_OVERRIDES})"

.PHONY: %_nvidia
%_nvidia:
	$(MAKE) build-$* \
		ADDITIONAL_ARGS="--instance-type p2.xlarge$(if $(ADDITIONAL_ARGS),$(SPACE)$(ADDITIONAL_ARGS))" \
		ADDITIONAL_OVERRIDES=overrides/nvidia.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)$(ADDITIONAL_OVERRIDES))
		VERBOSITY=$(VERBOSITY) \
		BUILD_DRY_RUN=$(BUILD_DRY_RUN)

# Centos 7 AWS
#
.PHONY: centos79-offline-nvidia
centos79-offline-nvidia:
	$(MAKE) aws-centos-7.9_offline-nvidia

.PHONY: centos79
centos79:
	$(MAKE) build-aws-centos-7.9

.PHONY: centos79-offline
centos7-offline:
	$(MAKE) aws-centos-7.9_offline

.PHONY: centos79-fips
centos79-fips:
	$(MAKE) aws-centos-7.9_fips

.PHONY: centos79-nvidia
centos7-nvidia:
	$(MAKE) aws-centos-7.9_nvidia

.PHONY: centos7
centos7: centos79

.PHONY: centos7-fips
centos7-fips: centos79-fips

.PHONY: centos7-offline
centos7-offline: centos79-offline

.PHONY: centos7-nvidia
centos7-nvidia: centos79-nvidia

# Centos 7 GCP
#
.PHONY: centos79
centos79-gcp:
	$(MAKE) build-gcp-centos-7.9

.PHONY: centos79-fips-gcp
centos79-fips-gcp:
	$(MAKE) gcp-centos-7.9_fips

# Centos 7 Azure
#

.PHONY: centos79-azure
centos79-azure:
	$(MAKE) build-azure-centos-7.9

.PHONY: centos7-azure
centos7-azure: centos79-azure

.PHONY: centos79-fips-azure
centos79-fips-azure:
	$(MAKE) azure-centos-7.9_fips

.PHONY: centos7-fips-azure
centos7-fips-azure: centos79-fips-azure

.PHONY: centos7-offline-azure
centos79-offline-azure:
	$(MAKE) azure-centos-7.9_offline

.PHONY: centos7-offline-azure
centos7-offline-azure: centos79-offline-azure

.PHONY: centos79-nvidia-azure
centos7-nvidia-azure:
	$(MAKE) azure-centos-7.9_nvidia

.PHONY: centos7-nvidia-azure
centos7-nvidia-azure: centos79-nvidia-azure

.PHONY: flatcar
flatcar:
	$(MAKE) build-aws-flatcar

.PHONY: flatcar-nvidia
flatcar-nvidia:
	$(MAKE) aws-flatcar_nvidia

# Flatcar Azure
.PHONY: flatcar-azure
flatcar-azure:
	$(MAKE) build-azure-flatcar

.PHONY: flatcar-nvidia-azure
flatcar-nvidia-azure:
	$(MAKE) azure-flatcar_nvidia

# Oracle 7 AWS
.PHONY: oracle79
oracle79:
	$(MAKE) build-aws-oracle-7.9

# old make targets to map to new ones
.PHONY: oracle7
oracle7: oracle79

# Oracle 7 Azure
.PHONY: oracle79-azure
oracle79-azure:
	$(MAKE) build-azure-oracle-7.9

# old make targets to map to new ones
.PHONY: oracle7-azure
oracle7-azure: oracle79-azure

# Oracle 8 AWS
.PHONY: oracle84
oracle84:
	$(MAKE) build-aws-oracle-8.4

# old make targets to map to new ones
.PHONY: oracle8
oracle8: oracle84

# Oracle 8 Azure
.PHONY: oracle84-azure
oracle84-azure:
	$(MAKE) build-azure-oracle-8.4

# old make targets to map to new ones
.PHONY: oracle8-azure
oracle8: oracle84-azure

# RHEL 7.9 AWS
.PHONY: rhel79
rhel79:
	$(MAKE) build-aws-rhel-7.9

.PHONY: rhel79-nvidia
rhel79-nvidia:
	$(MAKE) aws-rhel-7.9_nvidia

.PHONY: rhel79-offline-nvidia
rhel79-offline-nvidia:
	$(MAKE) aws-rhel-7.9_offline-nvidia

.PHONY: rhel79-fips
rhel79-fips:
	$(MAKE) aws-rhel-7.9_fips

.PHONY: rhel79-fips-offline
rhel79-fips-offline:
	$(MAKE) aws-rhel-7.9_offline-fips

# RHEL 7.9 Azure
.PHONY: rhel79-azure
rhel7-azure:
	$(MAKE) build-azure-rhel-7.9

.PHONY: rhel79-nvidia-azure
rhel7-nvidia-azure:
	$(MAKE) azure-rhel-7.9_nvidia

.PHONY: rhel79-fips-azure
rhel7-fips-azure:
	$(MAKE) azure-rhel-7.9_fips

.PHONY: rhel79-fips-offline-azure
rhel7-fips-offline-azure:
	$(MAKE) azure-rhel7.9_offline-fips

# old make targets to map to new ones
.PHONY: rhel7-azure
rhel7-azure: rhel79-azure

.PHONY: rhel7-nvidia-azure
rhel7-nvidia-azure: rhel79-nvidia-azure

.PHONY: rhel7-fips-azure
rhel7-fips-azure: rhel79-fips-azure

.PHONY: rhel7-fips-offline-azure
rhel7-fips-offline-azure: rhel79-fips-offline-azure

# RHEL 7.9 vSphere
.PHONY: rhel79-ova
rhel79-ova:
	$(MAKE) build-ova-rhel-7.9

.PHONY: rhel79-ova-offline
rhel79-ova-offline:
	$(MAKE) ova-rhel-7.9_offline

.PHONY: rhel79-ova-fips
rhel79-ova-fips:
	$(MAKE) ova-rhel-7.9_fips

.PHONY: rhel79-ova-fips-offline
rhel79-ova-fips-offline:
	$(MAKE) ova-rhel-7.9_offline-fips

# RHEL 8.2 AWS
.PHONY: rhel82
rhel82:
	$(MAKE) build-aws-rhel-8.2

.PHONY: rhel82-nvidia
rhel82-nvidia:
	$(MAKE) aws-rhel-8.2_nvidia

.PHONY: rhel82-fips
rhel82-fips:
	$(MAKE) aws-rhel-8.2_fips

.PHONY: rhel82-fips-offline
rhel82-fips-offline:
	$(MAKE) aws-rhel-8.2_offline-fips

# RHEL 8.4 AWS
.PHONY: rhel84
rhel84:
	$(MAKE) build-aws-rhel-8.4

.PHONY: rhel84-nvidia
rhel84-nvidia:
	$(MAKE) aws-rhel-8.4_nvidia

.PHONY: rhel84-offline-nvidia
rhel84-offline-nvidia:
	$(MAKE) aws-rhel-8.4_offline-nvidia

.PHONY: rhel84-fips
rhel84-fips:
	$(MAKE) aws-rhel-8.4_fips

.PHONY: rhel84-fips-offline
rhel84-fips-offline:
	$(MAKE) aws-rhel-8.4_offline-fips

# RHEL 8 Azure
.PHONY: rhel84-azure
rhel84-azure:
	$(MAKE) build-azure-rhel-8.4

.PHONY: rhel84-nvidia-azure
rhel84-nvidia-azure:
	$(MAKE) azure-rhel-8.4_nvidia

.PHONY: rhel84-fips-azure
rhel84-fips-azure:
	$(MAKE) azure-rhel-8.4_fips

.PHONY: rhel84-fips-offline-azure
rhel84-fips-offline-azure:
	$(MAKE) azure-rhel-8.4_offline-fips

# old make targets so this continues to work
.PHONY: rhel8-azure
rhel8-azure: rhel84-azure

.PHONY: rhel8-nvidia-azure
rhel8-nvidia-azure: rhel84-nvidia-azure

.PHONY: rhel8-fips-azure
rhel8-fips-azure: rhel84-fips-azure

.PHONY: rhel8-fips-offline-azure
rhel8-fips-offline-azure: rhel84-fips-offline-azure

# RHEL 8.4 vSphere
.PHONY: rhel84-ova
rhel84-ova:
	$(MAKE) build-ova-rhel-8.4

.PHONY: rhel84-ova-offline
rhel84-ova-offline:
	$(MAKE) ova-rhel-8.4_offline

.PHONY: rhel84-ova-fips
rhel84-ova-fips:
	$(MAKE) ova-rhel-8.4_fips

.PHONY: rhel82-ova-fips-offline
rhel84-ova-fips-offline:
	$(MAKE) ova-rhel-8.4_offline-fips

# SLES 15 AWS
.PHONY: sles15
sles15:
	$(MAKE) build-aws-sles-15

.PHONY: sles15-nvidia
sles15-nvidia:
	$(MAKE) aws-sles-15_nvidia

# SLES 15 Azure
.PHONY: sles15-azure
sles15-azure:
	$(MAKE) build-azure-sles-15

.PHONY: sles15-nvidia-azure
sles15-nvidia-azure:
	$(MAKE) azure-sles-15_nvidia

# Ubuntu 18(04) AWS
.PHONY: ubuntu18
ubuntu18:
	$(MAKE) build-aws-ubuntu-18

# Ubuntu 18(04) Azure
# preserve old make targets for CI
.PHONY: ubuntu18-azure
ubuntu18-azure: ubuntu1804-azure

.PHONY: ubuntu20-azure
ubuntu20-azure: ubuntu2004-azure
.PHONY: ubuntu1804-azure
ubuntu1804-azure:
	$(MAKE) build-azure-ubuntu-1804

# Ubuntu 18(04) GCP
.PHONY: ubuntu1804-gcp
ubuntu1804-gcp:
	$(MAKE) build-gcp-ubuntu-1804

# Ubuntu 20(04) AWS
.PHONY: ubuntu20
ubuntu20:
	$(MAKE) build-aws-ubuntu-20

.PHONY: ubuntu20-nvidia
ubuntu20-nvidia:
	$(MAKE) aws-ubuntu-20_nvidia

# Ubuntu 20(04) Azure
.PHONY: ubuntu2004-azure
ubuntu2004-azure:
	$(MAKE) build-azure-ubuntu-2004

# Ubuntu 20(04) GCP
.PHONY: ubuntu2004-gcp
ubuntu2004-gcp:
	$(MAKE) build-gcp-ubuntu-2004

.PHONY: ubuntu20-nvidia-azure
ubuntu20-nvidia-azure:
	$(MAKE) azure-ubuntu-20_nvidia
