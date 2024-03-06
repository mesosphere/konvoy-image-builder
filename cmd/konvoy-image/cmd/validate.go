package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var validateFlags app.ValidateFlags

// validateCmd runs validations against nodes to provision.
var validateCmd = &cobra.Command{
	SilenceUsage: true, SilenceErrors: true,
	Use:   "validate",
	Short: "validate existing infrastructure",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := app.Validate(validateFlags); err != nil {
			return fmt.Errorf("failed running validate %w", err)
		}

		return nil
	},
}

func init() {
	flagSet := validateCmd.Flags()
	ipAddUsed := "ip addresses used"
	flagSet.StringVar(&validateFlags.Inventory, "inventory-file", "inventory.yaml", "an ansible inventory defining your infrastructure")
	flagSet.StringVar(&validateFlags.ServiceSubnet, "service-subnet", "10.96.0.0/12", ipAddUsed+
		" for the service subnet")
	flagSet.StringVar(&validateFlags.PodSubnet, "pod-subnet", "192.168.0.0/16", ipAddUsed+
		" for the pod subnet")
	flagSet.IntVar(&validateFlags.APIServerPort, "apiserver-port", 6443, "apiserver port")
}
