#!/usr/bin/env bash
set -o errexit
set -o pipefail
set -o nounset

# shellcheck disable=SC2155
export AWS_VPC_ID="$("${TERRAFORM_BIN}" -chdir=.local/infra/aws output -raw vpc_id)"
# shellcheck disable=SC2155
export AWS_SECURITY_GROUP_ID="$("${TERRAFORM_BIN}" -chdir=.local/infra/aws output -raw security_group_id)"
# shellcheck disable=SC2155
export AWS_SUBNET_ID="$("${TERRAFORM_BIN}" -chdir=.local/infra/aws output -raw public_subnets)"
envsubst < .local/infra/aws/packer-custom-vpc-override.yaml.tmpl >> packer-custom-vpc-override.yaml
