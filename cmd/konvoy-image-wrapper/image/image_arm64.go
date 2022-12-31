//go:build EMBED_DOCKER_IMAGE_arm64
// +build EMBED_DOCKER_IMAGE_arm64

package image

//nolint:golint // blank import for embed support.
import (
	"bytes"
	_ "embed"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

const Repository = "mesosphere/konvoy-image-builder"

//go:embed konvoy-image-builder_linux_arm64.tar.gz
var konvoyImageTar []byte // memory is cheap, right?

func LoadImage() error {
	image := Tag()
	present, err := imageLoaded(image)
	if err != nil {
		return errors.Wrap(err, "error querying docker for images")
	}

	if !present {
		cmd := exec.Command("docker", "image", "load")
		cmd.Stdin = bytes.NewReader(konvoyImageTar)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
	return nil
}
