#!/bin/bash
# shellcheck disable=SC2001
echo "${SSH_BASTION_KEY_CONTENTS}" | sed 's/\\n/\n/g' >> vsphere-tests.pem
chmod 600 vsphere-tests.pem
"${ENVSUBST_ASSETS}"/envsubst < test/infra/vsphere/packer-vsphere-airgap.yaml.tmpl >> "$1"
