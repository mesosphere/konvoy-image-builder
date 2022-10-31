import os

import pytest
import testinfra.utils.ansible_runner

testinfra_hosts = testinfra.utils.ansible_runner.AnsibleRunner(
    os.environ["MOLECULE_INVENTORY_FILE"]
).get_hosts("all")


def test_containerd_running_and_enabled(host):
    containerd = host.service("containerd")
    assert containerd.is_enabled


def test_kubelet_running_and_enabled(host):
    kubelet = host.service("kubelet")
    assert kubelet.is_enabled


def test_kubectl_avail(host):
    # the path is only set on interactive shell. SO lets append it here
    cmd = host.run("bash -c 'PATH=$PATH:/opt/bin type kubectl'")
    assert cmd.succeeded is True


def test_kubeadm_avail(host):
    # the path is only set on interactive shell. SO lets append it here
    cmd = host.run("bash -c 'PATH=$PATH:/opt/bin type kubeadm'")
    assert cmd.succeeded is True


def test_cloudinit_feature_flags(host):
    """
    ubuntu 18.04: does not need the feature flag
    all other except flatcar: expect feature overrides
    """
    distro = host.system_info.distribution
    release = host.system_info.release

    # for flatcar we can skip
    if distro == "flatcar":
        pytest.skip("no changes on flatcar")

    # if cloud-init is lower than 20.0 we can skip
    cloud_init_version = host.run("cloud-init --version")
    assert cloud_init_version.succeeded

    cloud_init_version_str = cloud_init_version.stdout.strip("\n")
    if not cloud_init_version_str:
        cloud_init_version_str = cloud_init_version.stderr.strip("\n")

    assert cloud_init_version_str

    cloud_init_version_str_version_part = cloud_init_version_str.split(" ")[-1]
    major_version = cloud_init_version_str_version_part.split(".")[0]

    if int(major_version) < 20:
        pytest.skip("cloud-init major version ({}) below 20".format(major_version))

    if distro != "ubuntu":
        cmd = host.run(
            "python3 -c \"import sysconfig; print(sysconfig.get_path('purelib'))\""
        )
        assert cmd.succeeded

        featurefile = host.file(
            "{}/cloudinit/feature_overrides.py".format(cmd.stdout.strip("\n"))
        )
        assert featurefile.exists
        assert b"ERROR_ON_USER_DATA_FAILURE = False" in featurefile.content
    # ubuntu 18.04 still supported and no need for this feature flag
    elif distro == "ubuntu" and not release == "18.04":
        featurefile = host.file(
            "/usr/lib/python3/dist-packages/cloudinit/feature_overrides.py"
        )
        assert featurefile.exists
        assert b"ERROR_ON_USER_DATA_FAILURE = False" in featurefile.content
    else:
        assert True
