# Konvoy Image Builder

The goal of Konvoy Image Builder (KIB) is to produce a common operating surface to run Konvoy across heterogeneous infrastructure. KIB relies on ansible to install software, configure, and sanitize systems for running Konvoy. Packer is used to build images for cloud environments. Goss is used to validate system’s are capable of running Konvoy.

## Supported OS Families

Presently, KIB supports four OS families:

- Debian
- Red Hat
- Flatcar
- SUSE

## KIB Repository Layout

- `ansible`: contains the ansible playbooks, roles, and default variables
- `images`: contains image definitions for supported platforms. Presently, we provide AMI image definitions and generic image definitions. Generic image definitions are useful for preprovisioned infrastructure
- `overrides`: contains variable overrides for Nvidia and FIPS. Unless adding an overlay feature, these files can safely be ignored.

## Quickstart

## `konvoy-image` CLI

### Example Usage

```sh
konvoy-image build images/ami/centos-79.yaml
```

### CLI documentation

See [`konvoy-image`](docs/cli/konvoy-image.md)

## Example Images

| aws-region | base-os  | ami-id                | image params                                           |
|------------|----------|-----------------------|--------------------------------------------------------|
| us-west-2  | centos 7 | ami-0bc38a003a647b084 | [`images/ami/centos-7.yaml`](images/ami/centos-7.yaml) |

## Development

### Devkit Container

A devkit is provided to quickly allow usage and development. To build and
launch the devkit run:

```sh
make devkit.run
```

By default, the `devkit.run` target will run a shell, to specify another
command, set the `WHAT` variable. For example to run `make build` in the
devkit run:

```sh
make devkit.run WHAT='make build'
```

### Linting

The tooling consists of several languages, the main wrapper code is written in
`go` which is linted with `golangci-lint`. To lint the `go` files run:

```sh
make lint
```

Other languages are linted with the help of
[`super-linter`](https://github.com/github/super-linter). To lint everything
else run:

```sh
make super-lint
```

*NOTE* Konvoy Image Builder makes use of the `embed` feature of `go` 1.16.
`super-linter` currently uses `go` 1.15. It is expected that the `go` linter
will fail under `super-linter`, and is skipped for
[CI](.github/workflows/lint.yml).

### Testing

To run a specific end-to-end test, use a subset of the commands used to run the complete set of end-to-end tests in CI.

In this example, we run the end-to-end test against the latest version of Flatcar Linux:

```sh
WHAT="make flatcar-version.yaml" make devkit.run
make devkit.run \
    WHAT="./bin/konvoy-image build images/ami/flatcar.yaml --overrides flatcar-version.yaml -v 5" \
    INTERACTIVE=""
make docker.clean-latest-ami
```

### Building

To build the CLI command run:

```sh
make build
```

### Building the Wrapper

To build the wrapper for testing.

```sh
make build-wrapper
```
creates `./bin/konvoy-image-wrapper` binary for testing using konvoy image wrapper.

For further development, see the [Dev Docs](docs/dev).
