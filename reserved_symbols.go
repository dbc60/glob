// See LICENSE.txt for copyright and licensing information about this file.

package glob

type reservedSymbolSet map[rune]struct{}

var (
	exists          = struct{}{}
	reservedSymbols reservedSymbolSet
)

// hasReservedSymbol returns true if the string contains a character that is a
// member of reservedSymbols (a set of symbols not allowed glob patterns).
func hasReservedSymbol(pattern []rune) (bool, int) {
	for i, symbol := range pattern {
		if isReservedSymbol(symbol) {
			return true, i
		}
	}
	return false, 0
}

// isReservedSymbol returns true if the symbol is contained in the set of
// symbols not allowed in a file path.
func isReservedSymbol(symbol rune) bool {
	_, ok := reservedSymbols[symbol]
	return ok
}
