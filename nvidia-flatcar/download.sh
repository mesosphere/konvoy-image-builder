#!/bin/bash
set -ex

modules=/opt/modules  # Adjust this writable storage location as needed.
sudo mkdir -p "${modules}" "${modules}.wd"
sudo mount \
    -o "lowerdir=/usr/lib64/modules,upperdir=${modules},workdir=${modules}.wd" \
    -t overlay overlay /usr/lib64/modules

. /usr/share/flatcar/release
. /usr/share/flatcar/update.conf
url="https://${GROUP:-stable}.release.flatcar-linux.net/${FLATCAR_RELEASE_BOARD}/${FLATCAR_RELEASE_VERSION}/flatcar_developer_container.bin.bz2"
curl -f -L -O https://www.flatcar.org/security/image-signing-key/Flatcar_Image_Signing_Key.asc
gpg2 --import Flatcar_Image_Signing_Key.asc
curl -L "${url}" |
    tee >(bzip2 -d > flatcar_developer_container.bin) |
    gpg2 --verify <(curl -Ls "${url}.sig") -

mkdir flatcar_container
sudo mount -o ro,loop,offset=2097152 flatcar_developer_container.bin flatcar_container
sudo tar -cp --one-file-system -C flatcar_container/ .  | docker import - faiq/flatcar-dev:${FLATCAR_RELEASE_VERSION}
docker build -t faiq/flatcar-nvidia:dev --build-arg=FLATCAR_RELEASE_VERSION=${FLATCAR_RELEASE_VERSION} container/
