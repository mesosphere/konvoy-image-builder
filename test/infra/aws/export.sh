#!/usr/bin/env bash
set -o errexit
set -o pipefail
set -o nounset

# shellcheck disable=SC2155
export AWS_VPC_ID="$("${TERRAFORM_BIN}" -chdir="${INFRA_MODULES_DIR}/aws"  output -raw vpc_id)"
# shellcheck disable=SC2155
export AWS_SECURITY_GROUP_ID="$("${TERRAFORM_BIN}" -chdir="${INFRA_MODULES_DIR}/aws"  output -raw security_group_id)"
# shellcheck disable=SC2155
export AWS_SUBNET_ID="$("${TERRAFORM_BIN}" -chdir="${INFRA_MODULES_DIR}/aws" output -raw public_subnets)"
"${ENVSUBST_ASSETS}"/envsubst < "${INFRA_MODULES_DIR}"/aws/packer-aws-offline-override.yaml.tmpl >> "$1"
