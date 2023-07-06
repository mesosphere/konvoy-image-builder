variable "datastore_name" {
  description = "The datastore name"
  default     = "ci-kib"
}

variable "bastion_base_template" {
  description = "base template name"
  default     = "os-qualification-templates/d2iq-base-RockyLinux-9.1"
}

variable "resource_pool_name" {
  description = "The resource pool name"
  default     = "cluster-api"
}

variable "root_volume_size" {
  description = "Disk size for root volume"
  default     = 80
}

variable "vsphere_folder" {
  description = "folder name to store the VM"
  default     = "cluster-api"
}

variable "ssh_user" {
  description = "The root user"
  default     = "centos"
}


variable "bastion_owner" {
  description = "The root user"
  default     = "Konvoy image builder"
}

variable "ssh_public_key" {
  description = "Path to the SSH Public key. for example: ~/.ssh/id_rsa.pub"
  type        = string
}

variable "vsphere_network" {
  description = "vsphere network"
  type        = string
  default     = "Airgapped"
}
