package app

import (
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

type RootFlags struct {
	Verbosity logging.Verbosity
	RootFlagsCommon
}

func (r *RootFlags) LogDebug() bool {
	return logging.VerbosityToDebug(r.Verbosity)
}

type RootFlagsCommon struct {
	Color bool
}
