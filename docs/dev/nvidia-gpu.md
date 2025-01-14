# Nvidia GPU

`konvoy-image-builder` supports building images with installed Nvidia host components
that can be used to run GPU workloads on machines that come with a dedicated GPU
hardware.

To enable GPU build, add the following override:

```yaml
---
gpu:
  types:
    - nvidia
build_name_extra: "-nvidia"
```

There is an existing prepared override file in `overrides/nvidia.yaml` that will
enable installing Nvidia host components to the built image.

Example:

```sh
build --region us-west-2 --source-ami=ami-12345abcdef images/ami/centos-79.yaml \
    --overrides overrides/nvidia.yaml \
    --instance-type=g4dn.2xlarge
```

Example of building on AWS with Ubuntu 20.04 as the operating system and `p2.xlarge` as the AWS machines you are going to use:

```sh
konvoy-image build aws --region us-west-2 images/ami/ubuntu-2004.yaml \
    --overrides overrides/nvidia.yaml \
    --instance-type=p2.xlarge
```

**NOTE:** Providing the instance type is required for building GPU images.

## Supported images

| base os   | nvidia             |
|-----------|--------------------|
| rhel-79   | :white_check_mark: |
| rhel-84   | :white_check_mark: |
| rhel-86   | :white_check_mark: |
| ubuntu-20 | :white_check_mark: |
| flatcar   | :white_check_mark: |
