# flake8: noqa
# pylint: skip-file
# Cloud-Init Datasource for VMware Guestinfo

"""
A cloud init datasource for VMware GuestInfo.
"""

import base64
import collections
import copy
import ipaddress
import json
import os
import socket
import time
import zlib
from distutils.spawn import find_executable
from distutils.util import strtobool

import netifaces
from cloudinit import log as logging
from cloudinit import safeyaml, sources, util

# from cloud-init >= 20.3 subp is in its own module
try:
    from cloudinit.subp import ProcessExecutionError, subp
except ImportError:
    from cloudinit.util import ProcessExecutionError, subp

LOG = logging.getLogger(__name__)
NOVAL = "No value found"
VMWARE_RPCTOOL = find_executable("vmware-rpctool")
VMX_GUESTINFO = "VMX_GUESTINFO"
GUESTINFO_EMPTY_YAML_VAL = "---"
LOCAL_IPV4 = "local-ipv4"
LOCAL_IPV6 = "local-ipv6"
CLEANUP_GUESTINFO = "cleanup-guestinfo"
WAIT_ON_NETWORK = "wait-on-network"
WAIT_ON_NETWORK_IPV4 = "ipv4"
WAIT_ON_NETWORK_IPV6 = "ipv6"


class NetworkConfigError(Exception):
    """
    NetworkConfigError is raised when there is an issue getting or
    applying network configuration.
    """

    pass


class DataSourceVMwareGuestInfo(sources.DataSource):
    """
    This cloud-init datasource was designed for use with CentOS 7,
    which uses cloud-init 0.7.9. However, this datasource should
    work with any Linux distribution for which cloud-init is
    avaialble.

    The documentation for cloud-init 0.7.9"s datasource is
    available at http://bit.ly/cloudinit-datasource-0-7-9. The
    current documentation for cloud-init is found at
    https://cloudinit.readthedocs.io/en/latest/.

    Setting the hostname:
        The hostname is set by way of the metadata key "local-hostname".

    Setting the instance ID:
        The instance ID may be set by way of the metadata key "instance-id".
        However, if this value is absent then then the instance ID is
        read from the file /sys/class/dmi/id/product_uuid.

    Configuring the network:
        The network is configured by setting the metadata key "network"
        with a value consistent with Network Config Versions 1 or 2,
        depending on the Linux distro"s version of cloud-init:

            Network Config Version 1 - http://bit.ly/cloudinit-net-conf-v1
            Network Config Version 2 - http://bit.ly/cloudinit-net-conf-v2

        For example, CentOS 7"s official cloud-init package is version
        0.7.9 and does not support Network Config Version 2. However,
        this datasource still supports supplying Network Config Version 2
        data as long as the Linux distro"s cloud-init package is new
        enough to parse the data.

        The metadata key "network.encoding" may be used to indicate the
        format of the metadata key "network". Valid encodings are base64
        and gzip+base64.
    """

    dsname = "VMwareGuestInfo"

    def __init__(self, sys_cfg, distro, paths, ud_proc=None):
        sources.DataSource.__init__(self, sys_cfg, distro, paths, ud_proc)
        if not get_data_access_method():
            LOG.error("Failed to find vmware-rpctool")

    def get_data(self):
        """
        This method should really be _get_data in accordance with the most
        recent versions of cloud-init. However, because the datasource
        supports as far back as cloud-init 0.7.9, get_data is still used.

        Because of this the method attempts to do some of the same things
        that the get_data functions in newer versions of cloud-init do,
        such as calling persist_instance_data.
        """
        data_access_method = get_data_access_method()
        if not data_access_method:
            LOG.error("vmware-rpctool is required to fetch guestinfo value")
            return False

        # Get the metadata.
        self.metadata = load_metadata()

        # Get the user data.
        self.userdata_raw = guestinfo("userdata")

        # Get the vendor data.
        self.vendordata_raw = guestinfo("vendordata")

        # Check to see if any of the guestinfo data should be removed.
        if data_access_method == VMWARE_RPCTOOL and CLEANUP_GUESTINFO in self.metadata:
            clear_guestinfo_keys(self.metadata[CLEANUP_GUESTINFO])

        if self.metadata or self.userdata_raw or self.vendordata_raw:
            return True
        else:
            return False

    def setup(self, is_new_instance):
        """setup(is_new_instance)

        This is called before user-data and vendor-data have been processed.

        Unless the datasource has set mode to "local", then networking
        per "fallback" or per "network_config" will have been written and
        brought up the OS at this point.
        """

        host_info = wait_on_network(self.metadata)
        LOG.info("got host-info: %s", host_info)

        # Reflect any possible local IPv4 or IPv6 addresses in the guest
        # info.
        advertise_local_ip_addrs(host_info)

        # Ensure the metadata gets updated with information about the
        # host, including the network interfaces, default IP addresses,
        # etc.
        self.metadata = merge_dicts(self.metadata, host_info)

        # Persist the instance data for versions of cloud-init that support
        # doing so. This occurs here rather than in the get_data call in
        # order to ensure that the network interfaces are up and can be
        # persisted with the metadata.
        try:
            self.persist_instance_data()
        except AttributeError:
            pass

    @property
    def network_config(self):
        if "network" in self.metadata:
            LOG.debug("using metadata network config")
        else:
            LOG.debug("using fallback network config")
            self.metadata["network"] = {
                "config": self.distro.generate_fallback_config(),
            }
        return self.metadata["network"]["config"]

    def get_instance_id(self):
        # Pull the instance ID out of the metadata if present. Otherwise
        # read the file /sys/class/dmi/id/product_uuid for the instance ID.
        if self.metadata and "instance-id" in self.metadata:
            return self.metadata["instance-id"]
        with open("/sys/class/dmi/id/product_uuid", "r") as id_file:
            self.metadata["instance-id"] = str(id_file.read()).rstrip().lower()
            return self.metadata["instance-id"]

    def get_public_ssh_keys(self):
        public_keys_data = ""
        if "public-keys-data" in self.metadata:
            public_keys_data = self.metadata["public-keys-data"].splitlines()

        public_keys = []
        if not public_keys_data:
            return public_keys

        for public_key in public_keys_data:
            public_keys.append(public_key)

        return public_keys


