package version

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// Build information. Populated at build-time.
//
//nolint:gochecknoglobals // Version globals are set at build time.
var (
	version    string
	major      string
	minor      string
	patch      string
	revision   string
	branch     string
	commitDate string
	goVersion  = runtime.Version()
)

// Print returns version information.
func Print(program string) string {
	m := map[string]string{
		"program":    program,
		"version":    version,
		"major":      major,
		"minor":      minor,
		"patch":      patch,
		"revision":   revision,
		"branch":     branch,
		"commitDate": commitDate,
		"goVersion":  goVersion,
		"platform":   findSetting("GOOS", runtime.GOOS) + "/" + findSetting("GOARCH", runtime.GOARCH),
	}
	t := template.Must(template.New("version").Parse(`
	{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
		build date:       {{.commitDate}}
		go version:       {{.goVersion}}
		platform:         {{.platform}}
	`))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

// Info returns version, branch, revision, and git tree state information.
func Info() string {
	return fmt.Sprintf("(version=%s, branch=%s, revision=%s)", version, branch, revision)
}

// BuildContext returns goVersion, and commitDate information.
func BuildContext() string {
	return fmt.Sprintf("(go=%s, date=%s)", goVersion, commitDate)
}

func Version() string {
	return version
}

func Semver() *semver.Version {
	return semver.New(fmt.Sprintf("%s.%s.%s", major, minor, patch))
}

// findSetting attempts to find settings value from binary. These settings are embedded in binary at build time.
// if the settings is not found or attmpt to get buildinfo fails, return defaultValue back
// more info: https://pkg.go.dev/runtime/debug#BuildInfo
func findSetting(setting string, defaultValue string) string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return defaultValue
	}
	for _, s := range buildInfo.Settings {
		if s.Key == setting {
			return s.Value
		}
	}
	return defaultValue
}
