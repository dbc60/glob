// See LICENSE.txt for copyright and licensing information about this file.

package glob

type globError string

const (
	ErrGlobZeroLength     = globError("zero-length pattern")
	ErrGlobNoLeftBracket  = globError("class must start with a left bracket ('[')")
	ErrGlobTruncated      = globError("pattern truncated")
	ErrGlobInvalidEscape  = globError("invalid escape sequence")
	ErrGlobInvalidRange   = globError("invalid range")
	ErrGlobReservedSymbol = globError("glob error: reserved symbol found in pattern")
)

// satisfy the error interface
func (err globError) Error() string {
	return string(err)
}
