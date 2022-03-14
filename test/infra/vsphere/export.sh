#!/bin/bash
#CURDIR is populated from the make target
echo "${SSH_BASTION_KEY_CONTENTS}" | sed 's/\\n/\n/g' >> vsphere-tests.pem
chmod 600 vsphere-tests.pem
envsubst < test/infra/vsphere/packer-vsphere-airgap.yaml.tmpl >> packer-vsphere-airgap.yaml
