data "amazon-ami" "flatcar" {
  filters = {
    virtualization-type = "hvm"
    name                = "Flatcar-stable-*-hvm"
  }
  owners      = ["075585003325"]
  most_recent = true
}
local "flatcar_version" {
  expression = format(
    "packer: {distribution_version: %q, source_ami: %q}",
    element(split("-", data.amazon-ami.flatcar.name), 2),
    data.amazon-ami.flatcar.id
    )
}
source "file" "basic-example" {
  content =  local.flatcar_version
  target =  "flatcar-version.yaml"
}
build {
  sources = ["sources.file.basic-example"]
}
