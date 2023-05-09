SHELL:=/bin/bash
.DEFAULT_GOAL := help

OS := $(shell uname -s)

INTERACTIVE := $(shell [ -t 0 ] && echo 1)

root_mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
export REPO_ROOT_DIR := $(patsubst %/,%,$(dir $(root_mkfile_path)))
export REPO_REV ?= $(shell cd $(REPO_ROOT_DIR) && git describe --abbrev=12 --tags --match='v*' HEAD)

UID ?= $(shell id -u)
GID ?= $(shell id -g)
USER_NAME ?= $(shell id -u -n)
GROUP_NAME ?= $(shell id -g -n)

COVERAGE ?= $(REPO_ROOT_DIR)/coverage

VERBOSITY ?= 0

INVENTORY_FILE ?= $(REPO_ROOT_DIR)/inventory.yaml
COMMA:=,

export CGO_ENABLED=0
export GO_VERSION := $(shell cat go.mod | grep "go " -m 1 | cut -d " " -f 2)
GOLANG_IMAGE := golang:$(GO_VERSION)
ARCH := $(shell uname -m)

BUILDARCH ?= $(shell echo $(ARCH) | sed 's/x86_64/amd64/g')

export CI ?= no

export DOCKER_SOCKET ?= /var/run/docker.sock
ifeq ($(OS),Darwin)
export DOCKER_SOCKET_GID ?= $(shell /usr/bin/stat -f "%g" $(DOCKER_SOCKET))
else
export DOCKER_SOCKET_GID ?= $(shell stat -c %g $(DOCKER_SOCKET))
endif

export DOCKER_DEVKIT_GO_ENV_ARGS ?= \
	--env GOCACHE=/kib/.cache/go-build \
	--env GOMODCACHE=/kib/.cache/go-mod \
	--env GOLANGCI_LINT_CACHE=/kib/.cache/golangci-lint \

export DOCKER_DEVKIT_ENV_ARGS ?= \
	--env CI \
	--env GITHUB_TOKEN \
	--env BUILD_DRY_RUN \
	$(DOCKER_DEVKIT_GO_ENV_ARGS)

export DOCKER_DEVKIT_AWS_ARGS ?= \
	--env AWS_PROFILE \
	--env AWS_SECRET_ACCESS_KEY \
	--env AWS_SESSION_TOKEN \
	--env AWS_DEFAULT_REGION \
	--volume "$(HOME)/.aws":"/home/$(USER_NAME)/.aws"

DOCKER_GCP_CREDENTIALS_ARGS=--volume "$(HOME)/.gcloud":"/home/$(USER_NAME)/.gcloud" \
	                             --env GOOGLE_APPLICATION_CREDENTIALS=/home/$(USER_NAME)/.gcloud/credentials.json

export DOCKER_DEVKIT_GCP_ARGS ?= \
	$(DOCKER_GCP_CREDENTIALS_ARGS)

export DOCKER_DEVKIT_AZURE_ARGS ?= \
	--env AZURE_LOCATION \
	--env AZURE_CLIENT_ID \
	--env AZURE_CLIENT_SECRET \
	--env AZURE_SUBSCRIPTION_ID \
	--env AZURE_TENANT_ID \
	--volume "$(HOME)/.azure":"/home/$(USER_NAME)/.azure"

export DOCKER_DEVKIT_VSPHERE_ARGS ?= \
	--env VSPHERE_SERVER \
	--env VSPHERE_USERNAME \
	--env VSPHERE_PASSWORD \
	--env RHSM_USER \
	--env RHSM_PASS

export DOCKER_DEVKIT_BASTION_ARGS ?= \
	--env SSH_BASTION_USERNAME \
	--env SSH_BASTION_HOST \
	--env SSH_BASTION_KEY_CONTENTS

ifneq ($(wildcard $(DOCKER_SOCKET)),)
	export DOCKER_SOCKET_ARGS ?= \
		--volume "$(DOCKER_SOCKET)":/var/run/docker.sock
endif

export DOCKER_DEVKIT_PUSH_ARGS ?= \
	--volume "$(HOME)/.docker":"/home/$(USER_NAME)/.docker" \
	--env DOCKER_PASS \
	--env DOCKER_CLI_EXPERIMENTAL

# ulimit arg is a workaround for golang's "suboptimal" bug workaround that
# manifests itself in alpine images, resulting in packer plugins simply dying.
#
# On LTS distros like Ubuntu, kernel bugs are backported, so the kernel version
# may seem old even though it is not vulnerable. Golang ignores it and just
# looks at the distro+kernel combination to determine if it should panic or
# not. This results in packer silently failing when running in devkit
# container, as it is using Alpine linux. See the issue below for more details:
# https://github.com/docker-library/golang/issues/320
export DOCKER_ULIMIT_ARGS ?= \
	--ulimit memlock=67108864:67108864

