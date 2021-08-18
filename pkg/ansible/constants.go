package ansible

import (
	"os"
)

const (
	AnsibleCallbackWhiteListVerbose = "profile_tasks"
)

var (
	PythonPath  = os.Getenv("PYTHON_PATH")
	AnsiblePath = os.Getenv("ANSIBLE_PATH")
)
