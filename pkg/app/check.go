package app

import (
	"fmt"
	"io"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/constants"
	"github.com/mesosphere/konvoy-image-builder/pkg/executor"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

const (
	extraVarsTemplate = "%s=%s"
)

type CheckOptions struct {
	Out                 io.Writer
	ErrOut              io.Writer
	Verbosity           logging.Verbosity
	DryRun              bool
	ClusterName         string
	ServiceSubnet       string
	PodSubnet           string
	CalicoEncapsulation string
	CloudProvider       string
	KubernetesVersion   string
	TargetRAMMB         int
	ErrorsToIgnore      string
	PreflightChecks     string
}

func CheckPreflight(opts CheckOptions) error {
	ansibleExecutor, inventory, err := getParseOptionsGetExecutorAndInventory(opts)
	if err != nil {
		return err
	}

	playbookOptions := checkOptionsToPlaybookOptions(opts)
	return ansibleExecutor.PreflightsFull(inventory, playbookOptions)
}

func CheckNodes(opts CheckOptions) error {
	ansibleExecutor, inventory, err := getParseOptionsGetExecutorAndInventory(opts)
	if err != nil {
		return err
	}

	// The idea is that we want the checks to always run to confirm
	// that nodes didn’t get modified and lost the requirements from the
	// preflights.

	playbookOptions := checkOptionsToPlaybookOptions(opts)
	err = ansibleExecutor.PreflightsNode(inventory, playbookOptions)
	if err != nil {
		return err
	}

	return ansibleExecutor.CheckNodes(inventory, playbookOptions)
}

func CheckKubernetes(opts CheckOptions) error {
	ansibleExecutor, inventory, err := getParseOptionsGetExecutorAndInventory(opts)
	if err != nil {
		return err
	}

	playbookOptions := checkOptionsToPlaybookOptions(opts)
	return ansibleExecutor.CheckKubernetes(inventory, playbookOptions)
}

func getParseOptionsGetExecutorAndInventory(opts CheckOptions) (executor.Executor, string, error) {
	executorOpts := executor.Options{
		AnsibleDirectory: constants.AnsibleDir,
		Verbosity:        opts.Verbosity,
		DryRun:           opts.DryRun,
	}

	ansibleExecutor, err := executor.NewExecutor(opts.Out, opts.Out, executorOpts)
	if err != nil {
		return nil, "", err
	}

	inventory, err := executor.InventoryFilename()
	if err != nil {
		return nil, "", err
	}

	return ansibleExecutor, inventory, nil
}

func checkOptionsToPlaybookOptions(opts CheckOptions) *ansible.PlaybookOptions {
	playbookOptions := &ansible.PlaybookOptions{}
	values := []string{
		fmt.Sprintf(extraVarsTemplate, "service_subnet", opts.ServiceSubnet),
		fmt.Sprintf(extraVarsTemplate, "pod_subnet", opts.PodSubnet),
		fmt.Sprintf(extraVarsTemplate, "calico_encapsulation", opts.CalicoEncapsulation),
		fmt.Sprintf(extraVarsTemplate, "cloud_provider", opts.CloudProvider),
		fmt.Sprintf(extraVarsTemplate, "kubernetes_version", opts.KubernetesVersion),
	}

	if playbookOptions.ExtraVars == nil {
		playbookOptions.ExtraVars = make([]string, 0)
	}

	playbookOptions.ExtraVars = append(playbookOptions.ExtraVars, values...)

	return playbookOptions
}