export DOCKER_DEVKIT_USER_ARGS ?= \
	--user $(UID):$(GID) \
	--group-add $(DOCKER_SOCKET_GID)

export DOCKER_DEVKIT_SSH_ARGS ?= \
	--env SSH_AUTH_SOCK=/run/ssh-agent.sock \
	--volume $(SSH_AUTH_SOCK):/run/ssh-agent.sock

export DOCKER_DEVKIT_ARGS ?= \
	$(DOCKER_ULIMIT_ARGS) \
	$(DOCKER_DEVKIT_USER_ARGS) \
	--volume $(REPO_ROOT_DIR):/kib \
	--workdir /kib \
	$(DOCKER_SOCKET_ARGS) \
	$(DOCKER_DEVKIT_AWS_ARGS) \
	$(DOCKER_DEVKIT_GCP_ARGS) \
	$(DOCKER_DEVKIT_AZURE_ARGS) \
	$(DOCKER_DEVKIT_BASTION_ARGS) \
	$(DOCKER_DEVKIT_VSPHERE_ARGS) \
	$(DOCKER_DEVKIT_PUSH_ARGS) \
	$(DOCKER_DEVKIT_ENV_ARGS) \
	$(DOCKER_DEVKIT_SSH_ARGS)

export DOCKER_DEVKIT_DEFAULT_ARGS ?= \
	--rm \
	$(if $(INTERACTIVE),--tty) \
	--interactive

ifneq ($(shell git status --porcelain 2>/dev/null; echo $$?), 0)
	export GIT_TREE_STATE := dirty
else
	export GIT_TREE_STATE :=
endif


include make/ci.mk
include make/images.mk

# envsubst
# ---------------------------------------------------------------------
export ENVSUBST_VERSION ?= v1.2.0
export ENVSUBST_URL = https://github.com/a8m/envsubst/releases/download/$(ENVSUBST_VERSION)/envsubst-$(shell uname -s)-$(shell uname -m)
export ENVSUBST_ASSETS ?= $(CURDIR)/.local/envsubst/${ENVSUBST_VERSION}

.PHONY: install-envsubst
install-envsubst: ## install envsubst binary
install-envsubst: $(ENVSUBST_ASSETS)/envsubst

$(ENVSUBST_ASSETS)/envsubst:
	$(call print-target,install-envsubst)
	mkdir -p $(ENVSUBST_ASSETS)
	curl -Lf $(ENVSUBST_URL) -o $(ENVSUBST_ASSETS)/envsubst
	chmod +x $(ENVSUBST_ASSETS)/envsubst

include hack/pip-packages/Makefile
include test/infra/aws/Makefile
include test/infra/vsphere/Makefile

github-token.txt:
	echo $(GITHUB_TOKEN) >> github-token.txt


.PHONY: buildx
buildx:
buildx:
	 docker buildx create --use --name=konvoy-image-builder || true
	 docker run --privileged --rm tonistiigi/binfmt --install all &>/dev/null || true

###### Devkit container image
DEVKIT_IMAGE_DOCKERFILE ?= Dockerfile.devkit
DEVKIT_IMAGE_NAME ?= mesosphere/konvoy-image-builder-devkit
DEVKIT_IMAGE_TAG ?= $(shell cat ${DEVKIT_IMAGE_DOCKERFILE} requirements.txt requirements-devkit.txt  | sha256sum | cut -d" " -f 1)

.PHONY: devkit-image
## first tries to pull an image, if doesn't exist build and push the image
devkit-image:
	$(call print-target)
	docker image inspect $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG) &>/dev/null || \
	docker pull $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG) || \
	$(MAKE) devkit-image-amd64 devkit-image-arm64

.PHONY: devkit-image-amd64
devkit-image-amd64:
	$(MAKE) devkit-image-build-push BUILDARCH=amd64

.PHONY: devkit-image-arm64
devkit-image-arm64:
	$(MAKE) devkit-image-build-push BUILDARCH=arm64

