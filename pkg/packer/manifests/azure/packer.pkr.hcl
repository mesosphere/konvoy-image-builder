packer {
  required_plugins {
    azure = {
      version = ">= 1.3.1"
      source = "github.com/hashicorp/azure"
    }
    ansible = {
      version = ">= 1.1.0"
      source  = "github.com/hashicorp/ansible"
    }
    goss = {
      version = ">=3.1.5"
      source = "github.com/supershal/goss"
    }
  }
}

variable "ansible_override_files" {
  type    = list(string)
  default = []
}

variable "build_name" {
  type    = string
  default = ""
}

variable "build_name_extra" {
  type    = string
  default = ""
}

variable "capture_container_name" {
  type    = string
  default = "dkp-vhds"
}

variable "client_id" {
  type    = string
  default = "${env("AZURE_CLIENT_ID")}"
}

variable "client_secret" {
  type    = string
  default = "${env("AZURE_CLIENT_SECRET")}"
}

variable "cloud_environment_name" {
  type    = string
  default = "Public"
}

variable "distribution" {
  type = string
}

variable "distribution_version" {
  type = string
}

variable "disk_size" {
  type = number
  default = 80
}

variable "existing_ansible_ssh_args" {
  type    = string
  default = "${env("ANSIBLE_SSH_ARGS")}"
}

variable "extra_vars_file" {
  type    = string
  default = ""
}

variable "gallery_image_locations" {
  type = string
}

