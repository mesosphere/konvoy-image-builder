package cmd

import (
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one of [artifacts]",
	Args:  cobra.NoArgs,
}

func init() {
	uploadCmd.AddCommand(artifactsCmd)
}
