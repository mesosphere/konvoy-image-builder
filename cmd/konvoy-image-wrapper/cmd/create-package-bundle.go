package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"text/template"

	"github.com/blang/semver"
	"gopkg.in/yaml.v2"

	terminal "golang.org/x/term"
)

const (
	createPackageBundleCmd = "create-package-bundle"
	generatedDirName       = "generated"
)

type OSConfig struct {
	configDir      string
	containerImage string
}

var osToConfig = map[string]OSConfig{
	"centos-7.9": {
		configDir:      "bundles/centos7.9",
		containerImage: "docker.io/mesosphere/centos:7.9.2009.minimal",
	},
	"redhat-7.9": {
		configDir:      "bundles/redhat7.9",
		containerImage: "registry.access.redhat.com/ubi7/ubi:7.9",
	},
	"redhat-8.4": {
		configDir:      "bundles/redhat8.4",
		containerImage: "registry.access.redhat.com/ubi8/ubi:8.4",
	},
	"redhat-8.6": {
		configDir:      "bundles/redhat8.6",
		containerImage: "registry.access.redhat.com/ubi8/ubi:8.6",
	},
	"redhat-8.8": {
		configDir:      "bundles/redhat8.8",
		containerImage: "registry.access.redhat.com/ubi8/ubi:8.8",
	},
	"rocky-9.1": {
		configDir:      "bundles/rocky9.1",
		containerImage: "docker.io/library/rockylinux:9.1",
	},
	"ubuntu-20.04": {
		configDir:      "bundles/ubuntu20.04",
		containerImage: "docker.io/library/ubuntu:20.04",
	},
}

func getKubernetesVerisonFromAnsible() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working dir: %w", err)
	}
	bytes, err := os.ReadFile(path.Join(pwd, "ansible", "group_vars", "all", "defaults.yaml"))
	if err != nil {
		return "", fmt.Errorf("failed to read ansible defaults file %w", err)
	}
	var config map[string]interface{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return "", fmt.Errorf("failed to unmarshal ansible defaults %w", err)
	}
	kubeVersion, ok := config["kubernetes_version"].(string)
	if !ok {
		return "", fmt.Errorf("unable to parse kubernetes_version from ansible defaults")
	}
	_, err = semver.ParseTolerant(kubeVersion)
	if err != nil {
		return "", err
	}
	return kubeVersion, nil
}

func (r *Runner) CreatePackageBundle(args []string) error {
	var (
		osFlag                  string
		kubernetesVersionFlag   string
		fipsFlag                bool
		eusReposFlag            bool
		satelliteFlag           string
		subscriptionManagerFlag bool
		outputDirectoy          string
		containerImage          string
		fetchKernelHeaders      bool
	)
	flagSet := flag.NewFlagSet(createPackageBundleCmd, flag.ExitOnError)
	flagSet.StringVar(
		&osFlag,
		"os",
		"",
		fmt.Sprintf("The target OS you wish to create a package bundle for. Must be one of %v", getKeys(osToConfig)),
	)
	flagSet.StringVar(
		&kubernetesVersionFlag,
		"kubernetes-version",
		"",
		"The version of kubernetes to download packages for.",
	)
	flagSet.BoolVar(
		&fipsFlag,
		"fips",
		false,
		"If the package bundle should include fips packages.",
	)
	flagSet.BoolVar(
		&eusReposFlag,
		"enable-eus-repos",
		false,
		"If enabled fetches packages from EUS repositories when creating RHEL package bundles. Disabled by default.",
	)
	flagSet.StringVar(
		&satelliteFlag,
		"satellite-server-url",
		"",
		//nolint:lll // it is ok to have long help texts
		"If set, registers with and fetches packages from a Red Hat Satellite. All required repositories must be available in the Red Hat Satellite. Example: --satellite-server-url=\"https://satellite.nutanix.sh\"",
	)
	flagSet.BoolVar(
		&subscriptionManagerFlag,
		"skip-subscription-manager",
		false,
		//nolint:lll // it is ok to have long help texts
		"If enabled, skips authenticating with subscription-manager and fetching from the pre-configured official RHEL repositories when creating RHEL package bundles. Disabled by default.",
	)
	flagSet.StringVar(
		&outputDirectoy,
		"output-directory",
		"artifacts",
		"The directory to place the bundle in.",
	)
	flagSet.StringVar(
		&containerImage,
		"container-image",
		"",
		"A container image to use for building the package bundles",
	)
	flagSet.BoolVar(
		&fetchKernelHeaders,
		"fetch-kernel-headers",
		false,
		//nolint:lll // its ok to have long help texts
		"If enabled fetches kernel headers for the target operating system. To modify the version, edit the file at bundles/{OS_NAME}{VERSION}/packages.txt.gotmpl directly eg: bundles/redhat8.8/packages.txt.gotmpl. This is required for operating systems that will use NVIDIA GPU drivers.",
	)
	err := flagSet.Parse(args)
	if err != nil {
		return err
	}
	if osFlag == "" || outputDirectoy == "" {
		return errors.New("--os and --output-directory all must be set")
	}
	image := containerImage
	if containerImage == "" {
		image, err = getContainerImage(osFlag)
		if err != nil {
			return err
		}
	}
	kubernetesVersion := kubernetesVersionFlag
	if kubernetesVersion == "" {
		kubernetesVersion, err = getKubernetesVerisonFromAnsible()
		if err != nil {
			return err
		}
	}
	if eusReposFlag {
		//nolint:goconst // it is ok to not use const here
		r.env["EUS_REPOS"] = "true"
	}
	if satelliteFlag != "" {
		r.env["SATELLITE_SERVER_URL"] = satelliteFlag
	}
	if subscriptionManagerFlag {
		r.env["SKIP_SUBSCRIPTION_MANAGER"] = "true"
	}
	r.setHTTPProxyEnv()
	bundleCmd := "./bundle.sh"
	absPathToOutput := outputDirectoy
	if !path.IsAbs(outputDirectoy) {
		dir := r.workingDir
		absPathToOutput = path.Join(dir, outputDirectoy)
	}
	reposList, err := templateObjects(osFlag, kubernetesVersion, absPathToOutput, fipsFlag, fetchKernelHeaders)
	if err != nil {
		return err
	}
	config, found := osToConfig[osFlag]
	if !found {
		return fmt.Errorf("buildOS %s is invalid must be one of %v", osFlag, getKeys(osToConfig))
	}
	configDir := config.configDir
	dir := r.workingDir
	base := path.Join(dir, configDir)
	return startContainer(r.containerEngine, image, base, bundleCmd, absPathToOutput, reposList, r.env)
}

