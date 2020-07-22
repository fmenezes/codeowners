package checkers_test

import (
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
	"github.com/fmenezes/codeowners/checkers"
)

func TestNoOwnerCheck(t *testing.T) {
	input := struct {
		lineNo  int
		pattern string
		owners  []string
	}{
		lineNo:  1,
		pattern: "filepattern",
		owners:  []string{},
	}
	want := []codeowners.CheckResult{
		{
			LineNo:    1,
			Message:   "No owners specified",
			Severity:  codeowners.Error,
			CheckName: "NoOwner",
		},
	}

	checker := checkers.NoOwner{}
	got := checker.CheckLine(input.lineNo, input.pattern, input.owners...)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}
