// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"errors"
	"fmt"
	"testing"
)

// Verify getClass is working correctly
func TestGetClass(t *testing.T) {
	testIO := []struct {
		class    []rune
		expected []rune
		err      error
		negated  bool
	}{
		{
			class:    []rune{},
			expected: []rune{},
			err:      ErrGlobZeroLength,
		},
		{
			class:    []rune("["),
			expected: []rune{},
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune("]"),
			expected: []rune("]"),
			err:      ErrGlobNoLeftBracket,
		},
		{
			class:    []rune(`[\]`),
			expected: []rune(`[\]`),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune(`[abc\]`),
			expected: []rune(`[abc\]`),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune(`[abc\]xyz`),
			expected: []rune(`[abc\]xyz`),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune("abc]"),
			expected: []rune{'a'},
			err:      ErrGlobNoLeftBracket,
		},
		{
			class:    []rune("[abc"),
			expected: []rune("[abc"),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune{'[', escapeCharacter, ']'},
			expected: []rune{'[', escapeCharacter, ']'},
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune(`[abc\]`),
			expected: []rune(`[abc\]`),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune("[]"),
			expected: []rune("["),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune("[]]"),
			expected: []rune("[]]"),
		},
		{
			class:    []rune("[.0-\\]"),
			expected: []rune("[.0-\\]"),
			err:      ErrGlobInvalidRange,
		},
		{
			class:    []rune("[[]"),
			expected: []rune("[[]"),
		},
		{
			class:    []rune{'[', escapeCharacter, 'x', ']'},
			expected: []rune("["),
			err:      ErrGlobInvalidEscape,
		},
		{
			class:    []rune("[\\x]"),
			expected: []rune("["),
			err:      ErrGlobInvalidEscape,
		},
		{
			class:    []rune("[az\\"),
			expected: []rune("["),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune("[][!]"),
			expected: []rune("[][!]"),
		},
		{
			class:    []rune(`[\]]`),
			expected: []rune(`[\]]`),
		},
		{
			class:    []rune(`[bc]`),
			expected: []rune(`[bc]`),
		},
		{
			class:    []rune(`[\!abc\\xyz]`),
			expected: []rune(`[\!abc\\xyz]`),
		},
		{
			class:    []rune("[abc]"),
			expected: []rune("[abc]"),
		},
		{
			class:    []rune("[-]"),
			expected: []rune("[-]"),
		},
		{
			class:    []rune("[-x]"),
			expected: []rune("[-x]"),
		},
		{
			class:    []rune("[.-0]"),
			expected: []rune(""),
			err:      ErrGlobInvalidRange,
		},
		{
			class:    []rune("[!.-0]"),
			expected: []rune(""),
			err:      ErrGlobInvalidRange,
		},
		{
			class:    []rune("[x-]"),
			expected: []rune("[x-]"),
		},
		{
			class:    []rune(`[0\-9]`),
			expected: []rune(`[0\-9]`),
		},
		{
			class:    []rune("[0-9]"),
			expected: []rune("[0-9]"),
		},
		{
			class:    []rune(`[!]`),
			expected: []rune(`[!`),
			err:      ErrGlobTruncated,
		},
		{
			class:    []rune(`[![]`),
			expected: []rune(`[![]`),
			negated:  true,
		},
		{
			class:    []rune(`[![\]]`),
			expected: []rune(`[![\]]`),
			negated:  true,
		},
		{
			class:    []rune("[!abc]"),
			expected: []rune("[!abc]"),
			negated:  true,
		},
		{
			class:    []rune("[-Helo, 世界! abc?*0-9()[\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"),
			expected: []rune("[-Helo, 世界! abc?*0-9()[\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"),
			negated:  false,
		},
		{
			class:    []rune("[!-Helo, 世界! abc?*0-9()[\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"),
			expected: []rune("[!-Helo, 世界! abc?*0-9()[\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"),
			negated:  true,
		},
		{
			class:    []rune("[-Helo, 世界! abc?*0-9()[\\\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"),
			expected: []rune("[-Helo, 世界! abc?*0-9()[\\\\]"),
			negated:  false,
		},
		{
			class:    []rune("[!-Helo, 世界! abc?*0-9()[\\\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"),
			expected: []rune("[!-Helo, 世界! abc?*0-9()[\\\\]"),
			negated:  true,
		},
		{
			class:    []rune{'[', '!', GlobSeparator, ']'},
			expected: []rune{'[', '!', GlobSeparator, ']'},
			negated:  true,
		},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			count, err := classIsValid((test.class))
			t.Logf("%s classIsValid: %v, %d", name, err, count)
			if !errors.Is(err, test.err) {
				t.Errorf("Test ok %s (%s): Expected %s. Actual %s.", name, string(test.class), test.err, err)
			}

			if err == nil && test.err == nil {
				negated, actual := getClass(test.class)

				if test.negated != negated {
					t.Errorf("Test negated %s (%s): Expected %t. Actual %t.", name, string(test.class), test.negated, negated)
				}

				if len(test.expected) != len(actual) {
					// lengths don't match
					t.Errorf("Test results %s (%s): Expected=%s. Actual=%s", name, string(test.class), string(test.expected), string(actual))
				} else {
					// verify character-by-character
					for i := range actual {
						if actual[i] != test.expected[i] {
							t.Errorf("Test results %s (%s): Expected[%d]=%#U. Actual[%d]=%#U", name, string(test.class), i, test.expected[i], i, actual[i])
						}
					}
				}
			}
		})
	}
}

