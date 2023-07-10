resource "random_id" "build_id" {
  byte_length = 8
}
module "bastion_node" {
  source = "github.com/mesosphere/vcenter-tools/modules/vmclone"

  node_name      = "kib-vcd-bastion-${random_id.build_id.hex}"
  ssh_public_key =  file(var.ssh_public_key)

  datastore_name       = var.datastore_name
  datastore_is_cluster = false
  vm_template_name = var.bastion_base_template
  resource_pool_name = var.vcd_bastion_resource_pool_name
  root_volume_size = var.root_volume_size #80
  vsphere_folder = var.vcd_bastion_vsphere_folder
  ssh_user = var.ssh_user
  custom_attribute_owner      = var.bastion_owner
  custom_attribute_expiration = "8h"
}

resource "null_resource" "set_lower_mtu" {
  provisioner "remote-exec" {
    inline = [
      "sudo ip link set ens192 mtu 1400",
    ]
  }
  connection {
    host = module.bastion_node.nat_address
    port = module.bastion_node.nat_ssh_port
    user = var.ssh_user
    agent = true
  }
}

resource "null_resource" "copy_upload_template_script" {
   provisioner "file" {
    source = "${path.module}/upload-template.sh"
    destination =  "/home/${var.ssh_user}/upload-template.sh"
  }
  connection {
    host = module.bastion_node.nat_address
    port = module.bastion_node.nat_ssh_port
    user = var.ssh_user
    agent = true
  }
}

resource "null_resource" "execute_upload_template_script" {
  provisioner "remote-exec" {
    inline = [
      "export PROD_VSPHERE_URL=${var.vsphere_url}",
      "export PROD_VSPHERE_USERNAME=${var.vsphere_username}",
      "export PROD_VSPHERE_PASSWORD=${var.vsphere_password}",
      "export VCD_URL=${var.vcd_url}",
      "export VCD_USERNAME=${var.vcd_org_username}",
      "export VCD_PASSWORD=${var.vcd_org_password}",
      "export PROD_TEMPLATE_NAME=${var.vm_template_name_to_upload}",
      "set -o errexit",
      "chmod +x /home/${var.ssh_user}/upload-template.sh",
      "/home/${var.ssh_user}/upload-template.sh >> /home/${var.ssh_user}/upload-template.log"
    ]
  }
  provisioner "remote-exec" {
    inline = [
      "cat /home/${var.ssh_user}/upload-template.log"
    ]
  }
  connection {
    host = module.bastion_node.nat_address
    port = module.bastion_node.nat_ssh_port
    user = var.ssh_user
    agent = true
  }
  depends_on = [
    null_resource.set_lower_mtu,
    null_resource.copy_upload_template_script
   ]
}


output "bastion_node_default_ip_address" {
  value = module.bastion_node.default_ip_address
}

output "bastion_node_ssh_nat_address" {
  value = "${module.bastion_node.nat_address}:${module.bastion_node.nat_ssh_port}"
}
