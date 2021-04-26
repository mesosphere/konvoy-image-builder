## konvoy-image generate

generate files relating to building images

```
konvoy-image generate <image.yaml> [flags]
```

### Examples

```
generate --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-7.yaml
```

### Options

```
      --ami-regions stringArray          a list of regions to publish amis
      --containerd-version string        the version of containerd to install
  -h, --help                             help for generate
      --kubernetes-version string        the version of kubernetes to install
      --overrides stringArray            a list of override YAML files
      --region string                    the aws region to run the builder
      --source-ami string                a specific ami available in the builder region to source from
      --source-ami-filter-name string    a ami name filter on for selecting the source image
      --source-ami-filter-owner string   only search AMIs belonging to this owner id
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image](konvoy-image.md)	 - Create, provision, and customize images for running Konvoy

