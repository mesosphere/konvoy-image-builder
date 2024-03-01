package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var generateDocsCmd = &cobra.Command{
	Use:     "generate-docs <PATH>",
	Short:   "generate docs in path",
	Example: "generate-docs /tmp/docs",
	Args:    cobra.ExactArgs(1),

	RunE: func(_ *cobra.Command, args []string) error {
		path := args[0]

		if _, err := os.Stat(path); os.IsNotExist(err) {
			if mkErr := os.Mkdir(path, os.ModeDir); mkErr != nil {
				bail(fmt.Sprintf("error creating %s", path), mkErr, 2)
			}
		} else if err != nil {
			bail(fmt.Sprintf("error checking %s", path), err, 3)
		}

		if err := doc.GenMarkdownTree(rootCmd, path); err != nil {
			bail("error generating docs", err, 4)
		}

		return nil
	},
}
