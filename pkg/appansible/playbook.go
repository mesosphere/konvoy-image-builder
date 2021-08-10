package appansible

import (
	"path"

	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/constants"
)

type Playbook struct {
	Name      string
	Inventory string

	PlaybookOptions *ansible.PlaybookOptions
}

func NewPlaybook(name, inventory string, playbookOptions *ansible.PlaybookOptions) *Playbook {
	return &Playbook{
		Name:            name,
		Inventory:       inventory,
		PlaybookOptions: playbookOptions,
	}
}

func (p *Playbook) Run(runOptions RunOptions) error {
	io, err := NewIOConfig(p.Name, runOptions)
	if err != nil {
		return err
	}
	defer io.CloseFiles()

	runner := ansible.NewRunner(io.RunDir, io.LoggerConfig)
	if err := runner.StartPlaybook(p.Filename(), p.Inventory, p.PlaybookOptions); err != nil {
		return errors.Wrap(err, "error starting playbook")
	}

	if err := runner.WaitPlaybook(); err != nil {
		return errors.Wrap(err, "error waiting for playbook")
	}

	return nil
}

func (p *Playbook) Filename() string {
	return path.Join(constants.AnsiblePlaybookPath, p.Name+".yaml")
}
