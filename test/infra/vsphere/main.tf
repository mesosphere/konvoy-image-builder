// FIXME: https://d2iq.atlassian.net/browse/D2IQ-99451
terraform {
  required_providers {
    vsphere = {
      source = "hashicorp/vsphere"
      version = "~> 2.4.0"
    }
  }
}

resource "random_id" "build_id" {
  byte_length = 8
}
module "bastion_node" {
  source = "github.com/mesosphere/vcenter-tools/modules/vmclone"

  node_name      = "konvoy-image-builder-bastion-${random_id.build_id.hex}"
  ssh_public_key = file(var.ssh_public_key)

  datastore_name              = var.datastore_name
  datastore_is_cluster        = false
  vm_template_name            = var.bastion_base_template
  resource_pool_name          = var.resource_pool_name
  root_volume_size            = var.root_volume_size #80
  vsphere_folder              = var.vsphere_folder
  ssh_user                    = var.ssh_user
  custom_attribute_owner      = var.bastion_owner
  custom_attribute_expiration = "4h"
  vsphere_network             = var.vsphere_network
}

output "bastion_node_ssh_user" {
  value = var.ssh_user
}

output "bastion_node_ssh_nat_address" {
  value = module.bastion_node.nat_address
}

output "bastion_node_ssh_nat_port" {
  value = module.bastion_node.nat_ssh_port
}
