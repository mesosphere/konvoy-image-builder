# Build Base OS image using packer
This document outlines process to create bast image template from RedHat DVD ISO file

## Prerequisites
- Hashicorp Packer: Please download it for your operating system
- Access to vSphere: Make sure you have network access to vSphere.
- vSphere credentials exported:

```bash
export VSPHERE_USERNAME="<USERNAME>"
export VSPHERE_PASSWORD="<PASSWORD>"
export VSPHERE_SERVER="VSPHERE_SERVER_URL>"
```
- RedHat subscription: visit [RedHat Developer site](https://developers.redhat.com/) to register

## Build base image
Following steps builds creates base OS vsphere templates

1. Download DVD ISO from RedHat [download site](https://developers.redhat.com/products/rhel/download) for RHEL 7.9
you must login to redhat in order to download the DVD ISO file.

1. create .env file using .env.sample and configure required parameters
```bash
    cp .env.sample .env
```

1. Run packer with configured parameters

```bash
./run.sh

```
## How it works
The base image template is created from the RedHat DVD ISO file. The DVD ISO file has minimum core packages avaialble in it.
The packer vspher-iso build will perform following steps
1. Creates a VM in vSphere and mounts the DVD ISO (`/dev/sr0`) and kickstart files (`/dev/sr1`) as cdrom drives. The kickstart files are mounted from `./linux/rhel/<version>/ks.cfg`

1. Runs the kickstarts file at boot time using boot command `<up><tab> inst.text inst.ks=hd:sr1:linux/rhel/http/8/ks.cfg` [Reference documentation](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/installation_guide/sect-kickstart-howto#sect-kickstart-installation-starting-automatic)

1. The kickstart file installs core packages and openssh packages and creates a user without any password.

1. The post installation steps in the kickstart file adds provided ssh public key to the authorized keys for the user.

1. The packer is provided with the private key at runtime so that packer can ssh to the VM post installation and configure the VM

1. Packer shutsdown the VM using remote SSH command and converts it to template.

The VMs created with the template can only be accessed using the correosponding private key provided during base template installations.


## References:
- [Kickstart syntax reference](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/installation_guide/sect-kickstart-syntax)
- [Centos community kickstart files](https://github.com/CentOS/Community-Kickstarts)

