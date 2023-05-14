import os
import pytest
import testinfra.utils.ansible_runner

testinfra_hosts = testinfra.utils.ansible_runner.AnsibleRunner(
    os.environ["MOLECULE_INVENTORY_FILE"]
).get_hosts("all")

def test_kubelet_kubectl_installed(host):
    """
    we expect kubectl and kubelet package to be installed

    flatcar: skip no packages
    """
    distro = host.system_info.distribution
    if distro == "flatcar":
        pytest.skip("no packages on flatcar")

    assert host.package("kubectl").is_installed
    assert host.package("kubelet").is_installed

def test_kubeadm_installed(host):
    """
    we expect kubeadm package to be installed

    flatcar: skip no packages
    """
    distro = host.system_info.distribution
    if distro == "flatcar":
        pytest.skip("no packages on flatcar")

    assert host.package("kubeadm").is_installed

def test_kube_cmd_path(host):
    """
    kubelet, kubeadm and kubectl must be in path
    """
    distro = host.system_info.distribution
    if distro == "flatcar":
        pytest.skip("flatcar uses different PATH for non-interactive")

    assert host.exists("kubelet")
    assert host.exists("kubeadm")
    assert host.exists("kubectl")

def test_kube_cmd_path_flatcar(host):
    distro = host.system_info.distribution
    if distro != "flatcar":
        pytest.skip("flatcar uses different PATH for non-interactive")
    # the path is only set on interactive shell. SO lets append it here
    cmd = host.run("bash -c 'PATH=$PATH:/opt/bin type kubectl'")
    assert cmd.succeeded is True
    cmd = host.run("bash -c 'PATH=$PATH:/opt/bin type kubelet'")
    assert cmd.succeeded is True
    cmd = host.run("bash -c 'PATH=$PATH:/opt/bin type kubeadm'")
    assert cmd.succeeded is True
