package codeowners

import (
	"bufio"
	"io"
)

// Decoder providers functionality to read CODEOWNERS data
type Decoder struct {
	scanner *bufio.Scanner
	line    string
	lineNo  int
	done    bool
}

// NewDecoder generates a new Decoder instance. The reader should contain the contents of the CODEOWNERS file
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner: bufio.NewScanner(r),
		line:    "",
		lineNo:  0,
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
	d.lineNo++
	if len(line) == 0 && !d.done {
		d.peek()
	}
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
func (d *Decoder) Token() (Token, int) {
	pattern, owners := ParseLine(d.line)

	return Token{
		path:   pattern,
		owners: owners,
	}, d.lineNo
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
