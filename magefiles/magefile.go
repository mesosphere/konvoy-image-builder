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
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"gopkg.in/yaml.v2"
)

var (
	buildOSEnvVar            = "BUILD_OS"
	buildConfigurationEnvVar = "BUILD_CONFIG"
	buildInfraEnvVar         = "BUILD_INFRA"
	wrapperCmd               = "bin/konvoy-image-wrapper"
	baseURL                  = "https://downloads.d2iq.com/dkp"
	containerdURL            = "https://packages.d2iq.com/dkp/containerd"
	nvidiaURL                = "https://download.nvidia.com/XFree86/Linux-x86_64"
)

var (
	basic         = "basic"
	fips          = "fips"
	nvidia        = "nvidia"
	offline       = "offline"
	offlineFIPS   = "offline-fips"
	offlineNvidia = "offline-nvidia"

	validOS = []string{
		"centos 7.9",
		"rhel 7.9",
		"rhel 8.4",
		"rhel 8.6",
		"sles 15",
		"oracle 7.9",
		"flatcar",
		"ubuntu 18.04",
		"ubuntu 20.04",
	}

	validBuildConfig = []string{
		basic,
		fips,
		nvidia,
		offline,
		offlineFIPS,
		offlineNvidia,
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

func BuildWrapper() error {
	return sh.RunV("make", "build-wrapper")
}

// Runs E2e for images.
func RunE2e(buildOS, buildConfig, buildInfra string) error {
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
		fullOverride := fmt.Sprintf("--overrides=overrides/%s", override)
		overrideFlagForCmd = append(overrideFlagForCmd, fullOverride)
	}
	if buildConfig == offline || buildConfig == offlineNvidia || buildConfig == offlineFIPS {
		infraOverride := getInfraOverride(buildInfra)
		fullOverride := fmt.Sprintf("--overrides=%s", infraOverride)
		overrideFlagForCmd = append(overrideFlagForCmd, fullOverride)

		//TODO: @faiq - move this to mage
		if err := sh.RunV("make", infraOverride); err != nil {
			return fmt.Errorf("failed to create offline infra with override %s %v", infraOverride, err)
		}
		// we need to fetch the proper os-bundle
		// pip packages
		// image bundle
		// nvidia
		// containerd
		//TODO: @faiq - move this to mage
		if err := sh.RunV("make", "download-pip-packages"); err != nil {
			return fmt.Errorf("failed to download pip packages %v", err)
		}
		kubeVersion, err := getKubernetesVerisonForBuild()
		if err != nil {
			return fmt.Errorf("failed to read kubernetes version %w", err)
		}
		switch buildConfig {
		case offline, offlineNvidia:
			if err := fetchOSBundle(buildOS, kubeVersion, false); err != nil {
				return fmt.Errorf("failed to fetch OS bundle %w", err)
			}
			if err := fetchImageBundle(kubeVersion, false); err != nil {
				return fmt.Errorf("failed to fetch Image bundle %w", err)
			}
			if err := fetchContainerd(buildOS, false); err != nil {
				return fmt.Errorf("failed to fetch containerd %w", err)
			}
			if buildConfig == offlineNvidia {
				if err := fetchNvidiaRunFile(); err != nil {
					return fmt.Errorf("failed to fetch nvidiaRunFile %w", err)
				}
			}
		case offlineFIPS:
			if err := fetchOSBundle(buildOS, kubeVersion, false); err != nil {
				return fmt.Errorf("failed to fetch OS bundle %w", err)
			}
			if err := fetchImageBundle(kubeVersion, false); err != nil {
				return fmt.Errorf("failed to fetch Image bundle %w", err)
			}
			if err := fetchContainerd(buildOS, false); err != nil {
				return fmt.Errorf("failed to fetch containerd %w", err)
			}
		}
	}
	args := []string{"build", buildPath}
	args = append(args, overrideFlagForCmd...)
	return sh.RunV(wrapperCmd, args...)
}

// Clean up after yourself.
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("artifacts")
}