def decode(key, enc_type, data):
    """
    decode returns the decoded string value of data
    key is a string used to identify the data being decoded in log messages
    ----
    In py 2.7:
    json.loads method takes string as input
    zlib.decompress takes and returns a string
    base64.b64decode takes and returns a string
    -----
    In py 3.6 and newer:
    json.loads method takes bytes or string as input
    zlib.decompress takes and returns a bytes
    base64.b64decode takes bytes or string and returns bytes
    -----
    In py > 3, < 3.6:
    json.loads method takes string as input
    zlib.decompress takes and returns a bytes
    base64.b64decode takes bytes or string and returns bytes
    -----
    Given the above conditions the output from zlib.decompress and
    base64.b64decode would be bytes with newer python and str in older
    version. Thus we would covert the output to str before returning
    """
    LOG.debug("Getting encoded data for key=%s, enc=%s", key, enc_type)

    raw_data = None
    if enc_type == "gzip+base64" or enc_type == "gz+b64":
        LOG.debug("Decoding %s format %s", enc_type, key)
        raw_data = zlib.decompress(base64.b64decode(data), zlib.MAX_WBITS | 16)
    elif enc_type == "base64" or enc_type == "b64":
        LOG.debug("Decoding %s format %s", enc_type, key)
        raw_data = base64.b64decode(data)
    else:
        LOG.debug("Plain-text data %s", key)
        raw_data = data

    if isinstance(raw_data, bytes):
        return raw_data.decode("utf-8")
    return raw_data


def get_none_if_empty_val(val):
    """
    get_none_if_empty_val returns None if the provided value, once stripped
    of its trailing whitespace, is empty or equal to GUESTINFO_EMPTY_YAML_VAL.

    The return value is always a string, regardless of whether the input is
    a bytes class or a string.
    """

    # If the provided value is a bytes class, convert it to a string to
    # simplify the rest of this function"s logic.
    if isinstance(val, bytes):
        val = val.decode()

    val = val.rstrip()
    if len(val) == 0 or val == GUESTINFO_EMPTY_YAML_VAL:
        return None
    return val


