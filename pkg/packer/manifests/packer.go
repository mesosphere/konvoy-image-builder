package manifests

//nolint:golint // blank import for embed support.
import _ "embed"

// PackerAmazon storage for AWS packer template.
//
//go:embed aws/packer.pkr.hcl.tmpl
var PackerAmazon []byte

// PackerAzure storage for Azurepacker template.
//
//go:embed azure/packer.pkr.hcl
var PackerAzure []byte

// PackerGCP storage for GCP packer base template.
//
//go:embed gcp/packer.pkr.hcl
var PackerGCP []byte

// PackerVsphere storage for vSphere packer base template.
//
//go:embed vsphere/packer.pkr.hcl
var PackerVsphere []byte
