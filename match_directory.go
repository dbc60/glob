// See LICENSE.txt for copyright and licensing information about this file.

package glob

// matchDirectory accepts the simple pattern that follows the asterisk in a
// directory pattern and a candidate path. It returns true or false depending on
// whether the path matched
func matchDirectory(head, tail, path []rune) (bool, int) {
	var matched bool
	var total int
	var simpleCount int
	var simpleMatch bool

	// keep trying until a match with the path consumed, or the next character is
	// a separator.
	more := true
	for more {
		more = false
		if len(head) == 0 {
			// Trivial case: asterisk-only wildcard matches all but a path separator,
			// so if there are no path separators, it's a match.
			var i int
			for i = 0; i < len(path) && path[i] != Separator; i++ {
			}
			if i == len(path) {
				matched = true
				total = i
			}

			return matched, total
		}

		// Match path against the head pattern. If it fails, restart the match one
		// character later in path. Repeat until success or the target is consumed.
		// If successful, return true and the number of characters matched.
		// for-loop is like
		// do {...} while(!simpleMatch && len(path) > 0 && simpleCount > 0);
		retry := true
		var subtotal int
		for retry {
			simpleMatch, simpleCount = matchSimple(head, path)
			if !simpleMatch {
				if len(path) > 1 && path[0] != Separator {
					simpleCount = 1
					path = path[1:]
					subtotal++
				} else {
					simpleCount = 0
				}
			} else {
				subtotal += simpleCount
			}
			retry = !simpleMatch && len(path) > 0 && simpleCount > 0
		}

		if simpleMatch {
			if len(tail) == 0 && len(path) > simpleCount {
				// head is the last directory pattern, but there's more target to match.
				// Shift the path one character and restart the match
				if len(path) > 1 && path[0] != Separator {
					total += 1 + subtotal - simpleCount
					path = path[1:]
					more = true
				} else {
					// we can't shift the path any more, so the match fails
					simpleMatch = false
					total += subtotal
				}
			} else {
				total += subtotal
			}
		}
	}

	matched = simpleMatch
	return matched, total
}
