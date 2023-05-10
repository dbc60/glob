// See LICENSE.txt for copyright and licensing information about this file.

// Package glob implements functions for matching file paths against glob
// patterns. The Linux implementation expects directory separators to be
// forward slashes, while the Windows implementation expects backslashes.
// This is a library package meant to be a part of other programs. The main
// point of entry is the Match function:
//
//	Match(patternString, pathString string) (bool, int).
package glob

import (
	"strings"
)

// The escape character is used to enable interpreting characters used in glob
// patterns as literal characters.
const escapeCharacter = '\\'

// glob patterns always use '/' as a path separator, regardless of the OS.
const GlobSeparator = '/'
const GlobSeparatorString = "/"

// Match accepts a string representing a glob pattern and another string
// representing a path in the file system. It returns true and the number of
// runes in the path that were matched if the path is matched by the glob
// pattern or false and the number of runes matched if the match failed.
//
// This function may be used repeatedly to check multiple paths against a single
// glob pattern. As such, it is assumed that the glob represents a valid pattern.
//
// The goal of this algorithm is to use a set of glob patterns to match against
// paths in a file system where patterns can match paths either at a fixed-depth
// or recursively. A glob-pattern string may contain literal characters,
// character classes to match sets of characters against a single character in
// the path, or wildcards that may match a single character, zero or more
// characters except a path separator, or zero or more characters including a
// path separator (for recursively matching a paths).
//
// It helps to separate glob patterns into three different types.
//
// ## Type 1 Glob Patterns, Simple Patterns.
//
// Type 1 glob patterns are the most simple. They contain any of the following
// patterns and terminate at either the end of the pattern string or just before
// the first asterisk encountered:
//
//   - any literal character that is valid for a file or directory name. These are
//     dependent on the file system in use.
//   - a literal path separator (usually `\` or `/`). This is dependent on the
//     file system in use.
//   - `?`: match any single character, except a path separator.
//   - `[abc]`: match the next character to any character in the set delineated by
//     the square brackets.
//   - `[a-z0-9]`: match the next character to any character in the range. In this
//     example, there are actually two ranges — one from `a` to `z`
//     and the other from `0` to `9`.
//   - `[!abc]`: match the next character to any character NOT in the set.
//   - `[!a-z]`: match the next character to any character NOT in the range.
//
// Matching a Type 1 glob pattern succeeds when the pattern is consumed. If this
// is the only pattern in the glob pattern string, then the path must also be
// completely consumed. It fails on the first mismatch or when the path is
// completely consumed before the pattern is consumed.
//
// ## Type 2 Glob Pattern, Directory Pattern
//
// Type 2 glob patterns start with a single asterisk (`*`) followed by a Type 1
// pattern and terminates at either the end of the pattern string or the next
// asterisk encountered in that pattern string. The asterisk matches zero or
// more characters not including path separators. The Type 1 pattern that follows
// the initial asterisk is the "follower".
//
// The asterisk matches characters in the path in a non-greedy fashion. That is,
// it matches a character only when the follower fails to match the path either
// before the end of the path or before the next path separator _that doesn't
// match a literal path separator in the pattern_. Matching succeeds when the
// follower is consumed, and fails when either the path is consumed before the
// follower is consumed or a path separator in the path is encountered and it
// doesn't match a literal path separator in the follower. If this is the only
// pattern in the glob pattern string, then both the pattern and the path must
// be completely consumed for the match to succeed.
//
// ## Type 3 Glob Pattern, Recursive Pattern
//
// Type 3 glob patterns start with two asterisks (`**`) followed by a Type 1
// pattern and zero or more Type 2 patterns. Note that three or more consecutive
// asterisks are reduced to a single pair (`**`). A Type 3 pattern terminates
// either at the end of the glob pattern string or the next sequence of two or
// more asterisks in the pattern string. The double asterisk matches zero or
// more characters _including path separators_.
//
// Like Type 2 patterns, Type 3 patterns matches characters in the path in a non-
// greedy fashion. The difference is that its follower may need to be segmented
// into two or more followers, where each segment after the first is a Type 2
// pattern.
//
// The matching algorithm for a Type 3 glob pattern has two cases, the follower
// is either a Type 1 or a Type 2 pattern. In the first case there is only one
// segment. This case is processed similarly to the Type 2 pattern, except for
// the fact that a path separator may always be consumed from the path before
// the matching process restarts.
//
// The second case is a Type 1 segment followed by one or more Type 2 segments.
// In this case, matching must be restarted if _any of the segments fail_. Like
// before, the double asterisk consumes one character (which may be a path
// separator), and matching resumes with the first segment and this new, shorter
// path. It succeeds when all segments succeed, and fails when any segment fails
// and the path is consumed.
//
// ## Example Glob Patterns
//
// Here are some example glob patterns:
//
//	a       : match the letter "a"
//	*a      : match any sequence of characters except a path separator up to
//	          and including the first letter "a".
//	**a     : match any sequence of characters including any path separator up
//	          to and including the first letter "a".
//	/Users/*/Documents
//	        : match a "Documents" file or directory in any subdirectory of the
//			  "Users" directory.
//	/usr/**/[bc]a[!a-qsu-z]/?*.txt
//	        : match any file or directory that is a subdirectory of "/usr/" at
//	          any level, where the directory name starts with 'b' or 'c' followed
//	          by 'a' and another letter not in the set {a-q, s, u-z}, and the
//	          file or directory starts with any character ('?') followed by zero
//	          or more characters and end with ".txt". Some possible matches are
//	          "/usr/bat/x.txt", "/usr/foo/bar/baz/file.txt",
//			  "/usr/one/two/three/car/note.txt", and "/usr/ba世/界.txt"
func Match(patternString, pathString string) (bool, int) {
	var matchCount int
	var count int
	var matched bool

	patternMatched := true

	pattern := []rune(patternString)
	path := []rune(pathString)

	// Get the next chunk of the glob pattern, and the pattern type. Note that
	// there can be only one chunk that is simple, and it will be the first one.
	// If there are any subsequent chunks, they will be either one that matches
	// any file or directory name or one that does that recursively.
	head, tail, kind := nextPattern(pattern)

	if kind == patternSimple {
		patternMatched, matchCount = matchSimple(head, path)
		if patternMatched {
			path = path[matchCount:]
			if len(tail) == 0 {
				// no more glob patterns, so return what we have
				matched = len(path) == 0
				return matched, matchCount
			}
		}

		pattern = tail
		head, tail, kind = nextPattern(pattern)
	}

	for kind == patternDirectory && patternMatched && (len(head) > 0 || len(tail) > 0 || len(path) > 0) {
		if len(head) == 0 {
			// pattern is just a "*" wildcard
			matched = !strings.Contains(string(path), SeparatorString)
			if matched {
				matchCount += len(path)
			} else {
				matchCount += strings.Index(string(path), SeparatorString)
			}
			return matched, matchCount
		}

		patternMatched, count = matchDirectory(head, tail, path)
		if patternMatched {
			matchCount += count
			path = path[count:]
			pattern = tail
			head, tail, kind = nextPattern(pattern)
		}
	}

	for kind == patternRecursive && patternMatched && (len(head) > 0 || len(tail) > 0 || len(path) > 0) {
		if len(head) == 0 {
			pattern = tail
			head, tail, kind = nextPattern(pattern)
			if len(head) == 0 {
				// pattern is just "**" wildcard, which matches everything
				matched = true
				matchCount += len(path)
				return matched, matchCount
			}
		}

		patternMatched, count = matchRecursively(pattern, head, tail, path)
		if patternMatched {
			matchCount += count
			path = path[count:]
			pattern = tail
			head, tail, kind = nextPattern(pattern)
		}
	}

	matched = patternMatched && len(path) == 0
	return matched, matchCount
}
