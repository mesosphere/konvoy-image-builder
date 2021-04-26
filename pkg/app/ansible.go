package app

import (
	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
)

type ProvisionFlags struct {
	RootFlags

	ExtraVars []string
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
		RunsDirectory: AnsibleRunsDirectory,
		Verbosity:     flags.Verbosity,
	}
}
