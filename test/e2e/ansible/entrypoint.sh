#!/bin/bash
set -e

mkdir -p /home/sshuser/.ssh
if [[ -f "/setup/authorized_keys" ]]; then
    cp /setup/authorized_keys /home/sshuser/.ssh
fi
service ssh start

"$@"
