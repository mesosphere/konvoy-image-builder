#!/bin/bash
# Flatcar provides a basic free OS, supported for 9 months, and premium OS's.
# All of these, except for the most recent free OS, are provided through an
# AMI Marketplace subscription.  This prevents generation of a public AMI from
# these images.  Without further configuration, we must use the latest Flatcar
# release.  Hence, we run a Packer build to identify the latest Flatcar AMI
# and provide an override file to select it.

set -e -x

if [ -z "${AWS_DEFAULT_REGION}" ]
then
    # If running in an AWS region, use it
    AWS_DEFAULT_REGION=$(curl -sSL http://169.254.169.254/latest/meta-data/placement/region) || true
fi
if [ -z "${AWS_DEFAULT_REGION}" ]
then
    echo "Need AWS_DEFAULT_REGION to be set" >&2
    exit 1
fi

# Directory containing this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

packer build "${SCRIPT_DIR}/flatcar-ami.pkr.hcl"
