package packer

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

const DefaultPath string = "packer"

var supportedOnErrorActions = []string{
	"cleanup",
	"abort",
	"run-cleanup-provisioner",
}

type CLIRunner struct {
	Path string

	Out    io.Writer
	OutErr io.Writer
}

type BuildFlags struct {
	Force            bool
	Debug            bool
	Color            bool
	OnError          string
	VarFiles         []string
	Vars             map[string]string
	Bastion         string
	BastionPassword string
}

var ErrCLI = errors.New("error running packer")

func CLIError(op string) error {
	return fmt.Errorf("%w: %s", ErrCLI, op)
}

func (r *CLIRunner) Build(manifest string, flags BuildFlags) (*exec.Cmd, error) {
	a := []string{"build"}

	if flags.Force {
		a = append(a, "-force")
	}
	if flags.Debug {
		a = append(a, "-debug")
	}
	if !flags.Color {
		a = append(a, "-color=false")
	}

	if flags.OnError != "" {
		found := false
		for _, action := range supportedOnErrorActions {
			if action == flags.OnError {
				found = true
				break
			}
		}
		if !found {
			return nil, CLIError(
				fmt.Sprintf("packer clean up action is not valid, must be one of: %s", supportedOnErrorActions))
		}
		a = append(a, fmt.Sprintf("-on-error=%s", flags.OnError))
	}

	for _, f := range flags.VarFiles {
		a = append(a, fmt.Sprintf("-var-file=%s", f))
	}

	for k, v := range flags.Vars {
		a = append(a, "-var", fmt.Sprintf("'%s=%s'", k, v))
	}

	a = append(a, manifest)
	// TODO: log(debug): command
	return r.run(a...)
}

//nolint:gosec // private function, should not be abused
func (r *CLIRunner) run(args ...string) (*exec.Cmd, error) {
	c := exec.Command(r.Path, args...)
	c.Stdout = r.Out
	c.Stderr = r.OutErr
	err := c.Run()
	if err != nil {
		return c, fmt.Errorf("error running command: %w", err)
	}
	return c, nil
}
