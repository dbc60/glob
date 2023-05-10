// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"testing"
)

func TestPatternTypes(t *testing.T) {
	var p patternType
	for p = patternSimple; p <= patternUnknown; p++ {
		actual := p.String()
		switch p {
		case patternSimple:
			if actual != simpleString {
				t.Errorf("Simple pattern type: expected \"%s\", actual \"%s\".", simpleString, actual)
			}
		case patternDirectory:
			if actual != directoryString {
				t.Errorf("Simple pattern type: expected \"%s\", actual \"%s\".", directoryString, actual)
			}
		case patternRecursive:
			if actual != recursiveString {
				t.Errorf("Simple pattern type: expected \"%s\", actual \"%s\".", recursiveString, actual)
			}
		case patternUnknown:
			if actual != unknownPatternType {
				t.Errorf("Simple pattern type: expected \"%s\", actual \"%s\".", unknownPatternType, actual)
			}
		}
	}
}
