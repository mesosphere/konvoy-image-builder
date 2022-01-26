package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapMerge(t *testing.T) {
	m1 := map[string]interface{}{
		"one": "abc",
		"two": "xyz",
		"three": map[string]interface{}{
			"ten": 1234,
		},
	}
	m2 := map[string]interface{}{
		"one": "ABC",
		"three": map[string]interface{}{
			"ten":    12345,
			"eleven": "done",
		},
	}

	m3 := map[string]interface{}{
		"four": 4,
	}

	fmt.Printf("m1: %v\n", m1)

	if err := mergeMapsOverwrite(m1, m2, m3); err != nil {
		t.Fatalf("error merging: %v", err)
	}

	fmt.Printf("m1: %v\n", m1)
}

func TestEnrichKubernetesFullVersion(t *testing.T) {
	cases := []struct {
		name                         string
		config                       map[string]interface{}
		userDefinedKubernetesVersion string
		expectedErr                  error
		expectedK8sVersion           string
	}{
		{
			"KubernetesVersion only in config",
			map[string]interface{}{
				kubernetesVersionKey: interface{}("1.21.6"),
			},
			"",
			nil,
			"1.21.6",
		},
		{
			"KubernetesVersion only provided by user",
			map[string]interface{}{
				"some-unrelated-config": interface{}(1),
			},
			"1.21.6",
			nil,
			"1.21.6",
		},
		{
			"KubernetesVersion both in config and provided by user",
			map[string]interface{}{
				kubernetesVersionKey: interface{}("1.21.4"),
			},
			"1.21.6",
			nil,
			"1.21.6",
		},
		{
			"KubernetesVersion not at all provided",
			map[string]interface{}{
				"not-a-kubernetes-version": interface{}(1),
				"this-isnt-either":         interface{}(2),
			},
			"",
			ErrKubernetesVersionMissing,
			"1.21.6",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testErr := enrichKubernetesFullVersion(testCase.config, testCase.userDefinedKubernetesVersion)
			assert.ErrorIs(t, testErr, testCase.expectedErr, "Expected error: %s, but got: %s", testCase.expectedErr, testErr)
			if testErr == nil {
				assert.Contains(t, testCase.config, kubernetesFullVersionKey, "Expected key: %s to be in resulting config", kubernetesFullVersionKey)
				assert.Equal(t, testCase.expectedK8sVersion, getString(testCase.config, kubernetesFullVersionKey),
					"Expected version: %s but got %s", testCase.expectedK8sVersion, getString(testCase.config, kubernetesFullVersionKey))
			}
		})
	}
}

// TODO: tests
//func TestInitConfig(t *testing.T) {
//	type testCase struct {
//		image     string
//		overrides []string
//	}
//	bc := Builder{
//		CommonConfigPath: filepath.Join("testdata", "common.yaml"),
//		Image:            filepath.Join("testdata", "arch.yaml"),
//	}
//
//	if err := bc.InitConfig(); err != nil {
//		t.Fatalf("error running initConfig: %v", err)
//	}
//}
