# Konvoy Image Builder

Molecule Test: ![Teamcity Status](https://teamcity.mesosphere.io/app/rest/builds/buildType:%28id:ClosedSource_KonvoyImageBuilder_MoleculeTest%29/statusIcon.svg)

This repository contains tools for producing system images for the purpose of
running Konvoy.

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

*NOTE*: Konvoy Image Builder makes use of the `embed` feature of `go` 1.16.
`super-linter` currently uses `go` 1.15. It is expected that the `go` linter
will fail under `super-linter`, and is skipped for
[CI](.github/workflows/lint.yml).

### Building

To build the CLI command run:

```sh
make build
```
