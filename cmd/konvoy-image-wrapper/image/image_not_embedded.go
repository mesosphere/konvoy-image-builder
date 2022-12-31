//go:build !EMBED_DOCKER_IMAGE_amd64 && !EMBED_DOCKER_IMAGE_arm64
// +build !EMBED_DOCKER_IMAGE_amd64,!EMBED_DOCKER_IMAGE_arm64

package image

import (
	"os/exec"

	"github.com/pkg/errors"
)

const Repository = "mesosphere/konvoy-image-builder"

func LoadImage() error {
	image := Tag()
	found, err := imageLoaded(image)
	if err != nil {
		return errors.Wrap(err, "error querying docker for images")
	}
	if found {
		return nil
	}
	//nolint:gosec // this is necessary
	cmd := exec.Command("docker", "pull", Tag())
	return cmd.Run()
}
