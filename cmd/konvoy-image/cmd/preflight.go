package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

var preflightFlags app.PreflightFlags

var preflightCmd = &cobra.Command{
	Use:     "preflight inventory.yaml",
	Short:   "run preflights based on inventory.yaml",
	Example: "pre-flight ./inventory.yaml",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preflightFlags.RootFlags = rootFlags
		if err := app.Preflight(args[0], preflightFlags); err != nil {
			return errors.Wrap(err, "error running preflight")
		}
		return nil
	},
}
