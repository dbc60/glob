// See LICENSE.txt for copyright and licensing information about this file.

package glob

// matchRecursively compares a path to glob pattern that may match sequences of
// zero or more characters (including path separators). In other words,
// matchRecursively may match zero or more directory levels. If the match
// succeeds, return true and the number of characters it matched in the path. If
// the match fails, return false and the number of characters that were matched
// before a mismatch occurred.
//
// Note that head is a simple pattern, but represents a recursive pattern, and
// tail is zero or more directory patterns.
func matchRecursively(pattern, head, tail, path []rune) (bool, int) {
	var match bool
	var total, subtotal int

	more := true
	// Compare path to pattern. If head is the last pattern (i.e., tail is
	// empty), repeat until the path is consumed or there is no match.
	for more {
		more = false
		match, subtotal = matchRecursivePattern(pattern, head, tail, path)
		if match {
			total += subtotal
		}
	}

	return match, total
}

func matchRecursivePattern(pattern, head, tail, path []rune) (bool, int) {
	if len(head) == 0 {
		// Trivial case: ** wildcard matches any path
		return true, len(path)
	}

	var match bool
	var total int
	current := path
	shifted := path

	// Loop while the pattern hasn't matched and the path hasn't been consumed.
	// Use matchSimple to compare the path to the head pattern. If it fails,
	// restart the comparison starting with the next character in the path.
	// Repeat until all patterns have been matched to some sequence in the path
	// (success), or the path is consumed before the patterns are.
	//
	// N.B.: We don't have to consume the entire path. We only need to match
	// this recursive pattern and all subsequent directory "sub-patterns" here.
	// if nextPattern returns a recursive or simple pattern, we have
	// successfully matched this recursive pattern and all subsequent directory
	// patterns.
	//
	// If the next pattern is a simple pattern, it must be empty. If the path is
	// also consumed, then the path has matched the pattern. If the next pattern
	// is another recursive pattern, then we return to the caller which will
	// attempt to match any remaining path against that new recursive pattern.
	//
	// The for-loop is like "do {...} while(!matched && shiftCount > 0)"
	more := true
	for more {
		var kind patternType
		var subtotal int = 1 // 1 gets the for-loop started
		more = false

		// match head to some or all of the path
		for !match && len(current) > 0 && subtotal > 0 {
			match, subtotal = matchSimple(head, current)
			if !match {
				if len(current) > 1 {
					subtotal = 1
					current = current[1:]
				} else {
					subtotal = 0
				}
			}
			total += subtotal
		}

		// Match the rest of the patterns in tail (all directory patterns), if
		// any to the rest of the path.
		head, tail, kind = nextPattern(tail)
		for kind == patternDirectory && match &&
			(len(head) > 0 || len(tail) > 0 || len(current) > 0) {
			current = current[subtotal:]
			match, subtotal = matchDirectory(head, tail, current)
			if !match {
				// The directory pattern failed to match a part of the current
				// target. Start over with the recursive pattern, but shift the
				// starting position in path to the next match point for the
				// recursive pattern, and reset total to the number of
				// characters matched so far by the recursive pattern.
				head, tail, kind = nextPattern(pattern)
				shifted = shifted[1:]
				total = len(path) - len(shifted)
				current = shifted
				more = true
			} else {
				current = current[subtotal:]
				total += subtotal
				subtotal = 0
				head, tail, kind = nextPattern(tail)
			}
		}
	}

	// The total matched is the number of characters matched from a simple
	// match, zero-or-more directory matches, and zero-or-more recursive shifts.
	return match, total
}
