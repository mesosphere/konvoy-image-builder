#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

PROD_VSPHERE_URL=${PROD_VSPHERE_URL:=""}
PROD_VSPHERE_USERNAME=${PROD_VSPHERE_USERNAME:=""}
PROD_VSPHERE_PASSWORD=${PROD_VSPHERE_PASSWORD:=""}
PROD_VSPHERE_DATACENTER=${PROD_VSPHERE_DATACENTER:="dc1"}
PROD_TEMPLATE_NAME=${PROD_TEMPLATE_NAME:=""}


VCD_URL=${VCD_URL:=""}
VCD_USERNAME=${VCD_USERNAME:=""}
VCD_PASSWORD=${VCD_PASSWORD:=""}
VCD_ORG=${VCD_ORG:=""}
VCD_ORG_CATALOG=${VCD_ORG_CATALOG:=""}




GOVC_VERSION=${GOVC_VERSION:="v0.30.5"}
TMPDIR=${TMPDIR:=$SCRIPT_DIR}
# extract govc binary to /usr/local/bin
# note: the "tar" command must run with root permissions
GOVC_INSTALL_DIR="${TMPDIR}/govc-${GOVC_VERSION}"
mkdir -p "${GOVC_INSTALL_DIR}"
curl -L -o "${GOVC_INSTALL_DIR}"/govc.tar.gz  "https://github.com/vmware/govmomi/releases/download/${GOVC_VERSION}/govc_$(uname -s)_$(uname -m).tar.gz"
tar -xvzf "${GOVC_INSTALL_DIR}"/govc.tar.gz -C "${GOVC_INSTALL_DIR}"


echo "Exporting template ""${PROD_TEMPLATE_NAME}"" from production"
GOVC_URL=${PROD_VSPHERE_URL} GOVC_USERNAME=${PROD_VSPHERE_USERNAME} GOVC_PASSWORD=${PROD_VSPHERE_PASSWORD} "${GOVC_INSTALL_DIR}"/govc export.ovf -f=true -sha 256 -dc="${PROD_VSPHERE_DATACENTER}" -vm "${PROD_TEMPLATE_NAME}" "${SCRIPT_DIR}"

echo "fix the file size in the file descriptor"
# shellcheck disable=SC2012
VMDK_SIZE=$(ls -l "${SCRIPT_DIR}/${PROD_TEMPLATE_NAME}/${PROD_TEMPLATE_NAME}-disk-0.vmdk"  | awk '{print $5}')
sed -i -E "s/(ovf:size=\")[0-9].+(\"\/>)/\1${VMDK_SIZE}\2/g" "${SCRIPT_DIR}/${PROD_TEMPLATE_NAME}/${PROD_TEMPLATE_NAME}.ovf"

sudo dnf update -y
sudo dnf install python3 python3-pip -y
pip3 install vcd-cli
VCD_CLI="${HOME}/.local/bin/vcd"
${VCD_CLI} login "${VCD_URL}" "${VCD_ORG}" "${VCD_USERNAME}" --password "${VCD_PASSWORD}" --disable-warnings  --no-verify-ssl-certs -V 37.1

echo "upload template to VCD organization's catalog"
cd "${SCRIPT_DIR}/${PROD_TEMPLATE_NAME}" && tar -cvf "${SCRIPT_DIR}/${PROD_TEMPLATE_NAME}.ova" "${PROD_TEMPLATE_NAME}".mf "${PROD_TEMPLATE_NAME}".ovf "${PROD_TEMPLATE_NAME}"-disk-0.vmdk
${VCD_CLI} catalog upload "${VCD_ORG_CATALOG}" -i "${PROD_TEMPLATE_NAME}" "${SCRIPT_DIR}/${PROD_TEMPLATE_NAME}.ova"
