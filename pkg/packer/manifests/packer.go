package manifests

//nolint:golint // blank import for embed support.
import _ "embed"

//go:embed aws/packer.json.tmpl
//PackerAmazon storage for AWS packer template.
var PackerAmazon []byte

//go:embed azure/packer.json.tmpl
//PackerAzure storage for Azurepacker template.
var PackerAzure []byte

//go:embed vsphere/packer.json.tmpl
//PackerVsphere storage for vSphere packer base template.
var PackerVsphere []byte
