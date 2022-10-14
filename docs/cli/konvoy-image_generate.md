## konvoy-image generate

generate files relating to building images

### Synopsis

Generate files relating to building images. Specifying AWS arguments is deprecated and will be removed in a future version. Use the `aws` subcommand instead.

```
konvoy-image generate <image.yaml> [flags]
```

### Options

```
      --ami-groups stringArray           a list of AWS groups which are allowed use the image, using 'all' result in a public image
      --ami-regions stringArray          a list of regions to publish amis
      --ami-users stringArray            a list AWS user accounts which are allowed use the image
      --containerd-version string        the version of containerd to install
      --extra-vars strings               flag passed Ansible's extra-vars
  -h, --help                             help for generate
      --instance-type string             instance type used to build the AMI; the type must be present in the region in which the AMI is built
      --kubernetes-version string        The version of kubernetes to install. Example: 1.21.6
      --overrides strings                a comma separated list of override YAML files
      --region string                    the region in which to build the AMI
      --source-ami string                the ID of the AMI to use as the source; must be present in the region in which the AMI is built
      --source-ami-filter-name string    restricts the set of source AMIs to ones whose Name matches filter
      --source-ami-filter-owner string   restricts the source AMI to ones with this owner ID
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image](konvoy-image.md)	 - Create, provision, and customize images for running Konvoy
* [konvoy-image generate aws](konvoy-image_generate_aws.md)	 - generate files relating to building aws images
* [konvoy-image generate azure](konvoy-image_generate_azure.md)	 - generate files relating to building azure images
* [konvoy-image generate vsphere-iso](konvoy-image_generate_vsphere-iso.md)	 - generate files relating to building vsphere images from ISO

