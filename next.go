// See LICENSE.txt for copyright and licensing information about this file.

package glob

// nextPattern breaks down a glob pattern into one of the glob types. It accepts
// a slice of runes representing a glob pattern and returns a slice containing
// just the next sub-pattern and its type. It assumes the glob pattern is valid.
// If the sub-pattern is directory or recursive, it removes the leading
// asterisk(s).
func nextPattern(pattern []rune) (head, tail []rune, kind patternType) {
	var start int
	var end int
	kind = patternSimple

	// If the pattern isn't empty and the first character is an asterisk, then
	// it's either a directory or recursive pattern.
	if len(pattern) > 0 && pattern[0] == '*' {
		start++
		kind = patternDirectory

		// check if it's a recursive pattern
		if start < len(pattern) && pattern[start] == '*' {
			start++
			kind = patternRecursive

			// skip over extra asterisks
			for ; start < len(pattern) && pattern[start] == '*'; start++ {
			}
		}
	}

	var inClass bool // true if parsing inside a character class

	// find the end of the sub-pattern. An asterisk starts another chunk unless
	// it is found within a character class.
endPattern:
	for end = start; end < len(pattern); end++ {
		token := pattern[end]

		switch token {
		case escapeCharacter:
			// token is an escape character only if it's in a class, otherwise
			// it's a path separator for Windows and just another character for
			// Linux.
			if inClass {
				// skip the escaped character so we don't miss a real ']'
				end++
			}
		case '[':
			inClass = true
		case ']':
			inClass = false
		case '*':
			if !inClass {
				// end of current sub-pattern; '*' starts another one
				break endPattern
			}
		}
	}

	head = pattern[start:end]
	tail = pattern[end:]
	return
}
