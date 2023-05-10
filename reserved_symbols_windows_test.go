// See LICENSE.txt for copyright and licensing information about this file.

//go:build windows

package glob

import (
	"testing"
)

// Verify hasReservedSymbol and isReservedSymbol are working.
func TestHasReservedSymbol(t *testing.T) {
	testIO := []struct {
		name     string
		value    []rune
		expected bool
	}{
		{
			name:     "null rune is reserved",
			value:    []rune{0x00},
			expected: true,
		},
		{
			name:     "double quote is reserved",
			value:    []rune{'"'},
			expected: true,
		},
		{
			name:     "colon is reserved",
			value:    []rune{':'},
			expected: true,
		},
		{
			name:     "less than is reserved",
			value:    []rune{'<'},
			expected: true,
		},
		{
			name:     "greater than is reserved",
			value:    []rune{'>'},
			expected: true,
		},
		{
			name:     "vertical bar is reserved",
			value:    []rune{'|'},
			expected: true,
		},
		{
			name:     "empty string contains no reserved characters",
			value:    []rune{},
			expected: false,
		},
		{
			name:     "string contains no reserved symbols",
			value:    []rune("Hello, 世界"),
			expected: false,
		},
		{
			name:     "embedded reserved symbols",
			value:    []rune("Hello, 世界: 1 < 2"),
			expected: true,
		},
	}

	for _, test := range testIO {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := hasReservedSymbol(test.value)

			if actual != test.expected {
				t.Errorf("Expected: %t. Actual: %t.", true, actual)
			}
		})
	}

	t.Run("ascii control codes are reserved", func(t *testing.T) {
		for i := 0; i < 0x20; i++ {
			value := rune(i)
			if has := isReservedSymbol(value); !has {
				t.Errorf("For \"%d\" expected true. Actual: false", i)
			}
		}
	})

	t.Run("ascii del is reserved", func(t *testing.T) {
		value := rune(0x7F)
		if has := isReservedSymbol(value); !has {
			t.Error("Expected true. Actual false.")
		}
	})
}
