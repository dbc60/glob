// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"fmt"
	"testing"
)

// Verify matchDirectory is working correctly.
func TestMatchDirectory(t *testing.T) {
	testIO := []struct {
		pattern []rune   // test pattern (implicit leading '*')
		paths   []string // one or more paths to test
		matched []bool   // one expected success/failure match for each path
		counts  []int    // the expected number of runes matched on success/failure
	}{
		{
			pattern: []rune{},
			paths:   []string{"a", "ab", "abc", SeparatorString, SeparatorString + "a", "a" + SeparatorString},
			matched: []bool{true, true, true, false, false, false},
			counts:  []int{1, 2, 3, 0, 0, 0},
		},
		{
			pattern: []rune("z"),
			paths:   []string{"z", "z1", "z12", "12z", SeparatorString, "z" + SeparatorString, SeparatorString + "z"},
			matched: []bool{true, false, false, true, false, false, false},
			counts:  []int{1, 1, 1, 3, 0, 1, 0},
		},
		{
			pattern: []rune{GlobSeparator},
			paths:   []string{"x", "x" + SeparatorString, "xyz" + SeparatorString, SeparatorString, SeparatorString + "xyz", "abc" + SeparatorString + "xyz"},
			matched: []bool{false, true, true, true, false, false},
			counts:  []int{0, 2, 4, 1, 1, 4},
		},
		{
			pattern: []rune("?"),
			paths:   []string{"x", "x" + SeparatorString, "xyz" + SeparatorString, SeparatorString, "abc" + SeparatorString + "xyz"},
			matched: []bool{true, false, false, false, false},
			counts:  []int{1, 1, 3, 0, 3},
		},
		{
			pattern: []rune("aa?"),
			paths:   []string{"aaaaa"},
			matched: []bool{true},
			counts:  []int{5},
		},
		{
			pattern: []rune("."),
			paths:   []string{"a."},
			matched: []bool{true},
			counts:  []int{2},
		},
		{
			// {"*", "abc", true},
			pattern: []rune{},
			paths:   []string{"abc"},
			matched: []bool{true},
			counts:  []int{3},
		},
		{
			// {"*", SeparatorString + "abc", false},
			pattern: []rune{},
			paths:   []string{SeparatorString + "abc"},
			matched: []bool{false},
			counts:  []int{0},
		},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			if len(test.paths) != len(test.counts) || len(test.paths) != len(test.matched) {
				t.Errorf("Test %s: Invalid setup. The number of test paths (%d), expected counts (%d), and expected matches (%d) must be the same", name, len(test.paths), len(test.counts), len(test.matched))
			}

			for i, path := range test.paths {
				matched, count := matchDirectory(test.pattern, []rune{}, []rune(path))
				if matched != test.matched[i] {
					t.Errorf("Test %s [%d of %d] (%s, %s): Expected %t. Actual %t.", name, i+1, len(test.counts), string(test.pattern), path, test.matched[i], matched)
					break
				}
				if count != test.counts[i] {
					t.Errorf("Test %s [%d of %d] (%s, %s): Expected %d. Actual %d.", name, i+1, len(test.counts), string(test.pattern), path, test.counts[i], count)
				}
			}
		})
	}
}
