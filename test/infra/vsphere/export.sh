#!/bin/bash
set -eu
# shellcheck disable=SC2001
echo "${SSH_BASTION_KEY_CONTENTS}" | sed 's/\\n/\n/g' >> vsphere-tests.pem
chmod 600 vsphere-tests.pem

# shellcheck disable=SC2155
export SSH_BASTION_HOST="$("${TERRAFORM_BIN}" -chdir=.local/infra/vsphere output -raw bastion_node_ssh_nat_address)"
# shellcheck disable=SC2155
export SSH_BASTION_PORT="$("${TERRAFORM_BIN}" -chdir=.local/infra/vsphere output -raw bastion_node_ssh_nat_port)"
# shellcheck disable=SC2155
export SSH_BASTION_USERNAME="$("${TERRAFORM_BIN}" -chdir=.local/infra/vsphere output -raw bastion_node_ssh_user)"


"${ENVSUBST_ASSETS}"/envsubst < test/infra/vsphere/packer-vsphere-airgap.yaml.tmpl >> "$1"
