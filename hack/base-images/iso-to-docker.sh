#!/bin/bash

set -eo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

ISO_URL=$1
DOCKER_IMAGE=$2
CENTOS_ISO_URL="${ISO_URL:=http://centos.mirror.ndchost.com/7.9.2009/isos/x86_64/CentOS-7-x86_64-Minimal-2009.iso}"
DOCKER_IMAGE="${DOCKER_IMAGE:=centos:7.9.2009.minimal}"

PACKAGES=( docker squashfs-tools curl )
for element in "${PACKAGES[@]}"
  do
     TEST=`rpm -q --whatprovides $element`
     if [ "$?" -ne 0 ]
     then echo "RPM $element is missing"
     exit 1
     fi
done

filename=`basename $CENTOS_ISO_URL`
if [ ! -f $SCRIPT_DIR/$filename ]
then
  curl -s -o $filename $CENTOS_ISO_URL
fi

mkdir -p rootfs unsquashfs temp-unsquashfs
mount -o loop $filename rootfs
unsquashfs -f -d temp-unsquashfs rootfs/LiveOS/squashfs.img
mount -o loop temp-unsquashfs/LiveOS/rootfs.img unsquashfs

echo "building docker image $DOCKER_IMAGE"
tar -C unsquashfs -c . | docker import --change "CMD /bin/sh" - $DOCKER_IMAGE

docker images
umount rootfs unsquashfs
