// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"testing"
)

// Verify escaped backslashes in glob patterns match directory separators in file paths.
func TestNextPatternWindows(t *testing.T) {
	testIO := []struct {
		name         string
		pattern      []rune
		expectedHead []rune
		expectedType patternType
		tailLength   int
	}{
		{
			name:         "directory separator",
			pattern:      []rune("*\\\\"),
			expectedHead: []rune("\\\\"),
			expectedType: patternDirectory,
		},
		{
			name:         "recursive separator",
			pattern:      []rune("**\\\\"),
			expectedHead: []rune("\\\\"),
			expectedType: patternRecursive,
		},
	}

	for _, test := range testIO {
		t.Run(test.name, func(t *testing.T) {
			head, tail, kind := nextPattern(test.pattern)

			if len(head) != len(test.expectedHead) {
				t.Errorf("Test %s: Expected \"%s\". Actual \"%s\". Pattern lengths do not match.", test.name, string(test.expectedHead), string(head))
			}

			if len(tail) != test.tailLength {
				t.Errorf("Test %s: Expected \"%s\"(%d). Actual \"%s\"(%d). Next pattern lengths do not match.", test.name, string(test.pattern[len(head):]), test.tailLength, string(tail), len(tail))
			}

			if kind != test.expectedType {
				t.Errorf("Test %s Expected %s. Actual %s. Pattern types do not match.", test.name, test.expectedType, kind)
			}

			if string(head) != string(test.expectedHead) {
				t.Errorf("Test %s: Expected %s. Actual %s.", test.name, string(test.expectedHead), string(head))
			}
		})
	}
}
