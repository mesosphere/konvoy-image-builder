package app

import (
	"fmt"

	"github.com/pkg/errors"
)

// errors.
var (
	ErrInitConfig               = errors.New("init configuration failure")
	ErrBuild                    = errors.New("build failure")
	ErrKeyNotString             = errors.New("key is not a string")
	ErrKubernetesVersionMissing = errors.New("necessary kubernetes_version key missing")
	ErrPathNotString            = errors.New("path value is not a string")
	ErrPathNotSlice             = errors.New("path value is not a slice")
	ErrConfigRequired           = errors.New("config value is required")
)

func InitConfigError(op string, err error) error {
	if err != nil {
		return fmt.Errorf("%w: %s: %v", ErrInitConfig, op, err)
	}

	return fmt.Errorf("%w: %s", ErrInitConfig, op)
}

func BuildError(op string) error {
	return fmt.Errorf("%w: %s", ErrBuild, op)
}
