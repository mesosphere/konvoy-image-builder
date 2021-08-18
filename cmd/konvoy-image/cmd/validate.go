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
	Short:         "Run checks on the health of infrastructure",
	Args:          cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// preflight checks (fail fast here)
		if err := app.Validate(args[0], validateFlags); err != nil {
			return errors.Wrap(err, "error running provision")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	flagSet := rootCmd.Flags()
	flagSet.StringArrayVar(&validateFlags.ExtraVars, "extra-vars", []string{}, "flag passed Ansible's extra-vars")
	flagSet.StringVar(&validateFlags.ServiceSubnet, "service-subnet", "10.96.0.0/12", "ip addresses used"+
		" for the service subnet")
	flagSet.StringVar(&validateFlags.PodSubnet, "pod-subnet", "192.168.0.0/16", "ip addresses used"+
		" for the pod subnet")
	flagSet.StringVar(&validateFlags.CalicoEncapsulation, "calico-encapsulation", "vxlan", "calico "+
		"encapsulation")
	flagSet.StringVar(&validateFlags.CloudProvider, "cloud-provider", "aws", "cloud provider")
	flagSet.StringVar(&validateFlags.KubernetesVersion, "kubernetes-version", "1.22", "kubernetes version")
	flagSet.IntVar(&validateFlags.TargetRAMMB, "target-ram", 4096, "target ram size on a node in MB")
	flagSet.StringVar(&validateFlags.ErrorsToIgnore, "errors-to-ignore", "", "comma separated "+
		"list of errors to ignore")
	flagSet.StringVar(&validateFlags.PreflightChecks, "preflight-checks", "", "comma separated list "+
		"of preflight checks to run")
}
