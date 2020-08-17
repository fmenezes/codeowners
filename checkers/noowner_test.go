package checkers_test

import (
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
	"github.com/fmenezes/codeowners/checkers"
)

func TestNoOwnerCheck(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 0,
				EndLine:     1,
				EndColumn:   0,
			},
			Message:   "No owners specified",
			Severity:  codeowners.Error,
			CheckName: "NoOwner",
		},
	}

	checker := checkers.NoOwner{}
	got := checker.CheckLine("CODEOWNERS", input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}
