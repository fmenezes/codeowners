// Package codeowners provides funcionality to evaluate CODEOWNERS file.
package codeowners // import "github.com/fmenezes/codeowners"

import (
	"strings"
	"unicode"
)

// DefaultLocations provides default locations for the CODEOWNERS file
var DefaultLocations = [...]string{"CODEOWNERS", "docs/CODEOWNERS", ".github/CODEOWNERS"}

// ParseLine parses a CODEOWNERS line into file pattern and owners
func ParseLine(line string) (string, []string) {
	line = sanitiseLine(line)

	var previousRune rune
	data := strings.FieldsFunc(line, func(r rune) bool {
		result := unicode.IsSpace(r) && previousRune != '\\'
		previousRune = r
		return result
	})

	if len(data) > 1 {
		return data[0], data[1:]
	} else if len(data) == 1 {
		return data[0], nil
	} else {
		return "", nil
	}
}
