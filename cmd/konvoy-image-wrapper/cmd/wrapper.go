package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/mitchellh/go-homedir"
	errors2 "github.com/pkg/errors"
	terminal "golang.org/x/term"

	"github.com/mesosphere/konvoy-image-builder/cmd/konvoy-image-wrapper/image"
	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

const (
	envAWSCredentialsFile = "AWS_SHARED_CREDENTIALS_FILE" //nolint:gosec // environment var set by user

	envAWSDefaultRegion = "AWS_DEFAULT_REGION"
	envAWSRegion        = "AWS_REGION"
	envAWSProfile       = "AWS_PROFILE"
	envAWSAccessKeyID   = "AWS_ACCESS_KEY_ID"

	//nolint:gosec // environment var set by user
	envAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	//nolint:gosec // environment var set by user
	envAWSSessionToken = "AWS_SESSION_TOKEN"

	envAWSSTSRegionalEndpoints = "AWS_STS_REGIONAL_ENDPOINTS"
	envAWSSDKLoadConfig        = "AWS_SDK_LOAD_CONFIG"
	envAWSCABundle             = "AWS_CA_BUNDLE"

	envAzureLocation = "AZURE_LOCATION"

	envVSphereServer                     = "VSPHERE_SERVER"
	envVSphereUser                       = "VSPHERE_USERNAME"
	envVSpherePassword                   = "VSPHERE_PASSWORD"
	envRedHatSubscriptionManagerUser     = "RHSM_USER"
	envRedHatSubscriptionManagerPassword = "RHSM_PASS"

	//nolint:gosec // environment var set by user
	envGCPApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"

	envHTTPSProxy = "HTTPS_PROXY"
	envHTTPProxy  = "HTTP_PROXY"
	envNoProxy    = "NO_PROXY"

	containerWorkingDir = "/tmp/kib"
)

var ErrEnv = errors.New("manifest not support")

func ENvError(o string) error {
	return fmt.Errorf("%w: %s", ErrEnv, o)
}

type Runner struct {
	version string
	tty     bool

	usr                   *user.User
	usrGroup              *user.Group
	homeDir               string
	supplementaryGroupIDs []int
	tempDir               string
	workingDir            string
	env                   map[string]string
	volumes               []volume
}

type volume struct {
	kind   string
	source string
	target string
}

func NewRunner() *Runner {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error getting user home directory: %v", err)
	}

	return &Runner{
		tty:                   terminal.IsTerminal(int(os.Stdout.Fd())),
		homeDir:               home,
		supplementaryGroupIDs: []int{},
		env:                   map[string]string{},
	}
}

func (r *Runner) setUserAndGroups() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	r.usr = usr

	// If LookupGroupId fails (e.g., when Active Directory is used on
	// macOS), we still proceed. This is a best effort attempt. If the
	// group name cannot be retrieved, we'll use the name "konvoy" in
	// the container.
	usrGroup, err := user.LookupGroupId(r.usr.Gid)
	if err == nil {
		r.usrGroup = usrGroup
	}

	// Find supplementary groups IDs to set for the container.
	r.supplementaryGroupIDs = append(r.supplementaryGroupIDs, 0)

	// Add `docker` to the supplementary group of the container.
	// Ignore this step if `docker` group does not exist.
	dockerGroup, err := user.LookupGroup("docker")
	if err == nil {
		gid, err := strconv.Atoi(dockerGroup.Gid)
		if err != nil {
			return ENvError(fmt.Sprintf("docker gid '%s' is not an int", dockerGroup.Gid))
		}
		r.supplementaryGroupIDs = append(r.supplementaryGroupIDs, gid)
	}

	return nil
}

func (r *Runner) setAWSEnv() error {
	awsEnvVars := []string{
		envAWSDefaultRegion,
		envAWSRegion,
		envAWSProfile,
		envAWSAccessKeyID,
		envAWSSecretAccessKey,
		envAWSSessionToken,
		envAWSSTSRegionalEndpoints,
		envAWSSDKLoadConfig,
		envAWSCABundle,
	}

	for _, env := range awsEnvVars {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}

	// the homedir is already mounted, only mount if not the default
	credentialsFile, found := os.LookupEnv(envAWSCredentialsFile)
	if found {
		if fi, err := os.Stat(credentialsFile); err == nil {
			if fi.IsDir() {
				return ENvError(fmt.Sprintf("env %s must be set to a file", envAWSCredentialsFile))
			}
			r.addBindVolume(
				credentialsFile,
				filepath.Join(r.homeDir, ".aws", "credentials"),
			)
		} else if !os.IsNotExist(err) {
			return errors2.Wrap(err, fmt.Sprintf("could not determine if %q exists", credentialsFile))
		}
	}

	// find AWS_CA_BUNDLE and mount the file into the container if exists
	// https://docs.aws.amazon.com/cli/latest/topic/config-vars.html#general-options
	caBundle, err := awsCABundleFromEnv()
	if err != nil {
		return err
	}
	if caBundle != "" {
		r.addBindVolume(caBundle, caBundle)
	}

	return nil
}

