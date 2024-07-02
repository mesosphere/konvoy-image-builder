packer {
  required_plugins {
    vsphere = {
      version = ">= 1.0.8"
      source  = "github.com/hashicorp/vsphere"
    }
    ansible = {
      version = ">= 1.1.0"
      source  = "github.com/hashicorp/ansible"
    }
    sshkey = {
      version = ">= 1.0.1"
      source  = "github.com/ivoronin/sshkey"
    }
  }
}

variable "ansible_override_files" {
  type    = list(string)
  default = []
}

variable "build_name_extra" {
  type    = string
  default = ""
}

variable "build_name" {
  type    = string
  default = ""
}

variable "cluster" {
  type    = string
  default = ""
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
  default = ""
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

variable "existing_ansible_ssh_args" {
  type    = string
  default = "${env("ANSIBLE_SSH_ARGS")}"
}

variable "export_manifest" {
  type    = string
  default = "none"
}

variable "firmware" {
  type    = string
  default = "bios"
}

variable "folder" {
  type    = string
  default = ""
}

variable "gpu" {
  type    = string
  default = "false"
}

variable "gpu_nvidia_version" {
  type    = string
  default = ""
}

variable "gpu_types" {
  type    = string
  default = ""
}

variable "guest_os_type" {
  type = string
}

variable "insecure_connection" {
  type    = string
  default = "false"
}

variable "kubernetes_full_version" {
  type    = string
  default = ""
}

variable "linked_clone" {
  type    = string
  default = "true"
}

variable "manifest_output" {
  type    = string
  default = "manifest.json"
}

variable "memory" {
  type    = string
  default = "8192"
}

variable "ssh_agent_auth" {
  type    = string
  default = "false"
}

variable "ssh_password" {
  type    = string
  default = env("SSH_PASSWORD")
  sensitive = true
}

variable "ssh_private_key_file" {
  type    = string
  default = env("SSH_PRIVATE_KEY_FILE")
  sensitive = true
}

variable "ssh_public_key" {
  type    = string
  default = env("SSH_PUBLIC_KEY")
  sensitive = true
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
}

variable "vsphere_password" {
  type    = string
  default = "${env("VSPHERE_PASSWORD")}"
}

variable "vsphere_username" {
  type    = string
  default = "${env("VSPHERE_USERNAME") == "" ? env("VSPHERE_USER") : env("VSPHERE_USERNAME") }"
}

variable "vsphere_datacenter" {
  type = string
  default = "${env("VSPHERE_DATACENTER")}"
}

variable "ssh_bastion_host" {
  type = string
  default = ""
}

variable "ssh_bastion_port" {
  type = string
  default = "22"
}

variable "ssh_bastion_password" {
  type = string
  default = ""
}
variable "ssh_bastion_private_key_file" {
  type = string
  default = ""
}
variable "ssh_bastion_username" {
  type = string
  default = ""
}

variable "http_proxy" {
  type = string
  default = ""
}
variable "https_proxy" {
  type = string
  default = ""
}
variable "no_proxy" {
  type = string
  default = ""
}


variable "template" {
  type = string
  default = ""
}
variable "konvoy_image_builder_version" {
  type = string
  default = ""
}

variable "resource_pool" {
  type = string
  default = ""
}
variable "containerd_version" {
  type = string
  default = ""
}
variable "datacenter" {
  type = string
  default = ""
}
variable "network" {
  type = string
  default = ""
}
variable "host" {
  type = string
  default = ""
}
variable "custom_role" {
  type = string
  default = ""
}
variable "distro_arch" {
  type = string
  default = ""
}
variable "distro_name" {
  type = string
  default = ""
}
variable "distro_version" {
  type = string
  default = ""
}

variable "kubernetes_cni_semver" {
  type = string
  default = ""
}

variable "kubernetes_semver" {
  type = string
  default = ""
}

variable "kubernetes_source_type" {
  type = string
  default = ""
}

variable "kubernetes_typed_version" {
  type = string
  default = ""
}

variable "os_display_name" {
  type = string
  default = ""
}

variable "dry_run" {
  type    = bool
  default = false
}

variable "remote_folder" {
  type    = string
  default = "/tmp"
}

data "sshkey" "kibkey" {
  name = "konvoy-image-builder-tmpkey"
}

# "timestamp" template function replacement
locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

# All locals variables are generated from variables that uses expressions
# that are not allowed in HCL2 variables.
# Read the documentation for locals blocks here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/locals
locals {
  build_timestamp              = local.timestamp
  ssh_bastion_host             = var.ssh_bastion_host
  ssh_bastion_port             = var.ssh_bastion_port
  ssh_bastion_password         = var.ssh_bastion_password
  ssh_bastion_private_key_file = var.ssh_bastion_private_key_file
  ssh_bastion_username         = var.ssh_bastion_username
  vm_name                      = "konvoy-${var.build_name}-${var.kubernetes_full_version}-${local.build_timestamp}"

  # if only a public key is given we expect the private key to be loaded into ssh-agent
  ssh_agent_auth = var.ssh_agent_auth  != "false" ? true : var.ssh_private_key_file == "" && var.ssh_public_key != ""

  # inject generated key if no agent auth or private key is given
  ssh_private_key_file = var.ssh_private_key_file != "" ? var.ssh_private_key_file : local.ssh_agent_auth ? "" : data.sshkey.kibkey.private_key_path
  # when ssh_private_key_file uses the generated key inject its public key
  ssh_public_key = local.ssh_private_key_file == data.sshkey.kibkey.private_key_path ? data.sshkey.kibkey.public_key : chomp(var.ssh_public_key)
  ssh_password_hash = var.ssh_password != "" ? bcrypt(var.ssh_password): ""
  # prepare cloud-init
  cloud_init = <<EOF
#cloud-config
users:
  - name: ${var.ssh_username}
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: sudo, wheel
    lock_passwd: true
    ssh_authorized_keys:
      - ${local.ssh_public_key}
EOF
  ignition_config = <<EOF
{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.3.0"
  },
  "networkd": {},
  "passwd": {
    "users": [
      {
        "groups": [
          "wheel",
          "sudo",
          "docker"
        ],
        "name": "${var.ssh_username}",
        "passwordHash": "${local.ssh_password_hash}",
        "sshAuthorizedKeys": [
          "${local.ssh_public_key}"
        ]
      }
    ]
  },
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      },
      {
        "mask": true,
        "name": "update-engine.service"
      },
      {
        "mask": true,
        "name": "locksmithd.service"
      }
    ]
  }
}
EOF

  configuration_parameters_cloud_init = local.ssh_public_key != "" ? {
    "guestinfo.userdata" = base64encode(local.cloud_init),
    "guestinfo.userdata.encoding" = "base64",
    "guestinfo.metadata" = ""
    "guestinfo.metadata.encoding" = "base64"
  } : {}
  configuration_parameters_ignition = {
    "guestinfo.ignition.config.data" = base64encode(local.ignition_config),
    "guestinfo.ignition.config.data.encoding" = "base64",
  }
  configuration_parameters = var.distribution == "flatcar" ? local.configuration_parameters_ignition : local.configuration_parameters_cloud_init
  ansible_override_file_list = flatten([for override in var.ansible_override_files : concat(["--extra-vars"], [override])])
}

