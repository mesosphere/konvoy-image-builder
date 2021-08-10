package ansible

import (
	"io/ioutil"

	"github.com/mesosphere/konvoy/clientapis/pkg/apis/konvoy"
	clientapisconstants "github.com/mesosphere/konvoy/clientapis/pkg/constants"
	"gopkg.in/yaml.v2"
)

// FromFile returns the inventory defined in the file.
func FromFile(file string) (*konvoy.Inventory, error) {
	inventoryFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	inventory := &konvoy.Inventory{}
	if err := yaml.Unmarshal(inventoryFile, inventory); err != nil {
		return nil, err
	}

	// Annotate hosts with their node pool type and set the default node_pool
	for k := range inventory.ControlPlane.Hosts {
		h := inventory.ControlPlane.Hosts[k]
		h.ControlPlane = true
		setNodePoolIfEmpty(&h, clientapisconstants.DefaultControlPlanPoolName)
		inventory.ControlPlane.Hosts[k] = h
	}
	for k := range inventory.Bastion.Hosts {
		h := inventory.Bastion.Hosts[k]
		h.Bastion = true
		inventory.Bastion.Hosts[k] = h
	}

	// Add default node_pool to node hosts
	for k := range inventory.Node.Hosts {
		h := inventory.Node.Hosts[k]
		setNodePoolIfEmpty(&h, clientapisconstants.DefaultNodePoolName)
		inventory.Node.Hosts[k] = h
	}

	return inventory, nil
}

// RewriteWithDefaults will read the src file, apply defaults and write it out to dst.
func RewriteWithDefaults(src string, dst string) error {
	inventory, err := FromFile(src)
	if err != nil {
		return err
	}
	return ToFile(inventory, dst)
}

// ToFile writes the inventory to the file.
func ToFile(i *konvoy.Inventory, file string) error {
	data, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0600)
}

// setNodePoolIfEmpty will set the NodePool to the default value if empty.
func setNodePoolIfEmpty(host *konvoy.InventoryHost, nodePool string) {
	// hosts that are also control-plane hosts should not use the default label
	if host.NodePool == "" {
		host.NodePool = nodePool
	}
}