def advertise_local_ip_addrs(host_info):
    """
    advertise_local_ip_addrs gets the local IP address information from
    the provided host_info map and sets the addresses in the guestinfo
    namespace
    """
    if not host_info:
        return

    # Reflect any possible local IPv4 or IPv6 addresses in the guest
    # info.
    local_ipv4 = host_info.get(LOCAL_IPV4)
    if local_ipv4:
        set_guestinfo_value(LOCAL_IPV4, local_ipv4)
        LOG.info("advertised local ipv4 address %s in guestinfo", local_ipv4)

    local_ipv6 = host_info.get(LOCAL_IPV6)
    if local_ipv6:
        set_guestinfo_value(LOCAL_IPV6, local_ipv6)
        LOG.info("advertised local ipv6 address %s in guestinfo", local_ipv6)


def handle_returned_guestinfo_val(key, val):
    """
    handle_returned_guestinfo_val returns the provided value if it is
    not empty or set to GUESTINFO_EMPTY_YAML_VAL, otherwise None is
    returned
    """
    val = get_none_if_empty_val(val)
    if val:
        return val
    LOG.debug("No value found for key %s", key)
    return None


def get_guestinfo_value(key):
    """
    Returns a guestinfo value for the specified key.
    """
    LOG.debug("Getting guestinfo value for key %s", key)

    data_access_method = get_data_access_method()

    if data_access_method == VMX_GUESTINFO:
        env_key = ("vmx.guestinfo." + key).upper().replace(".", "_", -1)
        return handle_returned_guestinfo_val(key, os.environ.get(env_key, ""))

    if data_access_method == VMWARE_RPCTOOL:
        try:
            (stdout, stderr) = subp([VMWARE_RPCTOOL, "info-get guestinfo." + key])
            if stderr == NOVAL:
                LOG.debug("No value found for key %s", key)
            elif not stdout:
                LOG.error("Failed to get guestinfo value for key %s", key)
            else:
                return handle_returned_guestinfo_val(key, stdout)
        except ProcessExecutionError as error:
            if error.stderr == NOVAL:
                LOG.debug("No value found for key %s", key)
            else:
                util.logexc(
                    LOG, "Failed to get guestinfo value for key %s: %s", key, error
                )
        except Exception:
            util.logexc(
                LOG,
                "Unexpected error while trying to get guestinfo value for key %s",
                key,
            )

    return None


def set_guestinfo_value(key, value):
    """
    Sets a guestinfo value for the specified key. Set value to an empty string
    to clear an existing guestinfo key.
    """

    # If value is an empty string then set it to a single space as it is not
    # possible to set a guestinfo key to an empty string. Setting a guestinfo
    # key to a single space is as close as it gets to clearing an existing
    # guestinfo key.
    if value == "":
        value = " "

    LOG.debug("Setting guestinfo key=%s to value=%s", key, value)

    data_access_method = get_data_access_method()

    if data_access_method == VMX_GUESTINFO:
        return True

    if data_access_method == VMWARE_RPCTOOL:
        try:
            subp([VMWARE_RPCTOOL, ("info-set guestinfo.%s %s" % (key, value))])
            return True
        except ProcessExecutionError as error:
            util.logexc(
                LOG, "Failed to set guestinfo key=%s to value=%s: %s", key, value, error
            )
        except Exception:
            util.logexc(
                LOG,
                "Unexpected error while trying to set guestinfo key=%s to value=%s",
                key,
                value,
            )

    return None


def clear_guestinfo_keys(keys):
    """
    clear_guestinfo_keys clears guestinfo of all of the keys in the given list.
    each key will have its value set to "---". Since the value is valid YAML,
    cloud-init can still read it if it tries.
    """
    if not keys:
        return
    if not type(keys) in (list, tuple):
        keys = [keys]
    for key in keys:
        LOG.info("clearing guestinfo.%s", key)
        if not set_guestinfo_value(key, GUESTINFO_EMPTY_YAML_VAL):
            LOG.error("failed to clear guestinfo.%s", key)
        LOG.info("clearing guestinfo.%s.encoding", key)
        if not set_guestinfo_value(key + ".encoding", ""):
            LOG.error("failed to clear guestinfo.%s.encoding", key)


