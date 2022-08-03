package app

import (
	"fmt"
	"strings"

	"github.com/imdario/mergo"
)

func mergeUserOverridesToMap(overridesFilesList []string, m map[string]interface{}) error {
	overrides, err := getOverrides(overridesFilesList)
	if err != nil {
		return fmt.Errorf("error getting overrides: %w", err)
	}

	if err = MergeMapsOverwrite(m, overrides...); err != nil {
		return fmt.Errorf("error merging overrides: %w", err)
	}
	return nil
}

// recursively merges maps into orig, orig is modified.
func MergeMapsOverwrite(orig map[string]interface{}, maps ...map[string]interface{}) error {
	for _, m := range maps {
		if err := mergo.Merge(&orig, m, mergo.WithOverride); err != nil {
			return fmt.Errorf("error merging: %w", err)
		}
	}

	return nil
}

func getOverrides(paths []string) ([]map[string]interface{}, error) {
	overrides := make([]map[string]interface{}, 0, len(paths))

	for _, path := range paths {
		data, err := loadYAML(path)
		if err != nil {
			return nil, fmt.Errorf("error loading override: %w", err)
		}

		overrides = append(overrides, data)
	}

	return overrides, nil
}

func addExtraVarsToMap(extraVars []string, m map[string]interface{}) error {
	extraVarSet := make(map[string]interface{})
	for _, extraVars := range extraVars {
		set := strings.Split(extraVars, "=")
		//nolint:gomnd // the code is splitting on the equal
		if len(set) == 2 {
			k := set[0]
			v := set[1]
			extraVarSet[k] = v
		}
	}

	if err := MergeMapsOverwrite(m, extraVarSet); err != nil {
		return fmt.Errorf("error merging overrides: %w", err)
	}
	return nil
}
