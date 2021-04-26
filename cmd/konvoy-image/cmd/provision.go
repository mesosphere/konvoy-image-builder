package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var provisionFlags app.ProvisionFlags

var provisionCmd = &cobra.Command{
	Use:     "provision <inventory.yaml|hostname,>",
	Short:   "provision to an inventory.yaml or hostname, note the comma at the end of the hostname",
	Example: "build ./inventory.yaml",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provisionFlags.RootFlags = rootFlags

		if err := app.Provision(args[0], provisionFlags); err != nil {
			return errors.Wrap(err, "error running provision")
		}

		return nil
	},
}

func init() {
	fs := provisionCmd.Flags()
	fs.StringArrayVar(&provisionFlags.ExtraVars, "extra-vars", []string{}, "flag passed Ansible's extra-vars")
}
