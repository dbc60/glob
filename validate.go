// See LICENSE.txt for copyright and licensing information about this file.

package glob

// Validate the given glob pattern. Return true and the length of the pattern if
// it's valid. Otherwise, return false and the zero-based index of the rune
// where validation failed.
func Validate(pattern string) (int, error) {
	var index int

	head, tail, err := nextValidPattern([]rune(pattern))
	index += len(head)
	for err == nil && len(head) > 0 {
		head, tail, err = nextValidPattern(tail)
		index += len(head)
	}

	return index, err
}

func nextValidPattern(pattern []rune) (head, tail []rune, err error) {
	var start int
	var end int
	var isEscaped bool

	// If the pattern isn't empty and the first character is an asterisk, then
	// it's either a directory or recursive pattern.
	if len(pattern) > 0 && pattern[0] == '*' {
		// skip over additional asterisks
		for ; start < len(pattern) && pattern[start] == '*'; start++ {
		}
	}

	// find the end of the sub-pattern. An asterisk starts another chunk unless
	// it is found within a character class.
endPattern:
	for end = start; err == nil && end < len(pattern); end++ {
		token := pattern[end]

		// reserved symbols cannot be found in a path, so reject the pattern
		if isReservedSymbol(token) {
			err = ErrGlobReservedSymbol
			break
		}

		// Only wildcards and the escape character can be escaped.
		if isEscaped {
			isEscaped = false
			if token != '?' && token != '[' && token != '*' && token != escapeCharacter {
				err = ErrGlobInvalidEscape
				break
			} else {
				continue
			}
		}

		switch token {
		case escapeCharacter:
			isEscaped = true
		case '?':
			continue
		case '[':
			var count int
			count, err = classIsValid([]rune(pattern[end:]))
			end += count - 1
		case '*':
			// end of current sub-pattern; '*' starts another one
			break endPattern
		}
	}

	head = pattern[:end]
	tail = pattern[end:]

	return
}
