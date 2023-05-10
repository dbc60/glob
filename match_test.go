// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"fmt"
	"testing"
)

func TestPatterns(t *testing.T) {
	testIO := []struct {
		pattern  string
		path     string
		expected bool
	}{
		{"", "", true},
		{"", "a", false},
		{"", "aa", false},
		{"x", "", false},
		{"x", "x", true},
		{"*", "abc", true},
		{"*", SeparatorString + "abc", false},
		{"**" + GlobSeparatorString + "a" + GlobSeparatorString + "?*.txt", "a" + SeparatorString + "a" + SeparatorString + ".txt", false},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			_, err := Validate(test.pattern)
			if err != nil {
				t.Errorf("Test %s(%s, %s): pattern %s is invalid: %s", name, test.pattern, test.path, test.pattern, err.Error())
			} else {
				matched, _ := Match(test.pattern, test.path)
				if matched != test.expected {
					t.Errorf("Test %s(%s, %s): Expected %t. Actual %t.", name, test.pattern, test.path, test.expected, matched)
				}
			}
		})
	}
}

// Verify Match is working correctly.
func TestMatchAll(t *testing.T) {
	testIO := []struct {
		pattern string   // test pattern
		paths   []string // one or more paths to test
		matched []bool   // one expected success/failure match for each path
		counts  []int    // the expected number of runes matched on success/failure
	}{
		{
			pattern: "abc",
			paths:   []string{"abc"},
			matched: []bool{true},
			counts:  []int{3},
		},
		{
			pattern: "*",
			paths:   []string{"abc", "a" + SeparatorString + "b"},
			matched: []bool{true, false},
			counts:  []int{3, 1},
		},
		{
			pattern: "a*",
			paths:   []string{"a1b2c3d4e"},
			matched: []bool{true},
			counts:  []int{9},
		},
		{
			pattern: "a*a*", // a simple pathological test case
			paths:   []string{"aaaaaaaaaa"},
			matched: []bool{true},
			counts:  []int{10},
		},
		{
			pattern: "a*a*a*a*a*a*a*a*a*a*",
			paths:   []string{"aaaaaaaaaa"},
			matched: []bool{true},
			counts:  []int{10},
		},
		{
			pattern: "a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*",
			paths:   []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			matched: []bool{true},
			counts:  []int{100},
		},
		{
			pattern: "**",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "c"},
			matched: []bool{true},
			counts:  []int{5},
		},
		{
			pattern: "**c",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "c"},
			matched: []bool{true},
			counts:  []int{5},
		},
		{
			pattern: "**a**b**c",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "c"},
			matched: []bool{true},
			counts:  []int{5},
		},
		{
			pattern: "**a" + GlobSeparatorString + "**b" + GlobSeparatorString + "**c",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "c"},
			matched: []bool{true},
			counts:  []int{5},
		},
		{
			pattern: "**c*t",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "cat", "a" + SeparatorString + "b" + SeparatorString + "cut", "a" + SeparatorString + "b" + SeparatorString + "caught"},
			matched: []bool{true, true, true},
			counts:  []int{7, 7, 10},
		},
		{
			pattern: "**c*t",
			paths:   []string{"a" + SeparatorString + "b" + SeparatorString + "c" + SeparatorString + "ut"},
			matched: []bool{false},
			counts:  []int{0},
		},
		{
			// The path to match is "\Users\foo\ba世\界.txt" for Windows and "/Users/foo/ba世/界.txt" for Linux
			pattern: "/Users/**/[bc]a[!a-qsu-z]/?*.txt",
			paths:   []string{SeparatorString + "Users" + SeparatorString + "foo" + SeparatorString + "ba世" + SeparatorString + "界.txt"},
			matched: []bool{true},
			counts:  []int{20},
		},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			if len(test.paths) != len(test.counts) || len(test.paths) != len(test.matched) {
				t.Errorf("Test %s(%s): Invalid test setup. The number of test paths (%d), expected counts (%d), and expected matches (%d) must be the same", name, string(test.pattern), len(test.paths), len(test.counts), len(test.matched))
			}

			for i, path := range test.paths {
				count, err := Validate(test.pattern)
				if err != nil {
					t.Errorf("Test %s[%02d] pattern %s is invalid: %s", name, i+1, test.pattern, err.Error())
				}

				if count != len(test.pattern) {
					t.Errorf("Test %s[%02d] Validate(%s) returned %d, expected %d", name, i+1, test.pattern, count, len(test.pattern))
				}

				if err == nil {
					matched, count := Match(test.pattern, path)
					if matched != test.matched[i] {
						t.Errorf("Test %s[%02d]: Expected %t. Actual %t. Test %d of %d (%s, %s).", name, i+1, test.matched[i], matched, i+1, len(test.counts), string(test.pattern), test.paths[i])
					}

					if count != test.counts[i] {
						t.Errorf("Test %s[%02d]: Expected %d. Actual %d. Test %d of %d (%s, %s).", name, i+1, test.counts[i], count, i+1, len(test.counts), string(test.pattern), test.paths[i])
					}
				}
			}
		})
	}
}
