# Infra related goals
# ---------------------------------------------------------------------
export TERRAFORM_VERSION ?= 0.15.3
export TERRAFORM_ASSETS ?= $(CURDIR)/.local/terraform/$(TERRAFORM_VERSION)/bin/
export TERRAFORM_URL ?= https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_$(GOOS)_$(GOARCH).zip
export SCRIPT_DIR ?= $(CURDIR)/hack
export TERRAFORM_BIN = $(TERRAFORM_ASSETS)/terraform

TEST_MODULES_DIR = $(CURDIR)/test/infra
TEST_INFRA_STATE_DIR ?= $(CURDIR)/.local/infra

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
	mkdir -p $(TEST_INFRA_STATE_DIR)/aws
	$(TERRAFORM_BIN) -chdir=$(TEST_INFRA_STATE_DIR)/aws init -from-module=$(TEST_MODULES_DIR)/aws -force-copy || true
	cd $(TEST_INFRA_STATE_DIR)/aws && $(TEST_INFRA_STATE_DIR)/aws/init.sh

.PHONY: infra.aws.create
infra.aws.create: ## Create infrastructure
infra.aws.create: install-terraform
	$(call print-target)
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(TEST_INFRA_STATE_DIR)/aws apply -auto-approve

packer-custom-vpc-override.yaml:
packer-custom-vpc-override.yaml:
	@$(TEST_INFRA_STATE_DIR)/aws/export.sh

.PHONY: infra.aws.destroy
infra.aws.destroy: ## Destroy infrastructure
infra.aws.destroy: install-terraform
	$(call print-target)
	TF_LOG=INFO $(TERRAFORM_BIN) -chdir=$(TEST_INFRA_STATE_DIR)/aws destroy -auto-approve
	rm -r $(TEST_INFRA_STATE_DIR)/aws
