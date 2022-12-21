#!/bin/bash
set -o errexit -o pipefail -x

function usage() {
  echo "usage: $0 \\"
  echo '--push - pushes artiacts to github and dockerhub'
}

function main () {
  if  [ "${push}" = false ]; then
	  DOCKER_BUILDKIT=1 goreleaser --parallelism=1 --rm-dist --snapshot --timeout=2h --debug
    exit 0
  fi
  make docker-build-amd64
  make docker-build-arm64
  make push-manifest
  DOCKER_BUILDKIT=1 goreleaser release --snapshot --skip-publish --rm-dist --parallelism=1 --timeout=2h
  exit 0
}

push="false"

while [ "$1" != "" ]; do
  case $1 in
  -push | --push)
    shift
    push="true"
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
