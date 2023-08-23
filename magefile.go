//go:build mage
// +build mage

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"gopkg.in/yaml.v2"
)

var (
	wrapperCmd      = "bin/konvoy-image-wrapper"
	baseURL         = "https://downloads.d2iq.com/dkp"
	containerdURL   = "https://packages.d2iq.com/dkp/containerd"
	nvidiaURL       = "https://download.nvidia.com/XFree86/Linux-x86_64"
	overrideDirName = "overrides"
)

var (
	basic         = "basic"
	fips          = "fips"
	nvidia        = "nvidia"
	offline       = "offline"
	offlineFIPS   = "offline-fips"
	offlineNvidia = "offline-nvidia"
	rhck          = "rhck"

	validOS = []string{
		"centos 7.9",
		"redhat 7.9",
		"redhat 8.4",
		"redhat 8.6",
		"sles 15",
		"oracle 7.9",
		"flatcar",
		"ubuntu 18.04",
		"ubuntu 20.04",
		"rocky 9.0",
		"rocky 9.1",
	}

	validBuildConfig = []string{
		basic,
		fips,
		nvidia,
		offline,
		offlineFIPS,
		offlineNvidia,
		rhck,
	}

	aws   = "aws"
	azure = "azure"
	ova   = "ova"
	gcp   = "gcp"

	validInfra = []string{
		aws,
		azure,
		ova,
		gcp,
	}
	oldImageFormatVersion = 24
)

var (
	dryRunFlag         = "--dry-run"
	p2Instance         = "p2.xlarge"
	g4dnInstance       = "g4dn.2xlarge"
	azureGPUInstance   = "Standard_NC6s_v3"
	azureBuildInstance = "Standard_B2ms"
)

func BuildWrapper() error {
	return sh.RunV("make", "build-wrapper")
}

