# Image Builder Ansible
The included Ansible roles are being used to provision the image with needed software and drivers.


## Roles
Following roles are included

### `packages` - install main konvoy packages
This role installs packages like kubelet or kubeadm and sets the right
repositories to be used

### `gpu` - maintains driver packages for GPU support
Currently `gpu.types` only support nvidia.
It can be activated when `gpu.types = ['nvidia',...]`

### `config` - sets special config options
Being used to set expected config values

### `networking` - setup konvoy networking
set bridge utils and iptables

### `providers` - prepare cloud provider tooling
Installs and configures tooling for the cloud provider being used to generate
this image

### `images` - prepare image cache
This roll pulls a list of images when `download_images` is set

### `sysprep` - prepare image for distribution
truncate logs, system id, host key etc.


## Development
For development and testing molecule was prepared with scenarios and tests.

### Setup
Use `./ansible` as workdir and if not already done install python `3.9.8`

For better maintanability we create a virtual env.

```bash
cd ./ansible
asdf install python
python -m venv ./venv
source ./venv/bin/activate
```

The molecule dependencies are stored in `./requirements.txt`

```bash
pip install -r requirements.txt
source ./venv/bin/activate

```

### Working with molecule
molecule knows different deployment states which can be used at once
(`molecule test`) or as single steps.

For development its useful to use them one by one:

#### create
`molecule create` prepares the platforms ( or base images ) and creates the
instances.

#### prepare
`molecule prepare` is not used at the moment.

#### converge
`molecule converge` executes the actual roles or plays. The default scenario
uses the `provision.yaml` play to executes every role expcept `sysprep` as this
would intercept communication to the image

#### verify
`molecule verify` executes python testinfra to run tests against the systems.

#### destroy
`molecule destroy` destroys all platforms and the security groups being used.

#### test - all in one
`molecule test` without any arguments or options is more meant for CI runs.
It will destroy the instances after a run no matter if something failed or not.

If you do not want to execute all steps by hand you can run `molecule test --destroy never`.
Now if something fails and you want to try your fix you can rerun `molecule converge`

### Scenarios
Currently we have two main scenarios for molecule all based on ec2

```bash
molecule test
# means actually
molecule test -s default

molecule test -s ec2_gpu

# execute all scenarios
molecule test --all --parallel

# list all scenarios
molecule list
```

#### default
the default scenario just installs konvoy main packages and checks if they are
available.

#### ec2_gpu
the `ec2_gpu` scenario leaves out all other rules but the gpu specific one. It
uses gpu based instances so its test can check if they are available and being
ready for usage.

### tests
Tests are using python testinfra.

Check `molecule/default/tests/test_default.py` for example.