.PHONY: devkit-image-build-push
devkit-image-build-push: github-token.txt buildx
	$(call print-target)
	docker image inspect $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-$(BUILDARCH) &>/dev/null || \
	docker pull $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-$(BUILDARCH) || \
	docker buildx build \
	--pull \
	--push \
	--build-arg BUILDARCH=$(BUILDARCH) \
	--secret id=githubtoken,src=github-token.txt \
	--provenance=false \
	--platform linux/$(BUILDARCH) \
	--file $(DEVKIT_IMAGE_DOCKERFILE) \
	--tag=$(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-$(BUILDARCH) \
	$(dir $(DEVKIT_IMAGE_DOCKERFILE))

.PHONY: devkit-image-push-manifest
devkit-image-push-manifest: ## pushes the devkit-image
devkit-image-push-manifest: devkit-image-amd64 devkit-image-arm64
	$(call print-target)
	docker manifest create \
	$(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG) \
	--amend $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-amd64 \
	--amend $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-arm64
	docker manifest push $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)

##### Build KIB container image
export DOCKER_REPOSITORY ?= mesosphere/konvoy-image-builder
export DOCKER_IMG ?= $(DOCKER_REPOSITORY):$(REPO_REV)-$(BUILDARCH)

.PHONY: kib-image-build-amd64
kib-image-build-amd64:
	$(MAKE) kib-image-build BUILDARCH=amd64

.PHONY: kib-image-build-arm64
kib-image-build-arm64:
	$(MAKE) kib-image-build BUILDARCH=arm64

# we need to push these devkit images up to dockerhub.
# buildx does not pull the base images specified in Dockerfile from local docker cache.
# buildx always attempts to pull it from dockerhub.
# This behavior is not documented. see https://github.com/moby/moby/issues/42893#issuecomment-1241274246

# The latest Docker Engine version 23.0 defaults to buildx for building images.
# The caching improvements in this release might fix this problem.
# see release notes: https://docs.docker.com/engine/release-notes/23.0/#2301

# Always build the kib image to embed `konvoy-imag` binary from current revision
# devkit-image-$(BUILDARCH): will build and push arch specific devkit image to docker hub.
# konvoy-image-$(BUILDARCH): will create arch specific binary in ./bin directory.
# kib-image-build: will create docker image used by konvoy-image-wrapper.
#                  kib container image is built using './Dockerfile' that has two main dependencies
#				   (1) The devkit image reference from docker hub.
#					   The buildx always pulls base images from dockerhub. so this image has to be pushed to registry.
#				   (2) Arch specific konvoy-image binary file thats get copied in the image

# TODO: revisit this target when moving to docker 23.0.x.
# use 'docker build --platform' instead of 'docker buildx build' to avoid exporing images from buildx cache to docker cache.
.PHONY: kib-image-build
kib-image-build: devkit-image-$(BUILDARCH) konvoy-image-$(BUILDARCH)
	docker buildx build \
		--file $(REPO_ROOT_DIR)/Dockerfile \
		--build-arg BUILDARCH=$(BUILDARCH) \
		--build-arg BASE=$(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-$(BUILDARCH) \
		--platform linux/$(BUILDARCH) \
		--pull \
		--load \
		--tag=$(DOCKER_REPOSITORY):$(REPO_REV)-$(BUILDARCH) \
		$(REPO_ROOT_DIR)

.PHONY: kib-image-push-amd64
kib-image-push-amd64: kib-image-build-amd64
	docker push $(DOCKER_REPOSITORY):$(REPO_REV)-amd64

.PHONY: kib-image-push-arm64
kib-image-push-arm64: kib-image-build-arm64
	docker push $(DOCKER_REPOSITORY):$(REPO_REV)-arm64

# The arch specific images must be pushed to docker registry in order to create manifest file.
# A manifest file can not be created using locally cached images. see: https://github.com/docker/cli/issues/3350
# TODO: Build and push multi arch image using single command: docker buildx build --platform linux/amd64,linux/arm64 --output=type=registry
.PHONY: kib-image-push-manifest
kib-image-push-manifest: kib-image-push-amd64 kib-image-push-arm64
	docker manifest create \
		$(DOCKER_REPOSITORY):$(REPO_REV) \
		--amend $(DOCKER_REPOSITORY):$(REPO_REV)-arm64 \
		--amend $(DOCKER_REPOSITORY):$(REPO_REV)-amd64
	docker manifest push $(DOCKER_REPOSITORY):$(REPO_REV)

WHAT ?= bash

.PHONY: provision
provision: build
provision:
	./bin/konvoy-image provision --inventory-file $(INVENTORY_FILE)  \
	-v ${VERBOSITY} \
	$(if $(ADDITIONAL_OVERRIDES),--overrides=${ADDITIONAL_OVERRIDES}) \
	$(if $(EXTRA_OVERRIDE_VARS), --extra-vars=${EXTRA_OVERRIDE_VARS})

.PHONY: dev
dev: ## dev build
dev: clean generate build lint test mod-tidy build.snapshot

.PHONY: ci
ci: ## CI build
ci: dev diff

.PHONY: clean
clean: ## remove files created during build
	$(call print-target)
	rm -rf bin
	rm -rf dist
	rm -rf artifacts
	rm -rf "$(REPO_ROOT_DIR)/cmd/konvoy-image-wrapper/image/konvoy-image-builder.tar.gz"
	rm -f flatcar-version.yaml
	rm -f $(COVERAGE)*
	docker image rm $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG) || echo "image already removed"
	docker image rm $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-amd64 || echo "image already removed"
	docker image rm $(DEVKIT_IMAGE_NAME):$(DEVKIT_IMAGE_TAG)-arm64 || echo "image already removed"
	docker buildx rm konvoy-image-builder || echo "image already removed"

