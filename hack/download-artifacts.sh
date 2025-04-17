#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
set -x

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." &> /dev/null && pwd )"
PACKAGES_BASE_URL="https://packages.d2iq.com/konvoy/stable/linux/repos"
IMAGES_BASE_URL="https://downloads.d2iq.com/dkp/airgapped/kubernetes-images"

TARGET_ARTIFACTS_DIR="${TARGET_ARTIFACTS_DIR:-"/opt"}"
TARGET_KUBERNETES_IMAGES_DIR="${TARGET_KUBERNETES_IMAGES_DIR:-"${TARGET_ARTIFACTS_DIR}/kubernetes-images"}"

KUBERNETES_VERSION=$(awk -F': ' '/kubernetes_version/ {print $2}' "${PROJECT_DIR}/ansible/group_vars/all/defaults.yaml" | sed -n '2p' | xargs)
CRICTL_TOOLS_VERSION="$(echo "${KUBERNETES_VERSION}" | cut -d. -f1-2).0"
CNI_VERSION=$(awk -F': ' '/kubernetes_cni_version/ {print $2}' "${PROJECT_DIR}/ansible/group_vars/all/defaults.yaml" | sed -n '1p' | xargs)

# ensure target directories exist
mkdir -p "${TARGET_ARTIFACTS_DIR}"

# Download non-fips rpms
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubectl-${KUBERNETES_VERSION}-0.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-nokmem/x86_64/kubectl-${KUBERNETES_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubeadm-${KUBERNETES_VERSION}-0.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-nokmem/x86_64/kubeadm-${KUBERNETES_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubelet-${KUBERNETES_VERSION}-0.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-nokmem/x86_64/kubelet-${KUBERNETES_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/cri-tools-${CRICTL_TOOLS_VERSION}-0.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-nokmem/x86_64/cri-tools-${CRICTL_TOOLS_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubernetes-cni-${CNI_VERSION}-0.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-nokmem/x86_64/kubernetes-cni-${CNI_VERSION}-0.x86_64.rpm"
  
# Download FIPS RPMs
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubectl-${KUBERNETES_VERSION}-0-fips.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-fips/x86_64/kubectl-${KUBERNETES_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubeadm-${KUBERNETES_VERSION}-0-fips.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-fips/x86_64/kubeadm-${KUBERNETES_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubelet-${KUBERNETES_VERSION}-0-fips.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-fips/x86_64/kubelet-${KUBERNETES_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/cri-tools-${CRICTL_TOOLS_VERSION}-0-fips.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-fips/x86_64/cri-tools-${CRICTL_TOOLS_VERSION}-0.x86_64.rpm"
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/kubernetes-cni-${CNI_VERSION}-0-fips.rpm" "${PACKAGES_BASE_URL}/el/kubernetes-v${KUBERNETES_VERSION}-fips/x86_64/kubernetes-cni-${CNI_VERSION}-0.x86_64.rpm"

# download gpg key
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/d2iq-sign-authority-gpg-public-key" "${PACKAGES_BASE_URL}/d2iq-sign-authority-gpg-public-key"

# Download kubernetes images
mkdir -p "${TARGET_KUBERNETES_IMAGES_DIR}"
curl --fail -v -o "${TARGET_KUBERNETES_IMAGES_DIR}/kubernetes-images-${KUBERNETES_VERSION}-d2iq.1.tar" "${IMAGES_BASE_URL}/kubernetes-images-${KUBERNETES_VERSION}-d2iq.1.tar"
curl --fail -v -o "${TARGET_KUBERNETES_IMAGES_DIR}/kubernetes-images-${KUBERNETES_VERSION}-d2iq.1-fips.tar" "${IMAGES_BASE_URL}/kubernetes-images-${KUBERNETES_VERSION}-d2iq.1-fips.tar"

# Amazon SSM agent
curl --fail -v -o "${TARGET_ARTIFACTS_DIR}/amazon-ssm-agent.rpm" https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm
