PIP_PACKAGES_GETTER_IMAGE = mesosphere/pip-artifacts-getter

pip_artifact_dir = $(shell pwd)/artifacts/pip-packages

.PHONY: pip-packages-artifacts
pip-packages-artifacts: download-pip-packages artifacts/pip-packages.tar.gz

artifacts/pip-packages.tar.gz:
	cd artifacts/pip-packages && tar -czf ../pip-packages.tar.gz *

.PHONY: download-pip-packages
download-pip-packages:
	$(MAKE) pip-artifacts-getter.check
	$(MAKE) packages="awscli" download-pip-package

.PHONY: pip-artifacts-getter.check
pip-artifacts-getter.check:
	docker image inspect $(PIP_PACKAGES_GETTER_IMAGE):latest > /dev/null || $(MAKE) pip-artifacts-getter

.PHONY: pip-artifacts-getter
pip-artifacts-getter:
	@docker build \
	  -f hack/pip-packages/pip-getter/Dockerfile \
	  -t $(PIP_PACKAGES_GETTER_IMAGE):latest ./hack/pip-packages/pip-getter/.

$(pip_artifact_dir):
	mkdir -p $(pip_artifact_dir)

.PHONY: download-pip-package
download-pip-package: $(pip_artifact_dir)
	docker run --rm -v "$(pip_artifact_dir):/opt/dkp/pip/" -w /opt/dkp/pip $(PIP_PACKAGES_GETTER_IMAGE):latest \
		/bin/sh -c 'mkdir -p $(packages) && pip3 download --destination-directory=/opt/dkp/pip/$(packages) $(packages)'

