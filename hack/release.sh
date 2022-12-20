#!/bin/bash
set -o errexit -o pipefail -x

function usage() {
  echo "usage: $0 \\"
  echo '--push - pushes artiacts to github and dockerhub'
}

function main () {
  make docker-build-amd64 BUILDARCH=amd64 GOARCH=amd64
  make docker-build-arm64 BUILDARCH=arm64 GOARCH=arm64
  if ${push}; then
    make docker-push
    make push-manifest
	  DOCKER_BUILDKIT=1 goreleaser --parallelism=1 --rm-dist --debug --snapshot
    exit 0
  fi
	DOCKER_BUILDKIT=1 goreleaser release --snapshot --skip-publish --rm-dist --parallelism=1
}

push=false

while [ "$1" != "" ]; do
  case $1 in
  -push | --push)
    shift
    push=true
    ;;
  -h | --help)
    usage
    exit
    ;;
  *)
    echo "unknown parameter: $1"
    usage
    exit 1
    ;;
  esac
  shift
done

main
