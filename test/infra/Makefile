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
	cd $(INFRA_MODULES_DIR)/aws && $(INFRA_MODULES_DIR)/aws/init.sh

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

packer-ova-basic-override.yaml: install-envsubst
	"${ENVSUBST_ASSETS}"/envsubst < "$(INFRA_MODULES_DIR)"/vsphere/packer-vsphere-basic.yaml.tmpl > $@

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

