package codeowners_test

import (
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
)

func TestParseLine(t *testing.T) {
	testCases := []struct {
		input       string
		wantPattern string
		wantOwners  []string
	}{
		{
			input:       "* test@example.org",
			wantPattern: "*",
			wantOwners:  []string{"test@example.org"},
		},
		{
			input:       "filepattern test@example.org @owner @company/team",
			wantPattern: "filepattern",
			wantOwners:  []string{"test@example.org", "@owner", "@company/team"},
		},
		{
			input:       "file\\ with\\ spaces @owner",
			wantPattern: "file\\ with\\ spaces",
			wantOwners:  []string{"@owner"},
		},
		{
			input:       "    filepattern    ",
			wantPattern: "filepattern",
			wantOwners:  nil,
		},
		{
			input:       "filepattern @owner # comments",
			wantPattern: "filepattern",
			wantOwners:  []string{"@owner"},
		},
		{
			input:       "",
			wantPattern: "",
			wantOwners:  nil,
		},
		{
			input:       "# only comments on the line",
			wantPattern: "",
			wantOwners:  nil,
		},
	}

	for _, testCase := range testCases {
		gotPattern, gotOwners := codeowners.ParseLine(testCase.input)
		if gotPattern != testCase.wantPattern || !reflect.DeepEqual(gotOwners, testCase.wantOwners) {
			t.Errorf("Input: %s, Want: %s, %v, Got: %s, %v", testCase.input, testCase.wantPattern, testCase.wantOwners, gotPattern, gotOwners)
		}
	}
}
