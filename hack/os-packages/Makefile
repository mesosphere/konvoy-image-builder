RPM_PACKAGES_GETTER_IMAGE = mesosphere/rpm-artifacts-getter
RPM_FIPS_PACKAGES_GETTER_IMAGE = mesosphere/rpm-artifacts-getter-fips

KUBERNETES_RPM_PATCH = 0

.PHONY: os-packages-artifacts
os-packages-artifacts: download-all-rpms artifacts/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_x86_64_rpms.tar.gz

artifacts/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_x86_64_rpms.tar.gz:
	cd artifacts/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_el7_x86_64_rpms && tar -czf ../os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_x86_64_rpms.tar.gz *

.PHONY: download-all-rpms
download-all-rpms:
	$(MAKE) rhel_major_version=7 fips=0 bastion=0  download-rpms
#	$(MAKE) rhel_major_version=8 fips=0 bastion=0  download-rpms
#	$(MAKE) rhel_major_version=7 fips=0 bastion=1  download-rpms
#	$(MAKE) rhel_major_version=8 fips=0 bastion=1  download-rpms
#	$(MAKE) rhel_major_version=7 fips=1 bastion=0  download-rpms
#	$(MAKE) rhel_major_version=8 fips=1 bastion=0  download-rpms

.PHONY: download-rpms
download-rpms:
	$(MAKE) rhel_major_version=$(rhel_major_version) fips=$(fips) rpm-artifacts-getter.check
	$(MAKE) rhel_major_version=$(rhel_major_version) fips=$(fips) packages="audit ca-certificates conntrack-tools chrony curl ebtables open-vm-tools python3-pip python-netifaces python-requests socat sysstat yum-utils yum-plugin-versionlock nfs-utils NetworkManager" download-rpm
	$(MAKE) rhel_major_version=$(rhel_major_version) fips=$(fips) packages="container-selinux containerd.io-$(CONTAINERD_VERSION) container-selinux" download-rpm
	$(MAKE) rhel_major_version=$(rhel_major_version) fips=$(fips) packages="kubectl-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-$(KUBERNETES_RPM_PATCH) kubelet-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-$(KUBERNETES_RPM_PATCH) conntrack ebtables ethtool iproute iptables kubernetes-cni libnetfilter_cthelper libnetfilter_cttimeout libnetfilter_queue libblkid libmount libsmartcols libuuid socat util-linux kubeadm-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-$(KUBERNETES_RPM_PATCH) kubernetes-cni cri-tools" download-rpm
	curl https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm -o $(artifacts_dir)/amazon-ssm-agent.rpm
	$(MAKE) rhel_major_version=$(rhel_major_version) fips=$(fips) createrepo

ifeq ($(fips), 1)
artifacts_dir = $(shell pwd)/artifacts/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_el$(rhel_major_version)_x86_64_fips_rpms
rpm_getter_image = $(RPM_PACKAGES_GETTER_IMAGE)-fips-$(rhel_major_version)
getter_target = rpm-fips-artifacts-getter
endif
ifeq ($(fips), 0)
artifacts_dir = $(shell pwd)/artifacts/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_el$(rhel_major_version)_x86_64_rpms
rpm_getter_image = $(RPM_PACKAGES_GETTER_IMAGE)-$(rhel_major_version)
getter_target = rpm-artifacts-getter
endif

.PHONY: rpm-artifacts-getter.check
rpm-artifacts-getter.check:
	docker image inspect $(rpm_getter_image):$(DEFAULT_KUBERNETES_VERSION_SEMVER) > /dev/null || $(MAKE) fips=$(fips) rhel_major_version=$(rhel_major_version) rpm-artifacts-getter

.PHONY: rpm-artifacts-getter
rpm-artifacts-getter:
	@docker build \
	  -f hack/os-packages/$(getter_target)-$(rhel_major_version)/Dockerfile \
	  -t $(rpm_getter_image):$(DEFAULT_KUBERNETES_VERSION_SEMVER) ./hack/os-packages/$(getter_target)-$(rhel_major_version)/

$(artifact_dir):
	mkdir -p $(artifact_dir)

.PHONY: download-rpm
download-rpm: $(artifact_dir)
	docker run --rm -v "$(artifacts_dir):/opt/dkp/rpms/" -v /tmp/dnfcache:/var/cache/dnf -w /opt/dkp/rpms $(rpm_getter_image):$(DEFAULT_KUBERNETES_VERSION_SEMVER) \
		/bin/sh -c 'yumdownloader -y --setopt=skip_missing_names_on_install=False --destdir=. -x \*i686 --archlist=x86_64,noarch --resolve $(packages) && chown -R $(shell id -u):$(shell id -g) .'

.PHONY: createrepo
createrepo:
	docker run --rm -v "$(artifacts_dir):/opt/dkp/rpms/" -v /tmp/dnfcache:/var/cache/dnf -w /opt/dkp/rpms $(rpm_getter_image):$(DEFAULT_KUBERNETES_VERSION_SEMVER) \
		/bin/sh -c 'createrepo -v .'
