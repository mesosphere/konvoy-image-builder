package app

import (
	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/appansible"
)

type PreflightFlags struct {
	RootFlags
	ExtraVars []string
}

func Preflight(inventory string, flags PreflightFlags) error {
	playbook := appansible.NewPlaybook("preflight", inventory, &ansible.PlaybookOptions{
		ExtraVars: flags.ExtraVars,
	})

	if err := playbook.Run(NewRunOptions(flags.RootFlags)); err != nil {
		return errors.Wrap(err, "error running playbook")
	}
	return nil
}
