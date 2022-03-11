# Nvidia GPU

`konvoy-image-builder` supports building images with installed Nvidia host components
that can be used to run GPU workloads on machines that come with a dedicated GPU
hardware.

To enable GPU build, add the following override:

```yaml
gpu:
  type: nvidia
```

There is an existing prepared override file in `overrides/nvidia.yaml` that will
enable installing Nvidia host components to the built image.

Example:

```sh
build --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-7.yaml \
    --overrides overrides/nvidia.yaml
```

*NOTE* That for creating flatcar GPU images, the builder node has to have GPU support. `konvoy-image-builder` allows for providing the instance type to use:

```sh
--aws-instance-type <INSTANCE_TYPE_NAME>
```

Example:
```sh
build --region us-west-2 \
    --aws-instance-type p2.xlarge \
    images/ami/flatcar.yaml \
    --overrides overrides/nvidia.yaml
```

## Supported images

| base os   | nvidia             |
|-----------|--------------------|
| centos-7  | :white_check_mark: |
| centos-8  | :white_check_mark: |
| rhel-79   | :white_check_mark: |
| rhel-82   | :white_check_mark: |
| rhel-84   | :white_check_mark: |
| ubuntu-18 | :white_check_mark: |
| ubuntu-20 | :white_check_mark: |
| flatcar   | :white_check_mark: |
| sles-15   | :white_check_mark: |
