package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var validateFlags app.ValidateFlags

// validateCmd runs validations against nodes to provision.
var validateCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "validate",
	Short:         "validate existing infrastructure",
	Args:          cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := app.Validate(validateFlags); err != nil {
			return errors.Wrap(err, "error running validate")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	flagSet := validateCmd.Flags()
	flagSet.StringVar(&validateFlags.Inventory, "inventory-file", "inventory.yaml", "an ansible inventory defining your infrastructure")
	flagSet.StringVar(&validateFlags.ServiceSubnet, "service-subnet", "10.96.0.0/12", "ip addresses used"+
		" for the service subnet")
	flagSet.StringVar(&validateFlags.PodSubnet, "pod-subnet", "192.168.0.0/16", "ip addresses used"+
		" for the pod subnet")
	flagSet.IntVar(&validateFlags.APIServerPort, "apiserver-port", 6443, "apiserver port")
}
