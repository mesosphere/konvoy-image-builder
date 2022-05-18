## konvoy-image generate azure

generate files relating to building azure images

```
konvoy-image generate azure <image.yaml> [flags]
```

### Examples

```
azure --location westus2 --subscription-id <sub_id> images/azure/centos-79.yaml
```

### Options

```
      --client-id string                   the client id to use for the build
      --containerd-version string          the version of containerd to install
      --extra-vars strings                 flag passed Ansible's extra-vars
      --gallery-image-locations location   a list of locatins to publish the image (default same as location)
      --gallery-image-name string          the gallery image name to publish the image to
      --gallery-image-offer string         the gallery image offer to set (default "dkp")
      --gallery-image-publisher string     the gallery image publisher to set (default "dkp")
      --gallery-image-sku string           the gallery image sku to set
      --gallery-name string                the gallery name to publish the image in (default "dkp")
  -h, --help                               help for azure
      --instance-type string               the Instance Type to use for the build (default "Standard_D2ds_v5")
      --kubernetes-version string          The version of kubernetes to install. Example: 1.21.6
      --location string                    the location in which to build the image (default "westus2")
      --overrides strings                  a comma separated list of override YAML files
      --resource-group string              the resource group to create the image in (default "dkp")
      --subscription-id string             the subscription id to use for the build
      --tenant-id string                   the tenant id to use for the build
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image generate](konvoy-image_generate.md)	 - generate files relating to building images

