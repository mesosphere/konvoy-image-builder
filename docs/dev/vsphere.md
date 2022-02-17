# Create OS images for vsphere

Creating vsphere images is a two step process
- Create base OS image
- Create OS image for CAPI using the base OS image

## Create base OS image
----
Creating base OS image form DVD ISO files is a one time process. A base vsphere template will be created as result of building base image.
This base image template can be used to create more customized images.

[TODO: Add steps and screenshots]

## Create vSphere template image for CAPI
----
The process to create base OS image will create a vSphere template in Vsphere. This template is a linked clone of the VM so any VM created from template will be provisioned quickly.

**Environment Variables for vsphere:**
NOTE: use vsphere server URL without `http://` or `https://`
```bash
VSPHERE_SERVER=example.vsphere.url
VSPHERE_USERNAME=user@example.vsphere.url
VSPHERE_PASSWORD=example_password
```
**Environment Variables for RedHat subscription:**
```bash
RHSM_USER=example_user
RHSM_PASS=example_password
```
**Packer variables for vsphere:**
Following variables needed for the vsphere. Add this configuration in the `image.yaml`
```yaml
packer:
  cluster: "example_zone"
  datacenter: "example_datacenter"
  datastore: "example_datastore"
  folder: "example_folder"
  insecure_connection: "false"
  network: "example_network"
  resource_pool: "example_resource_pool"
  template: "example_base_OS_template_name"
  vsphere_guest_os_type: "exampple_rhel7_64Guest"
  guest_os_type: "example_rhel7-64"
  #goss params
  distribution: "example_RHEL"
  distribution_version: "example_7.9"
```

## Create template image on vsphere
```bash
konvoy-image build path/to/image.yaml
```
checkout example image configuration at: `<project_root>/images/ova/` directory