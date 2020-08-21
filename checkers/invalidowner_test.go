package checkers_test

import (
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
	"github.com/fmenezes/codeowners/checkers"
)

func TestInvalidOwnerCheckInvalidLongUsername(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern myusernamemyusernamemyusernamemyusernamemyusername",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   63,
			},
			Message:   "Owner 'myusernamemyusernamemyusernamemyusernamemyusername' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestInvalidOwnerCheckInvalidNoAt(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern invalid-owner",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   26,
			},
			Message:   "Owner 'invalid-owner' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestInvalidOwnerCheckInvalidHyphens(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @invalid--owner",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   28,
			},
			Message:   "Owner '@invalid--owner' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestInvalidOwnerCheckInvalidFormat(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @org/invalid/owner",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   31,
			},
			Message:   "Owner '@org/invalid/owner' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}
func TestInvalidOwnerCheckInvalidTrailingHyphen(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @invalid-owner-",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   28,
			},
			Message:   "Owner '@invalid-owner-' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestInvalidOwnerCheckMultipleInvalid(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern invalid-owner another-invalid-owner",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   26,
			},
			Message:   "Owner 'invalid-owner' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 27,
				EndLine:     1,
				EndColumn:   48,
			},
			Message:   "Owner 'another-invalid-owner' is invalid",
			Severity:  codeowners.Error,
			CheckName: "InvalidOwner",
		},
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestInvalidOwnerCheckPassUser(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @valid-owner",
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestInvalidOwnerCheckPassEmail(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern email@server.com",
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestInvalidOwnerCheckPassUserOrg(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @org/valid-owner",
	}

	checker := checkers.InvalidOwner{}
	validator := checker.NewValidator(codeowners.CheckerOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}
