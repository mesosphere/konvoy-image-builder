package executor

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mesosphere/konvoy-image-builder/pkg/ansible"
	"github.com/mesosphere/konvoy-image-builder/pkg/constants"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

// The Executor will carry out the installation plan.
type Executor interface {
	PreflightsFull(inventoryFile string, playbookOptions *ansible.PlaybookOptions, nodes ...string) error
	PreflightsBasic(inventoryFile string, playbookOptions *ansible.PlaybookOptions, nodes ...string) error
	PreflightsNode(inventoryFile string, playbookOptions *ansible.PlaybookOptions, nodes ...string) error
	CheckNodes(inventoryFile string, playbookOptions *ansible.PlaybookOptions) error
	CheckKubernetes(inventoryFile string, playbookOptions *ansible.PlaybookOptions) error
	RunPlaybook(inventoryFile string, playbookPath string, playbookOptions *ansible.PlaybookOptions) error
}

func InventoryFilename() (string, error) {
	inventory := filepath.Join(constants.WorkingDir, constants.DefaultInventoryFileName)
	f, err := os.Stat(inventory)
	if err != nil {
		if os.IsNotExist(err) {
			newErr := fmt.Errorf("required inventory file %q does not exist", inventory)
			return "", newErr
		}
		return "", fmt.Errorf("could not determine if inventory file %q exists", inventory)
	}
	if f.Size() == 0 {
		return "", fmt.Errorf("required inventory file %q is empty", inventory)
	}
	return inventory, nil
}

// Options are used to configure the executor.
type Options struct {
	// AnsibleDirectory is the location where the Ansible installation and playbooks are located
	AnsibleDirectory string
	// AssetsDirectory is the location where generated assets are to be stored
	AssetsDirectory string
	// RunsDirectory is where information about installation runs is kept
	RunsDirectory string
	// Verbosity of the executor
	Verbosity logging.Verbosity
	// DryRun determines if the executor should actually run the task
	DryRun bool
}

// NewExecutor returns an executor for performing installations according to the installation plan.
func NewExecutor(out io.Writer, errOut io.Writer, options Options) (Executor, error) {
	if options.AssetsDirectory == "" {
		options.AssetsDirectory = filepath.Join(constants.WorkingDir, "generated-assets/")
	}
	if options.RunsDirectory == "" {
		options.RunsDirectory = constants.RunsDir
	}

	return &ansibleExecutor{
		options:    options,
		out:        out,
		errOut:     errOut,
		ansibleDir: options.AnsibleDirectory,
	}, nil
}

type ansibleExecutor struct {
	options    Options
	out        io.Writer
	errOut     io.Writer
	ansibleDir string
}

type task struct {
	// name of the task used for the runs dir
	name string
	// the inventory of nodes to use
	inventoryFile string
	// the playbook filename
	playbook string
	// whether if is is a custom playbook
	customPlaybook bool
	// run the task on specific nodes
	limit []string
	// extraVars contain additional vars to pass to the runner
	playbookOptions *ansible.PlaybookOptions
}

// execute will run the given task, and setup all what's needed for us to run ansible.
func (ae *ansibleExecutor) execute(t task) error {
	runDirectory, err := ae.createRunDirectory(t.name)
	if err != nil {
		return fmt.Errorf("error creating working directory for %q: %w", t.name, err)
	}
	ansibleLogFilename := filepath.Join(runDirectory, "ansible.log")
	ansibleLogFile, err := os.Create(ansibleLogFilename)
	if err != nil {
		return fmt.Errorf("error creating ansible log file %q: %w", ansibleLogFilename, err)
	}
	defer ansibleLogFile.Close()
	runner := ae.ansibleRunner(ansibleLogFile)

	fmt.Fprintf(ae.out, "\nSTAGE [%s]\n", t.name)

	// Start running ansible with the given playbook
	if t.limit != nil && len(t.limit) != 0 {
		err = runner.StartPlaybookOnNode(pathFileName(t.playbook), t.inventoryFile, t.playbookOptions, t.limit...)
	} else {
		err = runner.StartPlaybook(pathFileName(t.playbook), t.inventoryFile, t.playbookOptions)
	}
	if err != nil {
		return fmt.Errorf("running ansible playbook failed: %w", err)
	}

	// Wait until ansible exits
	if err = runner.WaitPlaybook(); err != nil {
		return fmt.Errorf("ansible playbook execution failed: %w", err)
	}
	return nil
}

