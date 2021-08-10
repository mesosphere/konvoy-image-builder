#!/bin/sh

set -e

yum install -y python3
/usr/sbin/alternatives --set python /usr/bin/python3