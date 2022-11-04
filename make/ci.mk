#
# CI targets
#

# requires ANSIBLE_PATH, otherwise run `make ci.e2e.ansible`
e2e.ansible:
	make -C test/e2e/ansible e2e

ifeq ($(CI), true)
export DOCKER_DEVKIT_AWS_ARGS := --env AWS_ACCESS_KEY_ID --env AWS_SECRET_ACCESS_KEY
endif

# Run every E2E test in its own devkit container.
# All tests run in parallel. Adjust parallelism with --jobs.
# Output is interleaved when run in parallel. Use --output-sync=recurse to serialize output.
ci.e2e.build.all: ci.e2e.build.centos-7
ci.e2e.build.all: ci.e2e.build.ubuntu-18
ci.e2e.build.all: ci.e2e.build.ubuntu-20
ci.e2e.build.all: ci.e2e.build.rhel-8.6
ci.e2e.build.all: ci.e2e.build.sles-15
ci.e2e.build.all: ci.e2e.build.oracle-7
ci.e2e.build.all: ci.e2e.build.oracle-8
ci.e2e.build.all: ci.e2e.build.flatcar
ci.e2e.build.all: e2e.build.centos-7-offline
ci.e2e.build.all: e2e.build.rhel-7.9-offline-fips
ci.e2e.build.all: e2e.build.rhel-8.4-offline-fips
ci.e2e.build.all: e2e.build.rhel-8.6-offline
ci.e2e.build.all: ci.e2e.build.rhel-8.4-nvidia
ci.e2e.build.all: ci.e2e.build.rhel-8.6-nvidia
ci.e2e.build.all: ci.e2e.build.centos-7-nvidia
ci.e2e.build.all: ci.e2e.build.sles-15-nvidia
ci.e2e.build.all: ci.e2e.build.rhel-8.4-ova
ci.e2e.build.all: ci.e2e.build.rhel-7.9-ova

# Run an E2E test in its own devkit container.
ci.e2e.build.%:
	$(MAKE) e2e.build.$* ADDITIONAL_OVERRIDES=$(ADDITIONAL_OVERRIDES) VERBOSITY=$(VERBOSITY) ADDITIONAL_ARGS=$(ADDITIONAL_ARGS)

# AWS
e2e.build.centos-7: centos7

e2e.build.centos-7-offline: centos7-offline infra.aws.destroy

e2e.build.centos-7-offline-nvidia: centos79-offline-nvidia infra.aws.destroy

e2e.build.rhel-7.9-offline-fips: rhel79-fips-offline infra.aws.destroy

e2e.build.rhel-7.9-offline-nvidia: rhel79-offline-nvidia infra.aws.destroy

e2e.build.rhel-8.4-offline-fips: rhel84-fips-offline infra.aws.destroy

e2e.build.rhel-8.4-offline-nvidia: rhel84-offline-nvidia infra.aws.destroy

e2e.build.rhel-8.6-offline-nvidia: rhel86-offline-nvidia infra.aws.destroy

e2e.build.rhel-8.6-offline-fips: rhel86-fips-offline infra.aws.destroy

e2e.build.ubuntu-18: ubuntu18

e2e.build.ubuntu-20: ubuntu20

e2e.build.ubuntu-20-nvidia: ubuntu20-nvidia

e2e.build.sles-15: sles15

e2e.build.oracle-7: oracle7

e2e.build.oracle-8: oracle8

e2e.build.flatcar: flatcar

e2e.build.centos-7-nvidia: centos7-nvidia

e2e.build.sles-15-nvidia: sles15-nvidia

e2e.build.rhel-8.4-nvidia: rhel84-nvidia

e2e.build.rhel-8.6: rhel86

e2e.build.rhel-8.6-offline: rhel86-offline infra.aws.destroy

e2e.build.rhel-8.6-nvidia: rhel86-nvidia

# Azure
e2e.build.centos-7-azure: centos7-azure

e2e.build.flatcar-azure: flatcar-azure

e2e.build.oracle-7-azure: oracle7-azure

e2e.build.oracle-8-azure: oracle8-azure

e2e.build.sles-15-azure: sles15-azure

e2e.build.rhel-7-fips-azure: rhel7-fips-azure

e2e.build.rhel-8-fips-azure: rhel8-fips-azure

e2e.build.ubuntu-18-azure: ubuntu1804-azure

e2e.build.ubuntu-20-azure: ubuntu2004-azure

# vSphere
e2e.build.rhel-8.4-ova: rhel84-ova

e2e.build.rhel-7.9-ova: rhel79-ova

e2e.build.rhel-8.6-ova: rhel86-ova

# GCP
e2e.build.centos-7.9-gcp: centos79-gcp

e2e.build.centos-7.9-gcp-fips: centos79-fips-gcp

e2e.build.ubuntu-18.04-gcp: ubuntu1804-gcp

e2e.build.ubuntu-20.04-gcp: ubuntu2004-gcp

# use sibling containers to handle dependencies and avoid DinD
ci.e2e.ansible:
	make -C test/e2e/ansible e2e.setup
	WHAT="make -C test/e2e/ansible e2e.run" DOCKER_DEVKIT_DEFAULT_ARGS="--rm --net=host" make devkit.run
	make -C test/e2e/ansible e2e.clean
