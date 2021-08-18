package constants

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mesosphere/konvoy-image-builder/printerlib/pkg/printer"
)

const (
	AnsibleCallbackWhiteListVerbose = "profile_tasks"
)

var (
	AnsiblePlaybookPath  = "ansible"
	AnsibleRunsDirectory = "ansible-runs"

	AnsibleDir    = ansibleDir()
	ExecutableDir = executableDir()
	WorkingDir    = workingDir()
	PythonPath    = CalculatePythonPath(ansibleDir())
	RunsDir       = filepath.Join(WorkingDir, "runs/")
)

func ansibleDir() string {
	ap := os.Getenv("ANSIBLE_PATH")
	if ap == "" {
		ap = filepath.Join(ExecutableDir, "ansible")
	}

	return ap
}

// TODO return error but there is no way to recover.
func executableDir() string {
	ex, err := os.Executable()
	if err != nil {
		printer.PrintColor(os.Stderr, printer.Red, "Error: could not get executable directory %v", err)
		os.Exit(1)
	}
	return filepath.Dir(ex)
}

func workingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		printer.PrintColor(os.Stderr, printer.Red, "Error: could not get current working directory %v", err)
		os.Exit(1)
	}
	return wd
}

func CalculatePythonPath(ansibleDir string) string {
	pythonVersion, err := pythonVersion()
	if err != nil {
		return ""
	}
	return filepath.Join(ansibleDir, "lib", "python"+pythonVersion[:3], "site-packages")
}

func pythonVersion() (string, error) {
	pythonBinary, err := exec.LookPath("/usr/bin/python")
	if err != nil {
		return "", err
	}
	out, err := exec.Command(pythonBinary, "-V").CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.Split(string(out), " ")[1], nil
}