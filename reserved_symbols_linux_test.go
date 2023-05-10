// See LICENSE.txt for copyright and licensing information about this file.

//go:build linux

package glob

import (
	"testing"
)

func TestHasReservedSymbol(t *testing.T) {
	testIO := []struct {
		name     string
		value    []rune
		expected bool
		index    int
	}{
		{
			name:     "null is reserved",
			value:    []rune{rune(0x00)},
			expected: true,
			index:    0,
		},
	}

	for _, test := range testIO {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := hasReservedSymbol(test.value)

			if !actual {
				t.Errorf("Expected: %t. Actual: %t.", true, actual)
			}
		})
	}
}
