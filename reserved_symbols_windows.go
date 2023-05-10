// See LICENSE.txt for copyright and licensing information about this file.

//go:build windows

package glob

// Per the NTFS row in the "Comparison of filename limitations" table at
// https://en.wikipedia.org/wiki/Filename#Comparison_of_filename_limitations
// the reserved symbols are the ascii control codes (0x00-0x1F and 0x7F) and
// these characters: " * / : < > ? \ |. However, we need '\', '/', '?', and '*'
// for glob patterns:
//
//   - '\': is used as a path separator and an escape character so glob patterns
//     can contain literal characters that would otherwise be interpreted as
//     glob patterns ('[', ']', and '!' in a class range).
//   - '/': is a alternate path separator, because glob patterns are more fun
//     and flexible that way.
//   - '?': match any single character except a path separator
//   - '*': is used to match either zero or more characters except a path
//     separator (a single '*'), or zero or more characters including any and
//     all path separators (a sequence of two or more asterisks).
func init() {
	reservedSymbols = make(reservedSymbolSet)
	// ascii control codes 0x00 - 0x1F are reserved
	for i := 0; i < 0x20; i++ {
		reservedSymbols[rune(i)] = exists
	}
	reservedSymbols[rune(0x7F)] = exists
	reservedSymbols[rune('"')] = exists
	reservedSymbols[rune(':')] = exists
	reservedSymbols[rune('<')] = exists
	reservedSymbols[rune('>')] = exists
	reservedSymbols[rune('|')] = exists
}
