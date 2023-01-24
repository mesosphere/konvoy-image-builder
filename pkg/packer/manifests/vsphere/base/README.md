# Build Base OS image using packer

This document outlines process to create bast image template from a DVD ISO file

## Prerequisites

- Hashicorp Packer: Please download it for your operating system
- Access to vSphere: Make sure you have network access to vSphere.
- vSphere credentials exported:

```bash
export VSPHERE_USERNAME="<USERNAME>"
export VSPHERE_PASSWORD="<PASSWORD>"
export VSPHERE_SERVER="VSPHERE_SERVER_URL>"
```
- RedHat subscription: visit [RedHat Developer site](https://developers.redhat.com/) to register (to build RHEL OVAs)

## Build base RHEL image

Following steps creates base OS vsphere templates

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

## Build base Ubuntu image

Following steps creates base OS vsphere templates

1. Download [DVD ISO from Ubuntu](https://cdimage.ubuntu.com/ubuntu-legacy-server/releases/20.04/release/ubuntu-20.04.1-legacy-server-amd64.iso).

1. create .env file using .env.sample and configure required parameters
```bash
    cp .env.sample .env
```

1. Run packer with configured parameters, wait for it to start the VM.

```bash
./run.sh
```

2. Navigate to the vCenter console to perform the initial setup
   * Create a `builder` user with a secure password that you will use for the initial ssh connection.
   * Check to install Open SSH.
   * Wait for the install to complete and "Reboot Now" the machine when prompted.

3. SSH into the restarted VM with user `buider` and the password you create in the previous step. Run run the following steps (replacing the SSH public key placeholder with desired key):

```bash
sudo su 
# don't require password for sudo
echo 'builder ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/builder
chmod 440 /etc/sudoers.d/builder
# install open-vm-tools and then clean up apt cache
apt-get update
apt-get install open-vm-tools -y
apt-get purge --auto-remove -y
rm -rf /targe/var/lib/apt/lists/*
# add an SSH public key
mkdir -p /home/builder/.ssh/
echo '<CHANGE_ME>' >> /home/builder/.ssh/authorized_keys
chown -R builder:builder /home/builder/.ssh/
chmod 644 /home/builder/.ssh/authorized_keys
chmod 700 /home/builder/.ssh/
rm -f /etc/udev/rules.d/70-persistent-net.rules
shutdown -P now
```

4. In the vCenter console: right-click on the VM > Snapshots > Take Snapshot...

5. In the vCenter console: right-click on the VM > Template > Convert to Template

## How it works

The base image template is created from the RedHat DVD ISO file. The DVD ISO file has minimum core packages avaialble in it.
The packer vsphere-iso build will perform following steps
1. Creates a VM in vSphere and mounts the DVD ISO (`/dev/sr0`) and kickstart files (`/dev/sr1`) as cdrom drives. The kickstart files are mounted from `./linux/rhel/<version>/ks.cfg`

1. Runs the kickstarts file at boot time using boot command `<up><tab> inst.text inst.ks=hd:sr1:linux/rhel/http/8/ks.cfg` [Reference documentation](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/installation_guide/sect-kickstart-howto#sect-kickstart-installation-starting-automatic)

1. The kickstart file installs core packages and openssh packages and creates a user without any password.

1. The post installation steps in the kickstart file adds provided ssh public key to the authorized keys for the user.

1. The packer is provided with the private key at runtime so that packer can ssh to the VM post installation and configure the VM

1. Packer shutsdown the VM using remote SSH command and converts it to template.

The VMs created with the template can only be accessed using the correosponding private key provided during base template installations.


## References

- [Kickstart syntax reference](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/installation_guide/sect-kickstart-syntax)
- [Centos community kickstart files](https://github.com/CentOS/Community-Kickstarts)

