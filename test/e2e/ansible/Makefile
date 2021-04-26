ANSIBLE_IMAGE_REPO ?= konvoy-image-builder-ansible-e2e
ANSIBLE_IMAGE_TAG ?= latest

docker.build:
	docker build -t $(ANSIBLE_IMAGE_REPO):$(ANSIBLE_IMAGE_TAG) .

ANSIBLE_CONTAINER_NAME ?= konvoy-image-builder-ansible-e2e
SSH_PORT ?= 3001

E2E_SSH_KEY ?= generated/id_rsa
SSH_KEY ?= $(E2E_SSH_KEY)
PUBLIC_KEY ?= $(CURDIR)/$(E2E_SSH_KEY).pub
AUTHORIZED_KEYS_PATH ?= /home/sshuser/.ssh/authorized_keys
ifdef PUBLIC_KEY
PUBLIC_KEY_FLAG = -v $(PUBLIC_KEY):$(AUTHORIZED_KEYS_PATH)
endif

WHAT ?= bash

docker.run: docker.build
	docker run \
	--name $(ANSIBLE_CONTAINER_NAME) \
	--rm \
	-it \
	-p $(SSH_PORT):22 \
	$(PUBLIC_KEY_FLAG) \
	$(FLAGS) \
	$(ANSIBLE_IMAGE_REPO):$(ANSIBLE_IMAGE_TAG) $(WHAT)

docker.keygen-then-run:
	make keygen docker.run

docker.stop:
	docker stop $(ANSIBLE_CONTAINER_NAME)

ssh:
	ssh -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i $(SSH_KEY) sshuser@localhost -p $(SSH_PORT) $(SSH_WHAT)

keygen: keygen-clean
	ssh-keygen -t rsa -b 4096 -N "" -f $(E2E_SSH_KEY)

keygen-clean:
	rm $(E2E_SSH_KEY) $(E2E_SSH_KEY).pub || true

E2E_FLAGS ?= -d
# used to manually test if the container works
e2e.ssh:
	WHAT="" FLAGS="$(E2E_FLAGS)" make docker.keygen-then-run
	sleep 5
	SSH_WHAT=exit make ssh docker.stop keygen-clean
	WHAT="" FLAGS="$(E2E_FLAGS)" AUTHORIZED_KEYS_PATH=/setup/authorized_keys make docker.keygen-then-run
	sleep 5
	SSH_WHAT=exit make ssh docker.stop keygen-clean


.PHONY: e2e
e2e: e2e.setup e2e.run e2e.clean

e2e.setup:
	WHAT="" FLAGS="$(E2E_FLAGS)" make docker.keygen-then-run
	sleep 5

e2e.clean: docker.stop keygen-clean

e2e.run:
	go test -v -count=1 . # count=1 to break cache
