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
#export GOOS ?= linux
ARCH := $(shell uname -m)

BUILDARCH ?= $(shell echo $(ARCH) | sed 's/x86_64/amd64/g')

export CI ?= no
ifeq ($(CI),yes)
export TEAMCITY_EXTRA_MOUNT ?= /teamcity
export TEAMCITY_BUILD_ID ?= $(shell date +%s)
endif


export DOCKER_REPOSITORY ?= mesosphere/konvoy-image-builder
export DOCKER_SOCKET ?= /var/run/docker.sock
ifeq ($(OS),Darwin)
export DOCKER_SOCKET_GID ?= $(shell /usr/bin/stat -f "%g" $(DOCKER_SOCKET))
else
export DOCKER_SOCKET_GID ?= $(shell stat -c %g $(DOCKER_SOCKET))
endif

export DOCKER_ARCH_IMG ?= $(DOCKER_REPOSITORY):$(REPO_REV)-$(BUILDARCH)
export DOCKER_PHONY_FILE ?= .docker-$(shell echo '$(DOCKER_ARCH_IMG)' | tr '/:' '.')

export DOCKER_DEVKIT_IMG ?= $(DOCKER_REPOSITORY):latest-devkit-$(BUILDARCH)
export DOCKER_DEVKIT_PHONY_FILE ?= .docker-$(shell echo '$(DOCKER_DEVKIT_IMG)' | tr '/:' '.')
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

ifeq ($(strip $(TEAMCITY_EXTRA_MOUNT)),)
DOCKER_GCP_CREDENTIALS_ARGS=--volume "$(HOME)/.gcloud":"/home/$(USER_NAME)/.gcloud" \
	                             --env GOOGLE_APPLICATION_CREDENTIALS=/home/$(USER_NAME)/.gcloud/credentials.json
else
DOCKER_GCP_CREDENTIALS_ARGS=--volume $(TEAMCITY_EXTRA_MOUNT):$(TEAMCITY_EXTRA_MOUNT) \
								 --env GOOGLE_APPLICATION_CREDENTIALS=$(TEAMCITY_EXTRA_MOUNT)/$(TEAMCITY_BUILD_ID)-credentials.json
endif

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

# NOTE(jkoelker) Abuse ifeq and the junk variable to proxy docker image state
#                to the target file
ifneq ($(shell command -v docker),)
	ifeq ($(shell docker image ls --quiet "$(DOCKER_DEVKIT_IMG)"),)
		export junk := $(shell rm -rf $(DOCKER_DEVKIT_PHONY_FILE))
	endif
	ifeq ($(shell docker image ls --quiet "$(DOCKER_ARCH_IMG)"),)
		export junk := $(shell rm -rf $(DOCKER_PHONY_FILE))
	endif
endif

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


include make/ci.mk
include make/images.mk
include hack/pip-packages/Makefile
include test/infra/aws/Makefile
include test/infra/vsphere/Makefile


.PHONY: dev
dev: ## dev build
dev: clean generate build lint test mod-tidy bin/konvoy-image

.PHONY: generate
generate: ## go generate
	$(call print-target)
	go generate ./...

##########################
# helper target to run provision
##########################
.PHONY: provision
provision: build
provision:
	./bin/konvoy-image provision --inventory-file $(INVENTORY_FILE)  \
	-v ${VERBOSITY} \
	$(if $(ADDITIONAL_OVERRIDES),--overrides=${ADDITIONAL_OVERRIDES}) \
	$(if $(EXTRA_OVERRIDE_VARS), --extra-vars=${EXTRA_OVERRIDE_VARS})

##########################
# Build konvoy-image multi arch binaries
##########################

bin/konvoy-image_$(GOOS)_$(BUILDARCH): $(REPO_ROOT_DIR)/cmd
bin/konvoy-image_$(GOOS)_$(BUILDARCH): $(shell find $(REPO_ROOT_DIR)/cmd -type f -name '*'.go)
bin/konvoy-image_$(GOOS)_$(BUILDARCH): $(REPO_ROOT_DIR)/pkg
bin/konvoy-image_$(GOOS)_$(BUILDARCH): $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.go)
bin/konvoy-image_$(GOOS)_$(BUILDARCH): $(shell find $(REPO_ROOT_DIR)/pkg -type f -name '*'.tmpl)
bin/konvoy-image_$(GOOS)_$(BUILDARCH):
	$(call print-target)
	$(MAKE) docker GOOS=$(GOOS) GOARCH=$(BUILDARCH) WHAT="go build \
		-ldflags='-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=$(REPO_REV)' \
		-o ./bin/konvoy-image_$(GOOS)_$(BUILDARCH) ./cmd/konvoy-image/main.go"

