package app

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/constants"
)

const (
	extraVarsTemplate = "%s=%s"
)

type ProvisionFlags struct {
	RootFlags
	ClusterArgs

	ExtraVars []string
	Overrides []string
	Provider  string
	WorkDir   string
	Inventory string
}

type ValidateFlags struct {
	RootFlags

	ServiceSubnet       string
	PodSubnet           string
	APIServerEndpoint   string
	CalicoEncapsulation string
	CloudProvider       string
	ErrorsToIgnore      string

	ExtraVars []string
}

func Validate(inventory string, flags ValidateFlags) error {
	playbookOpts := validateFlagsToPlaybookOptions(flags)

	playbook := appansible.NewPlaybook("validate", inventory, playbookOpts)

	if err := playbook.Run(NewRunOptions(flags.RootFlags)); err != nil {
		return errors.Wrap(err, "error running playbook")
	}

	return nil
}

func Provision(inventory string, flags ProvisionFlags) error {
	playbook := appansible.NewPlaybook("provision", inventory, &ansible.PlaybookOptions{
		ExtraVars: flags.ExtraVars,
	})

	if err := playbook.Run(NewRunOptions(flags.RootFlags)); err != nil {
		return errors.Wrap(err, "error running playbook")
	}

	return nil
}

func NewRunOptions(flags RootFlags) appansible.RunOptions {
	return appansible.RunOptions{
		Out:           out,
		ErrOut:        errOut,
		RunsDirectory: constants.AnsibleRunsDirectory,
		Verbosity:     flags.Verbosity,
	}
}

func validateFlagsToPlaybookOptions(flags ValidateFlags) *ansible.PlaybookOptions {
	playbookOptions := &ansible.PlaybookOptions{}
	values := []string{
		fmt.Sprintf(extraVarsTemplate, "service_subnet", flags.ServiceSubnet),
		fmt.Sprintf(extraVarsTemplate, "pod_subnet", flags.PodSubnet),
		fmt.Sprintf(extraVarsTemplate, "calico_encapsulation", flags.CalicoEncapsulation),
		fmt.Sprintf(extraVarsTemplate, "cloud_provider", flags.CloudProvider),
		fmt.Sprintf(extraVarsTemplate, "apiserver_endpoint", flags.APIServerEndpoint),
		fmt.Sprintf(extraVarsTemplate, "errors_to_ignore", flags.ErrorsToIgnore),
	}

	if playbookOptions.ExtraVars == nil {
		playbookOptions.ExtraVars = make([]string, 0)
	}

	playbookOptions.ExtraVars = append(playbookOptions.ExtraVars, values...)

	return playbookOptions
}
