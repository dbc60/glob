// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"testing"
)

// Verify nextPattern is working correctly.
func TestNextPattern(t *testing.T) {
	testIO := []struct {
		name         string
		pattern      []rune
		expectedHead []rune
		expectedType patternType
		tailLength   int
	}{
		{
			name:         "empty",
			pattern:      []rune(""),
			expectedHead: []rune(""),
			expectedType: patternSimple,
		},
		{
			name:         "simple",
			pattern:      []rune("file"),
			expectedHead: []rune("file"),
			expectedType: patternSimple,
		},
		{
			name:         "class",
			pattern:      []rune("[fi\\]le]"),
			expectedHead: []rune("[fi\\]le]"),
			expectedType: patternSimple,
		},
		{
			name:         "directory",
			pattern:      []rune("*file"),
			expectedHead: []rune("file"),
			expectedType: patternDirectory,
		},
		{
			name:         "recursive",
			pattern:      []rune("**file"),
			expectedHead: []rune("file"),
			expectedType: patternRecursive,
		},
		{
			name:         "simple followed by directory",
			pattern:      []rune("file*more"),
			expectedHead: []rune("file"),
			expectedType: patternSimple,
			tailLength:   5,
		},
		{
			name:         "simple followed by recursive",
			pattern:      []rune("file**more"),
			expectedHead: []rune("file"),
			expectedType: patternSimple,
			tailLength:   6,
		},
		{
			name:         "directory followed by directory",
			pattern:      []rune("*file*more"),
			expectedHead: []rune("file"),
			expectedType: patternDirectory,
			tailLength:   5,
		},
		{
			name:         "directory followed by recursive",
			pattern:      []rune("*file**more"),
			expectedHead: []rune("file"),
			expectedType: patternDirectory,
			tailLength:   6,
		},
		{
			name:         "recursive followed by directory",
			pattern:      []rune("**file*more"),
			expectedHead: []rune("file"),
			expectedType: patternRecursive,
			tailLength:   5,
		},
		{
			name:         "recursive followed by recursive",
			pattern:      []rune("**file**more"),
			expectedHead: []rune("file"),
			expectedType: patternRecursive,
			tailLength:   6,
		},
		{
			name:         "recursive glob with 1 extra asterisk",
			pattern:      []rune("***file"),
			expectedHead: []rune("file"),
			expectedType: patternRecursive,
		},
		{
			name:         "recursive glob with 2 extra asterisks",
			pattern:      []rune("****file"),
			expectedHead: []rune("file"),
			expectedType: patternRecursive,
		},
		{
			name:         "recursive glob with many extra asterisks",
			pattern:      []rune("*********file"),
			expectedHead: []rune("file"),
			expectedType: patternRecursive,
		},
		{
			name:         "directory separator",
			pattern:      []rune("*\\" + SeparatorString),
			expectedHead: []rune("\\" + SeparatorString),
			expectedType: patternDirectory,
		},
		{
			name:         "recursive separator",
			pattern:      []rune("**\\" + SeparatorString),
			expectedHead: []rune("\\" + SeparatorString),
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
