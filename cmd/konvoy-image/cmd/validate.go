package cmd

import (
	"fmt"
	"net"

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
	Args:          cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isValidCIDR(validateFlags.PodSubnet) {
			return fmt.Errorf("pod-subnet %q was not a valid CIDR", validateFlags.PodSubnet)
		}

		if !isValidCIDR(validateFlags.ServiceSubnet) {
			return fmt.Errorf("service-subnet %q was not a valid CIDR", validateFlags.ServiceSubnet)
		}

		if err := app.Validate(args[0], validateFlags); err != nil {
			return errors.Wrap(err, "error running provision")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	flagSet := validateCmd.Flags()
	flagSet.StringVar(&validateFlags.ServiceSubnet, "service-subnet", "10.96.0.0/12", "ip addresses used"+
		" for the service subnet")
	flagSet.StringVar(&validateFlags.PodSubnet, "pod-subnet", "192.168.0.0/16", "ip addresses used"+
		" for the pod subnet")
	flagSet.IntVar(&validateFlags.APIServerPort, "apiserver-port", 6443, "apiserver port")
}

func isValidCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)

	return err == nil
}