func (r *Runner) setAzureEnv() {
	azureEnvVars := []string{
		envAzureLocation,
		app.AzureClientIDEnvVariable,
		app.AzureClientSecretEnvVariable,
		app.AzureSubscriptionIDEnvVariable,
		app.AzureTenantIDEnvVariable,
	}

	for _, env := range azureEnvVars {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}
}

func (r *Runner) setVSphereEnv() {
	for _, env := range []string{
		envVSphereServer,
		envVSphereUser,
		envVSpherePassword,
		envRedHatSubscriptionManagerUser,
		envRedHatSubscriptionManagerPassword,
	} {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}
}

func (r *Runner) setGCPEnv() error {
	path, found := os.LookupEnv(envGCPApplicationCredentials)
	if !found {
		// return early if env is not set, may be using a different provider
		return nil
	}
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() {
			return fmt.Errorf("env %s must be set to a file", envGCPApplicationCredentials)
		}
		r.env[envGCPApplicationCredentials] = path
		// bind to exact same path
		r.addBindVolume(path, path)
	} else if os.IsNotExist(err) {
		return fmt.Errorf("env %s ise set, but file %s does not exist", envGCPApplicationCredentials, path)
	} else {
		return fmt.Errorf("could not determine if %q exists: %w", path, err)
	}

	return nil
}

// awsCABundle will return the path to a custom AWS CA bundle from AWS_CA_BUNDLE env var.
func awsCABundleFromEnv() (string, error) {
	caBundle := os.Getenv(envAWSCABundle)
	if caBundle == "" {
		return caBundle, nil
	}

	fi, err := os.Stat(caBundle)
	if err != nil {
		return "", errors2.Wrap(err, "could not determine if file exists")
	}

	if fi.IsDir() {
		return "", ENvError(fmt.Sprintf("env %s must be set to a file", envAWSCABundle))
	}

	caBundle, err = filepath.Abs(caBundle)
	if err != nil {
		return "", errors2.Wrap(err, fmt.Sprintf("could not determine the absolute path for %q", caBundle))
	}

	return caBundle, nil
}

func (r *Runner) addBindVolume(source, target string) {
	r.volumes = append(r.volumes, volume{
		kind:   "bind",
		source: source,
		target: target,
	})
}

func (r *Runner) setHTTPProxyEnv() error {
	for _, env := range []string{envHTTPSProxy, envHTTPProxy, envNoProxy} {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}

	return nil
}

func (r *Runner) setAnsibleHostKeyChecking() {
	r.env["ANSIBLE_HOST_KEY_CHECKING"] = "false"
}

func (r *Runner) setupSSHAgent() {
	value, found := os.LookupEnv("SSH_AUTH_SOCK")
	if found {
		r.env["SSH_AUTH_SOCK"] = value
		r.addBindVolume(value, value)
	}
	value, found = os.LookupEnv("SSH_AGENT_PID")
	if found {
		r.env["SSH_AGENT_PID"] = value
	}
}

func (r *Runner) dockerRun(args []string) error {
	//nolint:gosec // running docker is inherently insecure
	cmd := exec.Command(
		"docker", "run",
		"--interactive",
		"--tty="+strconv.FormatBool(r.tty),
		"--rm",
		"--net=host",
		"-w", containerWorkingDir,
	)

	r.addBindVolume(r.workingDir, containerWorkingDir)

	if runtime.GOOS != "windows" {
		cmd.Args = append(cmd.Args, "-u", r.usr.Uid+":"+r.usr.Gid)

		r.addBindVolume(r.tempDir+"/passwd", "/etc/passwd")
		r.addBindVolume(r.tempDir+"/group", "/etc/group")
		r.addBindVolume(r.homeDir, r.homeDir)
		r.addBindVolume(r.tempDir+"/ssh_known_hosts", r.usr.HomeDir+"/.ssh/known_hosts")
	}

	for _, gid := range r.supplementaryGroupIDs {
		cmd.Args = append(cmd.Args, "--group-add", strconv.Itoa(gid))
	}

	// append custom envs after OS environment variables,to preserve existing behavior
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	cmd.Env = os.Environ()
	// Iterate through the environment variables map to add them as docker '-e ENV_VAR' arguments
	// and make sure they will be passed when invoking the command.
	for k, v := range r.env {
		cmd.Args = append(cmd.Args, "-e", k)
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	for _, v := range r.volumes {
		cmd.Args = append(cmd.Args, "--mount", "type="+v.kind+",source="+v.source+",target="+v.target)
	}

	cmd.Args = append(cmd.Args, image.Tag())
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// checkDockerVersion checks whether the docker version is greater than then
// minimum required.
func (r *Runner) checkDockerVersion() error {
	_, err := exec.Command("docker", "version", "-f", "{{.Client.Version}}").Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(ee.Stderr)
		}
		return err
	}
	return nil
}