//nolint:gocyclo,funlen // the function is relatively clear
func templateObjects(targetOS, kubernetesVersion, outputDir string, fips bool, fetchKernelHeaders bool) ([]string, error) {
	config, found := osToConfig[targetOS]
	if !found {
		return nil, fmt.Errorf("buildOS %s is invalid must be one of %v", targetOS, getKeys(osToConfig))
	}
	configDir := config.configDir
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory %w", err)
	}
	base := path.Join(dir, configDir)
	configDirFS := os.DirFS(base)
	l := make([]string, 0)
	generated := path.Join(base, generatedDirName)
	if err = os.MkdirAll(generated, 0o777); err != nil {
		return l, err
	}

	err = fs.WalkDir(configDirFS, ".", func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.Contains(filepath, "repo-templates") {
			newDir := path.Join(base, generatedDirName, "repos")
			if err := os.MkdirAll(newDir, 0o755); err != nil {
				return err
			}
		}

		if strings.Contains(filepath, "repo-templates") && strings.Contains(filepath, ".repo") &&
			!strings.Contains(filepath, "kubernetes.repo.gotmpl") {
			f, err := os.Open(path.Join(base, filepath))
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()
			baseName := path.Base(filepath)
			newFile := path.Join(base, generatedDirName, "repos", baseName)
			out, err := os.Create(newFile)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			_, err = io.Copy(out, f)
			if err != nil {
				return fmt.Errorf("failed to copy contents of repo file: %w", err)
			}
			l = append(l, out.Name())
		}

		//nolint:nestif // this if is not nested
		if strings.Contains(filepath, "kubernetes.repo.gotmpl") {
			kubernetesRepoTmpl, err := os.ReadFile(path.Join(base, filepath))
			if err != nil {
				return fmt.Errorf("failed to read template kubernetes repo file %w", err)
			}
			t, err := template.New("").Parse(string(kubernetesRepoTmpl))
			if err != nil {
				return fmt.Errorf("failed to parse go template: %w", err)
			}
			repoSuffix := "nokmem"
			if fips {
				repoSuffix = "fips"
			}
			templateInput := struct {
				RepoSuffix        string
				KubernetesVersion string
			}{
				RepoSuffix:        repoSuffix,
				KubernetesVersion: kubernetesVersion,
			}
			out, err := os.Create(path.Join(base, generatedDirName, "repos", "kubernetes.repo"))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			err = t.Execute(out, templateInput)
			if err != nil {
				return fmt.Errorf("failed to execute go template: %w", err)
			}
			l = append(l, out.Name())
		}

		if strings.Contains(filepath, "packages.txt.gotmpl") {
			imagesTmpl, err := os.ReadFile(path.Join(base, filepath))
			if err != nil {
				return fmt.Errorf("failed to read template images repo file %w", err)
			}
			t, err := template.New("").Parse(string(imagesTmpl))
			if err != nil {
				return fmt.Errorf("failed to parse go template: %w", err)
			}
			out, err := os.Create(path.Join(base, generatedDirName, "packages.txt"))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			templateInput := struct {
				KubernetesVersion  string
				FetchKernelHeaders bool
			}{
				KubernetesVersion:  kubernetesVersion,
				FetchKernelHeaders: fetchKernelHeaders,
			}
			err = t.Execute(out, templateInput)
			if err != nil {
				return fmt.Errorf("failed to execute go template: %w", err)
			}
		}

		//nolint:nestif // this if is not nested
		if strings.Contains(filepath, "bundle.sh.gotmpl") {
			outputBaseName := "/" + path.Base(outputDir)
			bundleTmpl, err := os.ReadFile(path.Join(base, filepath))
			if err != nil {
				return fmt.Errorf("failed to read template images repo file %w", err)
			}
			t, err := template.New("").Parse(string(bundleTmpl))
			if err != nil {
				return fmt.Errorf("failed to parse go template: %w", err)
			}
			out, err := os.OpenFile(path.Join(base, generatedDirName, "bundle.sh"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o777)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer out.Close()
			fipsSuffix := ""
			if fips {
				fipsSuffix = "_fips"
			}
			var criToolsVersion string
			var kubernetesMajorMinorVersion string
			if targetOS == "ubuntu-20.04" {
				kubernetesMajorMinorVersionNoV := strings.Join(strings.Split(kubernetesVersion, ".")[0:2], ".")
				// according to ansible.group_vars/all/defaults.yaml L18 this is always k8s major minor.0
				criToolsVersion = fmt.Sprintf("%s.0", kubernetesMajorMinorVersionNoV)
				kubernetesMajorMinorVersion = fmt.Sprintf("v%s", kubernetesMajorMinorVersionNoV)
			}
			templateInput := struct {
				KubernetesVersion string
				OutputDirectory   string
				FipsSuffix        string
				// these two get used by ubuntu
				CRIToolsVersion             string
				KubernetesMajorMinorVersion string
			}{
				KubernetesVersion:           kubernetesVersion,
				OutputDirectory:             outputBaseName,
				FipsSuffix:                  fipsSuffix,
				CRIToolsVersion:             criToolsVersion,
				KubernetesMajorMinorVersion: kubernetesMajorMinorVersion,
			}
			err = t.Execute(out, templateInput)
			if err != nil {
				return fmt.Errorf("failed to execute go template: %w", err)
			}
		}
		return nil
	})
	return l, err
}

