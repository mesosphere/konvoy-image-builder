#!/bin/bash
set -eu
# shellcheck disable=SC2001
echo "${SSH_BASTION_KEY_CONTENTS}" | sed 's/\\n/\n/g' >> vsphere-tests.pem
chmod 600 vsphere-tests.pem

SSH_BASTION_HOST="$("${TERRAFORM_BIN}" -chdir="${INFRA_MODULES_DIR}/vsphere" output -raw bastion_node_ssh_nat_address)"
export SSH_BASTION_HOST

SSH_BASTION_PORT="$("${TERRAFORM_BIN}" -chdir="${INFRA_MODULES_DIR}/vsphere" output -raw bastion_node_ssh_nat_port)"
export SSH_BASTION_PORT

SSH_BASTION_USERNAME="$("${TERRAFORM_BIN}" -chdir="${INFRA_MODULES_DIR}/vsphere"  output -raw bastion_node_ssh_user)"
export SSH_BASTION_USERNAME

"${ENVSUBST_ASSETS}"/envsubst < "${INFRA_MODULES_DIR}"/vsphere/packer-vsphere-airgap.yaml.tmpl >> "$1"