def guestinfo(key):
    """
    guestinfo returns the guestinfo value for the provided key, decoding
    the value when required
    """
    data = get_guestinfo_value(key)
    if not data:
        return None
    enc_type = get_guestinfo_value(key + ".encoding")
    return decode("guestinfo." + key, enc_type, data)


def load(data):
    """
    load first attempts to unmarshal the provided data as JSON, and if
    that fails then attempts to unmarshal the data as YAML. If data is
    None then a new dictionary is returned.
    """
    if not data:
        return {}
    try:
        return json.loads(data)
    except:
        return safeyaml.load(data)


def load_metadata():
    """
    load_metadata loads the metadata from the guestinfo data, optionally
    decoding the network config when required
    """
    data = load(guestinfo("metadata"))
    LOG.debug("loaded metadata %s", data)

    network = None
    if "network" in data:
        network = data["network"]
        del data["network"]

    network_enc = None
    if "network.encoding" in data:
        network_enc = data["network.encoding"]
        del data["network.encoding"]

    if network:
        LOG.debug("network data found")
        if isinstance(network, collections.Mapping):
            LOG.debug("network data copied to 'config' key")
            network = {"config": copy.deepcopy(network)}
        else:
            LOG.debug("network data to be decoded %s", network)
            dec_net = decode("metadata.network", network_enc, network)
            network = {
                "config": load(dec_net),
            }

        LOG.debug("network data %s", network)
        data["network"] = network

    return data


def get_datasource_list(depends):
    """
    Return a list of data sources that match this set of dependencies
    """
    return [DataSourceVMwareGuestInfo]


def get_default_ip_addrs():
    """
    Returns the default IPv4 and IPv6 addresses based on the device(s) used for
    the default route. Please note that None may be returned for either address
    family if that family has no default route or if there are multiple
    addresses associated with the device used by the default route for a given
    address.
    """
    gateways = netifaces.gateways()
    if "default" not in gateways:
        return None, None

    default_gw = gateways["default"]
    if netifaces.AF_INET not in default_gw and netifaces.AF_INET6 not in default_gw:
        return None, None

    ipv4 = None
    ipv6 = None

    gw4 = default_gw.get(netifaces.AF_INET)
    if gw4:
        _, dev4 = gw4
        addr4_fams = netifaces.ifaddresses(dev4)
        if addr4_fams:
            af_inet4 = addr4_fams.get(netifaces.AF_INET)
            if af_inet4:
                if len(af_inet4) > 1:
                    LOG.warn(
                        "device %s has more than one ipv4 address: %s", dev4, af_inet4
                    )
                elif "addr" in af_inet4[0]:
                    ipv4 = af_inet4[0]["addr"]

    # Try to get the default IPv6 address by first seeing if there is a default
    # IPv6 route.
    gw6 = default_gw.get(netifaces.AF_INET6)
    if gw6:
        _, dev6 = gw6
        addr6_fams = netifaces.ifaddresses(dev6)
        if addr6_fams:
            af_inet6 = addr6_fams.get(netifaces.AF_INET6)
            if af_inet6:
                if len(af_inet6) > 1:
                    LOG.warn(
                        "device %s has more than one ipv6 address: %s", dev6, af_inet6
                    )
                elif "addr" in af_inet6[0]:
                    ipv6 = af_inet6[0]["addr"]

    # If there is a default IPv4 address but not IPv6, then see if there is a
    # single IPv6 address associated with the same device associated with the
    # default IPv4 address.
    if ipv4 and not ipv6:
        af_inet6 = addr4_fams.get(netifaces.AF_INET6)
        if af_inet6:
            if len(af_inet6) > 1:
                LOG.warn("device %s has more than one ipv6 address: %s", dev4, af_inet6)
            elif "addr" in af_inet6[0]:
                ipv6 = af_inet6[0]["addr"]

    # If there is a default IPv6 address but not IPv4, then see if there is a
    # single IPv4 address associated with the same device associated with the
    # default IPv6 address.
    if not ipv4 and ipv6:
        af_inet4 = addr6_fams.get(netifaces.AF_INET)
        if af_inet4:
            if len(af_inet4) > 1:
                LOG.warn("device %s has more than one ipv4 address: %s", dev6, af_inet4)
            elif "addr" in af_inet4[0]:
                ipv4 = af_inet4[0]["addr"]

    return ipv4, ipv6


