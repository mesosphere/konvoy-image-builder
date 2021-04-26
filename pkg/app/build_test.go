package app

import (
	"fmt"
	"testing"
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