func validateBuildOS(buildOS string) error {
	if buildOS == "" {
		return fmt.Errorf("no buildOS found using %s", buildOSEnvVar)
	}
	found := false
	for _, valid := range validOS {
		if buildOS == valid {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("buildOS %s is invalid must be one of %v", buildOSEnvVar, validOS)
	}
	return nil
}

func validateBuildConfig(buildConfig string) error {
	if buildConfig == "" {
		return fmt.Errorf("no buildConfig found using %s", buildConfigurationEnvVar)
	}
	found := false
	for _, valid := range validBuildConfig {
		if buildConfig == valid {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("buildConfiguration %s is invalid must be one of %v", buildConfigurationEnvVar, validBuildConfig)
	}
	return nil
}

func validateBuildInfra(buildInfra string) error {
	if buildInfra == "" {
		return fmt.Errorf("no buildInfra found using %s", buildInfraEnvVar)
	}
	found := false
	for _, valid := range validInfra {
		if buildInfra == valid {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("buildInfra %s is invalid must be one of %v", buildInfraEnvVar, validInfra)
	}
	return nil

}

func getBuildPath(buildOS, buildInfra string) string {
	formattedOS := strings.TrimLeft(buildOS, " ")
	formattedOS = strings.TrimPrefix(formattedOS, " ")
	formattedOS = strings.Replace(formattedOS, " ", "-", 1)
	formattedOS = strings.Replace(formattedOS, ".", "", 1)

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
		return []string{"fips.yaml", "offline.yaml"}
	case offlineNvidia:
		return []string{"offline.yaml", "offline-nvidia.yaml"}
	}
	return nil
}

func getInfraOverride(buildInfra string) string {
	baseOfflineTemplate := "packer-%s-offline-override.yaml"
	switch buildInfra {
	case aws:
		return fmt.Sprintf(baseOfflineTemplate, aws)
	case ova:
		return fmt.Sprintf(baseOfflineTemplate, "vsphere")
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
	if err := os.MkdirAll("artifacts", 0775); err != nil {
		return fmt.Errorf("failed to make artifacts dir %w", err)
	}
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
		return err
	}
	imageBundleName := fmt.Sprintf("kubernetes-images-%s-d2iq.1", kubernetesVersion)
	v := semver.New(kubernetesVersion)
	ext := ".tar"
	if v.Minor < int64(oldImageFormatVersion) {
		imageBundleName = fmt.Sprintf("%s_images", kubernetesVersion)
		ext = ".tar.gz"
	}
	if fips {
		imageBundleName = imageBundleName + "_fips"
	}
	imageBundleName = imageBundleName + ext
	u.Path = path.Join(u.Path, "airgapped",
		"kubernetes-images",
		imageBundleName)
	fmt.Println("Downloading artifact from ", u.String())
	resp, err := fetchClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	outFile := path.Join("artifacts", "images", imageBundleName)
	if err := os.MkdirAll(path.Join("artifacts", "images"), 0775); err != nil {
		return err
	}
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func fetchContainerd(osName string, fips bool) error {
	bytes, err := os.ReadFile(path.Join("ansible", "group_vars", "all", "defaults.yaml"))
	if err != nil {
		return err
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
		return err
	}
	osInfo := strings.Split(osName, " ")
	osDist := osInfo[0]
	// TODO: improve this
	osMajorMinor := strings.Split(osInfo[1], ".")
	osMajor := osMajorMinor[0]
	osMinor := osMajorMinor[1]
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
	if err := os.MkdirAll("artifacts", 0775); err != nil {
		return fmt.Errorf("failed to create artifacts dir %w", err)
	}
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
	resp, err := fetchClient.Do(&http.Request{})
	if err != nil {
		return fmt.Errorf("failed to download runfile %w", err)
	}
	defer resp.Body.Close()
	outFile := path.Join("artifacts", runFile)
	if err := os.MkdirAll("artifacts", 0775); err != nil {
		return fmt.Errorf("failed to make artifacts dir %w", err)
	}
	out, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
