#!/bin/bash
envsubst < test/infra/vsphere/packer-vsphere-airgap.yaml.tmpl >> packer-vsphere-airgap.yaml
