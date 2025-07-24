# Create OS images for vSphere

Creating vSphere images is a two step process
- Create base OS image
- Create OS image for CAPI using the base OS image

## Create base OS image
----
Creating base OS image form DVD ISO files is a one time process. A base vSphere template will be created as a result of building a base image.
This base image template can be used to create more customized images.

<!-- [TODO: Add steps and screenshots] -->

## Create vSphere template image for CAPI
----
The process to create a base OS image will create a vSphere template in vSphere. This template is a linked clone of the VM so any VM created from template will be provisioned quickly.

**Environment Variables for vSphere:**
*NOTE:* use vSphere server URL without `http://` or `https://`

```bash
VSPHERE_SERVER=example.vsphere.url
VSPHERE_USERNAME=user@example.vsphere.url
VSPHERE_PASSWORD=example_password
```

**Environment Variables for RedHat subscription:**

For username and password:

```bash
RHSM_USER=example_user
RHSM_PASS=example_password
```

Or activation key and org ID:

```bash
RHSM_ACTIVATION_KEY=example_key
RHSM_ORG_ID=example_org_id
```

Setting pool_ids, environment, or consumer name is optional, but recommended for better subscription management:

```bash
RHSM_ENVIRONMENT=example_environment
RHSM_POOL_ID=example_pool_id
```

Or set a custom consumer name for easier identification:

```bash
RHSM_CONSUMER_NAME=example_consumer_name
RHSM_POOL_ID=example_pool_id
```

**Packer variables for vSphere:**
Following variables are needed for the vSphere. Add this configuration in the `image.yaml`

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

## Create template image on vSphere

```bash
konvoy-image build path/to/image.yaml
```

Checkout example image configuration at: [`<project_root>/images/ova/`](../../images/ova) directory.
