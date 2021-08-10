package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
	"github.com/mesosphere/konvoy-image-builder/pkg/version"
)

type rootCmdFlagsStruct struct {
	Verbosity int
	Verbose   bool
	app.RootFlagsCommon
}

var (
	rootCmdFlags rootCmdFlagsStruct
	rootFlags    app.RootFlags

	out    = os.Stdout
	errOut = os.Stderr
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "konvoy-image",
	Short: "Create, provision, and customize images for running Konvoy",
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		rootFlags, err = rootCmdFlagsToRootFlags(rootCmdFlags)
		if err != nil {
			return err
		}
		logging.SetLogLevel(rootFlags.Verbosity)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version for konvoy-image",
}

// TODO: needed for CLI docgen
//func RootCmd() *cobra.Command {
//	return rootCmd
//}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// set version string
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", version.Print("konvoy-image")))

	// also need to set Version to get cobra to print it
	rootCmd.Version = version.Info()

	// add a version command
	versionCmd.Run = func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(os.Stdout, "%s\n", version.Print("konvoy-image"))
	}

	if err := rootCmd.Execute(); err != nil {
		logging.Logger.Fatalf("error: %s", err)
	}
}

func bail(message string, err error, code int) {
	if err != nil {
		logging.Logger.Errorf("%v\n", errors.Wrap(err, message))
	} else {
		logging.Logger.Error(message + "\n")
	}
	os.Exit(code)
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(generateDocsCmd)
	rootCmd.AddCommand(provisionCmd)
	rootCmd.DisableAutoGenTag = true

	fs := rootCmd.PersistentFlags()
	fs.BoolVar(
		&rootCmdFlags.Verbose,
		"verbose",
		false,
		fmt.Sprintf(
			"enable debug level logging (same as --v %v)",
			logging.DebugLevel,
		),
	)
	fs.IntVarP(
		&rootCmdFlags.Verbosity,
		"v",
		"v",
		int(logging.UnknownLevel),
		fmt.Sprintf(
			"select verbosity level, should be between %d and %d",
			int(logging.UnknownLevel),
			int(logging.TraceLevel),
		),
	)
	fs.BoolVar(&rootCmdFlags.Color, "color", true, "enable color output")
}

var errVerboseCombo = errors.New(
	"'--verbose' and '--v' flags must not be set together, only use '--v' instead",
)

func rootCmdFlagsToRootFlags(rootCmdFlags rootCmdFlagsStruct) (app.RootFlags, error) {
	verbosity, err := logging.VerbosityFromFlags(rootCmdFlags.Verbose, rootCmdFlags.Verbosity)
	if err != nil {
		return app.RootFlags{}, errVerboseCombo
	}

	return app.RootFlags{
		Verbosity:       verbosity,
		RootFlagsCommon: rootCmdFlags.RootFlagsCommon,
	}, nil
}