variable "gallery_image_name" {
  type = string
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

variable "image_version" {
  type    = string
  default = "latest"
}

variable "image_publisher" {
  type    = string
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

variable "plan_image_offer" {
  type    = string
  default = ""
}

variable "plan_image_publisher" {
  type    = string
  default = ""
}

variable "plan_image_sku" {
  type    = string
  default = ""
}

variable "private_virtual_network_with_public_ip" {
  type    = bool
  default = true
}

variable "resource_group_name" {
  type    = string
  default = "dkp"
}

variable "shared_image_gallery_name" {
  type    = string
  default = "dkp"
}

variable "subscription_id" {
  type    = string
  default = "${env("AZURE_SUBSCRIPTION_ID")}"
}

variable "tenant_id" {
  type    = string
  default = "${env("AZURE_TENANT_ID")}"
}

variable "vm_size" {
  type    = string
  default = "Standard_D2s_v3"
}

variable "ssh_bastion_username" {
  type    = string
  default = ""
}

variable "ssh_bastion_password" {
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

variable "ssh_username" {
  type    = string
  default = ""
}

variable "gallery_image_publisher" {
  type    = string
  default = ""
}

variable "gallery_image_offer" {
  type    = string
  default = ""
}

variable "gallery_name" {
  type    = string
  default = ""
}

variable "https_proxy" {
  type    = string
  default = ""
}

variable "http_proxy" {
  type    = string
  default = ""
}

variable "no_proxy" {
  type    = string
  default = "true"
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

variable "virtual_network_name" {
  type = string
  default = ""
}

variable "virtual_network_resource_group_name" {
  type = string
  default = ""
}

variable "virtual_network_subnet_name" {
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

variable "location" {
  type = string
  default = "westus"
}

variable "gallery_image_sku" {
  type = string
}

# "timestamp" template function replacement
locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

# All locals variables are generated from variables that uses expressions
# that are not allowed in HCL2 variables.
# Read the documentation for locals blocks here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/locals
locals {
  build_timestamp                    = "${local.timestamp}"
  shared_image_gallery_image_version = formatdate("YYYY.MM.DDhhmmss", timestamp())
  gallery_image_locations = split(",", var.gallery_image_locations)
  generated_managed_image_name = "${var.gallery_image_name}-${local.build_timestamp}"
  # clean_resource_name https://github.com/hashicorp/packer-plugin-azure/blob/a56b4569b8e71781cbe56bd09cd0cdc7419e4d48/builder/azure/common/template_funcs.go#L23
  managed_image_name = trimsuffix(local.generated_managed_image_name, "-_.")
  ansible_override_file_list = flatten([for override in var.ansible_override_files : concat(["--extra-vars"], [override])])
}

# source blocks are generated from your builders; a source can be referenced in
# build blocks. A build block runs provisioner and post-processors on a
# source. Read the documentation for source blocks here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/source
# could not parse template for following block: "template: hcl2_upgrade:2: bad character U+0060 '`'"

source "azure-arm" "kib_image" {
  azure_tags = {
    build_date             = legacy_isotime("June 7, 7:22:43pm 2014") # json template isotime
    build_timestamp        = local.build_timestamp
    containerd_version     = var.containerd_version
    distribution           = var.distribution
    distribution_version   = var.distribution_version
    gpu                    = var.gpu
    gpu_nvidia_version     = var.gpu_nvidia_version
    gpu_types              = var.gpu_types
    image_builder_version  = var.konvoy_image_builder_version
    kubernetes_cni_version = var.kubernetes_cni_version
    kubernetes_version     = var.kubernetes_full_version
  }
  client_id                         = var.client_id
  client_secret                     = var.client_secret
  cloud_environment_name            = var.cloud_environment_name
  image_offer                       = var.distribution
  image_publisher                   = var.image_publisher
  image_sku                         = var.distribution_version
  image_version                     = var.image_version
  location                          = length(local.gallery_image_locations) > 0 ? element(local.gallery_image_locations, 0) : var.location
  managed_image_name                = local.managed_image_name
  managed_image_resource_group_name = var.resource_group_name
  os_type                           = "Linux"
  os_disk_size_gb                   = var.disk_size
  plan_info {
    plan_name      = var.plan_image_sku
    plan_product   = var.plan_image_offer
    plan_publisher = var.plan_image_publisher
  }
  private_virtual_network_with_public_ip = var.private_virtual_network_with_public_ip
  shared_image_gallery_destination {
    gallery_name        = var.gallery_name
    image_name          = var.gallery_image_name
    image_version       = local.shared_image_gallery_image_version
    replication_regions = local.gallery_image_locations
    resource_group      = var.resource_group_name
  }
  ssh_key_exchange_algorithms         = ["curve25519-sha256@libssh.org", "ecdh-sha2-nistp256", "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "diffie-hellman-group14-sha1", "diffie-hellman-group1-sha1"]
  ssh_username                        = "packer"
  subscription_id                     = var.subscription_id
  tenant_id                           = var.tenant_id
  virtual_network_name                = var.virtual_network_name
  virtual_network_resource_group_name = var.virtual_network_resource_group_name
  virtual_network_subnet_name         = var.virtual_network_subnet_name
  vm_size                             = var.vm_size

  skip_create_image = var.dry_run
}

# a build block invokes sources and runs provisioning steps on them. The
# documentation for build blocks can be found here:
# https://www.packer.io/docs/templates/hcl_templates/blocks/build
build {
  sources = ["source.azure-arm.kib_image"]

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
      PROVIDER = "azure"
    }
    version = var.goss_version
  }

  provisioner "shell" {
    inline = ["rm -r  ${var.remote_folder}/.goss-dir"]
  }

  post-processor "manifest" {
    custom_data = {
      containerd_version     = "${var.containerd_version}"
      distribution           = "${var.distribution}"
      distribution_version   = "${var.distribution_version}"
      kubernetes_cni_version = "${var.kubernetes_cni_semver}"
      kubernetes_version     = "${var.kubernetes_full_version}"
      compute_gallery_image_id = "/subscriptions/${var.subscription_id}/resourceGroups/${var.resource_group_name}/providers/Microsoft.Compute/galleries/${var.shared_image_gallery_name}/images/${var.gallery_image_name}/versions/${local.shared_image_gallery_image_version}"
    }
    output = "${var.manifest_output}"
  }
}