func (ae *ansibleExecutor) PreflightsFull(inventoryFile string,
	playbookOptions *ansible.PlaybookOptions,
	nodes ...string) error {
	if playbookOptions.ExtraVarsMap == nil {
		playbookOptions.ExtraVarsMap = make(map[string]interface{})
	}

	playbookOptions.ExtraVarsMap["preflight_type"] = "full"

	t := task{
		name:            constants.PreflightsName,
		playbook:        constants.PreflightsPlaybook,
		inventoryFile:   inventoryFile,
		limit:           nodes,
		playbookOptions: playbookOptions,
	}
	return ae.execute(t)
}

func (ae *ansibleExecutor) PreflightsNode(inventoryFile string,
	playbookOptions *ansible.PlaybookOptions,
	nodes ...string) error {
	if playbookOptions.ExtraVarsMap == nil {
		playbookOptions.ExtraVarsMap = make(map[string]interface{})
	}

	playbookOptions.ExtraVarsMap["preflight_type"] = "node"

	t := task{
		name:            constants.PreflightsName,
		playbook:        constants.PreflightsPlaybook,
		inventoryFile:   inventoryFile,
		limit:           nodes,
		playbookOptions: playbookOptions,
	}
	return ae.execute(t)
}

func (ae *ansibleExecutor) PreflightsBasic(inventoryFile string,
	playbookOptions *ansible.PlaybookOptions,
	nodes ...string) error {
	if playbookOptions.ExtraVarsMap == nil {
		playbookOptions.ExtraVarsMap = make(map[string]interface{})
	}
	playbookOptions.ExtraVarsMap["preflight_type"] = "basic"

	t := task{
		name:            constants.PreflightsName,
		playbook:        constants.PreflightsPlaybook,
		inventoryFile:   inventoryFile,
		limit:           nodes,
		playbookOptions: playbookOptions,
	}
	return ae.execute(t)
}

func (ae *ansibleExecutor) CheckNodes(inventoryFile string, playbookOptions *ansible.PlaybookOptions) error {
	t := task{
		name:            constants.CheckNodesName,
		playbook:        constants.CheckNodesPlaybook,
		inventoryFile:   inventoryFile,
		playbookOptions: playbookOptions,
	}
	return ae.execute(t)
}

func (ae *ansibleExecutor) CheckKubernetes(inventoryFile string, playbookOptions *ansible.PlaybookOptions) error {
	t := task{
		name:            constants.CheckKubernetesName,
		playbook:        constants.CheckKubernetesPlaybook,
		inventoryFile:   inventoryFile,
		playbookOptions: playbookOptions,
	}
	return ae.execute(t)
}

func (ae *ansibleExecutor) RunPlaybook(inventoryFile string, playbookPath string, playbookOptions *ansible.PlaybookOptions) error {
	t := task{
		name:            constants.RunPlaybookName,
		playbook:        playbookPath,
		customPlaybook:  true,
		inventoryFile:   inventoryFile,
		playbookOptions: playbookOptions,
	}
	return ae.execute(t)
}

func (ae *ansibleExecutor) createRunDirectory(runName string) (string, error) {
	start := time.Now()
	runDirectory := filepath.Join(ae.options.RunsDirectory, runName, start.Format("2006-01-02-15-04-05"))
	if err := os.MkdirAll(runDirectory, 0750); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}
	return runDirectory, nil
}

func (ae *ansibleExecutor) ansibleRunner(ansibleLog io.Writer) ansible.Runner {
	// Send stdout and stderr to ansibleOut

	loggerConfig := &ansible.LoggerConfig{
		Out:       ae.out,
		ErrOut:    ae.errOut,
		Log:       ansibleLog,
		Verbosity: int(ae.options.Verbosity),
	}

	return ansible.NewRunner(ae.ansibleDir, loggerConfig)
}

func pathFileName(name string) string {
	return path.Join(constants.AnsiblePlaybookPath, name)
}