// Runs E2e for images.
func RunE2e(buildOS, buildConfig, buildInfra string, dryRun bool) error {
	mg.Deps(BuildWrapper)
	err := validateBuildOS(buildOS)
	if err != nil {
		return err
	}
	err = validateBuildConfig(buildConfig)
	if err != nil {
		return err
	}
	err = validateBuildInfra(buildInfra)
	if err != nil {
		return err
	}
	buildPath := getBuildPath(buildOS, buildInfra)
	featureOverrides := getOverridesFromBuildConfig(buildConfig)
	overrideFlagForCmd := make([]string, 0, len(featureOverrides))
	for _, override := range featureOverrides {
		fullOverride := fmt.Sprintf("--overrides=%s/%s", overrideDirName, override)
		overrideFlagForCmd = append(overrideFlagForCmd, fullOverride)
	}

	// we need these extra overrides always for ova
	if buildConfig == offline || buildConfig == offlineNvidia || buildConfig == offlineFIPS {
		infraOverride := getInfraOverride(buildInfra)
		fullOverride := fmt.Sprintf("--overrides=%s", infraOverride)
		overrideFlagForCmd = append(overrideFlagForCmd, fullOverride)
		fmt.Printf("making infraOverride %s \n", infraOverride)
		// TODO: @faiq - move this to mage
		if err := sh.RunV("make", infraOverride); err != nil {
			return fmt.Errorf("failed to create offline infra with override %s %v", infraOverride, err)
		}

		defer func() {
			// TODO: @faiq - move this to mage
			if buildInfra == aws {
				if err := sh.RunV("make", "infra.aws.destroy"); err != nil {
					fmt.Printf("failed to delete offline infra %v\n", err)
				}
			}
		}()
		// we need to fetch the proper os-bundle
		// pip packages
		// image bundle
		// nvidia
		// containerd
		// TODO: @faiq - move this to mage
		if err := sh.RunV("make", "pip-packages-artifacts"); err != nil {
			return fmt.Errorf("failed to download pip packages %v", err)
		}
		kubeVersion, err := getKubernetesVerisonForBuild()
		if err != nil {
			return fmt.Errorf("failed to read kubernetes version %w", err)
		}

		// Fetch artifacts
		isFips := buildConfig == offlineFIPS
		if err := os.MkdirAll(path.Join("artifacts", "images"), 0775); err != nil {
			return fmt.Errorf("failed to create artifacts/images err %w", err)
		}
		if err := fetchOSBundle(buildOS, kubeVersion, isFips); err != nil {
			return fmt.Errorf("failed to fetch OS bundle %w", err)
		}
		if err := fetchImageBundle(kubeVersion, isFips); err != nil {
			return fmt.Errorf("failed to fetch Image bundle %w", err)
		}
		if err := fetchContainerd(buildOS, isFips); err != nil {
			return fmt.Errorf("failed to fetch containerd %w", err)
		}
		if buildConfig == offlineNvidia {
			if err := fetchNvidiaRunFile(); err != nil {
				return fmt.Errorf("failed to fetch nvidiaRunFile %w", err)
			}
		}
	}

	if buildOS == "flatcar" && buildInfra == ova {
		flatcarOverride := "packer-ova-flatcar-override.yaml"
		flatcarOverrideFlag := fmt.Sprintf("--overrides=%s", flatcarOverride)
		overrideFlagForCmd = append(overrideFlagForCmd, flatcarOverrideFlag)
		fmt.Printf("making flatcar override %s \n", flatcarOverride)
		// TODO: @faiq - move this to mage
		if err := sh.RunV("make", flatcarOverride); err != nil {
			return fmt.Errorf("failed to override for flatcar ova %s %v", flatcarOverride, err)
		}
	}

	args := []string{"build"}
	buildCmd := buildInfra
	if buildCmd == ova {
		buildCmd = "vsphere"
	}
	args = append(args, buildCmd)
	args = append(args, buildPath)
	// skip creating image
	if dryRun {
		args = append(args, dryRunFlag)
	} else {
		releaseOverrideFile, err := getReleaseOverride(buildConfig)
		if err != nil {
			return fmt.Errorf("failed to create release override file: %w", err)
		}
		overrideFlagForCmd = append(
			overrideFlagForCmd,
			fmt.Sprintf("--overrides=%s/%s", overrideDirName, releaseOverrideFile))
	}
	args = append(args, overrideFlagForCmd...)

	vmMachine := getVMForBuild(buildInfra, buildConfig)
	if vmMachine != "" {
		args = append(args, fmt.Sprintf("--instance-type=%s", vmMachine))
	}
	// extra args for gcp
	if buildInfra == gcp {
		args = append(args, "--project-id=eng-ksphere-platform-e2e", "--network=kib-ssh-network")
	}
	fmt.Printf("Running %s with args %v\n", wrapperCmd, args)
	return sh.RunV(wrapperCmd, args...)
}

// Clean up after yourself.
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("artifacts")
}

