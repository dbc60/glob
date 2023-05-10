// See LICENSE.txt for copyright and licensing information about this file.

package glob

// classIsValid accepts a pattern that starts with a left bracket and returns
// true and the length of the class if the class is valid, and false and the
// number of characters found before an error was encountered otherwise.
func classIsValid(pattern []rune) (int, error) {
	// A zero-length pattern is invalid
	if len(pattern) == 0 {
		return 0, ErrGlobZeroLength
	}

	// The pattern must start with '['
	if pattern[0] != '[' {
		return 1, ErrGlobNoLeftBracket
	}

	// The pattern must be at least 3 characters long, or if the first character
	// after the leading '[' is '!', it must be at least 4 characters long.
	if len(pattern) < 3 || pattern[1] == '!' && len(pattern) < 4 {
		// return length up to, but not including the end to indicate that not
		// only is the pattern invalid, but that it is too short. In effect, it
		// allows nextValidPattern to return a more intuitive character count.
		return len(pattern) - 1, ErrGlobTruncated
	}

	var i int = 1
	var err error
	var lo rune
	for done := false; !done && err == nil && i < len(pattern); i++ {
		token := pattern[i]
		switch token {
		case escapeCharacter:
			if i < len(pattern)-1 {
				// skip the escape character
				i++
				lo = pattern[i]
				switch lo {
				case '!', '-', ']', escapeCharacter:
					// These are the only valid escape sequences, but the class
					// can't end in "\]"
					if lo == ']' && i == len(pattern)-1 {
						err = ErrGlobTruncated
					}
				default:
					// invalid escape sequence
					err = ErrGlobInvalidEscape
				}
			} else {
				err = ErrGlobTruncated
			}
		case ']':
			// found end of class if the right bracket is not the first character, nor
			// is it the second character in a negated class
			if i != 1 || i != 2 && pattern[1] == '!' {
				done = true
			}
		case '-':
			// if the hyphen is the first character in the class, or the second
			// character in a negated class, or the last character just before
			// the terminating right bracket, then it's just a literal hyphen.
			if i == 1 || i == 2 && pattern[1] == '!' || i == len(pattern)-2 {
				continue
			}

			// We have a range. Verify it's valid (doesn't include '/')
			i++
			hi := pattern[i]
			if GlobSeparator > lo && GlobSeparator < hi || hi == escapeCharacter {
				err = ErrGlobInvalidRange
			}
		default:
			// capture the current token in case the next one starts a range
			lo = token
		}
	}

	if err == nil && pattern[i-1] != ']' {
		err = ErrGlobTruncated
	}

	return i, err
}

// matchClass accepts a class pattern and compares a character against
// characters defined by the pattern. It returns true if there is a match or
// false otherwise. A negated class pattern matches any character not included
// in the set, including the path separator.
//
// The string enclosed by the brackets cannot be empty; therefore ']' can be
// allowed between the brackets, provided that it is the first character. Thus,
// "[][!]" matches the three characters '[', ']', and '!'.).
func matchClass(pattern []rune, value rune, negated bool) bool {
	var escaped bool
	var matched bool

	// match a separator only if it's explicit in the class pattern and it's not
	// negated.
	if value == Separator {
		if !negated {
			// match the separator only if it is explicitly listed
			for _, token := range pattern {
				if token == GlobSeparator {
					return true
				}
			}
		} else {
			// negated classes never match the path separator, regardless of
			// whether it is explicitly listed.
			return false
		}
	}

	// Deal with the empty and negated empty classes.
	i := 1
	if negated {
		// move past the exclamation point
		i++
	}

	// loop through characters in the pattern attempting to match one of
	// them to the character. Stop when there is either a match, or the
	// pattern is consumed.
	var lo, hi rune

	for !matched && ((pattern[i] != ']' || i == 1) || escaped) {
		token := pattern[i]
		// If token is '\' and escaped is not true, set it to true, otherwise
		// ensure escaped is false.
		escaped = !escaped && (token == escapeCharacter)

		// skip past an escape character so literals such
		// as ']', '-', and '\' can be matched.
		if escaped {
			i++
			token = pattern[i]
			escaped = false
		}

		// initialize lo and hi to the same character
		lo = token
		hi = lo
		i++
		token = pattern[i]

		// if the next character in the pattern is '-' and it is not the last
		// character in the class, then we have a range. The value of "lo" is
		// already set, so capture the end of the range in "hi".
		if token == '-' && pattern[i+1] != ']' {
			i++
			token = pattern[i]
			escaped = !escaped && (token == escapeCharacter)

			// skip past an escape character so literals such
			// as ']', '-', and '\' can be matched.
			if escaped {
				i++
				token = pattern[i]
				escaped = false
			}
			hi = token
		}

		// Match if the character is in range.
		if value >= lo && value <= hi {
			matched = true
		}
	}

	// reverse the sense of matched if the class is negated.
	matched = matched != negated

	return matched
}

// getClass accepts a pattern that starts with a left bracket and returns a
// boolean indicating if the class is negated and a slice of runes containing
// the class pattern.
func getClass(pattern []rune) (negated bool, subpattern []rune) {
	// Assume the pattern is valid and starts with '[', for why else would we be here?
	var i int = 1

	negated = pattern[1] == '!'
	for done := false; !done; i++ {
		token := pattern[i]
		switch token {
		case ']':
			if i != 1 || i != 2 && negated {
				done = true
			}
		case escapeCharacter:
			// skip it
			i++
		}
	}

	subpattern = pattern[:i]
	negated = pattern[1] == '!'
	return
}