# create konvoy-image binary for current OS and architecture
bin/konvoy-image: bin/konvoy-image_$(GOOS)_$(BUILDARCH)
	$(call print-target)
	ln -sf ../bin/konvoy-image_$(GOOS)_$(BUILDARCH) bin/konvoy-image

# build konvoy-image binary for current OS and architecture
.PHONY: build
build: bin/konvoy-image

##########################
# Build devkit container
##########################
.PHONY: buildx
buildx:
	$(call print-target)
	docker run --privileged --rm tonistiigi/binfmt --install all || true

github-token.txt:
	$(call print-target)
	echo $(GITHUB_TOKEN) >> github-token.txt

$(DOCKER_DEVKIT_PHONY_FILE): github-token.txt buildx
$(DOCKER_DEVKIT_PHONY_FILE): Dockerfile.devkit install-envsubst
	$(call print-target)
	docker buildx build \
	--build-arg USER_ID=$(UID) \
	--build-arg GROUP_ID=$(GID) \
	--build-arg USER_NAME=$(USER_NAME) \
	--build-arg GROUP_NAME=$(GROUP_NAME) \
	--build-arg DOCKER_GID=$(DOCKER_SOCKET_GID) \
	--build-arg BUILDARCH=$(BUILDARCH) \
	$(if $(GITHUB_ACTION),--secret id=githubtoken$(COMMA)src=github-token.txt) \
	--platform linux/$(BUILDARCH) \
	--file $(REPO_ROOT_DIR)/Dockerfile.devkit \
	--tag=$(DOCKER_DEVKIT_IMG) \
	$(REPO_ROOT_DIR) && \
	touch $(DOCKER_DEVKIT_PHONY_FILE)

.PHONY: devkit
devkit:
	$(call print-target)
	$(MAKE) $(DOCKER_DEVKIT_PHONY_FILE)

.PHONY: devkit.run
devkit.run: ## run $(WHAT) in devkit
devkit.run: devkit
	$(call print-target)
	docker run \
		$(DOCKER_DEVKIT_DEFAULT_ARGS) \
		$(DOCKER_DEVKIT_ARGS) \
		"$(DOCKER_DEVKIT_IMG)" \
		$(WHAT)

##########################
# Build wrapper
##########################

# Docker image that gets embedded in the konvoy-image-wrapper has two dependencies
# 1. konvoy-image binary for linux OS and each amd64/arm64 arch
# 2. Devkit image used as base image per each amd64/arm64 arch
# based on BUILDARCH, build arm64 or amd64 compatible wrapper image
$(DOCKER_PHONY_FILE): $(DOCKER_DEVKIT_PHONY_FILE)
$(DOCKER_PHONY_FILE): Dockerfile
	$(call print-target)
	$(MAKE) bin/konvoy-image GOOS=linux BUILDARCH=$(BUILDARCH)
	DOCKER_BUILDKIT=1 docker buildx build \
		--file $(REPO_ROOT_DIR)/Dockerfile \
		--build-arg BUILDARCH=$(BUILDARCH) \
		--platform linux/$(BUILDARCH) \
		--build-arg BASE=docker.io/$(DOCKER_DEVKIT_IMG) \
		--tag=$(DOCKER_ARCH_IMG) \
		$(REPO_ROOT_DIR) \
	&& touch $(DOCKER_PHONY_FILE)

.PHONY: build-wrapper-image
build-wrapper-image:
	$(call print-target)
	$(MAKE) $(DOCKER_PHONY_FILE)
# builds konvoy-image-wrapper without go tag 'EMBED_DOCKER_IMAGE_arm64 or EMBED_DOCKER_IMAGE_amd64' to build wrapper without embedding. 
# this enables local testing faster.
# .goreleaser.yml file embeds the docker image using  EMBED_DOCKER_IMAGE_arm64 and EMBED_DOCKER_IMAGE_amd64 flags when releasing artifacts
bin/konvoy-image-wrapper: $(DOCKER_PHONY_FILE)
bin/konvoy-image-wrapper:
	$(call print-target)
	$(MAKE) docker WHAT="go build \
		-ldflags='-X github.com/mesosphere/konvoy-image-builder/pkg/version.version=$(REPO_REV)' \
		-o ./bin/konvoy-image-wrapper ./cmd/konvoy-image-wrapper/main.go"
	docker tag $(DOCKER_ARCH_IMG) $(DOCKER_REPOSITORY):$(REPO_REV)

