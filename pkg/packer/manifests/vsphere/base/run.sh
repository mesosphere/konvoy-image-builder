#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

# shellcheck source=/dev/null
if [ -f "${SCRIPT_DIR}"/.env ]; then
  set -a && source "${SCRIPT_DIR}"/.env && set +a
fi

cat "${SSH_PUBLIC_KEY_FILE}" > "${SCRIPT_DIR}"/linux/authorized_keys

PACKER_FILE=packer-base.json
OS_VARS_FILE_ARG="-var-file=${SCRIPT_DIR}/${BASE_OS}-base.json -var-file=${SCRIPT_DIR}/vsphere-base.json"
if [[ "${BASE_OS}" == ubuntu* ]]; then
  PACKER_FILE=packer-ubuntu-base.json
fi

if [[ "${BASE_OS}" == rocky* ]]; then
  PACKER_FILE=rockylinux.pkr.hcl
  OS_VARS_FILE_ARG="-var-file=${SCRIPT_DIR}/${BASE_OS}-base.pkrvars.hcl"
fi

packer build ${OS_VARS_FILE_ARG}  -on-error=abort "${SCRIPT_DIR}"/"${PACKER_FILE}"