// checkDockerRunning checks whether the docker daemon is running.
func (r *Runner) checkDockerRunning() error {
	out, err := exec.Command("docker", "info", "-f", "{{json .ServerVersion}}").Output()
	if err != nil || len(out) == 0 {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Stderr.Write(ee.Stderr)
		}
		return err
	}
	return nil
}

func (r *Runner) checkRequirements() error {
	err := r.checkDockerVersion()
	if err != nil {
		return err
	}

	err = r.checkDockerRunning()
	if err != nil {
		return err
	}
	return nil
}

func (r *Runner) setUserMapping() error {
	filePath := filepath.Join(r.tempDir, "passwd")
	content := fmt.Sprintf(
		"%s:!:%s:%s::%s:/bin/sh\n",
		r.usr.Username,
		r.usr.Uid,
		r.usr.Gid,
		r.homeDir,
	)
	//nolint:gosec // file must be world readable
	return ioutil.WriteFile(filePath, []byte(content), 0o644)
}

func (r *Runner) setGroupMapping() error {
	filePath := filepath.Join(r.tempDir, "group")

	// Default to "konvoy" as the group name in the container.
	groupName := "konvoy"
	if r.usrGroup != nil {
		groupName = r.usrGroup.Name
	}

	content := fmt.Sprintf(
		"%s::%s:\n",
		groupName,
		r.usr.Gid,
	)
	//nolint:gosec // file must be world readable
	return ioutil.WriteFile(filePath, []byte(content), 0o644)
}

func (r *Runner) Run(args []string) error {
	// Get the Konvoy image version for marker file
	var err error
	r.version = ""

	err = r.checkRequirements()
	if err != nil {
		return err
	}

	// Lookup for current user and its group ID.
	// This also look for supplementary group IDs to set.
	err = r.setUserAndGroups()
	if err != nil {
		return err
	}

	// Current dir on the host is used for working dir in the container,
	// which will become the default cluster name.
	r.workingDir, err = os.Getwd()
	if err != nil {
		return err
	}

	// Create a temporary dir to hold some files that need to be mounted to the container,
	// eg. /etc/passwd, /etc/group, etc.
	r.tempDir, err = ioutil.TempDir(r.workingDir, ".konvoy-image-tmp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(r.tempDir)

	// Setup the user and group mappings in the container so that uid and
	// gid on the host can be properly resolved in the container too.
	err = r.setUserMapping()
	if err != nil {
		return err
	}
	err = r.setGroupMapping()
	if err != nil {
		return err
	}

	// Mask the ssh config from the host. The ssh config format on OSX is
	// slightly different than that in Linux, which will cause Ansible to
	// fail sometimes.
	if _, err = os.Stat(r.usr.HomeDir + "/.ssh/config"); err == nil {
		f, ferr := os.Create(filepath.Join(r.tempDir, "ssh_config"))
		if ferr != nil {
			return ferr
		}

		// Make sure that that the file is only `rw` by the user.
		if ferr := f.Chmod(os.FileMode(0o600)); ferr != nil {
			return ferr
		}
		f.Close()

		r.addBindVolume(r.tempDir+"/ssh_config", r.usr.HomeDir+"/.ssh/config")
	}

	// Mask the ssh known_hosts file from the host.
	// This will prevent multiple runs from interfering with each other when targeting hosts with the same IPs.
	f, err := os.Create(filepath.Join(r.tempDir, "ssh_known_hosts"))
	if err != nil {
		return err
	}
	f.Close()

	err = r.setAWSEnv()
	if err != nil {
		return err
	}
	r.setAzureEnv()
	r.setVSphereEnv()
	err = r.setGCPEnv()
	if err != nil {
		return err
	}

	err = r.setHTTPProxyEnv()
	if err != nil {
		return err
	}

	err = image.LoadImage()
	if err != nil {
		return err
	}

	r.setAnsibleHostKeyChecking()
	r.setupSSHAgent()
	// Run the command in the konvoy docker container.
	return r.dockerRun(args)
}
