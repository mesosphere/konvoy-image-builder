package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var provisionFlags app.ProvisionFlags

var provisionCmd = &cobra.Command{
	Use:     "provision <inventory.yaml|hostname,>",
	Short:   "provision to an inventory.yaml or hostname, note the comma at the end of the hostname",
	Example: "provision --inventory-file inventory.yaml images/generic/centos-7.yaml",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provisionFlags.RootFlags = rootFlags
		builder := newBuilder()

		var workDir string
		var err error

		if buildFlags.workDir == "" {
			workDir, err = builder.InitConfig(newInitOptions(args[0], buildFlags.generateCLIFlags))
			if err != nil {
				bail("error rendering builder configuration", err, 2)
			}
		} else {
			workDir = buildFlags.workDir
			log.Printf("using workDir provided by --%s flag: %s", workDirFlagName, workDir)
		}

		return builder.Provision(workDir, provisionFlags)
	},
}

func init() {
	fs := provisionCmd.Flags()
	fs.StringArrayVar(&provisionFlags.ExtraVars, "extra-vars", []string{}, "flag passed Ansible's extra-vars")
	fs.StringVar(&provisionFlags.Provider, "provider", "", "specify a provider if you wish to install provider specific utilities")
	fs.StringVar(&provisionFlags.Inventory, "inventory-file", "", "an ansible inventory defining your infrastructure")

	addOverridesArg(fs, &provisionFlags.Overrides)
}
