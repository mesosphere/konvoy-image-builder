packer {
  required_plugins {
    vsphere = {
      version = ">= 1.0.8"
      source  = "github.com/hashicorp/vsphere"
    }
    ansible = {
      version = ">= 1.0.2"
      source  = "github.com/hashicorp/ansible"
    }
  }
}

variable "ansible_extra_vars" {
  type    = string
  default = ""
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

variable "ib_version" {
  type    = string
  default = "${env("IB_VERSION")}"
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

variable "ssh_bastion_host" {
  type = string
  default = ""
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

variable "goss_binary" {
  type = string
  default = "/usr/local/bin/goss-amd64"
}

variable "goss_entry_file" {
  type    = string
  default = null
}

variable "goss_inspect_mode" {
  type    = bool
  default = false
}

variable "goss_tests_dir" {
  type    = string
  default = null
}

variable "goss_url" {
  type    = string
  default = null
}

variable "goss_vars_file" {
  type    = string
  default = null
}
variable "goss_format" {
  type    = string
  default = null
}
variable "goss_format_options" {
  type    = string
  default = null
}
variable "goss_arch" {
  type    = string
  default = null
}
variable "goss_version" {
  type    = string
  default = null
}

variable "dry_run" {
  type    = bool
  default = false
}

variable "remote_folder" {
  type    = string
  default = "/tmp"
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
  ssh_bastion_password         = var.ssh_bastion_password
  ssh_bastion_private_key_file = var.ssh_bastion_private_key_file
  ssh_bastion_username         = var.ssh_bastion_username
  vm_name                      = "konvoy-${var.build_name}-${var.kubernetes_full_version}-${local.build_timestamp}"
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
  ssh_bastion_password         = local.ssh_bastion_password
  ssh_bastion_private_key_file = local.ssh_bastion_private_key_file
  ssh_bastion_username         = local.ssh_bastion_username
  ssh_key_exchange_algorithms  = ["curve25519-sha256@libssh.org", "ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "diffie-hellman-group14-sha1", "diffie-hellman-group1-sha1"]
  ssh_password                 = var.ssh_password
  ssh_private_key_file         = var.ssh_private_key_file
  ssh_timeout                  = "4h"
  ssh_username                 = var.ssh_username
  template                     = var.template
  username                     = var.vsphere_username
  vcenter_server               = var.vcenter_server
  vm_name                      = local.vm_name
  resource_pool                = var.resource_pool

  create_snapshot     = !var.dry_run
  convert_to_template = !var.dry_run
}

build {
  sources = ["source.vsphere-clone.kib_image"]

  provisioner "ansible" {
    ansible_env_vars = ["ANSIBLE_SSH_ARGS='${var.existing_ansible_ssh_args} -o IdentitiesOnly=yes -o HostkeyAlgorithms=+ssh-rsa -o PubkeyAcceptedAlgorithms=+ssh-rsa'", "ANSIBLE_REMOTE_TEMP='${var.remote_folder}/.ansible/'"]
    extra_arguments  = ["--extra-vars", "${var.ansible_extra_vars}"]
    playbook_file    = "./ansible/provision.yaml"
    user             = var.ssh_username
  }

  provisioner "shell" {
    inline = ["mkdir -p ${var.remote_folder}/.goss-dir"]
  }

  provisioner "file" {
    destination = "${var.remote_folder}/.goss-dir/goss"
    direction   = "upload"
    max_retries = "10"
    source      = var.goss_binary
  }


  provisioner "goss" {
    arch           = var.goss_arch
    download_path  = "${var.remote_folder}/.goss-dir/goss"
    format         = var.goss_format
    format_options = var.goss_format_options
    goss_file      = var.goss_entry_file
    inspect        = var.goss_inspect_mode
    skip_install   = true
    tests          = var.goss_tests_dir == null ? null : [var.goss_tests_dir]
    url            = var.goss_url
    use_sudo       = true
    vars_env = {
      HTTPS_PROXY = var.https_proxy
      HTTP_PROXY  = var.http_proxy
      NO_PROXY    = var.no_proxy
      http_proxy  = var.http_proxy
      https_proxy = var.https_proxy
      no_proxy    = var.no_proxy
    }
    vars_file = var.goss_vars_file
    vars_inline = {
      ARCH     = "amd64"
      OS       = lower(var.distribution)
      PROVIDER = "amazon"
    }
    version = var.goss_version
  }

  provisioner "shell" {
    inline = ["rm -r  ${var.remote_folder}/.goss-dir"]
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
      ib_version               = "${var.ib_version}"
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
    inline = ["echo 'destroying VM ${local.vm_name}' && govc vm.destroy ${local.vm_name}"]
    environment_vars =[
        "GOVC_URL=${var.vcenter_server}",
        "GOVC_USERNAME=${var.vsphere_username}",
        "GOVC_PASSWORD=${var.vsphere_password}"
    ]
  }
}
