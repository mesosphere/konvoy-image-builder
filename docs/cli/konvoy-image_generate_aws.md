## konvoy-image generate aws

generate files relating to building aws images

```
konvoy-image generate aws <image.yaml> [flags]
```

### Examples

```
aws --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-7.yaml
```

### Options

```
      --ami-groups stringArray           a list of AWS groups which are allowed use the image, using 'all' result in a public image
      --ami-regions stringArray          a list of regions to publish amis
      --ami-users stringArray            a list AWS user accounts which are allowed use the image
      --containerd-version string        the version of containerd to install
      --extra-vars strings               flag passed Ansible's extra-vars
  -h, --help                             help for aws
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

* [konvoy-image generate](konvoy-image_generate.md)	 - generate files relating to building images

