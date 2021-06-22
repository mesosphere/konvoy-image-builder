#!/usr/bin/env bash

set -ex

# How to compile custom kernel drivers:
# https://kinvolk.io/docs/flatcar-container-linux/latest/reference/developer-guides/kernel-modules/

MODULES_PATH="/opt/modules"

if [ ! -d "$MODULES_PATH" ]; then
    sudo mkdir -p "${MODULES_PATH}" "${MODULES_PATH}.wd"
    sudo mount \
        -o "lowerdir=/usr/lib64/modules,upperdir=${MODULES_PATH},workdir=${MODULES_PATH}.wd" \
        -t overlay overlay /usr/lib64/modules
fi

MODULES_MOUNT_UNIT="/etc/systemd/system/usr-lib64-modules.mount"
if [ ! -f "$MODULES_MOUNT_UNIT" ]; then
    cat <<EOF | sudo tee -a "$MODULES_MOUNT_UNIT"
[Unit]
Description=Custom Kernel Modules
Before=local-fs.target
ConditionPathExists=$MODULES_PATH

[Mount]
Type=overlay
What=overlay
Where=/usr/lib64/modules
Options=lowerdir=/usr/lib64/modules,upperdir=$MODULES_PATH,workdir=$MODULES_PATH.wd

[Install]
WantedBy=local-fs.target
EOF
fi

sudo systemctl enable usr-lib64-modules.mount

FLATCAR_IMAGE_FILE="flatcar_developer_container.bin"

if [ ! -f "$FLATCAR_IMAGE_FILE" ]; then
# Prepare flatcar container linux
    . /usr/share/flatcar/release
    . /usr/share/flatcar/update.conf
    FLATCAR_CONTAINER_URL="https://${GROUP:-stable}.release.flatcar-linux.net/${FLATCAR_RELEASE_BOARD}/${FLATCAR_RELEASE_VERSION}/flatcar_developer_container.bin.bz2"
    curl -L "${FLATCAR_CONTAINER_URL}" | bzip2 -d > "$FLATCAR_IMAGE_FILE"
fi

# `dev` func from https://gist.github.com/vbatts/9af92a341611751dc3a157f204a84973
# use to invoke a command in context of dev container
dev() {
    sudo systemd-nspawn \
        --bind=/usr/lib64/modules \
        --bind=/:/hostfs \
        --image=flatcar_developer_container.bin ${@}
}

# Prepare dev container for nvidia-driver installation
dev -- emerge-gitclone
dev -- emerge -gKv coreos-sources
dev -- emerge -gKv sys-fs/squashfs-tools
dev -- cp /usr/lib64/modules/$(ls /usr/lib64/modules)/build/.config /usr/src/linux/
dev -- make -C /usr/src/linux modules_prepare
# This is necessary in the 5.10+ kernel version
dev -- cp /usr/lib64/modules/$(ls /usr/lib64/modules)/build/Module.symvers /usr/src/linux/

KERNEL_VERSION=$(uname -r)
DRIVER_VERSION="460.84"
DRIVER_INSTALL_WORKDIR="/opt/nvidia"
DRIVER_PATH="$DRIVER_INSTALL_WORKDIR/NVIDIA-Linux-x86_64-$DRIVER_VERSION"

# Download nvidia driver
if [ ! -e "$DRIVER_PATH" ]; then
    sudo mkdir -p "$DRIVER_INSTALL_WORKDIR" && sudo chmod 777 "$DRIVER_INSTALL_WORKDIR"
    pushd "$DRIVER_INSTALL_WORKDIR"
    curl "https://us.download.nvidia.com/XFree86/Linux-x86_64/$DRIVER_VERSION/NVIDIA-Linux-x86_64-$DRIVER_VERSION.run" -o nvidia.run
    chmod +x nvidia.run
    ./nvidia.run -x -s
    popd
fi

# Run driver installation
if [ ! -e "$DRIVER_PATH/kernel/nvidia.ko" ]; then
    # This command will always exit with non-0 exit code because installer
    # tries to load module as part of installation process.
    dev --chdir="/hostfs$DRIVER_PATH" ./nvidia-installer -s -n \
        --kernel-name="$KERNEL_VERSION" \
        --kernel-source-path=/usr/src/linux \
        --no-check-for-alternate-installs \
        --no-opengl-files \
        --no-distro-scripts \
        --kernel-install-path="$DRIVER_PATH" \
        --log-file-name="$DRIVER_PATH"/nvidia-installer.log || true
fi

if [ ! -e "$DRIVER_PATH/kernel/nvidia.ko" ]; then
    echo "Failed to compile NVIDIA modules"
    cat "$DRIVER_PATH/nvidia-installer.log" && exit 1
fi

INSTALL_DIR="/usr/lib64/modules/nvidia-driver-$DRIVER_VERSION"
dev -- mkdir -p {"$INSTALL_DIR/bin","$INSTALL_DIR/lib64","$INSTALL_DIR/lib64/modules/$(uname -r)/kernel/drivers/video"}
dev -- find "/hostfs$DRIVER_PATH" -maxdepth 1 -name "*.so.*" \
    -exec cp {} "$INSTALL_DIR/lib64" \;
dev -- find "/hostfs$DRIVER_PATH" -maxdepth 1 -name "nvidia-*" -executable \
    -exec cp {} "$INSTALL_DIR/bin" \;
dev -- find "/hostfs$DRIVER_PATH/kernel" -maxdepth 1 -name "*.ko" \
    -exec cp {} "$INSTALL_DIR/lib64/modules/$(uname -r)/kernel/drivers/video" \;

sudo mkdir -p /etc/ld.so.conf.d/
echo "$INSTALL_DIR/lib64" | sudo tee -a /etc/ld.so.conf.d/nvidia.conf
sudo ldconfig

sudo rm /etc/modules-load.d/nvidia.conf
cat <<EOF | sudo tee -a /etc/modules-load.d/nvidia.conf
nvidia
nvidia-uvm
nvidia-modeset
EOF
