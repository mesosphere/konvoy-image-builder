//

package app_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mitchellh/pointerstructure"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

func TestConfigGet(t *testing.T) {
	config := app.Config{
		"string1": "abc",
		"string2": "xyz",
		"map": map[string]interface{}{
			"ten": 1234,
		},
	}

	assert.Equal(t, "abc", config.Get("string1"))
	assert.Equal(t, "xyz", config.Get("string2"))
	assert.Equal(t, "", config.Get("string3"))
	assert.Equal(t, "", config.Get(""))
}

func TestConfigGetWithError(t *testing.T) {
	config := app.Config{
		"string1": "abc",
		"string2": "xyz",
		"map": map[string]interface{}{
			"ten": 1234,
		},
	}

	value, err := config.GetWithError("string1")
	if assert.NoError(t, err) {
		assert.Equal(t, "abc", value)
	}

	value, err = config.GetWithError("string2")
	if assert.NoError(t, err) {
		assert.Equal(t, "xyz", value)
	}

	_, err = config.GetWithError("string3")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, pointerstructure.ErrNotFound)
	}

	_, err = config.GetWithError("map")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, app.ErrPathNotString)
	}
}

func TestConfigGetSlicetWithError(t *testing.T) {
	config := app.Config{
		"map": map[string]interface{}{
			"ten": 1234,
		},
		"slice": []string{
			"one",
			"two",
		},

		"interface": []interface{}{
			"three",
			"four",
		},
	}

	value, err := config.GetSliceWithError("slice")
	if assert.NoError(t, err) {
		assert.Equal(t, []string{"one", "two"}, value)
	}

	value, err = config.GetSliceWithError("interface")
	if assert.NoError(t, err) {
		assert.Equal(t, []string{"three", "four"}, value)
	}

	_, err = config.GetSliceWithError("slice1")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, pointerstructure.ErrNotFound)
	}

	_, err = config.GetSliceWithError("map")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, app.ErrPathNotSlice)
	}
}

func TestConfigSet(t *testing.T) {
	config := app.Config{
		"string1": "abc",
		"string2": "xyz",
	}

	assert.NoError(t, config.Set("string1", ""))

	value, err := config.GetWithError("string1")
	if assert.NoError(t, err) {
		assert.Equal(t, "abc", value)
	}

	assert.NoError(t, config.Set("string3", "nmo"))

	value, err = config.GetWithError("string3")
	if assert.NoError(t, err) {
		assert.Equal(t, "nmo", value)
	}

	assert.NoError(t, config.Set("string1", "jkl"))

	value, err = config.GetWithError("string1")
	if assert.NoError(t, err) {
		assert.Equal(t, "jkl", value)
	}

	err = config.Set("strings/error", "jkl")
	assert.Error(t, err)
	assert.ErrorIs(t, err, pointerstructure.ErrNotFound)
}

func TestConfigDelete(t *testing.T) {
	config := app.Config{
		"string1": "abc",
		"string2": "xyz",
	}

	assert.NoError(t, config.Delete("string1"))

	_, err := config.GetWithError("string1")
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, pointerstructure.ErrNotFound)
	}

	assert.Error(t, config.Delete("strings/error"))
	assert.ErrorIs(t, err, pointerstructure.ErrNotFound)
}

func TestBuildName(t *testing.T) {
	config := app.Config{}

	assert.NoError(t, config.Set(app.BuildNameKey, ""))
	assert.Equal(t, app.DefaultBuildName, app.BuildName(config))

	buildNameKey := gofakeit.Word()

	assert.NoError(t, config.Set(app.BuildNameKey, buildNameKey))
	assert.Equal(t, buildNameKey, app.BuildName(config))

	buildNameExtra := gofakeit.Word()
	assert.NoError(t, config.Set(app.BuildNameExtraKey, buildNameExtra))
	assert.Equal(t, fmt.Sprintf("%s%s", buildNameKey, buildNameExtra), app.BuildName(config))
}

type PackerStringer struct {
	Value string
}

func (stringer PackerStringer) String() string {
	return stringer.Value
}

type PackerDefault struct {
	Value string
}

