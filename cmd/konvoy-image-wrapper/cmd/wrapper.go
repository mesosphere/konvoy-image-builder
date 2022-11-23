package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
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
	envVSphereSSHUserName                = "SSH_USERNAME"
	envVSphereSSHPassword                = "SSH_PASSWORD"
	envVsphereSSHPrivatekeyFile          = "SSH_PRIVATE_KEY_FILE"

	//nolint:gosec // environment var set by user
	envGCPApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"

	envHTTPSProxy = "HTTPS_PROXY"
	envHTTPProxy  = "HTTP_PROXY"
	envNoProxy    = "NO_PROXY"

	containerWorkingDir = "/tmp/kib"
	windows             = "windows"
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
	kind        string
	source      string
	target      string
	bindOptions []string
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
	}

	for _, env := range awsEnvVars {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}

	_, found := os.LookupEnv(envAWSCredentialsFile)
	if !found {
		// fall-back to default location for aws credentials file
		os.Setenv(envAWSCredentialsFile, filepath.Join(r.usr.HomeDir, ".aws", "credentials"))
	}
	if err := r.mountFileEnv(envAWSCredentialsFile, filepath.Join(r.homeDir, ".aws", "credentials")); err != nil {
		return fmt.Errorf("unable to mount AWS credenttials file: %w", err)
	}

	if err := r.mountFileEnv(envAWSCABundle, ""); err != nil {
		return fmt.Errorf("unable to mount AWS CA bundle: %w", err)
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

func (r *Runner) setVSphereEnv() error {
	for _, env := range []string{
		envVSphereServer,
		envVSphereUser,
		envVSpherePassword,
		envRedHatSubscriptionManagerUser,
		envRedHatSubscriptionManagerPassword,
		envVSphereSSHUserName,
		envVSphereSSHPassword,
	} {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}

	if err := r.mountFileEnv(envVsphereSSHPrivatekeyFile, ""); err != nil {
		return fmt.Errorf("unable to mount ssh private key file for vsphere provider: %w", err)
	}
	return nil
}

func (r *Runner) setGCPEnv() error {
	// mount same path as the host path set by GOOGLE_APPLICATION_CREDENTIALS environment variable.
	return r.mountFileEnv(envGCPApplicationCredentials, "")
}

// mountFileEnv mounts absolute path to the file assigned by the environment variable.
// if the path to mount to the container is not provided, the absolute path of the host will be mounted to container.
func (r *Runner) mountFileEnv(envName string, containerPath string) error {
	filePath, found := os.LookupEnv(envName)
	if !found {
		return nil
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Ignore if file does not exists
			// The wrapper is common for all the provider so validation of the provider specific env variable should be moved to its subcommand
			return nil
		}
		return fmt.Errorf("error accessing file %q assigned to %s environment variable %w", filePath, envName, err)
	}

	if fi.IsDir() {
		return fmt.Errorf("env %s must be set to a file", envName)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("could not determine the absolute path for %s assigned to %s environment variable : %w", filePath, envName, err)
	}

	if containerPath == "" {
		containerPath = absFilePath
	}

	r.env[envName] = containerPath
	r.addBindVolume(absFilePath, containerPath, "readonly")
	return nil
}

func (r *Runner) addBindVolume(source, target string, options ...string) {
	r.volumes = append(r.volumes, volume{
		kind:        "bind",
		source:      source,
		target:      target,
		bindOptions: options,
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

	if runtime.GOOS != windows {
		cmd.Args = append(cmd.Args, "-u", r.usr.Uid+":"+r.usr.Gid)
		r.addBindVolume(r.tempDir, r.homeDir)
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
		bindOptions := ""
		if len(v.bindOptions) > 0 {
			bindOptions = fmt.Sprintf(",%s", strings.Join(v.bindOptions, ","))
		}
		cmd.Args = append(
			cmd.Args,
			"--mount",
			"type="+v.kind+",source="+v.source+",target="+v.target+bindOptions)
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
	if runtime.GOOS == windows {
		return nil
	}
	filePath := filepath.Join(r.tempDir, "passwd")
	content := fmt.Sprintf(
		"%s:!:%s:%s::%s:/bin/sh\n",
		r.usr.Username,
		r.usr.Uid,
		r.usr.Gid,
		r.homeDir,
	)
	//nolint:gosec // file must be world readable
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return err
	}
	r.addBindVolume(filePath, "/etc/passwd")
	return nil
}

func (r *Runner) setGroupMapping() error {
	if runtime.GOOS == windows {
		return nil
	}
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
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return err
	}
	r.addBindVolume(filePath, "/etc/group")
	return nil
}

// Mask the ssh config from the host. The ssh config format on OSX is
// slightly different than that in Linux, which will cause Ansible to
// fail sometimes.
func (r *Runner) maskSSHConfig() error {
	if runtime.GOOS == windows {
		return nil
	}
	_, err := os.Stat(r.usr.HomeDir + "/.ssh/config")
	if err != nil {
		if os.IsNotExist(err) {
			// ignore if the ~/.ssh/config file is not found
			return nil
		}
		return err
	}

	f, ferr := os.Create(filepath.Join(r.tempDir, "ssh_config"))
	if ferr != nil {
		return ferr
	}
	defer f.Close()

	// Make sure that that the file is only `rw` by the user.
	if ferr := f.Chmod(os.FileMode(0o600)); ferr != nil {
		return ferr
	}

	r.addBindVolume(r.tempDir+"/ssh_config", r.usr.HomeDir+"/.ssh/config")

	return nil
}

// Mask the ssh known_hosts file from the host.
// This will prevent multiple runs from interfering with each other when targeting hosts with the same IPs.
func (r *Runner) maskSSHKnownHosts() error {
	if runtime.GOOS == windows {
		return nil
	}
	f, err := os.Create(filepath.Join(r.tempDir, "ssh_known_hosts"))
	if err != nil {
		return err
	}
	f.Close()
	r.addBindVolume(f.Name(), r.usr.HomeDir+"/.ssh/known_hosts")
	return nil
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
	r.addBindVolume(r.workingDir, containerWorkingDir)

	// Create a temporary dir to hold some files that need to be mounted to the container,
	// eg. /etc/passwd, /etc/group, etc.
	r.tempDir, err = os.MkdirTemp(r.workingDir, ".konvoy-image-tmp")
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

	err = r.maskSSHConfig()
	if err != nil {
		return err
	}

	err = r.maskSSHKnownHosts()
	if err != nil {
		return err
	}

	err = r.setAWSEnv()
	if err != nil {
		return err
	}
	r.setAzureEnv()
	err = r.setVSphereEnv()
	if err != nil {
		return err
	}
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
