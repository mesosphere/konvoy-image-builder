# Konvoy Image Builder

Molecule Test: ![Teamcity Status](https://teamcity.mesosphere.io/app/rest/builds/buildType:%28id:ClosedSource_KonvoyImageBuilder_MoleculeTest%29/statusIcon.svg)

The goal of Konvoy Image Builder (KIB) is to produce a common operating surface to run konvoy across heterogeneous infrastructure. KIB relies on ansible to install software, configure, and sanitize systems for running konvoy. Packer is used to build images for cloud environments. Goss is used to validate system’s are capable of running Konvoy.

## Supported OS Families

Presently, KIB supports four OS families:

- Debian
- Red Hat
- Flatcar
- and SUSE

## KIB Repository Layout

- `ansible`: contains the ansible playbooks, roles, and default variables
- `images`: contains image definitions for supported platforms. Presently, we provide AMI image definitions and generic image definitions. Generic image definitions are useful for preprovisioned infrastructure
- `overrides`: contains variable overrides for nvidia and fips. Unless adding an overlay feature, these files can safely be ignored.

## Quickstart

## `konvoy-image` CLI

### Example Usage

```sh
konvoy-image build images/ami/centos-7.yaml
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

By default the `devkit.run` target will run a shell, to specify another
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

To run an specific end-to-end test, use a subset of the commands used to run the complete set of end-to-end tests in CI.

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

These are temporary instructions for building the wrapper for testing

```sh
make build.snapshot
```

Replace image tag with the version created by go releaser

```sh
docker save mesosphere/konvoy-image-builder:v1.0.0-alpha1-SNAPSHOT-e590962 \
 | gzip -c - > cmd/konvoy-image-wrapper/image/konvoy-image-builder.tar.gz
```

Build the wrapper

```sh
go build -tags EMBED_DOCKER_IMAGE \
 -ldflags="-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=v1.0.0-alpha1-SNAPSHOT-e590962" \
 -o ./bin/konvoy-image-wrapper ./cmd/konvoy-image-wrapper/main.go
```

For further development, see the [Dev Docs](docs/dev).