# patched socket.getfqdn() - see https://bugs.python.org/issue5004


def getfqdn(name=""):
    """Get fully qualified domain name from name.
    An empty argument is interpreted as meaning the local host.
    """
    name = name.strip()
    if not name or name == "0.0.0.0":
        name = socket.gethostname()
    try:
        addrs = socket.getaddrinfo(
            name, None, 0, socket.SOCK_DGRAM, 0, socket.AI_CANONNAME
        )
    except socket.error:
        pass
    else:
        for addr in addrs:
            if addr[3]:
                name = addr[3]
                break
    return name


def is_valid_ip_addr(val):
    """
    Returns false if the address is loopback, link local or unspecified;
    otherwise true is returned.
    """
    addr = None
    try:
        try:
            addr = ipaddress.ip_address(val)
        except ipaddress.AddressValueError:
            addr = ipaddress.ip_address(unicode(val))
    except:
        return False
    if addr.is_link_local or addr.is_loopback or addr.is_unspecified:
        return False
    return True


def get_host_info():
    """
    Returns host information such as the host name and network interfaces.
    """

    host_info = {
        "network": {
            "interfaces": {
                "by-mac": collections.OrderedDict(),
                "by-ipv4": collections.OrderedDict(),
                "by-ipv6": collections.OrderedDict(),
            },
        },
    }

    hostname = getfqdn(socket.gethostname())
    if hostname:
        host_info["hostname"] = hostname
        host_info["local-hostname"] = hostname
        host_info["local_hostname"] = hostname

    default_ipv4, default_ipv6 = get_default_ip_addrs()
    if default_ipv4:
        host_info[LOCAL_IPV4] = default_ipv4
    if default_ipv6:
        host_info[LOCAL_IPV6] = default_ipv6

    by_mac = host_info["network"]["interfaces"]["by-mac"]
    by_ipv4 = host_info["network"]["interfaces"]["by-ipv4"]
    by_ipv6 = host_info["network"]["interfaces"]["by-ipv6"]

    ifaces = netifaces.interfaces()
    for dev_name in ifaces:
        addr_fams = netifaces.ifaddresses(dev_name)
        af_link = addr_fams.get(netifaces.AF_LINK)
        af_inet4 = addr_fams.get(netifaces.AF_INET)
        af_inet6 = addr_fams.get(netifaces.AF_INET6)

        mac = None
        if af_link and "addr" in af_link[0]:
            mac = af_link[0]["addr"]

        # Do not bother recording localhost
        if mac == "00:00:00:00:00:00":
            continue

        if mac and (af_inet4 or af_inet6):
            key = mac
            val = {}
            if af_inet4:
                af_inet4_vals = []
                for ip_info in af_inet4:
                    if not is_valid_ip_addr(ip_info["addr"]):
                        continue
                    af_inet4_vals.append(ip_info)
                val["ipv4"] = af_inet4_vals
            if af_inet6:
                af_inet6_vals = []
                for ip_info in af_inet6:
                    if not is_valid_ip_addr(ip_info["addr"]):
                        continue
                    af_inet6_vals.append(ip_info)
                val["ipv6"] = af_inet6_vals
            by_mac[key] = val

        if af_inet4:
            for ip_info in af_inet4:
                key = ip_info["addr"]
                if not is_valid_ip_addr(key):
                    continue
                val = copy.deepcopy(ip_info)
                del val["addr"]
                if mac:
                    val["mac"] = mac
                by_ipv4[key] = val

        if af_inet6:
            for ip_info in af_inet6:
                key = ip_info["addr"]
                if not is_valid_ip_addr(key):
                    continue
                val = copy.deepcopy(ip_info)
                del val["addr"]
                if mac:
                    val["mac"] = mac
                by_ipv6[key] = val

    return host_info