func validateBuildOS(buildOS string) error {
	if buildOS == "" {
		return fmt.Errorf("no buildOS found using %s", buildOS)
	}
	found := false
	for _, valid := range validOS {
		if buildOS == valid {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("buildOS %s is invalid must be one of %v", buildOS, validOS)
	}
	return nil
}

func validateBuildConfig(buildConfig string) error {
	if buildConfig == "" {
		return fmt.Errorf("no buildConfig found using %s", buildConfig)
	}
	found := false
	for _, valid := range validBuildConfig {
		if buildConfig == valid {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("buildConfig %s is invalid must be one of %v", buildConfig, validBuildConfig)
	}
	return nil
}

func validateBuildInfra(buildInfra string) error {
	if buildInfra == "" {
		return fmt.Errorf("no buildInfra found using %s", buildInfra)
	}
	found := false
	for _, valid := range validInfra {
		if buildInfra == valid {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("buildInfra %s is invalid must be one of %v", buildInfra, validInfra)
	}
	return nil
}

func getBuildPath(buildOS, buildInfra string) string {
	formattedOS := strings.TrimLeft(buildOS, " ")
	formattedOS = strings.TrimPrefix(formattedOS, " ")
	formattedOS = strings.Replace(formattedOS, " ", "-", 1)
	formattedOS = strings.Replace(formattedOS, ".", "", 1)

	formattedOS = strings.Replace(formattedOS, "redhat", "rhel", 1)
	fileForOS := fmt.Sprintf("%s.yaml", formattedOS)

	infraDirForImage := buildInfra
	// we change aws for ami
	if infraDirForImage == "aws" {
		infraDirForImage = "ami"
	}
	return path.Join("images", infraDirForImage, fileForOS)
}

func getOverridesFromBuildConfig(buildConfig string) []string {
	switch buildConfig {
	case basic:
		return nil
	case fips:
		return []string{"fips.yaml"}
	case nvidia:
		return []string{"nvidia.yaml"}
	case offline:
		return []string{"offline.yaml"}
	case offlineFIPS:
		return []string{"offline-fips.yaml", "fips.yaml"}
	case offlineNvidia:
		return []string{"offline.yaml", "offline-nvidia.yaml"}
	case rhck:
		return []string{"rhck.yaml"}
	}
	return nil
}

// getReleaseOverride creates temporary release override file
// The release override file contains metadata about build name that
// get appended to the final image artifact.
// The konvoy e2e tests uses the released image by
// locating the image with kubernetes version and build metadata in the image name
func getReleaseOverride(buildConfig string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error finding current directory: %w", err)
	}
	buildNameExtra := "-release"
	switch buildConfig {
	case fips, offlineFIPS:
		// "-fips-release"
		buildNameExtra = fmt.Sprintf("-fips%s", buildNameExtra)
	case nvidia, offlineNvidia:
		// "-nvidia-release"
		buildNameExtra = fmt.Sprintf("-nvidia%s", buildNameExtra)
	}

	releaseFile := "release.yaml"
	releaseOverride := filepath.Join(currentDir, overrideDirName, releaseFile)
	content := []byte(fmt.Sprintf("---\nbuild_name_extra: %s\n", buildNameExtra))
	if err := os.WriteFile(releaseOverride, content, 0o644); err != nil {
		return "", fmt.Errorf("error creating release override file: %w", err)
	}
	return releaseFile, nil
}

func getInfraOverride(buildInfra string) string {
	baseOfflineTemplate := "packer-%s-offline-override.yaml"
	switch buildInfra {
	case aws:
		return fmt.Sprintf(baseOfflineTemplate, aws)
	case ova:
		return fmt.Sprintf(baseOfflineTemplate, ova)
	}
	return ""
}

func getKubernetesVerisonForBuild() (string, error) {
	bytes, err := os.ReadFile(path.Join("images", "common.yaml"))
	if err != nil {
		return "", err
	}
	var config map[string]interface{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return "", err
	}
	return config["kubernetes_version"].(string), nil
}

func fetchOSBundle(osName, kubernetesVersion string, fips bool) error {
	fetchClient := http.DefaultClient
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse url %w", err)
	}
	osInfo := strings.Split(osName, " ")
	osDist := osInfo[0]
	osMajor := strings.Split(osInfo[1], ".")[0]

	airgappedBundlePath := fmt.Sprintf("%s_%s_%s_x86_64", kubernetesVersion, osDist, osMajor)
	if fips {
		airgappedBundlePath = airgappedBundlePath + "_fips"
	}
	airgappedBundlePath = airgappedBundlePath + ".tar.gz"

	u.Path = path.Join(u.Path,
		"airgapped",
		"os-packages",
		airgappedBundlePath,
	)
	fmt.Println("Downloading artifact from ", u.String())
	resp, err := fetchClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
	})
	if err != nil {
		return fmt.Errorf("failed to download os bundle %w", err)
	}
	defer resp.Body.Close()
	outFile := path.Join("artifacts", airgappedBundlePath)
	out, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fetchImageBundle(kubernetesVersion string, fips bool) error {
	fetchClient := http.DefaultClient
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse url %w", err)
	}
	imageBundleName := fmt.Sprintf("kubernetes-images-%s-d2iq.1", kubernetesVersion)
	v := semver.New(kubernetesVersion)
	ext := ".tar"
	if v.Minor < int64(oldImageFormatVersion) {
		imageBundleName = fmt.Sprintf("%s_images", kubernetesVersion)
		ext = ".tar.gz"
	}
	if fips {
		imageBundleName = imageBundleName + "-fips"
	}
	imageBundleName = imageBundleName + ext
	u.Path = path.Join(u.Path, "airgapped",
		"kubernetes-images",
		imageBundleName)
	fmt.Println("Downloading image bundle from ", u.String())
	resp, err := fetchClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
	})
	if err != nil {
		return fmt.Errorf("failed to get image bundle %w", err)
	}
	defer resp.Body.Close()
	outFile := path.Join("artifacts", "images", imageBundleName)
	out, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fetchContainerd(osName string, fips bool) error {
	bytes, err := os.ReadFile(path.Join("ansible", "group_vars", "all", "defaults.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read file %w", err)
	}
	var config map[string]interface{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return fmt.Errorf("failed to unmarshal yaml %w", err)
	}
	containerdVersion, ok := config["containerd_version"].(string)
	if !ok {
		return fmt.Errorf("could not parse containerd version in bytes %v", bytes)
	}
	fetchClient := http.DefaultClient
	u, err := url.Parse(containerdURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL %w", err)
	}
	osInfo := strings.Split(osName, " ")
	osDist := osInfo[0]
	// TODO: improve this
	osMajorMinor := strings.Split(osInfo[1], ".")
	osMajor := osMajorMinor[0]
	osMinor := osMajorMinor[1]
	osDist = strings.Replace(osDist, "redhat", "rhel", 1)
	containerdPath := fmt.Sprintf("containerd-%s-d2iq.1-%s-%s.%s-x86_64", containerdVersion, osDist, osMajor, osMinor)
	if fips {
		containerdPath = containerdPath + "_fips"
	}
	containerdPath = containerdPath + ".tar.gz"
	u.Path = path.Join(u.Path, containerdPath)
	fmt.Println("fetching assets from ", u.String())
	resp, err := fetchClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
	})
	if err != nil {
		return fmt.Errorf("failed to get containerd %w", err)
	}
	defer resp.Body.Close()
	outFile := path.Join("artifacts", containerdPath)
	out, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fetchNvidiaRunFile() error {
	bytes, err := os.ReadFile(path.Join("ansible", "group_vars", "all", "defaults.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read file %w", err)
	}
	var config map[string]interface{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return fmt.Errorf("failed to unmarshal yaml %w", err)
	}
	nvidiaRunfileVersion, ok := config["nvidia_driver_version"].(string)
	if !ok {
		return fmt.Errorf("could not parse nvidia_driver_version version in bytes %v", bytes)
	}
	fetchClient := http.DefaultClient
	u, err := url.Parse(nvidiaURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL %w", err)
	}
	runFile := fmt.Sprintf("NVIDIA-Linux-x86_64-%s.run", nvidiaRunfileVersion)
	u.Path = path.Join(u.Path, nvidiaRunfileVersion, runFile)
	fmt.Println("Downloading artifact from ", u.String())
	resp, err := fetchClient.Do(&http.Request{
		URL:    u,
		Method: http.MethodGet,
	})
	if err != nil {
		return fmt.Errorf("failed to download runfile %w", err)
	}
	defer resp.Body.Close()
	outFile := path.Join("artifacts", runFile)
	out, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func getVMForBuild(buildInfra, buildConfig string) string {
	switch buildInfra {
	case aws:
		switch buildConfig {
		case offlineNvidia:
			return p2Instance
		case nvidia:
			return g4dnInstance
		default:
			return ""
		}
	case azure:
		switch buildConfig {
		case nvidia:
			return azureGPUInstance
		default:
			return azureBuildInstance

		}
	default:
		return ""
	}
}