func getKeys(m map[string]OSConfig) []string {
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func getContainerImage(targetOS string) (string, error) {
	config, found := osToConfig[targetOS]
	if !found {
		return "", fmt.Errorf("buildOS %s is invalid must be one of %v", targetOS, getKeys(osToConfig))
	}
	return config.containerImage, nil
}

func startContainer(containerEngine, containerImage,
	workingDir, runCmd, outputDir string,
	reposList []string, envs map[string]string,
) error {
	tty := terminal.IsTerminal(int(os.Stdout.Fd()))
	outputBaseName := path.Base(outputDir)
	//nolint:gosec // the input is sanatized and contained.
	cmd := exec.Command(
		containerEngine, "run",
		"--interactive",
		"--tty="+strconv.FormatBool(tty),
		"--rm",
		"-v", fmt.Sprintf("%s:/%s", outputDir, outputBaseName),
		"-v", fmt.Sprintf("%s:%s", workingDir, containerWorkingDir),
		"-w", fmt.Sprintf("%s/%s", containerWorkingDir, generatedDirName),
	)
	for _, repoFullPath := range reposList {
		repo := path.Base(repoFullPath)
		cmd.Args = append(
			cmd.Args,
			"-v",
			fmt.Sprintf("%s:%s/%s", repoFullPath, "/etc/yum.repos.d", repo),
		)
	}
	for k, v := range envs {
		cmd.Args = append(cmd.Args, "-e", k)
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Args = append(cmd.Args, []string{"--entrypoint", "/bin/sh", containerImage, "-c", runCmd}...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(c)
	go func() {
		for sig := range c {
			if signalErr := cmd.Process.Signal(sig); signalErr != nil {
				fmt.Fprintf(cmd.Stderr, "failed to relay signal %s %v\n", sig.String(), signalErr)
			}
		}
	}()
	fmt.Fprintf(os.Stdout, "running: %s \n", strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	return nil
}
