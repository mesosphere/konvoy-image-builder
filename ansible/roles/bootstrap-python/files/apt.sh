#!/bin/sh

set -e

apt-get update -qq
apt-get install -qq -y --no-install-recommends python2.7
ln -s /usr/bin/python2.7 /usr/bin/python
