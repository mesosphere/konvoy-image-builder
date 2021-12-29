#!/bin/bash

set -eo pipefail

#### Basic VAR definitions
USAGE="USAGE: $(basename "$0") kickstart ISO-URL"
KICKSTART="$1"
KSNAME=${KICKSTART%.*}
BUILDDATE=$(date +%Y%m%d)
BUILDROOT=/var/tmp/containers/$BUILDDATE/$KSNAME
CENTOS_ISO_URL="$2"
DOCKER_TAG="${DOCKER_TAG:=7.9.2009.minimal}"

#### Test for script requirements
# Did we get passed a kickstart
if [ "$#" -ne 2 ]; then
    echo "$USAGE"
    exit 1
fi

# Test for package requirements
PACKAGES=( anaconda-tui lorax yum-langpacks)
for Element in "${PACKAGES[@]}"
  do
    if ! rpm -q --whatprovides "$Element"; then
     echo "RPM $Element missing"
     exit 1
    fi
done

# Is the buildroot already present
if [ -d "$BUILDROOT" ]; then
    echo "The Build root, $BUILDROOT, already exists.  Would you like to remove it? [y/N] "
    read -r REMOVE
    if [ "$REMOVE" == "Y" ] || [ "$REMOVE" == "y" ]
      then
      if [ ! "$BUILDROOT" == "/" ]  
        then
        rm -rf "$BUILDROOT"
      fi
    else
      exit 1
    fi
fi

#download the ISO
ISO_FILE=$(basename "$CENTOS_ISO_URL")
if [ ! -f /tmp/"$ISO_FILE" ]
then
  echo "downloading: $CENTOS_ISO_URL"
  curl -s -o /tmp/"$ISO_FILE" "$CENTOS_ISO_URL"
fi

# Build the rootfs
time livemedia-creator --logfile=/tmp/"$KSNAME"-"$BUILDDATE".log \
     --no-virt --make-tar --ks "$KICKSTART" \
     --iso /tmp/"$ISO_FILE" \
     --image-name="$KSNAME"-docker.tar.xz --project "CentOS 7 Docker" \
     --releasever "7"

# Put the rootfs someplace
mkdir -p "$BUILDROOT"/docker
mv /var/tmp/"$KSNAME"-docker.tar.xz "$BUILDROOT"/docker/

# Create a Dockerfile to go along with the rootfs.

BUILDDATE_RFC3339="$(date -d "$BUILDDATE" --rfc-3339=seconds)"
cat << EOF > "$BUILDROOT"/docker/Dockerfile
FROM scratch
ADD $KSNAME-docker.tar.xz /
LABEL \\
    org.label-schema.schema-version="1.0" \\
    org.label-schema.name="CentOS Base Image" \\
    org.label-schema.vendor="CentOS" \\
    org.label-schema.license="GPLv2" \\
    org.label-schema.build-date="$BUILDDATE" \\
    org.opencontainers.image.title="CentOS Base Image" \\
    org.opencontainers.image.vendor="CentOS" \\
    org.opencontainers.image.licenses="GPL-2.0-only" \\
    org.opencontainers.image.created="$BUILDDATE_RFC3339"
CMD ["/bin/bash"]
EOF

# import docker image
docker import - centos:"$DOCKER_TAG" < "$BUILDROOT"/docker/"$KSNAME"-docker.tar.xz
