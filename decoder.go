package codeowners

import (
	"bufio"
	"io"
	"strings"
)

// Decoder providers functionality to read CODEOWNERS data
type Decoder struct {
	scanner *bufio.Scanner
	line    string
	done    bool
}

// NewDecoder generates a new CodeOwnersScanner instance. The reader should be a reader containing the contents of the CODEOWNERS file
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner: bufio.NewScanner(r),
		line:    "",
		done:    false,
	}
}

// peek will scan the next line
func (d *Decoder) peek() {
	if !d.scanner.Scan() {
		d.done = true
		return
	}
	d.line = d.scanner.Text()
	line := sanitiseLine(d.line)
	if len(line) == 0 && !d.done {
		d.peek()
	}
}

// sanitiseLine removes all empty space and comments from a given line
func sanitiseLine(line string) string {
	i := strings.Index(line, "#")
	if i >= 0 {
		line = line[:i]
	}
	return strings.Trim(line, " ")
}

// More returns true if there are available CODEOWNERS lines to be scanned.
// And also advances to the next line.
func (d *Decoder) More() bool {
	d.peek()
	return !d.done
}

// Token parses the next available line in the CODEOWNERS file.
// If More was never called it will return an empty token.
// After end of file Token will always return the last line.
func (d *Decoder) Token() Token {
	line := sanitiseLine(d.line)
	line = strings.ReplaceAll(line, "\\ ", "\\s")

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
