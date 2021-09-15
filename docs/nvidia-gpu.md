# Nvidia GPU

`konvoy-image-builder` supports building images with installed Nvidia host components
that can be used to run GPU workloads on machines that comes with a dedicated GPU
hardware.

To enable GPU build add the following override:

```yaml
gpu:
  type: nvidia
```

There is existing prepared override file in `overrides/nvidia.yaml` that will
enable installing Nvidia host components to the built image.

Example:

```sh
build --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-7.yaml \
    --overrides overrides/nvidia.yaml
```

Note, that for creating flatcar GPU images the builder node has to have GPU support. `konvoy-image-builder` allows for providing the instance type to use:

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

| base os | nvidia             |
|---------|--------------------|
| centos7 | :white_check_mark: |
| centos8 | :white_check_mark: |
| flatcar | :white_check_mark: |
| sles15  | :white_check_mark: |
