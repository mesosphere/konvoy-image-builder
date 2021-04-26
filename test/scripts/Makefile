SHELL := /bin/bash -euo pipefail

TEST_IMAGE_REPO ?= konvoy-image-builder-test
TEST_IMAGE_TAG ?= latest

docker.build:
	docker build -t $(TEST_IMAGE_REPO):$(TEST_IMAGE_TAG) .

TEST_CONTAINER_NAME ?= konvoy-image-builder-ansible-test
WORKDIR := /root/workdir

USER_NAME ?= $(shell id -u -n)

AWS_ARGS ?= \
	--env AWS_PROFILE \
	--env AWS_SECRET_ACCESS_KEY \
	--env AWS_SESSION_TOKEN \
	--env AWS_REGION \
	-v "$(HOME)/.aws":"/root/.aws"
ifneq ($(CI), true)
DEFAULT_FLAGS ?= -it
endif
FLAGS ?= -v $(CURDIR)/../../:$(WORKDIR) -w $(WORKDIR) $(AWS_ARGS)
WHAT ?= bash

docker.run: docker.build
	docker run \
	--name $(TEST_CONTAINER_NAME) \
	--rm \
	$(DEFAULT_FLAGS) \
	$(FLAGS) \
	$(TEST_IMAGE_REPO):$(TEST_IMAGE_TAG) $(WHAT)

docker.stop:
	docker stop $(TEST_CONTAINER_NAME)
