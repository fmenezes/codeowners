package codeowners_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
	_ "github.com/fmenezes/codeowners/checkers"
)

const dummyCheckerName string = "dummy"

type dummyChecker struct {
}

func (c dummyChecker) CheckLine(lineNo int, filePath string, owners ...string) []codeowners.CheckResult {
	return []codeowners.CheckResult{
		{
			LineNo:    1,
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
	if codeowners.Error.String() != "Error" {
		t.Errorf("codeowners.Error.String() should evaluate to 'Error'")
	}
	if codeowners.Warning.String() != "Warning" {
		t.Errorf("codeowners.Warning.String() should evaluate to 'Warning'")
	}
}

func TestSimpleCheck(t *testing.T) {
	input := "./test/data/pass"
	want := []codeowners.CheckResult{
		{
			LineNo:    1,
			Message:   "Dummy Error",
			Severity:  codeowners.Error,
			CheckName: dummyCheckerName,
		},
	}

	codeowners.RegisterChecker(dummyCheckerName, dummyChecker{})
	got, err := codeowners.Check(input, dummyCheckerName)
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Input %s, Want %v, Got %v", input, want, got)
	}
}

func TestNoProblemsFound(t *testing.T) {
	input := "./test/data/pass"
	got, err := codeowners.Check(input)
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if got != nil {
		t.Errorf("Input %s, Want %v, Got %v", input, nil, got)
	}
}

func TestCheckerNotFound(t *testing.T) {
	input := "./test/data/pass"
	_, err := codeowners.Check(input, "NonExistentChecker")
	if err == nil {
		t.Error("Should have errored")
	}
}

func TestNoCodeownersCheck(t *testing.T) {
	input := "./test/data"
	want := []codeowners.CheckResult{
		{
			LineNo:    0,
			Message:   "No CODEOWNERS file found",
			Severity:  codeowners.Error,
			CheckName: "NoCodeowners",
		},
	}

	got, err := codeowners.Check(input, dummyCheckerName)
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
			LineNo:    0,
			Message:   "Multiple CODEOWNERS files found (CODEOWNERS, docs/CODEOWNERS)",
			Severity:  codeowners.Warning,
			CheckName: "MultipleCodeowners",
		},
	}

	got, err := codeowners.Check(input, dummyCheckerName)
	if err != nil {
		t.Errorf("Input %s, Error %v", input, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Input %s, Want %v, Got %v", input, want, got)
	}
}

func ExampleCheck() {
	checks, err := codeowners.Check(".", codeowners.AvailableCheckers()...)
	if err != nil {
		panic(err)
	}
	for _, check := range checks {
		fmt.Printf("%d ::%s:: %s [%s]\n", check.LineNo, check.Severity, check.Message, check.CheckName)
	}
	//Output:
	//0 ::Error:: No CODEOWNERS file found [NoCodeowners]
}
