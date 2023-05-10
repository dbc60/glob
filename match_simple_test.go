// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"fmt"
	"testing"
)

// Verify matchSimple is working correctly
func TestMatchSimple(t *testing.T) {
	// This seems like a good sample of runes for testing matchClass
	testIO := []struct {
		pattern []rune   // test pattern
		paths   []string // one or more paths to test
		matched []bool   // one expected success/failure match for each path
		counts  []int    // the expected number of runes matched on success/failure
	}{
		{
			pattern: []rune{},
			paths:   []string{"x"},
			matched: []bool{false},
			counts:  []int{0},
		},
		{
			pattern: []rune{'x'},
			paths:   []string{"x", "y"},
			matched: []bool{true, false},
			counts:  []int{1, 0},
		},
		{
			pattern: []rune{escapeCharacter, ']'},
			paths:   []string{"]", "?"},
			matched: []bool{true, false},
			counts:  []int{1, 0},
		},
		{
			pattern: []rune{'?'},
			paths:   []string{SeparatorString, "x"},
			matched: []bool{false, true},
			counts:  []int{0, 1},
		},
		{
			pattern: []rune{GlobSeparator},
			paths:   []string{SeparatorString, "z"},
			matched: []bool{true, false},
			counts:  []int{1, 0},
		},
		{
			pattern: []rune("[ABC]"),
			paths:   []string{"B", "z"},
			matched: []bool{true, false},
			counts:  []int{1, 0},
		},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			if len(test.paths) != len(test.counts) || len(test.paths) != len(test.matched) {
				t.Errorf("Test %s: Invalid test setup. The number of test paths (%d), expected counts (%d), and expected matches (%d) must be the same", name, len(test.paths), len(test.counts), len(test.matched))
			}

			for i, path := range test.paths {
				matched, count := matchSimple(test.pattern, []rune(path))
				if matched != test.matched[i] {
					t.Errorf("Test %s [%d of %d] (%s, %s): Expected %t. Actual %t.", name, i+1, len(test.counts), string(test.pattern), path, test.matched[i], matched)
				}
				if count != test.counts[i] {
					t.Errorf("Test %s [%d of %d] (%s, %s): Expected %d. Actual %d.", name, i+1, len(test.counts), string(test.pattern), path, test.counts[i], count)
				}
			}
		})
	}
}
