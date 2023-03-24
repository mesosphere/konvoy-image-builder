#!/bin/sh

# Exit as soon as there is an unexpected error.
set -e

# local path of the vmware guest datasource for cloud init
VMWARE_DS_PATH="${VMWARE_DS_PATH:-/tmp/DataSourceVMwareGuestInfo.py}"


if ! command -v python >/dev/null 2>&1 && \
   ! command -v python3 >/dev/null 2>&1; then
  echo "python 2 or 3 is required" 1>&2
  exit 1
fi

# PYTHON_VERSION may be 2 or 3 and indicates which version of Python
# is used by cloud-init. This variable is not set until PY_MOD_CLOUD_INIT
# is resolved.
CLOUD_INIT_PYTHON=$(head -1 "$(command -v cloud-init || echo /dev/null)"  | cut -c 3-)
echo "cloud-init python: ${CLOUD_INIT_PYTHON}"
PYTHON_VERSION=$(${CLOUD_INIT_PYTHON} -c 'import sys;print(sys.version_info.major)' || echo "")

get_py_mod_dir() {
  _script='import os; import '"${1-}"'; print(os.path.dirname('"${1-}"'.__file__));'
  case "${PYTHON_VERSION}" in
  2)
    python -c ''"${_script}"'' 2>/dev/null || echo ""
    ;;
  3)
    python3 -c ''"${_script}"'' 2>/dev/null || echo ""
    ;;
  *)
    { python3 -c ''"${_script}"'' || python -c ''"${_script}"'' || echo ""; } 2>/dev/null
    ;;
  esac
}

# PY_MOD_CLOUD_INIT is set to the the "cloudinit" directory in either
# the Python2 or Python3 lib directory. This is also used to determine
# which version of Python is repsonsible for running cloud-init.
PY_MOD_CLOUD_INIT="$(get_py_mod_dir cloudinit)"
if [ -z "${PY_MOD_CLOUD_INIT}" ]; then
  echo "cloudinit is required" 1>&2
  exit 1
fi

if echo "${PY_MOD_CLOUD_INIT}" | grep -q python2; then
  PYTHON_VERSION=2
else
  PYTHON_VERSION=3
fi
echo "using python ${PYTHON_VERSION}"

# The following modules are required:
#   * netifaces
# netifaces module is already installed then it is assumed to be compatible.
if [ -z "$(get_py_mod_dir netifaces)" ]; then
   echo "netifaces is required" 1>&2
   exit 1
fi

# move the cloud init datasource into the cloud-init's "sources" directory.
if [ ! -f "${PY_MOD_CLOUD_INIT}/sources/DataSourceVMwareGuestInfo.py" ]; then
  mv "${VMWARE_DS_PATH}" "${PY_MOD_CLOUD_INIT}/sources/DataSourceVMwareGuestInfo.py"
fi

# Make sure that the datasource can execute without error on this host.
echo "validating datasource"
case "${PYTHON_VERSION}" in
2)
  python "${PY_MOD_CLOUD_INIT}/sources/DataSourceVMwareGuestInfo.py" 1>/dev/null
  ;;
3)
  python3 "${PY_MOD_CLOUD_INIT}/sources/DataSourceVMwareGuestInfo.py" 1>/dev/null
  ;;
esac
