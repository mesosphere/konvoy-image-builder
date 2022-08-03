package app_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	"github.com/mesosphere/konvoy-image-builder/pkg/app"
)

func TestMapMerge(t *testing.T) {
	m1 := map[string]interface{}{
		"one": gofakeit.Word(),
		"two": gofakeit.Word(),
		"three": map[string]interface{}{
			"ten": gofakeit.Int64(),
		},
	}
	m2 := map[string]interface{}{
		"one": gofakeit.Word(),
		"three": map[string]interface{}{
			"ten":    gofakeit.Int64(),
			"eleven": gofakeit.Word(),
		},
	}

	m3 := map[string]interface{}{
		"four": gofakeit.Int64(),
	}

	fmt.Printf("m1: %v\n", m1)

	err := app.MergeMapsOverwrite(m1, m2, m3)
	assert.NoError(t, err)

	if assert.Contains(t, m1, "one") {
		assert.Equal(t, m1["one"], m2["one"])
	}

	assert.Contains(t, m1, "two")

	if assert.Contains(t, m1, "three") {
		if assert.Contains(t, m1["three"], "ten") {
			if assert.Contains(t, m1["three"], "eleven") {
				assert.Equal(t, m1["three"], m2["three"])
			}
		}
	}

	if assert.Contains(t, m1, "four") {
		assert.Equal(t, m1["four"], m3["four"])
	}

}

func TestMapMergeForArtifacts(t *testing.T) {
	cases := []struct {
		caseName                     string
		mapFromOverridesAndExtraVars map[string]interface{}
		passedUserArgs               map[string]interface{}
		expectedOutput               map[string]interface{}
	}{
		{
			caseName: "if there is an exisiting field from the overrides",
			mapFromOverridesAndExtraVars: map[string]interface{}{
				"os_packages_local_bundle_file": "passed-artifacts/centos-7.9.tar.gz",
				"fips":                          true,
			},
			passedUserArgs: map[string]interface{}{
				"os_packages_local_bundle_file":  "artifacts/centos-7.9-fips.tar.gz",
				"containerd_local_bundle_file":   "artifacts/containerd-1.4.13-centos-7.9-d2iq.1.tar.gz",
				"pip_packages_local_bundle_file": "artifacts/pip-packages.tar.gz",
				"images_local_bundle_dir":        "artifacts/images",
			},
			expectedOutput: map[string]interface{}{
				"os_packages_local_bundle_file":  "artifacts/centos-7.9-fips.tar.gz",
				"containerd_local_bundle_file":   "artifacts/containerd-1.4.13-centos-7.9-d2iq.1.tar.gz",
				"pip_packages_local_bundle_file": "artifacts/pip-packages.tar.gz",
				"images_local_bundle_dir":        "artifacts/images",
				"fips":                           true,
			},
		},
		{
			caseName: "if there is no exisiting field from the overrides",
			mapFromOverridesAndExtraVars: map[string]interface{}{
				"fips": false,
			},
			passedUserArgs: map[string]interface{}{
				"os_packages_local_bundle_file":  "artifacts/centos-7.9-fips.tar.gz",
				"containerd_local_bundle_file":   "artifacts/containerd-1.4.13-centos-7.9-d2iq.1.tar.gz",
				"pip_packages_local_bundle_file": "artifacts/pip-packages.tar.gz",
				"images_local_bundle_dir":        "artifacts/images",
			},
			expectedOutput: map[string]interface{}{
				"os_packages_local_bundle_file":  "artifacts/centos-7.9-fips.tar.gz",
				"containerd_local_bundle_file":   "artifacts/containerd-1.4.13-centos-7.9-d2iq.1.tar.gz",
				"pip_packages_local_bundle_file": "artifacts/pip-packages.tar.gz",
				"images_local_bundle_dir":        "artifacts/images",
				"fips":                           false,
			},
		},
	}
	for _, testCase := range cases {
		err := app.MergeMapsOverwrite(testCase.mapFromOverridesAndExtraVars, testCase.passedUserArgs)
		assert.NoError(t, err)
		assert.Equal(t, testCase.mapFromOverridesAndExtraVars, testCase.expectedOutput)
	}
}
