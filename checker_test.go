package codeowners_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
)

const dummyCheckerName string = "dummy"

type dummyChecker struct {
}

type dummyCheckerValidator struct {
	codeownersFileLocation string
	directory              string
}

func (c dummyChecker) NewValidator(options codeowners.ValidatorOptions) codeowners.Validator {
	return dummyCheckerValidator{
		codeownersFileLocation: options.CodeownersFileLocation,
		directory:              options.Directory,
	}
}

func (c dummyCheckerValidator) ValidateLine(lineNo int, line string) []codeowners.CheckResult {
	return []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:  c.codeownersFileLocation,
				StartLine: lineNo,
				EndLine:   lineNo,
			},
			Message:   "Dummy Error",
			Severity:  codeowners.Error,
			CheckName: dummyCheckerName,
		},
	}
}

func TestRegisterChecker(t *testing.T) {
	err := codeowners.RegisterChecker(dummyCheckerName, dummyChecker{})
	if err != nil {
		t.Error(err)
	}
	found := false
	for _, checker := range codeowners.AvailableCheckers() {
		if checker == dummyCheckerName {
			found = true
		}
	}
	if !found {
		t.Errorf("%s not properly registered", dummyCheckerName)
	}
}

func TestRegisterCheckerAgain(t *testing.T) {
	codeowners.RegisterChecker(dummyCheckerName, dummyChecker{})
	err := codeowners.RegisterChecker(dummyCheckerName, dummyChecker{})
	if err == nil {
		t.Errorf("%s should be already registered, expecting an error", dummyCheckerName)
	}
}

func TestSeverityLevelLabels(t *testing.T) {
	if codeowners.Error.Name() != "Error" {
		t.Errorf("codeowners.Error.String() should evaluate to 'Error'")
	}
	if codeowners.Warning.Name() != "Warning" {
		t.Errorf("codeowners.Warning.String() should evaluate to 'Warning'")
	}
}

func TestPositionFormat(t *testing.T) {
	testCases := []struct {
		input codeowners.Position
		want  string
	}{
		{
			input: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 1,
				EndLine:     2,
				EndColumn:   2,
			},
			want: "CODEOWNERS 1:1-2:2",
		},
		{
			input: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 1,
				EndLine:     1,
				EndColumn:   2,
			},
			want: "CODEOWNERS 1:1-2",
		},
		{
			input: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 1,
				EndLine:     1,
				EndColumn:   1,
			},
			want: "CODEOWNERS 1:1",
		},
		{
			input: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 0,
				EndLine:     1,
				EndColumn:   0,
			},
			want: "CODEOWNERS 1",
		},
		{
			input: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 0,
				EndLine:     0,
				EndColumn:   0,
			},
			want: "CODEOWNERS 1",
		},
		{
			input: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   0,
				StartColumn: 0,
				EndLine:     0,
				EndColumn:   0,
			},
			want: "CODEOWNERS 0",
		},
	}

	for _, testCase := range testCases {
		got := testCase.input.Format()
		if got != testCase.want {
			t.Errorf("Input: %v, Want: %v, Got: %v", testCase.input, testCase.want, got)
		}
	}
}

func TestSimpleCheck(t *testing.T) {
	input := "./test/data/pass"
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 0,
				EndLine:     1,
				EndColumn:   0,
			},
			Message:   "Dummy Error",
			Severity:  codeowners.Error,
			CheckName: dummyCheckerName,
		},
	}

	codeowners.RegisterChecker(dummyCheckerName, dummyChecker{})
	got, err := codeowners.Check(codeowners.CheckOptions{
		Directory: input,
		Checkers:  []string{dummyCheckerName},
	})
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Input %s, Want %v, Got %v", input, want, got)
	}
}

func TestNoProblemsFound(t *testing.T) {
	input := "./test/data/pass"
	got, err := codeowners.Check(codeowners.CheckOptions{
		Directory: input,
	})
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if got != nil {
		t.Errorf("Input %s, Want %v, Got %v", input, nil, got)
	}
}

func TestCheckerNotFound(t *testing.T) {
	input := "./test/data/pass"
	_, err := codeowners.Check(codeowners.CheckOptions{
		Directory: input,
		Checkers:  []string{"NonExistentChecker"},
	})
	if err == nil {
		t.Error("Should have errored")
	}
}

func TestNoCodeownersCheck(t *testing.T) {
	input := "./test/data"
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath: "CODEOWNERS",
			},
			Message:   "No CODEOWNERS file found",
			Severity:  codeowners.Error,
			CheckName: "NoCodeowners",
		},
	}

	got, err := codeowners.Check(codeowners.CheckOptions{
		Directory: input,
		Checkers:  []string{dummyCheckerName},
	})
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Input %s, Want %v, Got %v", input, want, got)
	}
}

func TestMultipleCodeownersCheck(t *testing.T) {
	input := "./test/data/multiple_codeowners"
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath: "CODEOWNERS",
			},
			Message:   "Multiple CODEOWNERS files found (CODEOWNERS, docs/CODEOWNERS)",
			Severity:  codeowners.Warning,
			CheckName: "MultipleCodeowners",
		},
	}

	got, err := codeowners.Check(codeowners.CheckOptions{
		Directory: input,
		Checkers:  []string{dummyCheckerName},
	})
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Input %s, Want %v, Got %v", input, want, got)
	}
}

func ExampleCheck() {
	checks, err := codeowners.Check(codeowners.CheckOptions{
		Directory: ".",
		Checkers:  codeowners.AvailableCheckers(),
	})
	if err != nil {
		panic(err)
	}
	for _, check := range checks {
		fmt.Printf("%s ::%s:: %s [%s]\n", check.Position.Format(), check.Severity.Name(), check.Message, check.CheckName)
	}
	//Output:
	//CODEOWNERS 0 ::Error:: No CODEOWNERS file found [NoCodeowners]
}
