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
| us-west-2  | centos 7 | ami-0bc38a003a647b084 | [`images/ami/centos-79.yaml`](images/ami/centos-79.yaml) |

## Development

## Recommended Tools

* [Go](https://golang.org/doc/install)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Buildx](https://docs.docker.com/build/install-buildx/)
* [goreleaser](https://goreleaser.com/install/)
* [magefile](https://magefile.org/)

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

[Magefile](https://magefile.org/) tool is used to run konvoy image builder e2e tests.

In this example, we run the end-to-end test against the Centos 7.9 with air-gapped and fips configuration in AWS
```sh
runE2e "centos 7.9" "offline-fips" aws false
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
