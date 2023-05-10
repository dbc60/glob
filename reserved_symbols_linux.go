// See LICENSE.txt for copyright and licensing information about this file.

//go:build linux

package glob

// Per the "most UNIX file systems" row in the "Comparison of filename
// limitations" table at
// https://en.wikipedia.org/wiki/Filename#Comparison_of_filename_limitations
// the reserved symbols are nul (0x00) and /. However, we need '/' for glob
// pattern matching, so nul is the only reserved symbol.
func init() {
	reservedSymbols = make(reservedSymbolSet)
	reservedSymbols[rune(0x00)] = exists
}
