#!/bin/bash
set -ex
DRIVER_NAME=nvidia
DRIVER_INSTALL_DIR=/opt/drivers
DRIVER_CACHE_DIR="${DRIVER_INSTALL_DIR}/archive/${DRIVER_NAME}"
LD_ROOT=/
mkdir -p "${DRIVER_INSTALL_DIR}"
mkdir -p "${DRIVER_CACHE_DIR}"
rm -rf "opt/drivers/${DRIVER_NAME}"

docker run --rm -v /usr/lib64/modules:/usr/lib64/modules -v ${DRIVER_CACHE_DIR}:/out faiq/flatcar-nvidia:dev
ln -s "${DRIVER_CACHE_DIR}" "${DRIVER_INSTALL_DIR}/${DRIVER_NAME}"

mkdir -p "/etc/ld.so.conf.d"
echo "${DRIVER_INSTALL_DIR}/${DRIVER_NAME}/lib" >  /etc/ld.so.conf.d/${DRIVER_NAME}.conf
ldconfig -r /

depmod -b "${LD_ROOT}"
# This is an NVIDIA dep that is not specified in the module.dep file.
modprobe -d "${LD_ROOT}" ipmi_devintf || true
depmod -b   "${DRIVER_INSTALL_DIR}/${DRIVER_NAME}" || true
modprobe -d "${DRIVER_INSTALL_DIR}/${DRIVER_NAME}" nvidia
modprobe -d "${DRIVER_INSTALL_DIR}/${DRIVER_NAME}" nvidia-modeset
modprobe -d  "${DRIVER_INSTALL_DIR}/${DRIVER_NAME}" nvidia-uvm

# Verify installation
${DRIVER_INSTALL_DIR}/${DRIVER_NAME}/bin/nvidia-smi
${DRIVER_INSTALL_DIR}/${DRIVER_NAME}/bin/nvidia-modprobe -u -m -c 0
