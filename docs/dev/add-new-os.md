# Konvoy Image Builder - Adding Operating Systems

An OS family represents major branches in how machines are provisioned. Adding an OS family is an advanced task and will require detailed knowledge of the OS family being added and familiarity with ansible.

Differences within OS family distributions (Ubuntu vs Debian) and versions (Centos 7 vs Centos 8) are handled only as needed. In most cases, version and distribution differences are handled by ansible, and likely will not require any changes to the existing playbooks.

## Adding a New OS Family

Adding support for a new OS family will be implemented primarily in ansible. As ansible abstracts many common OS features, the process for adding a new OS is fairly painless. For example, this PR represents the addition of the Debian OS family and the addition of Ubuntu 20.04:

- [PR #100](https://github.com/mesosphere/konvoy-image-builder/pull/100)

- Repositories for Debian and Red Hat families are defined [here](https://github.com/mesosphere/konvoy-image-builder/blob/main/ansible/group_vars/all/system.yaml#L1-L33).
- Release links to the actual binaries are defined [here](https://github.com/mesosphere/konvoy-image-builder/blob/main/ansible/group_vars/all/defaults.yaml#L55-L62).

### Images

Image definitions for AMIs should contain variables which allow packer to properly discover the base image for the OS family and version of the target.

An example of an image definition for producing a centos-8 image is defined [here](https://github.com/mesosphere/konvoy-image-builder/blob/main/images/ami/centos-8.yaml).

#### Goss

When adding a new OS family, a GOSS spec is required.

For [PR #100](https://github.com/mesosphere/konvoy-image-builder/pull/100/files), all that was done to achieve this was to simply copy `goss/centos` to `goss/ubuntu`.

## Adding a New Supported Version

If the OS family is supported, adding a new image definition is straightforward. Here is an example of adding Ubuntu 18.04 support, which is part of the Debian OS family:

- [PR #104](https://github.com/mesosphere/konvoy-image-builder/pull/104)

All that's needed is to provide an image definition which instructs packer where to find the base image. You can make changes in the ansible roles to account for any differences.

### Image Builder Commands and Testing

#### Build

An AWS machine image can be built using the following following command:

```sh
./konvoy-image build images/ami/<image>.yaml
```

#### Generate

When developing ansible roles and tasks it’s much easier to test against a long running instance. KIB can be used to generate extravars for a given target.

```sh
./konvoy-image generate images/ami/<image>.yaml
```

This will create a directory in `./work` which will contain files needed to use ansible and packer directly.

To review the rest of the CLI commands, go to the [CLI docs](../cli).