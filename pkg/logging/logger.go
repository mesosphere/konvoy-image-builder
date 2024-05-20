package logging

import (
	"errors"

	"github.com/sirupsen/logrus"

	logrustooling "github.com/onrik/logrus/filename"
)

// Logger is the custom global logger object.
var Logger = logrus.New()

type Verbosity int

const (
	UnknownLevel   Verbosity = iota // Unknown verbosity, no logs in Ansible.
	TimestampLevel                  // Timestamp logs in Ansible, Ansible verbosity set to 1
	ErrorLevel                      // Ansible verbosity set to 2
	WarningLevel                    // Ansible verbosity set to 3
	InfoLevel                       // Ansible verbosity set to 4
	DebugLevel                      // Ansible verbosity set to 5, Packer debug flag on
	TraceLevel                      // Ansible verbosity set to 6
)

func VerbosityToDebug(verbosity Verbosity) bool {
	return verbosity >= DebugLevel
}

// SetLogLevel sets the log-level for the logger basing on the Verbosity
// provided. The log level string must be correct, otherwise function will
// panic.
func SetLogLevel(loglevel Verbosity) {
	var enableTimestamps bool
	var enableSourcecodeInfo bool

	// Can be done in less verbose way, but let's opt for explicitness and
	// clarity.
	switch loglevel {
	case UnknownLevel:
		Logger.SetLevel(logrus.ErrorLevel)
		enableTimestamps = false
		enableSourcecodeInfo = false
	case TimestampLevel:
		Logger.SetLevel(logrus.WarnLevel)
		enableTimestamps = true
		enableSourcecodeInfo = false
	case ErrorLevel:
		Logger.SetLevel(logrus.InfoLevel)
		enableTimestamps = true
		enableSourcecodeInfo = true
	case WarningLevel:
		Logger.SetLevel(logrus.DebugLevel)
		enableTimestamps = true
		enableSourcecodeInfo = true
	default:
		// InfoLevel and higher
		Logger.SetLevel(logrus.TraceLevel)
		enableTimestamps = true
		enableSourcecodeInfo = true
	}

	if enableTimestamps {
		(Logger.Formatter).(*logrus.TextFormatter).DisableTimestamp = false
		(Logger.Formatter).(*logrus.TextFormatter).FullTimestamp = true
	} else {
		(Logger.Formatter).(*logrus.TextFormatter).DisableTimestamp = true
		(Logger.Formatter).(*logrus.TextFormatter).FullTimestamp = false
	}

	// We want this function to be callable multiple times, hence we zero-out
	// the existing hooks before we add(readd) the filename hook
	Logger.ReplaceHooks(make(logrus.LevelHooks))
	if enableSourcecodeInfo {
		filenameHook := logrustooling.NewHook()
		filenameHook.Field = "src"
		Logger.AddHook(filenameHook)
	}
}

var ErrVerbosityConflict = errors.New(
	"verbose and verbosity can not both be set",
)

// VerbosityFromFlags returns the verbosity given "verbose" and "verbosity"
// flags.
func VerbosityFromFlags(verbose bool, verbosity int) (Verbosity, error) {
	if verbose {
		if verbosity != int(UnknownLevel) {
			return UnknownLevel, ErrVerbosityConflict
		}
		return DebugLevel, nil
	}
	return IntToVerbosity(verbosity), nil
}

func IntToVerbosity(verbosity int) Verbosity {
	if verbosity < int(UnknownLevel) {
		return UnknownLevel
	}
	if verbosity > int(TraceLevel) {
		return TraceLevel
	}
	return Verbosity(verbosity)
}
