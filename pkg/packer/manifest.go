package packer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"

	"github.com/mesosphere/konvoy-image-builder/pkg/packer/manifests"
)

type RenderOptions struct {
	SourceAMIDefined bool
	DryRun           bool
	Bastion          string
	UseBastion       bool
	BastionPassword  string
	BastionUser      string
}

var packerManifests = map[string][]byte{
	"amazon":  manifests.PackerAmazon,
	"vsphere": manifests.PackerVsphere,
}

var ErrManifestNotSupported = errors.New("manifest not support")

func ManifestNotSupportedError(manifest string) error {
	return fmt.Errorf("%w: %s", ErrManifestNotSupported, manifest)
}

func GetManifest(buildType string, options *RenderOptions) ([]byte, error) {
	m, ok := packerManifests[buildType]
	if !ok {
		return nil, ManifestNotSupportedError(buildType)
	}
	return renderPackerJSON(string(m), options)
}

func renderPackerJSON(t string, options *RenderOptions) ([]byte, error) {
	o, err := template.New("packer-json").Delims("((", "))").Parse(t)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse `packer-json` template")
	}

	var out []byte
	buf := bytes.NewBuffer(out)

	if err := o.Execute(buf, options); err != nil {
		return nil, errors.Wrap(err, "error executing template")
	}

	return buf.Bytes(), nil
}