func TestGenPackerVars(t *testing.T) {
	extraVarsPath := gofakeit.Word()
	k8sFullVersion := gofakeit.Word()

	packerBytes := []byte(gofakeit.Word())
	packerDefault := PackerDefault{
		Value: gofakeit.Word(),
	}
	packerString := gofakeit.Word()
	packerStringer := PackerStringer{
		Value: gofakeit.Word(),
	}

	config := app.Config{
		"packer": map[interface{}]interface{}{
			"string":   packerString,
			"bytes":    packerBytes,
			"stringer": packerStringer,
			"nil":      nil,
			"default":  packerDefault,
		},
		app.KubernetesFullVersionKey: k8sFullVersion,
		"gpu": map[interface{}]interface{}{
			"types": []string{"nvidia"},
		},
		"nvidia_driver_version": "1234",
	}

	jsonVars, err := app.GenPackerVars(config, extraVarsPath)
	assert.NoError(t, err)
	var packer map[string]interface{}

	err = json.Unmarshal(jsonVars, &packer)
	assert.NoError(t, err)
	if assert.Contains(t, packer, "string") {
		assert.Equal(t, packerString, packer["string"])
	}

	if assert.Contains(t, packer, "bytes") {
		assert.Equal(t, string(packerBytes), packer["bytes"])
	}

	if assert.Contains(t, packer, "stringer") {
		assert.Equal(t, packerStringer.Value, packer["stringer"])
	}

	if assert.Contains(t, packer, "nil") {
		assert.Equal(t, "", packer["nil"])
	}

	if assert.Contains(t, packer, "default") {
		assert.Equal(t, fmt.Sprintf("{%s}", packerDefault.Value), packer["default"])
	}

	if assert.Contains(t, packer, app.KubernetesFullVersionKey) {
		assert.Equal(t, k8sFullVersion, packer[app.KubernetesFullVersionKey])
	}

	if assert.Contains(t, packer, "gpu") {
		assert.Equal(t, "true", packer["gpu"])
	}

	if assert.Contains(t, packer, "gpu_types") {
		assert.Contains(t, "nvidia", packer["gpu_types"])
	}

	if assert.Contains(t, packer, "gpu_nvidia_version") {
		assert.Equal(t, "1234", packer["gpu_nvidia_version"])
	}

	if assert.Contains(t, packer, app.AnsibleExtraVarsKey) {
		assert.Equal(t, fmt.Sprintf("@%s", extraVarsPath), packer[app.AnsibleExtraVarsKey])
	}
}

func TestMergeUserArgs(t *testing.T) {
	config := app.Config{
		"packer": map[interface{}]interface{}{},
	}
	userArgs := app.UserArgs{
		ClusterArgs: app.ClusterArgs{
			ContainerdVersion: gofakeit.Word(),
			KubernetesVersion: gofakeit.Word(),
		},
	}

	assert.NoError(t, app.MergeUserArgs(config, userArgs))
	if assert.Contains(t, config, app.ContainerdVersionKey) {
		assert.Equal(t, userArgs.ContainerdVersion, config[app.ContainerdVersionKey])
	}

	if assert.Contains(t, config, app.KubernetesVersionKey) {
		assert.Equal(t, userArgs.KubernetesVersion, config[app.KubernetesVersionKey])
	}
}

func TestMergeAmazonUserArgs(t *testing.T) {
	config := app.Config{
		"packer": map[interface{}]interface{}{},
	}

	// NOTE(jkoelker) Test with `SourceAMI == ""` first so when set we can test that the filters
	//                are removed.
	amazonArgs := &app.AmazonArgs{
		SourceAMI:        "",
		AMIFilterName:    gofakeit.Word(),
		AMIFilterOwner:   gofakeit.Word(),
		AWSBuilderRegion: gofakeit.Word(),
		AMIRegions: []string{
			gofakeit.Word(),
		},
		AWSInstanceType: gofakeit.Word(),
		AMIUsers: []string{
			gofakeit.Word(),
		},
		AMIGroups: []string{
			gofakeit.Word(),
		},
	}

	err := app.MergeAmazonUserArgs(config, amazonArgs)
	assert.NoError(t, err)

	value, err := config.GetWithError(app.PackerBuilderRegionPath)
	assert.NoError(t, err)
	assert.Equal(t, amazonArgs.AWSBuilderRegion, value)

	_, err = config.GetWithError(app.PackerSourceAMIPath)
	assert.Error(t, err)

	value, err = config.GetWithError(app.PackerFilterNamePath)
	assert.NoError(t, err)
	assert.Equal(t, amazonArgs.AMIFilterName, value)

	value, err = config.GetWithError(app.PackerFilterOwnerPath)
	assert.NoError(t, err)
	assert.Equal(t, amazonArgs.AMIFilterOwner, value)

	value, err = config.GetWithError(app.PackerAMIRegionsPath)
	assert.NoError(t, err)
	assert.Equal(t, strings.Join(amazonArgs.AMIRegions, ","), value)

	value, err = config.GetWithError(app.PackerAMIUsersPath)
	assert.NoError(t, err)
	assert.Equal(t, strings.Join(amazonArgs.AMIUsers, ","), value)

	value, err = config.GetWithError(app.PackerAMIGroupsPath)
	assert.NoError(t, err)
	assert.Equal(t, strings.Join(amazonArgs.AMIGroups, ","), value)

	amazonArgs.SourceAMI = gofakeit.Word()

	err = app.MergeAmazonUserArgs(config, amazonArgs)
	assert.NoError(t, err)

	_, err = config.GetWithError(app.PackerFilterNamePath)
	assert.Error(t, err)

	_, err = config.GetWithError(app.PackerFilterOwnerPath)
	assert.Error(t, err)

	value, err = config.GetWithError(app.PackerSourceAMIPath)
	assert.NoError(t, err)
	assert.Equal(t, amazonArgs.SourceAMI, value)
}

