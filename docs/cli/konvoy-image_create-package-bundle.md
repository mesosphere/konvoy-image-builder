## konvoy-image create-package-bundle

build os package bundles for airgapped installs

```
konvoy-image create-package-bundle [flags]
```

### Examples

```
create-package-bundle --os redhat-8.4 --kubernetes-version=1.28.6 --output-directory=artifacts
```

### Options

```
      --container-image string      A container image to use for building the package bundles
      --fips                        If the package bundle should include fips packages.
  -h, --help                        help for create-package-bundle
      --kubernetes-version string   The version of kubernetes to download packages for. Example: 1.21.6
      --os string                   The target OS you wish to create a package bundle for. Must be one of [centos-7.9 redhat-7.9 redhat-8.4 redhat-8.6 redhat-8.8 oracle-7.9 rocky-9.1 ubuntu-18.04 ubuntu-20.04]
      --output-directory string     The directory to place the bundle in.
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image](konvoy-image.md)	 - Create, provision, and customize images for running Konvoy