# source blocks are generated from your builders; a source can be referenced in
# build blocks. A build block runs provisioner and post-processors on a
# source. Read the documentation for source blocks here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/source
source "vsphere-clone" "kib_image" {
  CPUs                         = var.cpu
  RAM                          = var.memory
  cluster                      = var.cluster
  communicator                 = "ssh"
  cpu_cores                    = var.cpu_cores
  datacenter                   = var.datacenter
  datastore                    = var.datastore
  folder                       = var.folder
  host                         = var.host
  insecure_connection          = var.insecure_connection
  linked_clone                 = var.linked_clone
  network                      = var.network
  password                     = var.vsphere_password
  ssh_agent_auth               = var.ssh_agent_auth
  ssh_bastion_host             = local.ssh_bastion_host
  ssh_bastion_port             = local.ssh_bastion_port
  ssh_bastion_password         = local.ssh_bastion_password
  ssh_bastion_private_key_file = local.ssh_bastion_private_key_file
  ssh_bastion_username         = local.ssh_bastion_username
  ssh_key_exchange_algorithms  = ["curve25519-sha256@libssh.org", "ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "diffie-hellman-group14-sha1", "diffie-hellman-group1-sha1"]
  ssh_password                 = var.ssh_password
  ssh_private_key_file         = local.ssh_private_key_file
  ssh_timeout                  = "4h"
  ssh_username                 = var.ssh_username
  template                     = var.template
  username                     = var.vsphere_username
  vcenter_server               = var.vcenter_server
  vm_name                      = local.vm_name
  resource_pool                = var.resource_pool

  cd_label = "cidata"
  cd_content = {
    "/user-data"       = local.cloud_init,
    "/user-data.txt"       = local.cloud_init,
    "/meta-data"       = "",
  }

  // try injecting cloud-init via guestinfo
  configuration_parameters = local.configuration_parameters

  create_snapshot     = !var.dry_run
  convert_to_template = !var.dry_run
}

build {
  sources = ["source.vsphere-clone.kib_image"]

  provisioner "ansible" {
    ansible_env_vars = ["ANSIBLE_SSH_ARGS='${var.existing_ansible_ssh_args} -o IdentitiesOnly=yes -o HostkeyAlgorithms=+ssh-rsa -o PubkeyAcceptedAlgorithms=+ssh-rsa'", "ANSIBLE_REMOTE_TEMP='${var.remote_folder}/.ansible/'"]
    extra_arguments  = concat(local.ansible_override_file_list, ["--scp-extra-args", "'-O'"])
    playbook_file    = "${path.cwd}/ansible/provision.yaml"
    user             = var.ssh_username
  }

  post-processor "manifest" {
    custom_data = {
      build_date               = "${timestamp()}"
      build_name               = "${var.build_name}"
      build_timestamp          = "${local.build_timestamp}"
      build_type               = "node"
      containerd_version       = "${var.containerd_version}"
      custom_role              = "${var.custom_role}"
      disk_size                = "${var.disk_size}"
      distro_arch              = "${var.distro_arch}"
      distro_name              = "${var.distro_name}"
      distro_version           = "${var.distro_version}"
      firmware                 = "${var.firmware}"
      gpu                      = "${var.gpu}"
      gpu_nvidia_version       = "${var.gpu_nvidia_version}"
      gpu_types                = "${var.gpu_types}"
      guest_os_type            = "${var.guest_os_type}"
      kubernetes_cni_version   = "${var.kubernetes_cni_semver}"
      kubernetes_version       = "${var.kubernetes_full_version}"
      kubernetes_source_type   = "${var.kubernetes_source_type}"
      kubernetes_typed_version = "${var.kubernetes_typed_version}"
      os_name                  = "${var.os_display_name}"
      vsphere_guest_os_type    = "${var.vsphere_guest_os_type}"
      distribution             = "${var.distribution}"
      distribution_version     = "${var.distribution_version}"
    }
    name       = "packer-manifest"
    output     = "${var.manifest_output}"
    strip_path = true
  }
  post-processor "shell-local" {
    inline = [ "if ${var.dry_run}; then echo 'destroying VM ${local.vm_name} with command: govc vm.destroy -dc=${var.vsphere_datacenter} ${local.vm_name}'; govc vm.destroy -dc=${var.vsphere_datacenter} ${local.vm_name}; fi"]
    environment_vars =[
        "GOVC_URL=${var.vcenter_server}",
        "GOVC_USERNAME=${var.vsphere_username}",
        "GOVC_PASSWORD=${var.vsphere_password}"
    ]
  }
}
