package codeowners_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/fmenezes/codeowners"
)

func exec(input string) ([][]string, int) {
	decoder := codeowners.NewDecoder(strings.NewReader(input))
	got := [][]string{}
	c := 0
	for decoder.More() {
		c++
		token := decoder.Token()
		got = append(got, append([]string{token.Path()}, token.Owners()...))
	}
	return got, c
}

func assert(t *testing.T, input string, want [][]string) {
	got, gotCount := exec(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
	if gotCount != len(want) {
		t.Errorf("Input: %v, Want: %v scans, Got: %v scans", input, len(want), gotCount)
	}
}

func TestSimple(t *testing.T) {
	assert(t, `* test@example.org`, [][]string{
		{"*", "test@example.org"},
	})
}

func TestMultipleOwners(t *testing.T) {
	assert(t, `* test@example.org @owner @company/team`, [][]string{
		{"*", "test@example.org", "@owner", "@company/team"},
	})
}

func TestFilesWithSpaces(t *testing.T) {
	assert(t, `file\ with\ spaces @owner`, [][]string{
		{"file with spaces", "@owner"},
	})
}

func TestMultipleLines(t *testing.T) {
	assert(t, `* test@example.org
file @owner`, [][]string{
		{"*", "test@example.org"},
		{"file", "@owner"},
	})
}

func TestEmptyFile(t *testing.T) {
	assert(t, ``, [][]string{})
}

func TestEmptyLines(t *testing.T) {
	assert(t, `* test@example.org



file @owner




`, [][]string{
		{"*", "test@example.org"},
		{"file", "@owner"},
	})
}

func TestIgnoreComments(t *testing.T) {
	assert(t, `* test@example.org # comment
# comment
file @owner`, [][]string{
		{"*", "test@example.org"},
		{"file", "@owner"},
	})
}

func TestNoOwners(t *testing.T) {
	assert(t, `*`, [][]string{
		{"*"},
	})
}

func TestLastToken(t *testing.T) {
	decoder := codeowners.NewDecoder(strings.NewReader(`filepattern @owner`))
	if !decoder.More() {
		t.Error("More should be true")
	}
	for i := 0; i < 3; i++ { //calling 3 times to prove it always returns the last line
		token := decoder.Token()
		if token.Path() != "filepattern" {
			t.Error("Path should be 'filepattern'")
		}
		if len(token.Owners()) != 1 || token.Owners()[0] != "@owner" {
			t.Error("Owners should match ['@owner']")
		}

		if decoder.More() {
			t.Error("More should be false")
		}
	}
}

func TestMoreNotCalled(t *testing.T) {
	decoder := codeowners.NewDecoder(strings.NewReader(`filepattern @owner`))
	token := decoder.Token()
	if token.Path() != "" {
		t.Error("Path should be empty")
	}
	if len(token.Owners()) != 0 {
		t.Error("Owners should be empty")
	}
}

func ExampleDecoder() {
	decoder := codeowners.NewDecoder(strings.NewReader(`* test@example.org
filepattern @owner`))
	for decoder.More() {
		token := decoder.Token()
		fmt.Printf("File Pattern: %s\n", token.Path())
		fmt.Printf("Owners: %v\n", token.Owners())
	}
	// Output:
	// File Pattern: *
	// Owners: [test@example.org]
	// File Pattern: filepattern
	// Owners: [@owner]
}
