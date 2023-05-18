package image

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/mesosphere/konvoy-image-builder/pkg/version"
)

func imageLoaded(containerEngine, image string) bool {
	cmd := exec.Command(containerEngine, "image", "inspect", image)
	stdErrBuf := bytes.NewBuffer(make([]byte, 0))
	cmd.Stderr = stdErrBuf
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func Tag() string {
	return fmt.Sprintf("%s:%s", Repository, version.Version())
}
