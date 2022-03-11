package cmd

import (
	"github.com/spf13/pflag"
)

func addOverridesArg(fs *pflag.FlagSet, overrides *[]string) {
	fs.StringSliceVar(overrides, "overrides", []string{}, "a comma separated list of override YAML files")
}

func addExtraVarsArg(fs *pflag.FlagSet, extraVars *[]string) {
	fs.StringSliceVar(extraVars, "extra-vars", []string{}, "flag passed Ansible's extra-vars")
}

func addClusterArgs(fs *pflag.FlagSet, kubernetesVersion, containerdVersion *string) {
	fs.StringVar(kubernetesVersion, "kubernetes-version", "", "The version of kubernetes to install. Example: 1.21.6")
	fs.StringVar(containerdVersion, "containerd-version", "", "the version of containerd to install")
}
