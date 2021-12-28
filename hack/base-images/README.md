# Create docker image from Centos ISO 

This page documents steps required to create minmal docker image from centos 7 ISO

The build scripts are slightly modified version of https://github.com/CentOS/sig-cloud-instance-build/tree/master/docker
Build scripts on the `sig-cloud-instance-build` repo are used to create official centos docker images.

### Prerequisites
All steps needed to be executed in a Centos VM.
- Create EC2 instance: t2.large with 40 GB disk

### Create docker image
- install required packages
```
sudo yum install -y docker lorax anaconda-tui yum-langpacks virt-install libvirt-python
sudo systemctl start docker
```
- Create rootfs tarball and import to docker
```
export DOCKER_TAG=7.9.2009.minimal
./containerbuild.sh centos-7-x86_64.ks http://centos.mirror.ndchost.com/7.9.2009/isos/x86_64/CentOS-7-x86_64-Minimal-2009.iso
  
```
Creates docker image `centos:7.9.2009.minimal`


## Alternate method of creating docker image from ISO
This method mounts ISO on loop device and extracts Liveos readonly squash file system to create uncompressed rootfs. 

```
./iso-to-docker.sh
```