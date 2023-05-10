// See LICENSE.txt for copyright and licensing information about this file.

package glob

// matchSimple compares a simple path string to a simple glob pattern. If any
// character sequence in path, starting with its first character, matches the
// entire pattern, then return true and the number of runes matched. Otherwise,
// return false and the number of runes matched before a mismatch occurred.
//
// Simple matching includes matching literal characters one-for-one, matching
// any single character to the '?' wildcard, and matching any single character
// to any character defined by a set (aka, character class).
//
// if a class match fails, exit the loop
func matchSimple(pattern, path []rune) (bool, int) {
	var index int
	var matchedCount int
	var isEscaped bool
	var matched = true

	if len(pattern) == 0 && len(path) > 0 {
		// handle special case of a zero length pattern
		return false, 0
	}

mismatch:
	for index = 0; index < len(pattern) && matchedCount < len(path) && matched; index++ {
		token := pattern[index]
		value := path[matchedCount]
		if isEscaped {
			isEscaped = false
			if token == value {
				matchedCount++
			} else {
				// consume a pattern-character and break out of the loop
				index++
				matched = false
				break mismatch
			}
		} else {
			switch token {
			case escapeCharacter:
				// skip past the escape character so literals, such as ']', '-',
				// '\', and '*' can be matched.
				isEscaped = true
			case '[':
				negated, classPattern := getClass(pattern[index:])

				matched = matchClass(classPattern, value, negated)
				if matched {
					matchedCount++
					// set index past the end of class character
					index += len(classPattern) - 1
				} else {
					// consume a pattern-character and break out of the loop
					index++
					break mismatch
				}
			case '?':
				// match any single character except a path separator
				if value != Separator {
					matchedCount++
				} else {
					// consume a pattern-character and break out of the loop
					index++
					matched = false
					break mismatch
				}
			default:
				// match a literal character in the pattern to one in the path.
				if token == value {
					matchedCount++
				} else {
					// consume a pattern-character and break out of the loop
					index++
					matched = false
					break mismatch
				}
			}
		}
	}

	// ensure the entire pattern was used
	matched = matched && index == len(pattern)
	return matched, matchedCount
}
