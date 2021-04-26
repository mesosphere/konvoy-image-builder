package ansible

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/term"

	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

var logger = logging.Logger

// OutputFormat is used for controlling the STDOUT format of the Ansible runner.
type OutputFormat string

// Runner for running Ansible playbooks.
type Runner interface {
	// TODO files should be some io interface to falicitate testing
	// StartPlaybook runs the playbook asynchronously with the given inventory and extra vars.
	StartPlaybook(playbookFileName, inventory string, playbookOptions *PlaybookOptions) error
	// WaitPlaybook blocks until the execution of the playbook is complete. If an error occurred,
	// it is returned. Otherwise, returns nil to signal the completion of the playbook.
	WaitPlaybook() error
}

type PlaybookOptions struct {
	ExtraVars    []string
	ExtraVarsMap map[string]interface{}
}

type LoggerConfig struct {
	// Out is the stdout writer for the Ansible process
	Out io.Writer
	// ErrOut is the stderr writer for the Ansible process
	ErrOut io.Writer
	// Output also gets logged to a file
	Log       io.Writer
	Verbosity int
}

func NewLoggerConfig(writer io.Writer, verbosity int) *LoggerConfig {
	return &LoggerConfig{
		Out:       writer,
		ErrOut:    writer,
		Log:       writer,
		Verbosity: verbosity,
	}
}

type runner struct {
	*LoggerConfig
	pythonPath   string
	ansiblePath  string
	runDir       string
	waitPlaybook func() error
}

// NewRunner returns a new runner for running Ansible playbooks.
func NewRunner(runDir string, loggerConfig *LoggerConfig) Runner {
	return &runner{
		LoggerConfig: loggerConfig,
		pythonPath:   PythonPath,
		ansiblePath:  AnsiblePath,
		runDir:       runDir,
	}
}

var ErrPlaybookNotStarted = errors.New(
	"playbook not started",
)

// WaitPlaybook blocks until the ansible process running the playbook exits.
// If the process exits with a non-zero status, it will return an error.
func (r *runner) WaitPlaybook() error {
	if r.waitPlaybook == nil {
		return errors.Wrap(ErrPlaybookNotStarted, "wait called")
	}
	execErr := r.waitPlaybook()

	if execErr != nil {
		return fmt.Errorf("error running ansible: %w", execErr)
	}
	return nil
}

func (r *runner) StartPlaybook(playbookFileName, inventory string, playbookOptions *PlaybookOptions) error {
	if playbookOptions == nil {
		playbookOptions = &PlaybookOptions{}
	}

	marshalledExtraVars, err := marshalledExtraVars(playbookOptions.ExtraVarsMap)
	if err != nil {
		return err
	}

	/* #nosec : G204: Subprocess launched with function call as argument or cmd arguments */
	cmd := exec.Command(path.Join(r.ansiblePath, "bin", "ansible-playbook"), playbookFileName,
		"-i", inventory,
		"--extra-vars", string(marshalledExtraVars),
	)
	cmd.Args = append(cmd.Args, r.runConditionalFlags(playbookOptions)...)
	cmd.Env = r.runEnv()
	PrintCmd(cmd, r.Log)

	// also Log to a file
	cmd.Stdout = io.MultiWriter(r.Log, r.Out)
	cmd.Stderr = io.MultiWriter(r.Log, r.ErrOut)
	cmd.Stdin = os.Stdin

	logger.SetOutput(r.Out)

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "error running playbook")
	}

	r.waitPlaybook = cmd.Wait

	return nil
}

func marshalledExtraVars(extraVars map[string]interface{}) ([]byte, error) {
	marshalledExtraVars, err := json.Marshal(extraVars)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal extra-vars")
	}
	return marshalledExtraVars, nil
}

func (r *runner) runEnv() []string {
	pythonPathEnvVar := "PYTHONPATH=" + r.pythonPath
	ansibleConfigEnvVar := "ANSIBLE_CONFIG=" + filepath.Join(r.ansiblePath, "playbooks", "ansible.cfg")
	env := append(os.Environ(), pythonPathEnvVar, ansibleConfigEnvVar)
	if r.Verbosity > 0 {
		ansibleCallbackWhitelistEnvVar := "ANSIBLE_CALLBACK_WHITELIST=" + AnsibleCallbackWhiteListVerbose
		env = append(env, ansibleCallbackWhitelistEnvVar)
	}
	// usually Ansible would automatically determine if it should output color based on TTY
	// because using an io.MultiWriter, Ansible won't print color
	// force it to print color but only with a TTY, to avoid printing extra characters when running as a pod in Kubernetes
	ansibleForceColor := fmt.Sprintf("ANSIBLE_FORCE_COLOR=%t", term.IsTerminal(int(os.Stdout.Fd())))

	return append(env, ansibleForceColor)
}

func (r *runner) runConditionalFlags(playbookOptions *PlaybookOptions) []string {
	args := make([]string, 0, 2)

	for _, extraVarsFilePath := range playbookOptions.ExtraVars {
		args = append(args, "--extra-vars", extraVarsFilePath)
	}

	if r.Verbosity > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", r.Verbosity)))
	}

	return args
}

func PrintCmd(cmd *exec.Cmd, log io.Writer) {
	fmt.Fprintln(log, "To run the Ansible command, you can run:")
	fmt.Fprintln(log, strings.Join(cmd.Args, " "))
	fmt.Fprintln(log, "\nComplete environment:\n"+strings.Join(cmd.Env, "\n"))
}