def wait_on_network(metadata):
    # Determine whether we need to wait on the network coming online.
    wait_on_ipv4 = False
    wait_on_ipv6 = False
    if WAIT_ON_NETWORK in metadata:
        wait_on_network = metadata[WAIT_ON_NETWORK]
        if WAIT_ON_NETWORK_IPV4 in wait_on_network:
            wait_on_ipv4_val = wait_on_network[WAIT_ON_NETWORK_IPV4]
            if isinstance(wait_on_ipv4_val, bool):
                wait_on_ipv4 = wait_on_ipv4_val
            else:
                wait_on_ipv4 = bool(strtobool(wait_on_ipv4_val))
        if WAIT_ON_NETWORK_IPV6 in wait_on_network:
            wait_on_ipv6_val = wait_on_network[WAIT_ON_NETWORK_IPV6]
            if isinstance(wait_on_ipv6_val, bool):
                wait_on_ipv6 = wait_on_ipv6_val
            else:
                wait_on_ipv6 = bool(strtobool(wait_on_ipv6_val))

    # Get information about the host.
    host_info = None
    while host_info is None:
        host_info = get_host_info()
        if wait_on_ipv4:
            ipv4_ready = False
            if "network" in host_info:
                if "interfaces" in host_info["network"]:
                    if "by-ipv4" in host_info["network"]["interfaces"]:
                        if len(host_info["network"]["interfaces"]["by-ipv4"]) > 0:
                            ipv4_ready = True
            if not ipv4_ready:
                LOG.info("ipv4 not ready")
                host_info = None
        if wait_on_ipv6:
            ipv6_ready = False
            if "network" in host_info:
                if "interfaces" in host_info["network"]:
                    if "by-ipv6" in host_info["network"]["interfaces"]:
                        if len(host_info["network"]["interfaces"]["by-ipv6"]) > 0:
                            ipv6_ready = True
            if not ipv6_ready:
                LOG.info("ipv6 not ready")
                host_info = None
        if host_info is None:
            LOG.info("waiting on network")
            time.sleep(1)

    return host_info


def get_data_access_method():
    if os.environ.get(VMX_GUESTINFO, ""):
        return VMX_GUESTINFO
    if VMWARE_RPCTOOL:
        return VMWARE_RPCTOOL
    return None


_MERGE_STRATEGY_ENV_VAR = "CLOUD_INIT_VMWARE_GUEST_INFO_MERGE_STRATEGY"
_MERGE_STRATEGY_DEEPMERGE = "deepmerge"


def merge_dicts(a, b):
    merge_strategy = os.getenv(_MERGE_STRATEGY_ENV_VAR)
    if merge_strategy == _MERGE_STRATEGY_DEEPMERGE:
        try:
            LOG.info("merging dictionaries with deepmerge strategy")
            return merge_dicts_with_deep_merge(a, b)
        except Exception as err:
            LOG.error("deep merge failed: %s" % err)
    LOG.info("merging dictionaries with stdlib strategy")
    return merge_dicts_with_stdlib(a, b)


def merge_dicts_with_deep_merge(a, b):
    from deepmerge import always_merger

    return always_merger.merge(a, b)


def merge_dicts_with_stdlib(a, b):
    for key, value in a.items():
        if isinstance(value, dict):
            node = b.setdefault(key, {})
            merge_dicts_with_stdlib(value, node)
        else:
            b[key] = value
    return b


def main():
    """
    Executed when this file is used as a program.
    """
    try:
        logging.setupBasicLogging()
    except Exception:
        pass
    metadata = {
        "wait-on-network": {"ipv4": True, "ipv6": "false"},
        "network": {"config": {"dhcp": True}},
    }
    host_info = wait_on_network(metadata)
    metadata = merge_dicts(metadata, host_info)
    print(util.json_dumps(metadata))


if __name__ == "__main__":
    main()
