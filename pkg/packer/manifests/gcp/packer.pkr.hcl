packer {
  required_plugins {
    googlecompute = {
      version = ">= 1.0.11"
      source = "github.com/hashicorp/googlecompute"
    }
    ansible = {
      version = ">= 1.1.0"
      source  = "github.com/hashicorp/ansible"
    }
  }
}

variable "account_file" {
  type    = string
  default = "${env("GOOGLE_APPLICATION_CREDENTIALS")}"
}

variable "build_name" {
  type    = string
  default = ""
}

variable "build_name_extra" {
  type    = string
  default = ""
}

variable "disk_size" {
  type    = string
  default = "80"
}

variable "disk_type" {
  type    = string
  default = "pd-standard"
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

variable "konvoy_image_builder_version" {
  type    = string
  default = "0.0.1"
}

variable "kubernetes_full_version" {
  type    = string
  default = ""
}

variable "manifest_output" {
  type    = string
  default = "manifest.json"
}

variable "network" {
  type    = string
  default = ""
}

variable "project_id" {
  type    = string
  default = ""
}

variable "region" {
  type    = string
  default = ""
}

variable "http_proxy" {
  type    = string
  default = ""
}

variable "https_proxy" {
  type    = string
  default = ""
}

variable "source_image" {
  type    = string
  default = ""
}

variable "ssh_bastion_host" {
  type    = string
  default = ""
}

variable "ssh_bastion_private_key_file" {
  type    = string
  default = ""
}

variable "ssh_bastion_username" {
  type    = string
  default = ""
}

variable "ssh_bastion_password" {
  type    = string
  default = ""
}

variable "distribution_family" {
  type    = string
  default = ""
}

variable "distribution" {
  type    = string
  default = ""
}

variable "distribution_version" {
  type    = string
  default = ""
}

variable "ssh_username" {
  type    = string
  default = ""
}

variable "no_proxy" {
  type    = string
  default = ""
}

variable "containerd_version" {
  type    = string
  default = ""
}

variable "kubernetes_cni_version" {
  type    = string
  default = ""
}

variable "kubernetes_cni_semver" {
  type    = string
  default = ""
}

variable "ansible_override_files" {
  type    = list(string)
  default = []
}

variable "existing_ansible_ssh_args" {
  type    = string
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

# "timestamp" template function replacement
locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

# All locals variables are generated from variables that uses expressions
# that are not allowed in HCL2 variables.
# Read the documentation for locals blocks here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/locals
locals {
  build_timestamp      = local.timestamp
  zone                 = "${var.region}-a"
  generated_image_name = "konvoy-${var.build_name}-${var.kubernetes_full_version}-${local.build_timestamp}"
  # clean_resource_name https://github.com/hashicorp/packer-plugin-googlecompute/blob/81d8d5a740c0d7fb0b02be93133ac17a11557f34/builder/googlecompute/template_funcs.go#L20
  image_name           = regex_replace(lower(local.generated_image_name), "[^-a-z0-9]", "-")
  ansible_override_file_list = flatten([for override in var.ansible_override_files : concat(["--extra-vars"], [override])])
}

# source blocks are generated from your builders; a source can be referenced in
# build blocks. A build block runs provisioner and post-processors on a
# source. Read the documentation for source blocks here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/source
# could not parse template for following block: "template: hcl2_upgrade:2: bad character U+0060 '`'"

source "googlecompute" "kib_image" {
  disk_size    = var.disk_size
  disk_type    = var.disk_type
  image_family = replace("konvoy-${var.build_name}-${var.kubernetes_full_version}", ".", "-")
  image_labels = {
    build_date             = lower(regex_replace(legacy_isotime("June 7, 7:22:43pm 2014"), "[.: ,]+", "-")) # json template isotime
    build_timestamp        = lower(regex_replace(local.build_timestamp, "[.: ,]+", "-"))
    containerd_version     = lower(regex_replace(var.containerd_version, "[.: ,]+", "-"))
    gpu                    = lower(regex_replace(var.gpu, "[.: ,]+", "-"))
    gpu_nvidia_version     = lower(regex_replace(var.gpu_nvidia_version, "[.: ,]+", "-"))
    gpu_types              = lower(regex_replace(var.gpu_types, "[.: ,]+", "-"))
    image_builder_version  = lower(regex_replace(var.konvoy_image_builder_version, "[.: ,]+", "-"))
    kubernetes_cni_version = lower(regex_replace(var.kubernetes_cni_version, "[.: ,]+", "-"))
    kubernetes_version     = lower(regex_replace(var.kubernetes_full_version, "[.: ,]+", "-"))
  }
  image_name                  = local.image_name
  network                     = var.network
  project_id                  = var.project_id
  region                      = var.region
  source_image                = var.source_image
  source_image_family         = var.distribution_family
  ssh_key_exchange_algorithms = ["curve25519-sha256@libssh.org", "ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "diffie-hellman-group14-sha1", "diffie-hellman-group1-sha1"]
  ssh_username                = var.ssh_username
  wait_to_add_ssh_keys        = "20s"
  zone                        = local.zone

  skip_create_image = var.dry_run
}

# a build block invokes sources and runs provisioning steps on them. The
# documentation for build blocks can be found here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/build
build {
  sources = ["source.googlecompute.kib_image"]

  provisioner "shell" {
    remote_folder = "${var.remote_folder}"
    environment_vars = ["HTTP_PROXY=${var.http_proxy}", "http_proxy=${var.http_proxy}", "HTTPS_PROXY=${var.https_proxy}", "https_proxy=${var.https_proxy}", "NO_PROXY=${var.no_proxy}", "no_proxy=${var.no_proxy}", "BUILD_NAME=${var.build_name}"]
    inline           = ["if [ $BUILD_NAME != \"ubuntu-1804\" ]; then exit 0; fi", "while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for cloud-init...'; sleep 1; done", "sudo apt-get -qq update && sudo DEBIAN_FRONTEND=noninteractive apt-get -qqy install python python-pip"]
  }

  provisioner "ansible" {
    ansible_env_vars = ["ANSIBLE_SSH_ARGS='${var.existing_ansible_ssh_args}'", "ANSIBLE_REMOTE_TEMP='${var.remote_folder}/.ansible/'"]
    extra_arguments  = concat(local.ansible_override_file_list, ["--scp-extra-args", "'-O'"])
    playbook_file    = "${path.cwd}/ansible/provision.yaml"
    user             = "${var.ssh_username}"
  }

  post-processor "manifest" {
    custom_data = {
      containerd_version     = "${var.containerd_version}"
      distribution           = "${var.distribution}"
      distribution_version   = "${var.distribution_version}"
      kubernetes_cni_version = "${var.kubernetes_cni_semver}"
      kubernetes_version     = "${var.kubernetes_full_version}"
    }
    output = "${var.manifest_output}"
  }
}
