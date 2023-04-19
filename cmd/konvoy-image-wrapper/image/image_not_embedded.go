//go:build !EMBED_DOCKER_IMAGE
// +build !EMBED_DOCKER_IMAGE

package image

import (
	"os/exec"
)

const Repository = "mesosphere/konvoy-image-builder"

func LoadImage(containerEngine string) error {
	image := Tag()
	if imageLoaded(containerEngine, image) {
		return nil
	}
	//nolint:gosec // this is necessary
	cmd := exec.Command(containerEngine, "pull", Tag())
	return cmd.Run()
}
