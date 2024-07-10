//go:build EMBED_DOCKER_IMAGE
// +build EMBED_DOCKER_IMAGE

package image

//nolint:golint // blank import for embed support.
import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
)

const Repository = "mesosphere/konvoy-image-builder"

//go:embed konvoy-image-builder.tar.gz
var konvoyImageTar []byte // memory is cheap, right?

func LoadImage(containerEngine string) error {
	image := Tag()
	if imageLoaded(containerEngine, image) {
		return nil
	}
	cmd := exec.Command(containerEngine, "image", "load")
	cmd.Stdin = bytes.NewReader(konvoyImageTar)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run cmd %s with error %w", cmd.Args, err)
	}
	return nil
}
