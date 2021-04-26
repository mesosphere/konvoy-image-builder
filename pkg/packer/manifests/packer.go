package manifests

//nolint:golint // blank import for embed support.
import _ "embed"

//go:embed aws/packer.json.tmpl
//PackerAmazon storage for AWS packer template.
var PackerAmazon []byte

// add more embedded files here
// packer/azure/packer.json.tmpl for example
// var PackerAzure []byte
