// See LICENSE.txt for copyright and licensing information about this file.

package main

import (
	"fmt"
	"log"
	"os"

	glob "dbc60/goglob"
)

// Exercise glob.Validate and glob.Match
func main() {
	if len(os.Args) != 3 {
		log.Fatal("expected a glob pattern and file path as arguments")
	}

	pattern := os.Args[1]
	path := os.Args[2]
	count, err := glob.Validate(pattern)

	if err != nil {
		log.Fatalf("Pattern %s is invalid at position %d: %s.", pattern, count+1, err)
	}

	matched, count := glob.Match(pattern, path)
	if matched {
		fmt.Printf("Pattern %s matches path %s\n", pattern, path)
	} else {
		fmt.Printf("Pattern %s mismatches at position %d in path %s\n", pattern, count+1, path)
	}
}
