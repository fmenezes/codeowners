package codeowners_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/fmenezes/codeowners"
)

func exec(input string) ([][]string, int) {
	scanner := codeowners.New(strings.NewReader(input))
	got := [][]string{}
	c := 0
	for scanner.More() {
		c++
		token := scanner.Token()
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
