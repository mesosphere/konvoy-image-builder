rpm_packages_getter_image = mesosphere/rpm-artifacts-getter
KUBERNETES_RPM_PATCH = 0
# possible values: https://docs.ansible.com/ansible/2.6/user_guide/playbooks_conditionals.html#ansible-distribution
# RedHat, CentOS, Ubuntu, Debian, SUSE, Amazon
# supported OS/version/fips matrix: https://docs.d2iq.com/dkp/konvoy/2.1/supported-operating-systems/
os_distribution ?= centos
os_distribution_major_version ?= 7
os_arch ?= x86_64
fips ?= 0

ifeq ($(fips), 1)
artifacts_dir = $(shell pwd)/artifacts/os-packages/$(os_distribution)/$(os_distribution_major_version)/$(os_arch)/fips/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_$(os_arch)
rpm_getter_image = $(rpm_packages_getter_image)-$(os_distribution)-fips:$(os_distribution_major_version)
getter_target = rpm-artifacts-getter/$(os_distribution)/$(os_distribution_major_version)/fips
endif
ifeq ($(fips), 0)
artifacts_dir = $(shell pwd)/artifacts/os-packages/$(os_distribution)/$(os_distribution_major_version)/$(os_arch)/os-packages_$(DEFAULT_KUBERNETES_VERSION_SEMVER)_$(os_arch)
rpm_getter_image = $(rpm_packages_getter_image)-$(os_distribution):$(os_distribution_major_version)
getter_target = rpm-artifacts-getter/$(os_distribution)/$(os_distribution_major_version)
endif

common_packages="audit ca-certificates conntrack-tools chrony curl ebtables open-vm-tools python3-pip python3-netifaces python3-requests socat sysstat yum-utils yum-plugin-versionlock nfs-utils NetworkManager tzdata nspr nss python3-chardet yum-utils nss-softokn-freebl nss-util nss-softokn nss nss-tools  NetworkManager-libnm libdrm libpciaccess python3-libxml2 python3-kitchen dnf-plugin-ovl"
containerd_packages="container-selinux containerd.io-$(CONTAINERD_VERSION)"
kubernetes_packages="kubectl-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-$(KUBERNETES_RPM_PATCH) kubelet-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-$(KUBERNETES_RPM_PATCH) conntrack ebtables ethtool iproute iptables kubernetes-cni libnetfilter_cthelper libnetfilter_cttimeout libnetfilter_queue libblkid libmount libsmartcols libuuid socat util-linux kubeadm-$(DEFAULT_KUBERNETES_VERSION_SEMVER)-$(KUBERNETES_RPM_PATCH) kubernetes-cni cri-tools"
centos_7_packages="python-netifaces python-requests libxml2-python python-kitchen python-chardet yum-plugin-ovl python2-dnf"
rhel_7_packages=$(centos_7_packages)
rhel_8_packages="python3-netifaces python3-requests python3-libxml2 python3-kitchen python3-chardet dnf-plugin-ovl nfs-utils keyutils-libs python36 platform-python-setuptools iproute-tc chkconfig"

.PHONY: os-packages-artifacts
os-packages-artifacts: $(artifacts_dir) archive-os-packages

$(artifacts_dir):
	mkdir -p $(artifacts_dir)

.PHONY: archive-os-packages
archive-os-packages:
	$(MAKE) download-rpms
	cd $(artifacts_dir) && tar -czf ../$(shell basename $(artifacts_dir)).tar.gz *

.PHONY: download-rpms
download-rpms:
	$(MAKE) rpm-artifacts-getter.check
	$(MAKE) packages="$(subst $\",,$(common_packages) $(containerd_packages) $(kubernetes_packages))" download-rpm
	$(MAKE) download-rpms-$(os_distribution)-$(os_distribution_major_version)
	curl https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm -o $(artifacts_dir)/amazon-ssm-agent.rpm
	$(MAKE) createrepo-$(os_distribution)-$(os_distribution_major_version)

download-rpms-redhat-7:
	$(MAKE) packages=$(rhel_7_packages) download-rpm
download-rpms-redhat-8:
	$(MAKE) packages=$(rhel_8_packages) download-rpm
download-rpms-centos-7:
	$(MAKE) packages=$(centos_7_packages) download-rpm
download-rpms%:    #default filtering
	@echo "no OS specific packages. skipping..."

ifndef rhel_major_version
rhel_major_version = 7
endif

.PHONY: rpm-artifacts-getter.check
rpm-artifacts-getter.check:
	docker image inspect $(rpm_getter_image) > /dev/null || $(MAKE) rpm-artifacts-getter

.PHONY: rpm-artifacts-getter
rpm-artifacts-getter:
	@docker build \
	  -f hack/os-packages/$(getter_target)/Dockerfile \
	  -t $(rpm_getter_image) ./hack/os-packages/$(getter_target)

.PHONY: download-rpm
download-rpm: $(artifacts_dir)
	docker run --rm -v "$(artifacts_dir):/opt/dkp/rpms/" -v /tmp/dnfcache:/var/cache/dnf -w /opt/dkp/rpms $(rpm_getter_image) \
	/bin/sh -c 'yumdownloader -y --setopt=skip_missing_names_on_install=False --destdir=. -x \*i686 --archlist=x86_64,noarch --resolve $(packages) && chown -R $(shell id -u):$(shell id -g) .'

.PHONY: createrepo-centos
createrepo-centos-%:
	docker run --rm -v "$(artifacts_dir):/opt/dkp/rpms/" -v /tmp/dnfcache:/var/cache/dnf -w /opt/dkp/rpms $(rpm_getter_image) \
		/bin/sh -c 'createrepo -v . && chown -R $(shell id -u):$(shell id -g) repodata/'

.PHONY: createrepo-redhat-7
createrepo-redhat-7: createrepo-centos-7

.PHONY: createrepo-redhat-8
createrepo-redhat-8:
	docker run --rm -v "$(artifacts_dir):/opt/dkp/rpms/" -v /tmp/dnfcache:/var/cache/dnf -w /opt/dkp/rpms $(rpm_getter_image) \
		/bin/sh -c 'createrepo -v . && chown -R $(shell id -u):$(shell id -g) repodata/ && repo2module . --module-name offline && createrepo_mod .'

.PHONY: clean-artifacts
clean-artifacts:
	rm -rf $(shell pwd)/artifacts
	for image in $$(\
		docker image ls \
			--filter=reference='$(RPM_PACKAGES_GETTER_IMAGE)*' \
			--format "{{.Repository}}:{{.Tag}}"); do \
		docker image rm $$image;\
	done