func TestMergeAzureUserArgs(t *testing.T) {
	config := app.Config{
		"packer":                     map[interface{}]interface{}{},
		app.KubernetesFullVersionKey: "1.24.6",
	}

	azureArgs := &app.AzureArgs{
		Location: gofakeit.Word(),
	}

	azureArgs.CloudEndpoint = &app.AzureCloudFlag{
		Endpoint: app.AzureCloudEndpointPublic,
	}

	if err := app.MergeAzureUserArgs(config, azureArgs); assert.NoError(t, err) {
		if value, err := config.GetSliceWithError(
			app.PackerAzureGalleryLocations,
		); assert.NoError(t, err) {
			assert.Equal(t, []string{azureArgs.Location}, value)
		}
	}

	azureArgs.GalleryImageLocations = []string{
		gofakeit.Word(),
		gofakeit.Word(),
		gofakeit.Word(),
	}

	azureArgs.GalleryImageName = gofakeit.Word()
	azureArgs.GalleryImageOffer = gofakeit.Word()
	azureArgs.GalleryImagePublisher = gofakeit.Word()
	azureArgs.GalleryImageSKU = gofakeit.Word()
	azureArgs.GalleryName = gofakeit.Word()
	azureArgs.ResourceGroupName = gofakeit.Word()
	azureArgs.SubscriptionID = gofakeit.Word()

	if err := app.MergeAzureUserArgs(config, azureArgs); assert.NoError(t, err) {
		if value, err := config.GetSliceWithError(
			app.PackerAzureGalleryLocations,
		); assert.NoError(t, err) {
			assert.Equal(t, azureArgs.GalleryImageLocations, value)
		}

		value, err := config.GetWithError(app.PackerAzureGalleryImageNamePath)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.GalleryImageName, value)

		value, err = config.GetWithError(app.PackerAzureGalleryImageOfferPath)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.GalleryImageOffer, value)

		value, err = config.GetWithError(app.PackerAzureGalleryImagePublisherPath)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.GalleryImagePublisher, value)

		value, err = config.GetWithError(app.PackerAzureGalleryImageSKU)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.GalleryImageSKU, value)

		value, err = config.GetWithError(app.PackerAzureLocation)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.Location, value)

		value, err = config.GetWithError(app.PackerAzureGalleryNamePath)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.GalleryName, value)

		value, err = config.GetWithError(app.PackerAzureResourceGroupNamePath)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.ResourceGroupName, value)

		value, err = config.GetWithError(app.PackerAzureSubscriptionIDPath)
		assert.NoError(t, err)
		assert.Equal(t, azureArgs.SubscriptionID, value)
	}
}

func TestMergeAzureGalleryImageName(t *testing.T) {
	azureArgs := &app.AzureArgs{
		Location: gofakeit.Word(),
		CloudEndpoint: &app.AzureCloudFlag{
			Endpoint: app.AzureCloudEndpointPublic,
		},
	}

	testcase := []struct {
		fullK8sVersion                string
		wantFullK8sVersionInImageName string
	}{
		{"1.24.6", "1.24.6"},
		{"1.24.6+fips.0", "1.24.6-fips.0"},
		{"1.24.6-nvidia", "1.24.6-nvidia"},
	}
	for _, tc := range testcase {
		config := app.Config{
			"packer":                     map[interface{}]interface{}{},
			app.KubernetesFullVersionKey: tc.fullK8sVersion,
		}
		if err := app.MergeAzureUserArgs(config, azureArgs); assert.NoError(t, err) {
			value, err := config.GetWithError(app.PackerAzureGalleryImageNamePath)
			assert.NoError(t, err)
			assert.Contains(t, value, tc.wantFullK8sVersionInImageName)
		}
	}
}

func TestConfigEnrichKubernetesFullVersion(t *testing.T) {
	config := app.Config{}

	assert.ErrorIs(
		t,
		app.EnrichKubernetesFullVersion(config, ""),
		app.ErrKubernetesVersionMissing,
	)

	k8sVersion := gofakeit.AppVersion()

	assert.NoError(t, config.Set(app.KubernetesVersionKey, k8sVersion))

	assert.NoError(t, app.EnrichKubernetesFullVersion(config, ""))
	assert.Equal(t, k8sVersion, config.Get(app.KubernetesFullVersionKey))

	userVersion := gofakeit.AppVersion()

	assert.NoError(t, app.EnrichKubernetesFullVersion(config, userVersion))
	assert.Equal(t, userVersion, config.Get(app.KubernetesFullVersionKey))

	metadata := gofakeit.AppVersion()

	assert.NoError(t, config.Set(app.KubernetesBuildMetadataKey, metadata))

	assert.NoError(t, app.EnrichKubernetesFullVersion(config, userVersion))
	assert.Equal(
		t,
		fmt.Sprintf("%s+%s", userVersion, metadata),
		config.Get(app.KubernetesFullVersionKey),
	)

	assert.NoError(t, app.EnrichKubernetesFullVersion(config, ""))
	assert.Equal(
		t,
		fmt.Sprintf("%s+%s", k8sVersion, metadata),
		config.Get(app.KubernetesFullVersionKey),
	)
}
