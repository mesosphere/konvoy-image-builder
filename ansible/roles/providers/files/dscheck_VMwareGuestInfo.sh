#!/bin/sh

# Cloud-Init Datasource for VMware Guestinfo
#
# This file should be installed to /usr/bin/dscheck_VMwareGuestInfo
# without the ".sh" extension. The extension only exists to make it easier
# to identify the file during development.
#
# This file provides cloud-init's ds-identify program a shell type that
# can be resolved with "type dscheck_VMwareGuestInfo" and used to validate
# where a datasource is installed and useable.
#
# Cloud-init's ds-identify program in /usr/lib/cloud-init includes functions
# to determine whether or not datasources can be used. Because the program
# is a shell script and uses "type dscheck_DATASOURCE_NAME" to determine
# if there is a matching bash type that can answer for the datasource,
# it's possible to respond with an external script. While other datasources
# have functions in ds-identify, the "type" command looks up types both
# in Bash's function table as well as script in the PATH. Therefore the
# ds-identify program, when looking up whether or not the datasource
# VMwareGuestInfo can be used, will defer to this file when it is in the
# PATH and named dscheck_VMwareGuestInfo.
#

if [ -n "${VMX_GUESTINFO}" ]; then
  if [ -n "${VMX_GUESTINFO_METADATA}" ] || \
     [ -n "${VMX_GUESTINFO_USERDATA}" ] || \
     [ -n "${VMX_GUESTINFO_VENDORDATA}" ]; then
     exit 0
  fi
fi

if ! command -v vmware-rpctool >/dev/null 2>&1; then
  exit 1
fi

if { vmware-rpctool "info-get guestinfo.metadata" || \
     vmware-rpctool "info-get guestinfo.userdata" || \
     vmware-rpctool "info-get guestinfo.vendordata"; } >/dev/null 2>&1; then
   exit 0
fi

exit 1
