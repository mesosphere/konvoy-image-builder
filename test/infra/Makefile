# Infra related goals
# ---------------------------------------------------------------------
export GOOS ?= $(shell go env GOOS 2>/dev/null)
export GOARCH ?= $(shell go env GOARCH 2>/dev/null)
export TERRAFORM_VERSION ?= 1.5.2
export TERRAFORM_ASSETS ?= $(CURDIR)/.local/terraform/$(TERRAFORM_VERSION)/bin/
export TERRAFORM_URL ?= https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_$(GOOS)_$(GOARCH).zip
export SCRIPT_DIR ?= $(CURDIR)/hack
export TERRAFORM_BIN = $(TERRAFORM_ASSETS)/terraform

INFRA_MODULES_DIR = $(CURDIR)/test/infra

# export for local testing
SSH_BASTION_PUBLIC_KEY_FILE ?= ""
VSPHERE_DATASTORE ?= "ci-kib"

ifeq ($(CI),true)
  $(shell echo "$(SSH_BASTION_PUBLIC_KEY_CONTENTS)" > $(CURDIR)/vsphere-bastion.pub)
  SSH_BASTION_PUBLIC_KEY_FILE = $(CURDIR)/vsphere-bastion.pub
endif

# VCD terraform inputs
VCD_TF_FLAGS ?= \
	-var='vsphere_url=$(VSPHERE_SERVER)' \
	-var='vsphere_username=$(VSPHERE_USERNAME)' \
	-var='vsphere_password=$(VSPHERE_PASSWORD)' \
	-var='datastore_name=$(VSPHERE_DATASTORE)' \
	-var='ssh_public_key=$(SSH_BASTION_PUBLIC_KEY_FILE)' \
	-var='vcd_url=$(VCD_SERVER)' \
	-var='vcd_org_username=$(VCD_ORG_USERNAME)' \
	-var='vcd_org_password=$(VCD_ORG_PASSWORD)' \
	-var='vm_template_name_to_upload=$(TEMPLATE_NAME_TO_UPLOAD)'

.PHONY: install-terraform
install-terraform: ## Download and unpack terraform binary
install-terraform: $(TERRAFORM_BIN)

$(TERRAFORM_BIN):
	$(call print-target,install-terraform)
	mkdir -p $(TERRAFORM_ASSETS)
	curl -o $(TERRAFORM_ASSETS)/terraform.zip -fsSL $(TERRAFORM_URL) && \
		unzip $(TERRAFORM_ASSETS)/terraform.zip -d $(TERRAFORM_ASSETS) && chmod +x $(TERRAFORM_BIN) && \
		rm -rf $(TERRAFORM_ASSETS)/terraform.zip

.PHONY: infra.aws.init
infra.aws.init: ## Initialize infrastructure metadata
infra.aws.init: install-terraform
	$(call print-target)
	$(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/aws init
	cd $(INFRA_MODULES_DIR)/aws && $(TEST_INFRA_STATE_DIR)/aws/init.sh

.PHONY: infra.vsphere.init
infra.vsphere.init: ## Initialize infrastructure metadata
infra.vsphere.init: install-terraform
	$(call print-target)
	$(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vsphere init

.PHONY: infra.aws.create
infra.aws.create: ## Create infrastructure
infra.aws.create: infra.aws.init
	$(call print-target)
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/aws plan
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/aws apply -auto-approve

PHONY: infra.vsphere.create
infra.vsphere.create: ## Create infrastructure
infra.vsphere.create: TF_VAR_ssh_public_key=$(SSH_BASTION_PUBLIC_KEY_FILE)
infra.vsphere.create: TF_VAR_datastore_name=$(VSPHERE_DATASTORE)
infra.vsphere.create: infra.vsphere.init
	$(call print-target)
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vsphere plan
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vsphere apply -auto-approve

packer-aws-offline-override.yaml: infra.aws.create install-envsubst
	@$(INFRA_MODULES_DIR)/aws/export.sh $@

packer-ova-offline-override.yaml: infra.vsphere.create install-envsubst
	@$(INFRA_MODULES_DIR)/vsphere/export.sh $@

.PHONY: infra.aws.destroy
infra.aws.destroy: ## Destroy infrastructure
infra.aws.destroy: install-terraform infra.aws.init
	$(call print-target)
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/aws destroy -auto-approve

.PHONY: infra.vsphere.destroy
infra.vsphere.destroy: ## Destroy infrastructure
infra.vsphere.destroy: TF_VAR_ssh_public_key=$(SSH_BASTION_PUBLIC_KEY_FILE)
infra.vsphere.destroy: TF_VAR_datastore_name=$(VSPHERE_DATASTORE)
infra.vsphere.destroy: infra.vsphere.init
	$(call print-target)
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vsphere destroy -auto-approve

.PHONY: infra.vcd.init
infra.vcd.init: ## Initialize VCD infrastructure
infra.vcd.init: install-terraform
	$(call print-target)
	$(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vcd init

PHONY: infra.vcd.create
infra.vcd.create: ## Create VCD infrastructure
infra.vcd.create: infra.vcd.init
	$(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vcd plan $(VCD_TF_FLAGS)
	$(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vcd apply -auto-approve $(VCD_TF_FLAGS)

PHONY: infra.vcd.destroy
infra.vcd.destroy: ## Destroy VCD infrastructure
infra.vcd.destroy: infra.vcd.init
	$(TERRAFORM_BIN) -chdir=$(INFRA_MODULES_DIR)/vcd destroy -auto-approve $(VCD_TF_FLAGS)
