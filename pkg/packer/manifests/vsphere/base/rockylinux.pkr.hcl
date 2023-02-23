packer {
  required_plugins {
    vsphere = {
      version = ">= 1.0.8"
      source = "github.com/hashicorp/vsphere"
    }
  }
}

variable "vsphere_user" {
  type    = string
  default = "${env("VSPHERE_USERNAME") == "" ? env("VSPHERE_USER") : env("VSPHERE_USERNAME") }"
}

variable "cluster" {
  type    = string
  default = env("VSPHERE_CLUSTER")
}

variable "cpu" {
  type    = string
  default = "4"
}

variable "cpu_cores" {
  type    = string
  default = "1"
}

variable "datastore" {
  type    = string
  default = env("VSPHERE_DATASTORE")
}


variable "disk_size" {
  type    = string
  default = "20480"
}

variable "distribution" {
  type    = string
  default = ""
}

variable "distribution_version" {
  type    = string
  default = ""
}

variable "firmware" {
  type    = string
  default = "bios"
}

variable "folder" {
  type    = string
  default = ""
}

variable "guest_os_type" {
  type = string
}

variable "ib_version" {
  type    = string
  default = "${env("IB_VERSION")}"
}

variable "insecure_connection" {
  type    = string
  default = "false"
}

variable "iso_url" {
  type    = string
  default = env("ISO_URL")
}

variable "iso_checksum" {
  type    = string
  default = env("ISO_SHA256_CHECKSUM")
}

variable "memory" {
  type    = string
  default = "8192"
}

variable "ssh_password" {
  type    = string
  default = ""
}

variable "ssh_timeout" {
  type    = string
  default = "60m"
}

variable "ssh_username" {
  type    = string
  default = env("SSH_USERNAME")
}

variable "vcenter_server" {
  type    = string
  default = "${env("VSPHERE_SERVER")}"
}

variable "vsphere_guest_os_type" {
  type = string
  default = ""
}

variable "vsphere_password" {
  type    = string
  default = "${env("VSPHERE_PASSWORD")}"
}

variable "network" {
  default = env("VSPHERE_NETWORK")
}
variable "resource_pool" {
  default = env("VSPHERE_RESOURCE_POOL")
}

variable "build_name" {
  default = "rockylinux-base"
}

variable "datacenter" {
  default=env("VSPHERE_DATACENTER")
}

variable "distro_arch" {
  default = ""
}
variable "distro_name" {
  default = ""
}
variable "distro_version" {
  default = ""
}

variable "os_iso_path" {
  type    = string
  default = ""
}

variable "private_key_path" {
  type    = string
  default = env("SSH_PRIVATE_KEY_FILE")
}

variable "public_key_contents" {
  type    = string
  default = env("SSH_PUBLIC_KEY_CONTENTS")
}
# "timestamp" template function replacement
locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

data "sshkey" "install" {}

