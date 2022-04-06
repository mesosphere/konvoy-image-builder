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

AIRGAPPED_BUNDLE_URL ?= konvoy-kubernetes-staging.s3.us-west-2.amazonaws.com
ARTIFACTS_DIR ?= artifacts/
DEFAULT_KUBERNETES_VERSION_SEMVER ?= $(shell \
	grep -E -e "kubernetes_version:" ansible/group_vars/all/defaults.yaml | \
	cut -d\" -f2 \
)

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

# NOTE(jkoelker) Convert the provider to an image subdir
image_dir = $(subst aws,ami,$(call provider, $(1)))

# NOTE(jkoelker) Extract the file from the last part (same as `cut -d- -f2-`)
#                and squashes major and minor, e.g 7.9 -> 79, 8.2 -> 82
image_file = $(subst .,,$(subst $(SPACE),-,$(wordlist 2, 3, $(subst -,$(SPACE),$(1)))))

$(ARTIFACTS_DIR):
	mkdir -p $(ARTIFACTS_DIR)

$(ARTIFACTS_DIR)/images:
	mkdir -p $(ARTIFACTS_DIR)/images

# TODO(jkoelker) UnPHONYify these targets
.PHONY: download-images-bundle
download-images-bundle: $(ARTIFACTS_DIR)/images
	curl -o $(ARTIFACTS_DIR)/images/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_images$(bundle_suffix).tar.gz -fsSL https://$(AIRGAPPED_BUNDLE_URL)/konvoy/airgapped/kubernetes-images/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_images$(bundle_suffix).tar.gz

.PHONY: download-os-packages-bundle
download-os-packages-bundle: $(ARTIFACTS_DIR)
	curl -o $(ARTIFACTS_DIR)/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_$(os_distribution)_$(os_distribution_major_version)_$(os_distribution_arch)$(bundle_suffix).tar.gz -fsSL https://$(AIRGAPPED_BUNDLE_URL)/konvoy/airgapped/os-packages/$(DEFAULT_KUBERNETES_VERSION_SEMVER)_$(os_distribution)_$(os_distribution_major_version)_$(os_distribution_arch)$(bundle_suffix).tar.gz

# NOTE(jkoelker) set no-op cleanup targets for providers that support `DryRun`.
.PHONY: aws-build-image-cleanup
aws-build-image-cleanup: ;

.PHONY: ova-build-image-cleanup
ova-build-image-cleanup: ;

# NOTE(jkoelker) The common build target every other target ends up calling.
.PHONY: build-image
build-image: build
build-image: $(IMAGE)
build-image: ## Build an image on a provider
	./bin/konvoy-image build $(subst ova,,$(PROVIDER)) $(IMAGE) \
	--dry-run=$(BUILD_DRY_RUN) \
	-v ${VERBOSITY} \
	$(if $(ADDITIONAL_OVERRIDES),--overrides=$(ADDITIONAL_OVERRIDES)) \
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
		ADDITIONAL_OVERRIDES=overrides/fips.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)$(ADDITIONAL_OVERRIDES)) \
		VERBOSITY=$(VERBOSITY) \
		BUILD_DRY_RUN=$(BUILD_DRY_RUN)

.PHONY: %_offline
%_offline:
	# NOTE(jkoelker) Fail fast if offline is not supported for provider
	$(MAKE) devkit.run WHAT="make packer-$(call provider,$*)-offline-override.yaml"
	$(MAKE) os_distribution=$(call os_distro,$(call distro,$*)) \
		os_distribution_major_version=$(call major_version,$(call version,$*)) \
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
		os_distribution_major_version=$(call major_version,$(call version,$*)) \
		os_distribution_arch=x86_64 \
		bundle_suffix=_fips \
		download-os-packages-bundle
	$(MAKE) pip-packages-artifacts
	$(MAKE) bundle_suffix=_fips download-images-bundle
	$(MAKE) devkit.run WHAT="make $*_fips \
		BUILD_DRY_RUN=${BUILD_DRY_RUN} \
		VERBOSITY=$(VERBOSITY) \
		ADDITIONAL_ARGS=\"$(ADDITIONAL_ARGS)\" \
		ADDITIONAL_OVERRIDES=overrides/offline-fips.yaml,packer-$(call provider,$*)-offline-override.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)${ADDITIONAL_OVERRIDES})"

.PHONY: %_nvidia
%_nvidia:
	$(MAKE) build-$* \
		ADDITIONAL_ARGS="--aws-instance-type p2.xlarge$(if $(ADDITIONAL_ARGS), $(SPACE)$(ADDITIONAL_ARGS))" \
		ADDITIONAL_OVERRIDES=overrides/nvidia.yaml$(if $(ADDITIONAL_OVERRIDES),$(COMMA)$(ADDITIONAL_OVERRIDES))
		VERBOSITY=$(VERBOSITY) \
		BUILD_DRY_RUN=$(BUILD_DRY_RUN)

# Centos 7
.PHONY: centos7
centos7:
	$(MAKE) build-aws-centos-7

.PHONY: centos7-fips
centos7-fips:
	$(MAKE) aws-centos-7_fips

.PHONY: centos7-offline
centos7-offline:
	$(MAKE) aws-centos-7_offline

.PHONY: centos7-nvidia
centos7-nvidia:
	$(MAKE) aws-centos-7_nvidia

.PHONY: flatcar
flatcar:
	$(MAKE) build-aws-flatcar ADDITIONAL_OVERRIDES=flatcar-version.yaml

.PHONY: flatcar-nvidia
flatcar-nvidia:
	$(MAKE) aws-flatcar_nvidia ADDITIONAL_OVERRIDES=flatcar-version.yaml

# Oracle 7
.PHONY: oracle7
oracle7:
	$(MAKE) build-aws-oracle-7

# Oracle 8
.PHONY: oracle8
oracle8:
	$(MAKE) build-aws-oracle-8

# RHEL 7.9
.PHONY: rhel79
rhel79:
	$(MAKE) build-aws-rhel-7.9

.PHONY: rhel79-nvidia
rhel79-nvidia:
	$(MAKE) aws-rhel-7.9_nvidia

.PHONY: rhel79-fips
rhel79-fips:
	$(MAKE) aws-rhel-7.9_fips

.PHONY: rhel79-fips-offline
rhel79-fips-offline:
	$(MAKE) aws-rhel7.9_offline-fips

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

# RHEL 8.2
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

# RHEL 8.4
.PHONY: rhel84
rhel84:
	$(MAKE) build-aws-rhel-8.4

.PHONY: rhel84-nvidia
rhel84-nvidia:
	$(MAKE) aws-rhel-8.4_nvidia

.PHONY: rhel84-fips
rhel84-fips:
	$(MAKE) aws-rhel-8.4_fips

.PHONY: rhel84-fips-offline
rhel84-fips-offline:
	$(MAKE) aws-rhel-8.4_offline-fips

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

# SLES 15
.PHONY: sles15
sles15:
	$(MAKE) build-aws-sles-15

.PHONY: sles15-nvidia
sles15-nvidia:
	$(MAKE) aws-sles-15_nvidia

# Ubuntu 18(04)
.PHONY: ubuntu18
ubuntu18:
	$(MAKE) build-aws-ubuntu-18

# Ubuntu 20(04)
.PHONY: ubuntu20
ubuntu20:
	$(MAKE) build-aws-ubuntu-20

.PHONY: ubuntu20-nvidia
ubuntu20-nvidia:
	$(MAKE) aws-ubuntu-20_nvidia
