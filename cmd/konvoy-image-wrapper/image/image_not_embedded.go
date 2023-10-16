//go:build !EMBED_DOCKER_IMAGE
// +build !EMBED_DOCKER_IMAGE

package image

import (
	"fmt"
	"os/exec"
)

const Repository = "mesosphere/konvoy-image-builder"

func LoadImage(containerEngine string) error {
	image := Tag()
	if imageLoaded(containerEngine, image) {
		return nil
	}

	//nolint:gosec // this is necessary
	cmd := exec.Command(containerEngine, "pull", fmt.Sprintf("docker.io/%s", Tag()))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to pull image %s using cmd %s with error %w", Tag(), cmd.Args, err)
	}
	return nil
}