// Verify matchClass is working correctly
func TestMatchClass(t *testing.T) {
	var testValues = []rune(testValueString)
	const testPatternClass = "[-Helo, 世界! abc?*0-9()[\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]"
	const testPatternClassNegated = "[!-Helo, 世界! abc?*0-9()[\\]{}XYZ'more\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สวัสดี]dfghijklnpqstuvwXYZ0123456789"
	testIO := []struct {
		pattern  []rune
		values   []rune
		expected []bool
	}{
		{
			pattern:  []rune("[][!]"),
			values:   []rune{']'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[][!]"),
			values:   []rune{'['},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[][!]"),
			values:   []rune{'!'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[?]"),
			values:   []rune{'?'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[\\]]"),
			values:   []rune{']'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[!\\]]"),
			values:   []rune{']'},
			expected: []bool{false},
		},
		{
			// match any character except  'a', 'b', 'c',
			// and '/' on Linux or '\' on Windows
			pattern:  []rune{'[', '!', 'a', 'b', GlobSeparator, 'c', ']'},
			values:   []rune{Separator},
			expected: []bool{false},
		},
		{
			// match any character except  'a', 'b', 'c', and '/' on Linux or
			// '\' on Windows, because negated classes never match a path
			// separator whether it is listed or not.
			pattern:  []rune("[!abc]"),
			values:   []rune{Separator},
			expected: []bool{false},
		},
		{
			// Match any character except a separator
			pattern:  []rune{'[', '!', GlobSeparator, ']'},
			values:   []rune{Separator},
			expected: []bool{false},
		},
		{
			pattern:  []rune{'[', '!', escapeCharacter, escapeCharacter, ']'},
			values:   testValues,
			expected: testsMatchNoSeparator(testValues),
		},
		{
			pattern:  []rune(testPatternClass),
			values:   testValues,
			expected: testsMatchNoSeparator(testValues),
		},
		{
			pattern:  []rune(testPatternClassNegated),
			values:   testValues,
			expected: make([]bool, len(testValues)),
		},
		{
			pattern:  []rune{'[', GlobSeparator, ']'},
			values:   []rune{Separator},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[好]"),
			values:   []rune{'好'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[好]"),
			values:   testValues,
			expected: make([]bool, len(testValues)),
		},
		{
			pattern:  []rune("[abc!wx-z]"),
			values:   []rune{'!'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[\\!abcwx-z]"),
			values:   []rune{'!'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[abc\\!wx-z]"),
			values:   []rune{'!'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[-abc!wxz]"),
			values:   []rune{'-'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[-\\]abc!wxz]"),
			values:   []rune{'-'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[abcwx-z-]"),
			values:   []rune{'-'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[abcwx\\-z]"),
			values:   []rune{'-'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[--.0abc!wxz]"),
			values:   []rune{'-'},
			expected: []bool{true},
		},
		{
			pattern:  []rune("[0-9]"),
			values:   []rune("a0123456789z"),
			expected: []bool{false, true, true, true, true, true, true, true, true, true, true, false},
		},
		{
			pattern:  []rune("[!0-9]"),
			values:   []rune("0123456789"),
			expected: make([]bool, 10),
		},
		{
			pattern:  []rune("[[?!*-]"),
			values:   []rune("[?!*-"),
			expected: testsMatchAll([]rune("[?!*-")),
		},
		{
			pattern:  []rune("[![?!*-]"),
			values:   []rune("[?!*-"),
			expected: make([]bool, len([]rune("[?!*-"))),
		},
		{
			pattern:  []rune("[.0-\\\\]"),
			values:   []rune("\\"),
			expected: []bool{true},
		},
		{
			pattern:  []rune("[\\\\-\\]]"),
			values:   []rune("]"),
			expected: []bool{true},
		},
		{
			pattern:  []rune("[-Hello, 世界! abc?*0-9()[\\]{}XYZ'MORE\"\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98สdf-z]"),
			values:   []rune("\\"),
			expected: []bool{false},
		},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			if len(test.values) != len(test.expected) {
				t.Errorf("Test %s(%s): Invalid test setup. The number of test value (%d), expected results (%d)", name, string(test.pattern), len(test.values), len(test.expected))
			}

			negated, pattern := getClass(test.pattern)
			for i, value := range test.values {
				matched := matchClass(pattern, value, negated)
				if matched != test.expected[i] {
					t.Errorf("Test %s[%02d] (\"%s\", %#U): Expected match %t. Actual match %t.", name, i+1, string(test.pattern), value, test.expected[i], matched)
				}
			}
		})
	}
}

// Initialize expected values for tests that always match
func testsMatchAll(values []rune) []bool {
	tests := make([]bool, len(values))

	for i := range tests {
		tests[i] = true
	}

	return tests
}

// Initialize expected values for tests that never match a path separator
func testsMatchNoSeparator(values []rune) []bool {
	tests := make([]bool, len(values))

	for i := range tests {
		if values[i] != Separator {
			tests[i] = true
		}
	}

	return tests
}
