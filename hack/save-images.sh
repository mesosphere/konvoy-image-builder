#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Some boiler plate
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Some constants
IMAGE_LIST_FILE="${1:-images.out}"
EXTRA_LIST_FILE="${2:-}"
TAR_FILE_NAME="${3:-airgapped.container-images.tar.gz}"

# check if the image list exists
if [ ! -f "$(pwd)/$IMAGE_LIST_FILE" ]; then
  echo "list of images not found, please run 'make list-images' to generate and/or check output file"
fi

echo "reading from: $(pwd)/${IMAGE_LIST_FILE}"
while read -r image; do
  images+=("$image")
done <"$(pwd)/${IMAGE_LIST_FILE}"

# check if the extra image list exists
if [ -f "$(pwd)/$EXTRA_LIST_FILE" ]; then
  echo "extra file found, adding from: $(pwd)/${EXTRA_LIST_FILE}"
  while read -r image; do
    images+=("$image")
  done <"$(pwd)/${EXTRA_LIST_FILE}"
elif [ -n "${EXTRA_LIST_FILE}" ]; then
  echo "not found: $(pwd)/${EXTRA_LIST_FILE}"
  exit 1
fi

# we wish to fail quickly here and not hide errors, docker pull can fail but in bash loop it will fail silently
# we ensure we capture these failures immediately and not continue.
for image in "${images[@]}"; do
  printf "%b\n" "${image}"
  docker pull "${image}"
  # shellcheck disable=SC2181
  if [ $? != 0 ]; then
    echo "error pulling ${image}"
    exit 1
  fi
done

printf "\n saving to tar file: %s\n" "$(pwd)/${TAR_FILE_NAME}"
# shellcheck disable=SC2068
docker save ${images[@]} | gzip >"$(pwd)/${TAR_FILE_NAME}"