source "vsphere-iso" "base-template" {
  CPUs                         = var.cpu
  RAM                          = var.memory
  cluster                      = var.cluster
  disk_controller_type         = ["pvscsi"]
  guest_os_type                = var.vsphere_guest_os_type
  network_adapters {
    network_card = "vmxnet3"
    network      = var.network
  }
  boot_wait = "10s"
  boot_command = ["<up><tab> inst.text inst.ks=hd:sr1:/bootfile.cfg<enter><wait>"]
  firmware = var.firmware

  cd_content = {
    "/bootfile.cfg" = local.kickstart_el9,
    # make it cloud-config compatible
    "/user-data" = local.kickstart_el9,
    "/meta-data" = ""
  }

  cd_label = "rockylinux-9.1"

  storage {
    disk_size             = 40000
  }

  communicator                 = "ssh"
  cpu_cores                    = "${var.cpu_cores}"
  datacenter                   = "${var.datacenter}"
  resource_pool                = "${var.resource_pool}"
  datastore                    = "${var.datastore}"
  folder                       = "${var.folder}"
  insecure_connection          = "${var.insecure_connection}"
  iso_url                      = "${var.iso_url}"
  iso_checksum                 = "${var.iso_checksum}"
  password                     = "${var.vsphere_password}"
  #ssh_bastion_host             = "${local.ssh_bastion_host}"
  #ssh_bastion_password         = "${local.ssh_bastion_password}"
  ssh_private_key_file         = "${var.private_key_path}"
  ssh_clear_authorized_keys    = false
  #ssh_bastion_username         = "${local.ssh_bastion_username}"
  ssh_key_exchange_algorithms  = ["curve25519-sha256@libssh.org", "ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "diffie-hellman-group14-sha1", "diffie-hellman-group1-sha1"]
  ssh_password                 = "${var.ssh_password}"
  ssh_timeout                  = "4h"
  ssh_username                 = "${var.ssh_username}"
  username                     = "${var.vsphere_user}"
  vcenter_server               = "${var.vcenter_server}"
  vm_name                      = "base-rocky-9.1"

  create_snapshot     = true
  convert_to_template = true
  iso_paths              = ["${var.os_iso_path}"]
}

build {
    sources = [
        "source.vsphere-iso.base-template",
    ]
}

# hardcoded text blobs
locals {
  kickstart_el9 = <<EOF
   repo --name="AppStream" --baseurl="http://download.rockylinux.org/pub/rocky/9/AppStream/x86_64/os/"
  cdrom
  # Use text install
  text

  # Don't run the Setup Agent on first boot
  firstboot --disabled
  eula --agreed

  # Keyboard layouts
  keyboard --vckeymap=us --xlayouts='us'

  # System language
  lang en_US.UTF-8

  # Network information
  network --bootproto=dhcp --onboot=on --ipv6=auto --activate --hostname=rockylinux

  # Lock Root account
  rootpw --lock

  # Create builder user
  #authselect --enableshadow --passalgo=sha512 --kickstart
  user --groups=wheel --name=builder --gecos="builder"
  sshkey --username=${var.ssh_username} "${var.public_key_contents}"

  # System services
  selinux --permissive
  firewall --disabled
  services --enabled="NetworkManager,sshd,chronyd"

  # System timezone
  timezone UTC

  # System booloader configuration
  bootloader --location=mbr --boot-drive=sda
  zerombr
  clearpart --all --initlabel --drives=sda
  part / --fstype="ext4" --grow --asprimary --label=slash --ondisk=sda

  skipx

  %packages --ignoremissing --excludedocs
  @^minimal-environment
  @core
  openssh-server
  open-vm-tools
  sudo
  sed
  python3

  # unnecessary firmware
  -aic94xx-firmware
  -atmel-firmware
  -b43-openfwwf
  -bfa-firmware
  -ipw2100-firmware
  -ipw2200-firmware
  -ivtv-firmware
  -iwl*-firmware
  -libertas-usb8388-firmware
  -ql*-firmware
  -rt61pci-firmware
  -rt73usb-firmware
  -xorg-x11-drv-ati-firmware
  -zd1211-firmware
  -cockpit
  -quota
  -alsa-*
  -fprintd-pam
  -intltool
  -microcode_ctl
  %end

  %addon com_redhat_kdump --disable
  %end

  reboot
  # Enable/disable the following services
  services --enabled=sshd
  %post --erroronfail --nochroot --logfile=/mnt/sysimage/var/log/ks-post.log

  update-ca-trust force-enable
  # Remove the package cache
  # Disable quiet boot and splash screen
  sed --follow-symlinks -i "s/ rhgb quiet//" /mnt/sysimage/etc/default/grub
  sed --follow-symlinks -i "s/ rhgb quiet//" /mnt/sysimage/boot/grub2/grubenv

  echo '${var.ssh_username} ALL=(ALL) NOPASSWD: ALL' >/etc/sudoers.d/${var.ssh_username} && chmod 440 /etc/sudoers.d/${var.ssh_username}

  yum -y clean all

  swapoff -a
  rm -f /swapfile
  sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab

  # Ensure on next boot that network devices get assigned unique IDs.
  #sed -i '/^\(HWADDR\|UUID\)=/d' /etc/sysconfig/network-scripts/ifcfg-*
  %end
  EOF
}
