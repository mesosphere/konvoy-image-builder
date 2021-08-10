package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/mesosphere/konvoy-image-builder/pkg/constants"
	"github.com/mesosphere/konvoy-image-builder/pkg/logging"
)

type commonCLIFlags struct {
	verbose   bool //nolint:structcheck // not used directly here
	verbosity int  //nolint:structcheck // not used directly here
	dryRun    bool //nolint:structcheck // not used directly here
}

func AddVerboseFlag(flagSet *pflag.FlagSet, verbose *bool) {
	flagSet.BoolVar(verbose, "verbose", false, "enable debug level logging")
}

func AddVerbosityFlag(flagSet *pflag.FlagSet, verbosity *int) {
	flagSet.IntVarP(verbosity, "v", "v",
		int(logging.UnknownLevel),
		fmt.Sprintf("select verbosity level, "+
			"should be between %d and %d", int(logging.UnknownLevel), int(logging.TraceLevel)))
}

func intToVerbosity(verbosity int) logging.Verbosity {
	if verbosity < int(logging.UnknownLevel) {
		return logging.UnknownLevel
	}
	if verbosity > int(logging.TraceLevel) {
		return logging.TraceLevel
	}
	return logging.Verbosity(verbosity)
}

func verbosityFromFlags(verbose bool, verbosity int) (int, error) {
	if verbose {
		// TODO: Print the warning once `--v` becomes a visible flag.
		// fmt.Fprintln(os.Stdout, "WARNING: '--verbose' is deprecated, use '--v 4' instead (other log levels are also accepted)")
		if verbosity != int(logging.UnknownLevel) {
			return int(logging.UnknownLevel), errors.New("'--verbose' and '--v' " +
				"must not be set together, only use '--v' instead")
		}
		return int(logging.DebugLevel), nil
	}
	return verbosity, nil
}

func AddClusterNameFlag(flagSet *pflag.FlagSet, clusterName *string) {
	flagSet.StringVar(clusterName, "cluster-name", filepath.Base(constants.WorkingDir), "name used to "+
		"prefix the cluster and all the created resources")
}
