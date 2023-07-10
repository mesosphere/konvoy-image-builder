
variable "ssh_public_key" {
  description = "Path to the SSH Public key. for example: ~/.ssh/id_rsa.pub"
  type        = string
}

variable "ssh_user" {
  description = "SSH Username"
  default     = "konvoy"
  type        = string
}

variable "bastion_base_template" {
  description = "Name of the template used for bastion"
  default     = "d2iq-base-templates/d2iq-base-RockyLinux-9.1"
  type        = string
}

variable "vcd_bastion_vsphere_folder" {
  description = "Name of the folder where the bastion vm will be created"
  default     = "users/users-ksphere-platform"
  type        = string
}

variable "datastore_name" {
  description = "Name of the datastore"
  type        = string
}

variable "vcd_bastion_resource_pool_name" {
  description = "Name of the vsphere resource pool"
  default     = "users-ksphere-platform"
  type        = string
}

variable "bastion_owner" {
  description = "owner name to tag the bastion host with"
  default     = "kib-e2e"
  type        = string
}

variable "root_volume_size" {
  description = "size of the root volume"
  default     = 80
  type        = number
}

variable "vsphere_url" {
  description = "URL production vsphere environment"
  type        = string
}

variable "vsphere_username" {
  description = "username for production vsphere environment"
  type        = string
}

variable "vsphere_password" {
  description = "password for production vSphere environment"
  type        = string
  #sensitive = true
}

variable "vcd_url" {
  description = "URL for cloud director environment"
  type        = string
}

variable "vcd_org_username" {
  description = "username for cloud director tenant organization"
  type        = string
}

variable "vcd_org_password" {
  description = "password for cloud director tenant organization"
  type        = string
  #sensitive = true
}

variable "vm_template_name_to_upload" {
  description = "name of the VM template to upload"
  type        = string
}