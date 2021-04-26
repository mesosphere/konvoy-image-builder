package packer

import "os/exec"

type Interface interface {
	Build(manifest string, flags BuildFlags) (*exec.Cmd, error)
}
