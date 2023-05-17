package cmd

import (
	"bytes"
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

	envVSphereServer     = "VSPHERE_SERVER"
	envVSphereUser       = "VSPHERE_USERNAME"
	envVSpherePassword   = "VSPHERE_PASSWORD"
	envVSphereDatacenter = "VSPHERE_DATACENTER"
	envVsphereDatastore  = "VSPHERE_DATASTORE"

	envRedHatSubscriptionManagerUser          = "RHSM_USER"
	envRedHatSubscriptionManagerPassword      = "RHSM_PASS"
	envRedHatSubscriptionManagerActivationKey = "RHSM_ACTIVATION_KEY"
	envRedHatSubscriptionManagerOrgID         = "RHSM_ORG_ID"

	envVSphereSSHUserName       = "SSH_USERNAME"
	envVSphereSSHPassword       = "SSH_PASSWORD"
	envVsphereSSHPrivatekeyFile = "SSH_PRIVATE_KEY_FILE"

	//nolint:gosec // environment var set by user
	envGCPApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"

	envHTTPSProxy = "HTTPS_PROXY"
	envHTTPProxy  = "HTTP_PROXY"
	envNoProxy    = "NO_PROXY"

	containerWorkingDir   = "/tmp/kib"
	windows               = "windows"
	containerEngineEnv    = "KIB_CONTAINER_ENGINE"
	containerEngineDocker = "docker"
	containerEnginePodman = "podman"
)

var ErrEnv = errors.New("manifest not support")

func EnvError(o string) error {
	return fmt.Errorf("%w: %s", ErrEnv, o)
}

type Runner struct {
	usr                   *user.User
	usrGroup              *user.Group
	homeDir               string
	supplementaryGroupIDs []int
	tempDir               string
	workingDir            string
	env                   map[string]string
	volumes               []volume
	containerEngine       string
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
	var containerEngine string
	switch p := os.Getenv(containerEngineEnv); p {
	case "":
		containerEngine = detectContainerEngine()
	case "podman":
		containerEngine = containerEnginePodman
	case "docker":
		containerEngine = containerEngineDocker
	default:
		log.Printf("ignoring unknown value %q for %s", p, containerEngineEnv)
	}

	return &Runner{
		homeDir:               home,
		supplementaryGroupIDs: []int{},
		env:                   map[string]string{},
		containerEngine:       containerEngine,
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
			return EnvError(fmt.Sprintf("docker gid '%s' is not an int", dockerGroup.Gid))
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
		envVSphereDatacenter,
		envVsphereDatastore,
		envRedHatSubscriptionManagerUser,
		envRedHatSubscriptionManagerPassword,
		envRedHatSubscriptionManagerActivationKey,
		envRedHatSubscriptionManagerOrgID,
		envVSphereSSHUserName,
		envVSphereSSHPassword,
		envVsphereSSHPrivatekeyFile,
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

func (r *Runner) setHTTPProxyEnv() {
	for _, env := range []string{envHTTPSProxy, envHTTPProxy, envNoProxy} {
		value, found := os.LookupEnv(env)
		if found {
			r.env[env] = value
		}
	}
}

func (r *Runner) setAnsibleHostKeyChecking() {
	r.env["ANSIBLE_HOST_KEY_CHECKING"] = "false"
}

func (r *Runner) setupSSHAgent() {
	value, found := os.LookupEnv("SSH_AUTH_SOCK")
	if found {
		r.env["SSH_AUTH_SOCK"] = value
		r.addBindVolume(value, value, "readonly")
	}
	value, found = os.LookupEnv("SSH_AGENT_PID")
	if found {
		r.env["SSH_AGENT_PID"] = value
	}
}

func (r *Runner) dockerRun(args []string) error {
	//nolint:gosec // we validate this
	cmd := exec.Command(
		r.containerEngine, "run",
		"--interactive",
		"--tty=false",
		"--rm",
		"--net=host",
		"-w", containerWorkingDir,
	)
	if r.containerEngine == containerEnginePodman {
		cmd.Args = append(cmd.Args, "--userns=keep-id", "--security-opt", "label=disable")
	}

	if runtime.GOOS != windows && r.containerEngine == containerEngineDocker {
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd failed %w", err)
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
	// Lookup for current user and its group ID.
	// This also look for supplementary group IDs to set.
	err = r.setUserAndGroups()
	if err != nil {
		return fmt.Errorf("failed to set groups %w", err)
	}

	// Current dir on the host is used for working dir in the container,
	// which will become the default cluster name.
	r.workingDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get wd %w", err)
	}
	r.addBindVolume(r.workingDir, containerWorkingDir)

	// Create a temporary dir to hold some files that need to be mounted to the container,
	// eg. /etc/passwd, /etc/group, etc.
	r.tempDir, err = os.MkdirTemp(r.workingDir, ".konvoy-image-tmp")
	if err != nil {
		return fmt.Errorf("failed to make temp %w", err)
	}
	defer os.RemoveAll(r.tempDir)

	// Setup the user and group mappings in the container so that uid and
	// gid on the host can be properly resolved in the container too.
	if r.containerEngine == containerEngineDocker {
		err = r.setUserMapping()
		if err != nil {
			return fmt.Errorf("failed to set user mapping %w", err)
		}
		err = r.setGroupMapping()
		if err != nil {
			return fmt.Errorf("failed to set group mapping %w", err)
		}
	}

	err = r.maskSSHConfig()
	if err != nil {
		return fmt.Errorf("failed to set mask ssh config %w", err)
	}

	err = r.maskSSHKnownHosts()
	if err != nil {
		return fmt.Errorf("failed to set mask ssh known_hosts %w", err)
	}

	err = r.setAWSEnv()
	if err != nil {
		return fmt.Errorf("failed to set aws env %w", err)
	}
	r.setAzureEnv()
	err = r.setVSphereEnv()
	if err != nil {
		return fmt.Errorf("failed to set vshpere env %w", err)
	}
	err = r.setGCPEnv()
	if err != nil {
		return fmt.Errorf("failed to set gcp env %w", err)
	}

	r.setHTTPProxyEnv()

	err = image.LoadImage(r.containerEngine)
	if err != nil {
		return fmt.Errorf("failed to load image %w", err)
	}

	r.setAnsibleHostKeyChecking()
	r.setupSSHAgent()
	// Run the command in the konvoy docker container.
	return r.dockerRun(args)
}

// detectContainerEngine determines which container engine should be used.
// if both docker and podman installed then docker engine takes precedence
// if none of them are detected then fallback to existing behavior of using
// docker as a container engine.
func detectContainerEngine() string {
	if isDockerAvailable() {
		return containerEngineDocker
	} else if isPodmanAvailable() {
		return containerEnginePodman
	}
	// fall back to current behavior for backward compatibility
	return containerEngineDocker
}

func isPodmanAvailable() bool {
	cmd := exec.Command("podman", "-v")
	var buff bytes.Buffer
	cmd.Stdout = &buff
	err := cmd.Run()
	if err != nil {
		return false
	}
	return strings.HasPrefix(buff.String(), "podman version")
}

func isDockerAvailable() bool {
	cmd := exec.Command("docker", "-v")
	var buff bytes.Buffer
	cmd.Stdout = &buff
	err := cmd.Run()
	if err != nil {
		return false
	}
	return strings.HasPrefix(buff.String(), "Docker version")
}
