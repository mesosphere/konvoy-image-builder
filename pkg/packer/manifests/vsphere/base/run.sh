#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

# shellcheck source=/dev/null
if [ -f "${SCRIPT_DIR}"/.env ]; then
  set -a && source "${SCRIPT_DIR}"/.env && set +a
fi

cat "${SSH_PUBLIC_KEY_FILE}" > "${SCRIPT_DIR}"/linux/authorized_keys

packer build -var-file "${SCRIPT_DIR}"/rhel-"${BASE_OS}"-base.json -var-file "${SCRIPT_DIR}"/vsphere-base.json -on-error=abort "${SCRIPT_DIR}"/packer-base.json

