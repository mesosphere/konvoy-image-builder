terraform {
  required_providers {
    vsphere =  {
      version = "1.15.0"
    }
  }
}
provider "vsphere" {
  allow_unverified_ssl = false

}
variable "datacenter_name" {
  description = "The datacenter name"
  default     = "dc1"
}

variable "datastore_name" {
  description = "The datastore name"
  default     = "ovh-nfs"
}

variable "resource_pool_name" {
  description = "The resource pool name"
  default     = "cluster-api"
}

variable "network_name_airgapped" {
  description = "The network name"
  default     = "Airgapped"
}

variable "network_name_public" {
  description = "The network name"
  default     = "Public"
}

variable "bastion_vm_template" {
  description = "The VM template name for the bastion machine"
  default     = "kib-builder-template"
}


variable "root_user" {
  description = "The root user"
  default     = "builder"
}

data "vsphere_datacenter" "dc" {
  name = var.datacenter_name
}

data "vsphere_datastore" "datastore" {
  name          = var.datastore_name
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_resource_pool" "pool" {
  name          = var.resource_pool_name
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "network_public" {
  name          = var.network_name_public
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "network_airgapped" {
  name          = var.network_name_airgapped
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "bastion_template" {
  name          = var.bastion_vm_template
  datacenter_id = data.vsphere_datacenter.dc.id
}

resource "vsphere_virtual_machine" "konvoy-e2e-bastion" {
  name             = "bastion-host"
  resource_pool_id = data.vsphere_resource_pool.pool.id
  datastore_id     = data.vsphere_datastore.datastore.id

  num_cpus = 4
  memory   = 6144
  guest_id = "centos7_64Guest"


  network_interface {
    network_id = data.vsphere_network.network_public.id
  }

  network_interface {
    network_id = data.vsphere_network.network_airgapped.id
  }


  clone {
    template_uuid = data.vsphere_virtual_machine.bastion_template.id
    linked_clone  = false
  }

  disk {
    label = "disk0"
    datastore_id     = data.vsphere_datastore.datastore.id
    size = 80
  }
}

output "bastion_ip" {
  value = vsphere_virtual_machine.konvoy-e2e-bastion.guest_ip_addresses
}
