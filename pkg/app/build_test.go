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
