#!/bin/bash
# Flatcar's update typically occurs 10 minutes after boot. We need to turn it
# off before then. Hence, we run this script before all other Flatcar scripts.
# We should check if this can be moved to Ansible and still work or if we can
# configure even earlier with a cloud-init / Container Linux Config.
set -e

[[ "$BUILD_NAME" != *"flatcar"* ]] && exit 0

# Prevent systemd from starting `update-engine` and `locksmithd`. The default
# system unit files check for `/usr/.noupdate`, but since `/usr` is a read-only
# filesystem we cannot create this. Add `/etc/.noupdate` to perform the same
# role.
for service in locksmithd update-engine
do
    mkdir -p /etc/systemd/system/${service}.service.d
    cat > /etc/systemd/system/${service}.service.d/10-noupdate.conf <<EOF
[Unit]
ConditionPathExists=!/etc/.noupdate
EOF
done
touch /etc/.noupdate
systemctl daemon-reload

# Stop the currently running service instances.
systemctl stop locksmithd
systemctl stop update-engine

# Disable automated reboots.  With `locksmithd` disabled, no reboot strategy is
# required. However, setting this ensures that we are conservative if
# `locksmithd` is re-enabled.
sed -i '/^REBOOT_STRATEGY=.*/d' /etc/flatcar/update.conf
echo 'REBOOT_STRATEGY=off' >> /etc/flatcar/update.conf
