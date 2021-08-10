package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

type checkCLIFlags struct {
	commonCLIFlags
	clusterName         string
	serviceSubnet       string
	podSubnet           string
	calicoEncapsulation string
	cloudProvider       string
	kubernetesVersion   string
	targetRAMMB         int
	errorsToIgnore      string
	preflightChecks     string
}

var checkFlags checkCLIFlags

// checkCmd runs checks against kubernetes and addons.
var checkCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "check",
	Short:         "Run checks on the health of the cluster",
	Args:          cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVerbosityAndSetLogs(); err != nil {
			return err
		}

		errors := make([]error, 0)
		opts := checkFlagsToOptions(checkFlags)

		// preflight checks (fail fast here)
		if err := app.CheckPreflight(opts); err != nil {
			return err
		}

		// check nodes
		if err := app.CheckNodes(opts); err != nil {
			errors = append(errors, err)
		}

		// check kubernetes
		if err := app.CheckKubernetes(opts); err != nil {
			errors = append(errors, err)
		}

		if len(errors) == 0 {
			return nil
		}

		var b bytes.Buffer
		for _, err := range errors {
			b.WriteString(fmt.Sprintf("- %v\n", err))
		}
		return fmt.Errorf(b.String())
	},
}

// checkPreflightCmd runs checks to see if machines are ready to run kubeadm.
var checkPreflightCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "preflight",
	Short:         "Run checks to validate machines are ready for installation",
	Args:          cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVerbosityAndSetLogs(); err != nil {
			return err
		}

		return app.CheckPreflight(checkFlagsToOptions(checkFlags))
	},
}

// checkNodesCmd runs checks on the nodes to see if they are healthy.
var checkNodesCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "nodes",
	Short:         "Run checks on the nodes",
	Args:          cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVerbosityAndSetLogs(); err != nil {
			return err
		}

		return app.CheckNodes(checkFlagsToOptions(checkFlags))
	},
}

// checkKubernetesCmd runs checks against kubernetes.
var checkKubernetesCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "kubernetes",
	Short:         "Run checks on the cluster components",
	Args:          cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkVerbosityAndSetLogs(); err != nil {
			return err
		}

		return app.CheckKubernetes(checkFlagsToOptions(checkFlags))
	},
}

func checkFlagsToOptions(flag checkCLIFlags) app.CheckOptions {
	return app.CheckOptions{
		Out:                 out,
		ErrOut:              errOut,
		Verbosity:           intToVerbosity(flag.verbosity),
		DryRun:              flag.dryRun,
		ClusterName:         flag.clusterName,
		ServiceSubnet:       flag.serviceSubnet,
		PodSubnet:           flag.podSubnet,
		CalicoEncapsulation: flag.calicoEncapsulation,
		CloudProvider:       flag.cloudProvider,
		KubernetesVersion:   flag.kubernetesVersion,
		TargetRAMMB:         flag.targetRAMMB,
		ErrorsToIgnore:      flag.errorsToIgnore,
		PreflightChecks:     flag.preflightChecks,
	}
}

func checkVerbosityAndSetLogs() error {
	var err error
	checkFlags.verbosity, err = verbosityFromFlags(checkFlags.verbose, checkFlags.verbosity)
	if err != nil {
		return fmt.Errorf("error choosing verbosity: %w", err)
	}
	logging.SetLogLevel(logging.Verbosity(checkFlags.verbosity))

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.AddCommand(checkPreflightCmd)
	checkCmd.AddCommand(checkNodesCmd)
	checkCmd.AddCommand(checkKubernetesCmd)

	addFlags(checkCmd)
	addFlags(checkPreflightCmd)
	addFlags(checkNodesCmd)
	addFlags(checkKubernetesCmd)
}

func addFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()
	AddVerboseFlag(flagSet, &checkFlags.verbose)
	AddVerbosityFlag(flagSet, &checkFlags.verbosity)
	AddClusterNameFlag(flagSet, &checkFlags.clusterName)
	addServiceSubnetFlag(flagSet, &checkFlags.serviceSubnet)
	addPodSubnetFlag(flagSet, &checkFlags.podSubnet)
	addCalicoEncapsulationFlag(flagSet, &checkFlags.calicoEncapsulation)
	addCloudProviderFlag(flagSet, &checkFlags.cloudProvider)
	addKubernetesVersionFlag(flagSet, &checkFlags.kubernetesVersion)
	addTargetRAMMBFlag(flagSet, &checkFlags.targetRAMMB)
	addErrorsToIgnoreFlag(flagSet, &checkFlags.errorsToIgnore)
	addPreflightChecksFlag(flagSet, &checkFlags.preflightChecks)
	if err := flagSet.MarkHidden("v"); err != nil {
		println("could not mark hidden the 'verbosity' flag")
	}
}

func addServiceSubnetFlag(flagSet *pflag.FlagSet, serviceSubnet *string) {
	flagSet.StringVar(serviceSubnet, "service-subnet", "10.96.0.0/12", "ip addresses used"+
		" for the service subnet")
}

func addPodSubnetFlag(flagSet *pflag.FlagSet, podSubnet *string) {
	flagSet.StringVar(podSubnet, "pod-subnet", "192.168.0.0/16", "ip addresses used"+
		" for the pod subnet")
}

func addCalicoEncapsulationFlag(flagSet *pflag.FlagSet, calicoEncapsulation *string) {
	flagSet.StringVar(calicoEncapsulation, "calico-encapsulation", "vxlan", "calico "+
		"encapsulation")
}

func addCloudProviderFlag(flagSet *pflag.FlagSet, cloudProvider *string) {
	flagSet.StringVar(cloudProvider, "cloud-provider", "aws", "cloud provider")
}

func addKubernetesVersionFlag(flagSet *pflag.FlagSet, kubernetesVersion *string) {
	flagSet.StringVar(kubernetesVersion, "kubernetes-version", "1.22", "kubernetes version")
}

func addTargetRAMMBFlag(flagSet *pflag.FlagSet, targetRAMMB *int) {
	flagSet.IntVar(targetRAMMB, "target-ram", 4096, "target ram size on a node in MB")
}

func addErrorsToIgnoreFlag(flagSet *pflag.FlagSet, errorsToIgnore *string) {
	flagSet.StringVar(errorsToIgnore, "errors-to-ignore", "", "comma separated "+
		"list of errors to ignore")
}

func addPreflightChecksFlag(flagSet *pflag.FlagSet, preflightChecks *string) {
	flagSet.StringVar(preflightChecks, "preflight-checks", "", "comma separated list "+
		"of preflight checks to run")
}
