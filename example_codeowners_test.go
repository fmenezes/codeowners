package codeowners_test

import (
	"fmt"
	"strings"

	"github.com/fmenezes/codeowners"
)

func Example() {
	codeownerDecoder := codeowners.NewDecoder(strings.NewReader(`* test@example.org`))
	for codeownerDecoder.More() {
		token := codeownerDecoder.Token()
		fmt.Printf("File Pattern: %s\n", token.Path())
		fmt.Printf("Owners: %v\n", token.Owners())
	}
	// Output:
	// File Pattern: *
	// Owners: [test@example.org]
}
