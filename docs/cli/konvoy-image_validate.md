## konvoy-image validate

validate existing infrastructure

```
konvoy-image validate [flags]
```

### Options

```
      --apiserver-port int      apiserver port (default 6443)
  -h, --help                    help for validate
      --pod-subnet string       ip addresses used for the pod subnet (default "192.168.0.0/16")
      --service-subnet string   ip addresses used for the service subnet (default "10.96.0.0/12")
```

### Options inherited from parent commands

```
      --color     enable color output (default true)
  -v, --v int     select verbosity level, should be between 0 and 6
      --verbose   enable debug level logging (same as --v 5)
```

### SEE ALSO

* [konvoy-image](konvoy-image.md)	 - Create, provision, and customize images for running Konvoy

