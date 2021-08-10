package ansible

import (
	random "math/rand"
	"time"

	"github.com/mesosphere/konvoy-image-builder/pkg/constants"
)

type ExtraVars struct {
	WorkingDirectory string                 `json:"working_directory" yaml:"working_directory"`
	ExtraRunnerVars  map[string]interface{} `json:"extraRunnerVars,omitempty" yaml:"extraRunnerVars,omitempty"`
	KeepalivedVrid   int                    `json:"keepalived_vrid" yaml:"keepalived_vrid"`
	RPMsTarFile      string                 `json:"rpms_tar_file" yaml:"rpms_tar_file"`
}

func extraVars(playbookOptions *PlaybookOptions) *ExtraVars {
	keepalivedVrid := random.New(random.NewSource(time.Now().UnixNano())).Intn(255) + 1 // range(1,255)

	extraVars := &ExtraVars{
		WorkingDirectory: constants.WorkingDir,
		ExtraRunnerVars:  playbookOptions.ExtraVarsMap,
		KeepalivedVrid:   keepalivedVrid,
	}

	return extraVars
}
