package packer

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/packer/manifests"
)

var packerManifests = map[string][]byte{
	"amazon":        manifests.PackerAmazon,
	"azure":         manifests.PackerAzure,
	"vsphere":       manifests.PackerVsphere,
	"googlecompute": manifests.PackerGCP,
}

var ErrManifestNotSupported = errors.New("manifest not support")

func ManifestNotSupportedError(manifest string) error {
	return fmt.Errorf("%w: %s", ErrManifestNotSupported, manifest)
}

func GetManifest(buildType string) ([]byte, error) {
	m, ok := packerManifests[buildType]
	if !ok {
		return nil, ManifestNotSupportedError(buildType)
	}
	return m, nil
}
