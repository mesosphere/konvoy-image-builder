package image

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mesosphere/konvoy-image-builder/pkg/version"
)

func imageLoaded(image string) (bool, error) {
	cmd := exec.Command("docker", "image", "inspect", image)
	stdErrBuf := bytes.NewBuffer(make([]byte, 0))
	cmd.Stderr = stdErrBuf

	if err := cmd.Run(); err != nil {
		stdErr := stdErrBuf.Bytes()
		if bytes.Index(stdErr, []byte("No such image")) > 0 {
			return false, nil
		}
		fmt.Fprintf(os.Stderr, "docker error: %s", stdErr)
		return false, err
	}
	return true, nil
}

func Tag() string {
	return fmt.Sprintf("%s:%s", Repository, version.Version())
}