.PHONY: generate
generate: ## go generate
	$(call print-target)
	go generate ./...

###### build arch specific konvoy-image binary
# TODO: refactor targets to remove duplication
bin/konvoy-image: $(REPO_ROOT_DIR)/cmd
bin/konvoy-image: $(shell find $(REPO_ROOT_DIR)/cmd -type f -name '*'.go)
bin/konvoy-image: $(REPO_ROOT_DIR)/pkg
bin/konvoy-image: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.go)
bin/konvoy-image: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.hcl)
bin/konvoy-image:
	$(call print-target)
	GOARCH=$(BUILDARCH) GOOS=$(GOOS) go build \
		-ldflags='-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=$(REPO_REV)' \
		-o ./dist/konvoy-image_$(GOOS)_$(GOARCH)/konvoy-image ./cmd/konvoy-image/main.go
	mkdir -p bin
	ln -sf ../dist/konvoy-image_$(GOOS)_$(GOARCH)/konvoy-image bin/konvoy-image

# Creates bin/konvoy-image-amd64 which will be copied to KIB container image for amd64. see Dockerfile
bin/konvoy-image-amd64: $(REPO_ROOT_DIR)/cmd
bin/konvoy-image-amd64: $(shell find $(REPO_ROOT_DIR)/cmd -type f -name '*'.go)
bin/konvoy-image-amd64: $(REPO_ROOT_DIR)/pkg
bin/konvoy-image-amd64: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.go)
bin/konvoy-image-amd64: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.hcl)
bin/konvoy-image-amd64:
	$(call print-target)
	GOARCH=amd64 GOOS=$(GOOS) go build \
		-ldflags='-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=$(REPO_REV)' \
		-o ./dist/konvoy-image_linux_amd64/konvoy-image ./cmd/konvoy-image/main.go
	mkdir -p bin
	ln -sf ../dist/konvoy-image_linux_amd64/konvoy-image bin/konvoy-image-amd64

# Creates bin/konvoy-image-arm64 which will be copied to KIB container image for arm64. see Dockerfile
bin/konvoy-image-arm64: $(REPO_ROOT_DIR)/cmd
bin/konvoy-image-arm64: $(shell find $(REPO_ROOT_DIR)/cmd -type f -name '*'.go)
bin/konvoy-image-arm64: $(REPO_ROOT_DIR)/pkg
bin/konvoy-image-arm64: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.go)
bin/konvoy-image-arm64: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.hcl)
bin/konvoy-image-arm64:
	$(call print-target)
	GOARCH=arm64 GOOS=$(GOOS) go build \
		-ldflags='-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=$(REPO_REV)' \
		-o ./dist/konvoy-image_linux_arm64/konvoy-image ./cmd/konvoy-image/main.go
	mkdir -p bin
	ln -sf ../dist/konvoy-image_linux_arm64/konvoy-image bin/konvoy-image-arm64

konvoy-image-amd64:
	$(MAKE) bin/konvoy-image-amd64 GOOS=linux GOARCH=amd64

konvoy-image-arm64:
	$(MAKE) bin/konvoy-image-arm64 GOOS=linux GOARCH=arm64

###### build konvoy image wrapper

bin/konvoy-image-wrapper: kib-image-build-$(BUILDARCH)
	$(call print-target)
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build \
		-ldflags='-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=$(REPO_REV)-$(BUILDARCH)' \
		-o ./bin/konvoy-image-wrapper ./cmd/konvoy-image-wrapper/main.go

dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image: $(REPO_ROOT_DIR)/cmd
dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image: $(shell find $(REPO_ROOT_DIR)/cmd -type f -name '*'.go)
dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image: $(REPO_ROOT_DIR)/pkg
dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.go)
dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image: $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.hcl)
dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image:
	$(call print-target)
	goreleaser build --snapshot --clean --id konvoy-image --single-target

