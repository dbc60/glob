// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"fmt"
	"testing"
)

// Verify Validate is working.
func TestValidate(t *testing.T) {
	testIO := []struct {
		pattern string
		err     error
		length  int
	}{
		{"", nil, 0},
		{"a", nil, 1},
		{"abc", nil, 3},
		{"?", nil, 1},
		{"a?", nil, 2},
		{"?b", nil, 2},
		{"\\?", nil, 2},
		{"a\\?", nil, 3},
		{"\\?b", nil, 3},
		{"\\[", nil, 2},
		{"a\\[", nil, 3},
		{"\\[b", nil, 3},
		{"\\*", nil, 2},
		{"a\\*", nil, 3},
		{"\\*b", nil, 3},
		{"\\\\", nil, 2},
		{"a\\\\", nil, 3},
		{"\\\\b", nil, 3},
		{"\\b", ErrGlobInvalidEscape, 1},
		{"a\\b", ErrGlobInvalidEscape, 2},
		{"[]", ErrGlobTruncated, 1},
		{"[!]", ErrGlobTruncated, 2},
		{"[a]", nil, 3},
		{"[!a]", nil, 4},
		{"[]]", nil, 3},
		{"[]![]", nil, 5},
		{"[-]", nil, 3},
		{"[a-]", nil, 4},
		{"[-b]", nil, 4},
		{"[a-z]", nil, 5},
		{"[?*\\\\]", nil, 6},
		{"[\\!\\-?*\\]]", nil, 10},
		{"*asdf", nil, 5},
		{"*asdf*", nil, 6},
		{string(rune(0x00)), ErrGlobReservedSymbol, 0},
	}
	for i, test := range testIO {
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			length, err := Validate(test.pattern)
			if test.err != err || test.length != length {
				t.Errorf("Test %s: pattern(%s), expected (%d, %s), actual (%d, %s)",
					name, test.pattern, test.length, test.err, length, err)
			}
		})
	}
}
