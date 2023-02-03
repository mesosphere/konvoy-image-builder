# Infra related goals
# ---------------------------------------------------------------------
export SCRIPT_DIR ?= $(CURDIR)/hack

TEST_SCRIPTS_DIR ?= $(CURDIR)/test/infra/vsphere

packer-ova-offline-override.yaml: install-envsubst
	eval $(ssh-agent)
	ssh-add vsphere-tests.pem
	@$(TEST_SCRIPTS_DIR)/export.sh $@