.PHONY: build
build: bin/konvoy-image
build: ## go build

.PHONY: build-wrapper
build-wrapper: bin/konvoy-image-wrapper

.PHONY: docs
docs: build
	$(REPO_ROOT_DIR)/bin/konvoy-image generate-docs $(REPO_ROOT_DIR)/docs/cli

.PHONY: docs.check
docs.check: docs
docs.check:
	@test -z "$(shell git status --porcelain -- $(REPO_ROOT_DIR)/docs)" \
		|| (echo ''; \
			echo 'Need docs update:'; \
			echo ''; \
			git status --porcelain -- "$(REPO_ROOT_DIR)/docs"; \
			echo ''; \
			echo 'Run `make docs` and commit the results'; \
			exit 1)

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run -c .golangci.yml --fix

# Add a convience alias
.PHONY: super-linter
super-linter: super-lint

.PHONY: super-lint
include $(REPO_ROOT_DIR)/.github/super-linter.env
export
export DOCKER_SUPER_LINTER_ARGS ?= \
	--env RUN_LOCAL=true \
	--env-file $(REPO_ROOT_DIR)/.github/super-linter.env \
	--volume $(REPO_ROOT_DIR):/tmp/lint
export DOCKER_SUPER_LINTER_VERSION ?= $(shell \
	grep 'uses: github/super-linter' $(REPO_ROOT_DIR)/.github/workflows/lint.yml | cut -d@ -f2 \
)
export DOCKER_SUPER_LINTER_IMG := github/super-linter:$(DOCKER_SUPER_LINTER_VERSION)

super-lint: ## run all linting with super-linter
	$(call print-target)
	docker run \
		--rm \
		$(if $(INTERACTIVE),--tty) \
		--interactive \
		$(DOCKER_SUPER_LINTER_ARGS) \
		$(DOCKER_SUPER_LINTER_IMG)

.PHONY: super-lint-shell
super-lint-shell: ## open a shell in the super-linter container
	$(call print-target)
	docker run \
		--rm \
		$(if $(INTERACTIVE),--tty) \
		--interactive \
		$(DOCKER_SUPER_LINTER_ARGS) \
		--workdir=/tmp/lint \
		--entrypoint="/bin/bash" \
		$(DOCKER_SUPER_LINTER_IMG) -l

.PHONY: test
test: ## go test with race detector and code coverage
	$(call print-target)
	CGO_ENABLED=1 go test $(shell go list ./... | grep -v e2e) -- -race -short -v  

.PHONY: integration-test
integration-test: ## go test with race detector for integration tests
	$(call print-target)
	CGO_ENABLED=1 go test -race -run Integration -v ./...

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	$(call print-target)
	go mod tidy

.PHONY: build.snapshot
build.snapshot: dist/konvoy-image_linux_amd64/konvoy-image
build.snapshot:
	$(call print-target)
	# NOTE(jkoelker) shenanigans to get around goreleaser and
	#                `make release-bundle` being able to share the same
	#                `Dockerfile`. Unfortunatly goreleaser forbids
	#                copying the dist folder into the temporary folder
	#                that it uses as its docker build context ;(.
	# NOTE (faiq): does anyone use this target?
	mkdir -p bin
	cp dist/konvoy-image_linux_$(BUILDARCH)/konvoy-image bin/konvoy-image
	goreleaser --parallelism=1 --skip-publish --snapshot --clean

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi



.PHONY: release
release:
	$(call print-target)
	DOCKER_BUILDKIT=1 goreleaser release --clean --parallelism=1 --timeout=2h

.PHONY: release-snapshot
release-snapshot:
	$(call print-target)
	DOCKER_BUILDKIT=1 goreleaser --parallelism=1 --clean --snapshot --timeout=2h

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	$(call print-target)
	go clean -r -i -cache -testcache -modcache

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef

cmd/konvoy-image-wrapper/image/konvoy-image-builder.tar.gz: kib-image-push-manifest
	# we need to build the appropriate image for the bundle we're creating
	# followed by saving it as just image name so that we can put in the release tar
	# the docker images are published before this by hack/release.sh, making this safe.
	docker pull $(DOCKER_REPOSITORY):$(REPO_REV)-$(BUILDARCH)
	docker tag $(DOCKER_REPOSITORY):$(REPO_REV)-$(BUILDARCH) $(DOCKER_REPOSITORY):$(REPO_REV)
	docker save $(DOCKER_REPOSITORY):$(REPO_REV) | gzip -c - > "$(REPO_ROOT_DIR)/cmd/konvoy-image-wrapper/image/konvoy-image-builder.tar.gz"
