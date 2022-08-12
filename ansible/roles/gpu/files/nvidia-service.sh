#!/bin/bash

ldconfig
modprobe -d / ipmi_devintf || true
depmod
modprobe nvidia
modprobe nvidia-modeset
modprobe nvidia-uvm
