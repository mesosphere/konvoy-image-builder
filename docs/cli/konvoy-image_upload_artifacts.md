## konvoy-image upload artifacts

upload artifacts to hosts defined in inventory-file

```
konvoy-image upload artifacts [flags]
```

### Options

```
      --container-images-dir string   path to container images for install on remote hosts.
  -h, --help                          help for artifacts
      --inventory-file string         an ansible inventory defining your infrastructure (default "inventory.yaml")
      --os-packages-bundle string     path to os-packages tar file for install on remote hosts.
      --pip-packages-bundle string    path to pip-packages tar filefor install on remote hosts.
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image upload](konvoy-image_upload.md)	 - Upload one of [artifacts]

