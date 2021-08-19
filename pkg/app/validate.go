package app

import (
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
)

const (
	extraVarsTemplate = "%s=%s"
)

type ValidateFlags struct {
	RootFlags

	Inventory     string
	ServiceSubnet string
	PodSubnet     string
	APIServerPort int
}

func Validate(flags ValidateFlags) error {
	playbookOpts := validateFlagsToPlaybookOptions(flags)

	if flags.Inventory == "" {
		return fmt.Errorf("inventory file cannot be empty")
	}

	if !isValidCIDR(flags.PodSubnet) {
		return fmt.Errorf("pod-subnet %q is not a valid CIDR", flags.PodSubnet)
	}

	if !isValidCIDR(flags.ServiceSubnet) {
		return fmt.Errorf("service-subnet %q is not a valid CIDR", flags.ServiceSubnet)
	}

	playbook := appansible.NewPlaybook("validate", flags.Inventory, playbookOpts)

	if err := playbook.Run(NewRunOptions(flags.RootFlags)); err != nil {
		return errors.Wrap(err, "error running playbook")
	}

	return nil
}

func validateFlagsToPlaybookOptions(flags ValidateFlags) *ansible.PlaybookOptions {
	playbookOptions := &ansible.PlaybookOptions{}
	values := []string{
		fmt.Sprintf(extraVarsTemplate, "service_subnet", flags.ServiceSubnet),
		fmt.Sprintf(extraVarsTemplate, "pod_subnet", flags.PodSubnet),
		fmt.Sprintf("%s=%d", "apiserver_port", flags.APIServerPort),
	}

	if playbookOptions.ExtraVars == nil {
		playbookOptions.ExtraVars = make([]string, 0)
	}

	playbookOptions.ExtraVars = append(playbookOptions.ExtraVars, values...)

	return playbookOptions
}

func isValidCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)

	return err == nil
}
