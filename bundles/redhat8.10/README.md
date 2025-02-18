# Updating packages.txt.gotmpl

The package list includes packages required by the NVIDIA GPU drivers, kernel- headers, and kernel-devel. These headers need to match the kernel version in the base image.

Periodically, the base image is updated, and we need to update the package versions. To find the expected kernel version, look for the error message, e.g.

```log
amazon-ebs.kib_image: fatal: [default]: FAILED! => {"changed": false, "failures": ["No package kernel-headers-4.18.0-553.34.1.el8_10.x86_64 available.", "No package kernel-devel-4.18.0-553.34.1.el8_10.x86_64 available."], "msg": "Failed to install some of the specified packages", "rc": 1, "results": []}
```

In this example, the expected kernel version is 4.18.0-553.34.1.