# builds konvoy image wrapper binary without embedding the docker contaienr file.
.PHONY: build-wrapper
build-wrapper: bin/konvoy-image-wrapper

# This target is used when building release artifacts using goreleaser. 
# goreleaser embeds the container image in the final konvoy-image-wrapper binary
cmd/konvoy-image-wrapper/image/konvoy-image-builder_linux_$(BUILDARCH).tar.gz: $(DOCKER_PHONY_FILE)
	$(call print-target)
	docker save $(DOCKER_ARCH_IMG) | gzip -c - > "$(REPO_ROOT_DIR)/cmd/konvoy-image-wrapper/image/konvoy-image-builder_linux_$(BUILDARCH).tar.gz"

##########################
# Relese Targets
##########################

.PHONY: push-wrapper-image
push-wrapper-image: build-wrapper-image
	$(call print-target)
	docker push $(DOCKER_REPOSITORY):$(REPO_REV)-$(BUILDARCH)

.PHONY: push-manifest
push-manifest:
	$(call print-target)
	docker manifest create \
		$(DOCKER_REPOSITORY):$(REPO_REV) \
		--amend $(DOCKER_REPOSITORY):$(REPO_REV)-arm64 \
		--amend $(DOCKER_REPOSITORY):$(REPO_REV)-amd64
	docker manifest push $(DOCKER_REPOSITORY):$(REPO_REV)

.PHONY: release
release:
	$(call print-target)
# set --parallelism=1 because the goreleaser pre executes hook will run pre execute hook
# cmd/konvoy-image-wrapper/image/konvoy-image-builder_linux_amd64.tar.gz in parallel for each linux, darwin and windows binary.
# this can corrupt content of the file.
	goreleaser --parallelism=1 --rm-dist --debug --snapshot --parallelism=1

.PHONY: release-snapshot
release-snapshot:
	$(call print-target)
# set --parallelism=1 because the goreleaser pre executes hook will run pre execute hook
# cmd/konvoy-image-wrapper/image/konvoy-image-builder_linux_amd64.tar.gz in parallel for each linux, darwin and windows binary.
# this can corrupt content of the file.
	goreleaser release --snapshot --skip-publish --rm-dist --parallelism=1 --debug

##########################
# docs targets
##########################
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

##########################
# linter targets
##########################
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

##########################
# Test targets
##########################
.PHONY: test
test: ## go test with race detector and code coverage
	$(call print-target)
	CGO_ENABLED=1 go-acc --covermode=atomic --output=$(COVERAGE).out --ignore=e2e ./... -- -race -short -v
ifneq ($(CI),)
	gocover-cobertura -by-files < $(COVERAGE).out > $(COVERAGE).xml
else
	go tool cover -html=$(COVERAGE).out -o $(COVERAGE).html
endif

.PHONY: integration-test
integration-test: ## go test with race detector for integration tests
	$(call print-target)
	CGO_ENABLED=1 go test -race -run Integration -v ./...

##########################
# clean targets
##########################
.PHONY: clean
clean: ## remove files created during build
	$(call print-target)
	rm -rf bin
	rm -rf dist
	rm -rf artifacts
	rm -rf "$(REPO_ROOT_DIR)/cmd/konvoy-image-wrapper/image/konvoy-image-builder.tar.gz"
	rm -f flatcar-version.yaml
	rm -f $(COVERAGE)*
	docker image rm $(DOCKER_DEVKIT_IMG) || echo "image already removed"
	docker buildx rm konvoy-image-builder || echo "image already removed"

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	$(call print-target)
	go clean -r -i -cache -testcache -modcache

##########################
# helper targets
##########################
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi


# TODO: remove this target and its references once all jobs are moved to github actions
# explicity install `go` in github actions runner so that we dont have to download docker container
# github actions can cache go-installer so installing `go` in runner will be quick.
WHAT ?= bash
.PHONEY: docker
docker:
	docker run \
	--rm \
	$(DOCKER_ULIMIT_ARGS) \
	--volume $(REPO_ROOT_DIR):/build \
	--workdir /build \
	--env GOOS \
	--env GOARCH \
	--env BUILDARCH \
	$(GOLANG_IMAGE) \
	/bin/bash -c "$(WHAT)"

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	$(call print-target)
	go mod tidy

.PHONY: ci
ci: ## CI build
ci: dev diff
