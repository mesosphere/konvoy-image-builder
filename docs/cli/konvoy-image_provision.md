## konvoy-image provision

provision to an inventory.yaml or hostname, note the comma at the end of the hostname

```
konvoy-image provision <inventory.yaml|hostname,> [flags]
```

### Examples

```
build ./inventory.yaml
```

### Options

```
      --extra-vars stringArray   flag passed Ansible's extra-vars
  -h, --help                     help for provision
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image](konvoy-image.md)	 - Create, provision, and customize images for running Konvoy

