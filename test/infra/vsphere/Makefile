# Infra related goals
# ---------------------------------------------------------------------
export SCRIPT_DIR ?= $(CURDIR)/hack

TEST_SCRIPTS_DIR ?= $(CURDIR)/test/infra/vsphere

packer-ova-offline-override.yaml: install-envsubst
	@$(TEST_SCRIPTS_DIR)/export.sh $@
