package appansible

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

type RunOptions struct {
	Out           io.Writer
	ErrOut        io.Writer
	RunsDirectory string
	Verbosity     logging.Verbosity
}

type IOConfig struct {
	RunName string
	RunDir  string

	LoggerConfig *ansible.LoggerConfig
	CloseFiles   func()
}

func NewIOConfig(runName string, runOptions RunOptions) (*IOConfig, error) {
	start := time.Now()
	runDir := filepath.Join(runOptions.RunsDirectory, runName, start.Format("2006-01-02-15-04-05"))
	if err := os.MkdirAll(runDir, 0o750); err != nil {
		return nil, errors.Wrap(err, "error creating directory")
	}

	ansibleLogFilename := filepath.Join(runDir, "ansible.log")
	ansibleLogFile, err := os.Create(ansibleLogFilename)
	if err != nil {
		return nil, fmt.Errorf("error creating ansible log file %q: %w", ansibleLogFilename, err)
	}

	return &IOConfig{
		RunName: runName,
		RunDir:  runDir,
		LoggerConfig: &ansible.LoggerConfig{
			Out:       runOptions.Out,
			ErrOut:    runOptions.ErrOut,
			Log:       ansibleLogFile,
			Verbosity: int(runOptions.Verbosity),
		},
		CloseFiles: func() {
			ansibleLogFile.Close()
		},
	}, nil
}
