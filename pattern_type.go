// See LICENSE.txt for copyright and licensing information about this file.

package glob

type patternType int

const (
	patternSimple patternType = iota
	patternDirectory
	patternRecursive
	patternUnknown
	simpleString       = "simple"
	directoryString    = "directory"
	recursiveString    = "recursive"
	unknownPatternType = "unknown"
)

func (p patternType) String() string {
	switch p {
	case patternSimple:
		return simpleString
	case patternDirectory:
		return directoryString
	case patternRecursive:
		return recursiveString
	default:
		return unknownPatternType
	}
}
