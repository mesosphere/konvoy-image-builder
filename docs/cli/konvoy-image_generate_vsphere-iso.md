## konvoy-image generate vsphere-iso

generate files relating to building vsphere images from ISO

```
konvoy-image generate vsphere-iso <image.yaml> [flags]
```

### Examples

```
vsphere-iso --vsphere-datacenter=dc1 --vsphere-cluster=zone1 --vsphere-network=Public --vsphere-datastore=datastore1 images/vsphere-iso/ubuntu-2004.yaml
```

### Options

```
      --containerd-version string      the version of containerd to install
      --extra-vars strings             flag passed Ansible's extra-vars
  -h, --help                           help for vsphere-iso
      --iso-checksum string            replace the templates iso checksum
      --iso-url string                 replace the templates iso url
      --kubernetes-version string      The version of kubernetes to install. Example: 1.21.6
      --overrides strings              a comma separated list of override YAML files
      --ssh-username string            specify the initial user which gets SUDO permissions
      --vcenter-server string          vCenter server address (or environment variable VSPHERE_SERVER)
      --vsphere-cluster string         vSphere cluster name
      --vsphere-datacenter string      vSphere data center name
      --vsphere-datastore string       vSphere datastore to be used
      --vsphere-folder string          place VM in vSphere folder
      --vsphere-insecure-connection    ignore SSL certificate errors
      --vsphere-network string         specify vSphere resource pool for the build VM
      --vsphere-password string        vSphere password (or environment variable VSPHERE_PASSWORD)
      --vsphere-resource-pool string   specify vSphere resource pool for the build VM
      --vsphere-user string            vSphere user (or environment variable VSPHERE_USER)
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image generate](konvoy-image_generate.md)	 - generate files relating to building images

