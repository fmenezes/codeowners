// Package codeowners provides funcionality to evaluate CODEOWNERS file.
package codeowners // import "github.com/fmenezes/codeowners"

import (
	"bufio"
	"io"
	"strings"
)

// DefaultLocations provides default locations for the CODEOWNERS file
var DefaultLocations = [...]string{"CODEOWNERS", "docs/CODEOWNERS", ".github/CODEOWNERS"}

// Decoder providers functionality to read CODEOWNERS data
type Decoder struct {
	scanner *bufio.Scanner
	line    *string
	done    bool
}

// New generates a new CodeOwnersScanner instance. The reader should be a reader containing the contents of the CODEOWNERS file
func New(r io.Reader) *Decoder {
	return &Decoder{
		scanner: bufio.NewScanner(r),
		line:    nil,
		done:    false,
	}
}

func (s *Decoder) peek() {
	if !s.scanner.Scan() {
		s.done = true
	}
	line := sanitiseLine(s.scanner.Text())
	s.line = &line
	if len(*s.line) == 0 && !s.done {
		s.peek()
	}
}

func sanitiseLine(line string) string {
	i := strings.Index(line, "#")
	if i >= 0 {
		line = line[:i]
	}
	return strings.Trim(line, " ")
}

// More returns if there are available CODEOWNERS lines to be scanned
func (s *Decoder) More() bool {
	s.peek()
	return !s.done
}

// Token parses the next available line in the CODEOWNERS file
func (s *Decoder) Token() Token {
	line := strings.ReplaceAll(*s.line, "\\ ", "\\s")

	data := strings.Split(line, " ")

	for i := range data {
		data[i] = strings.ReplaceAll(data[i], "\\s", " ")
	}

	return Token{
		path:   data[0],
		owners: data[1:],
	}
}

// Token providers reading capabilities for every CODEOWNERS line
type Token struct {
	path   string
	owners []string
}

// Path returns the file path pattern
func (t Token) Path() string {
	return t.path
}

// Owners returns the owners
func (t Token) Owners() []string {
	return t.owners
}
