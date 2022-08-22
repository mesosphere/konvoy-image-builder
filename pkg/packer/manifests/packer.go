package manifests

//nolint:golint // blank import for embed support.
import _ "embed"

// PackerAmazon storage for AWS packer template.
//
//go:embed aws/packer.json.tmpl
var PackerAmazon []byte

// PackerAzure storage for Azurepacker template.
//
//go:embed azure/packer.json.tmpl
var PackerAzure []byte

// PackerGCP storage for GCP packer base template.
//
//go:embed gcp/packer.json.tmpl
var PackerGCP []byte

// PackerVsphere storage for vSphere packer base template.
//
//go:embed vsphere/packer.json.tmpl
var PackerVsphere []byte
