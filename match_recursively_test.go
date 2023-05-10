// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"fmt"
	"testing"
)

// Verify matchRecursively is working correctly
func TestMatchRecursively(t *testing.T) {
	testIO := []struct {
		pattern string   // test pattern
		paths   []string // one or more paths to test
		matched []bool   // one expected success/failure match for each path
		counts  []int    // the expected number of runes matched on success/failure
	}{
		{
			pattern: "**",
			paths:   []string{"a", "ab", "abc", SeparatorString, SeparatorString + "a", "a" + SeparatorString},
			matched: []bool{true, true, true, true, true, true},
			counts:  []int{1, 2, 3, 1, 2, 2},
		},
		{
			pattern: "**c",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "c"},
			matched: []bool{true},
			counts:  []int{5},
		},
		{
			pattern: "**c*t",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "cat", "a" + SeparatorString + "b" + SeparatorString + "cut", "a" + SeparatorString + "b" + SeparatorString + "caught", "a" + SeparatorString + "b" + SeparatorString + "c" + SeparatorString + "ut", "abcdefghijklmnopqrst"},
			matched: []bool{true, true, true, false, true},
			counts:  []int{7, 7, 10, 0, 20},
		},
		{
			pattern: "**a*b*c",
			paths:   []string{"abcd" + SeparatorString + "xyz" + SeparatorString + "abc"},
			matched: []bool{true},
			counts:  []int{12},
		},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			if len(test.paths) != len(test.counts) || len(test.paths) != len(test.matched) {
				t.Errorf("Test %s: Invalid setup. The number of test paths (%d), expected counts (%d), and expected matches (%d) must be the same.", name, len(test.paths), len(test.counts), len(test.matched))
			}

			if len(test.pattern) < 2 {
				t.Errorf("Test %s: pattern \"%s\" is too short. It must start with '**'.", name, test.pattern)
			}

			if string(test.pattern[:2]) != "**" {
				t.Errorf("Test %s: pattern \"%s\" must start with '**'.", name, test.pattern)
			}

			for i, path := range test.paths {
				head, tail, kind := nextPattern([]rune(test.pattern))
				if kind != patternRecursive {
					t.Errorf("Test %s: expected pattern (%s) to be recursive; actual is %s", name, test.pattern, kind)
				}

				matched, count := matchRecursively([]rune(test.pattern), head, tail, []rune(path))
				if matched != test.matched[i] {
					t.Errorf("Test %s[%02d]: Expected %t. Actual %t. Test %d of %d (%s, %s).", name, i+1, test.matched[i], matched, i+1, len(test.counts), test.pattern, test.paths[i])
					break
				}
				if count != test.counts[i] {
					t.Errorf("Test %s[%02d]: Expected %d. Actual %d. Test %d of %d (%s, %s).", name, i+1, test.counts[i], count, i+1, len(test.counts), test.pattern, test.paths[i])
				}
			}
		})
	}
}
