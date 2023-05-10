// See LICENSE.txt for copyright and licensing information about this file.

package glob

import (
	"fmt"
	"testing"
)

// Verify getClass is working correctly
func TestErrors(t *testing.T) {
	testIO := []struct {
		err error
	}{
		{err: ErrGlobZeroLength},
		{err: ErrGlobNoLeftBracket},
		{err: ErrGlobTruncated},
		{err: ErrGlobInvalidEscape},
		{err: ErrGlobInvalidRange},
		{err: nil},
	}

	for i, test := range testIO {
		// There are less than 1000 tests, so each test can have a name of 3 digits
		// with leading zeros
		name := fmt.Sprintf("%03d", i+1)
		t.Run(name, func(t *testing.T) {
			switch test.err {
			case ErrGlobZeroLength:
				t.Logf("Error %s", test.err.Error())
			case ErrGlobNoLeftBracket:
				t.Logf("Error %s", test.err.Error())
			case ErrGlobTruncated:
				t.Logf("Error %s", test.err.Error())
			case ErrGlobInvalidEscape:
				t.Logf("Error %s", test.err.Error())
			case ErrGlobInvalidRange:
				t.Logf("Error %s", test.err.Error())
			case nil:
				t.Log("Error is nil")
			default:
				t.Errorf("Error %s", test.err.Error())
			}
		})
	}
}
