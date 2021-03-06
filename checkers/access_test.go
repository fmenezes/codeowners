package checkers_test

import (
	"reflect"
	"testing"

	"github.com/fmenezes/codeowners"
	"github.com/fmenezes/codeowners/checkers"
)

func TestAccessCheck(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @owner",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   19,
			},
			Message:   "Owner '@owner' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}

	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckError(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @ownerWithError",
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
			Message:   "Owner '@ownerWithError' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}

	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckPass(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @ownerWithAccess",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestAccessCheckPassInvalidOwners(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern owner",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestAccessCheckPassNoOwners(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestAccessCheckPassHttpGitConfig(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @ownerWithAccess",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              "http",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}
func TestAccessCheckNoTokenPass(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @ownerWithAccess",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              "bad",
		CodeownersFileLocation: "bad",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestAccessCheckInvalidDir(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @ownerWithAccess",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "bad",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   29,
			},
			Message:   "Owner '@ownerWithAccess' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              "bad",
		CodeownersFileLocation: "bad",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckNoCollaborator(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @noOwner",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   21,
			},
			Message:   "Owner '@noOwner' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckIsCollaboratorError(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @noOwnerWithError",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   30,
			},
			Message:   "Owner '@noOwnerWithError' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckTeamPass(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @github/justice-league",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}

func TestAccessCheckTeamDeny(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern @org/denyTeam",
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
			Message:   "Owner '@org/denyTeam' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckEmailPass(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern notfound@example.com",
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if got != nil {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, nil, got)
	}
}
func TestAccessCheckEmailDeny(t *testing.T) {
	input := struct {
		lineNo int
		line   string
	}{
		lineNo: 1,
		line:   "filepattern found@example.com",
	}
	want := []codeowners.CheckResult{
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   1,
				StartColumn: 13,
				EndLine:     1,
				EndColumn:   30,
			},
			Message:   "Owner 'found@example.com' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := validator.ValidateLine(input.lineNo, input.line)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}

func TestAccessCheckDenyMemo(t *testing.T) {
	input := []struct {
		lineNo int
		line   string
	}{
		{
			lineNo: 1,
			line:   "filepattern @org/denyTeam",
		},
		{
			lineNo: 2,
			line:   "filepattern2 @org/denyTeam",
		},
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
			Message:   "Owner '@org/denyTeam' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
		{
			Position: codeowners.Position{
				FilePath:    "CODEOWNERS",
				StartLine:   2,
				StartColumn: 14,
				EndLine:     2,
				EndColumn:   27,
			},
			Message:   "Owner '@org/denyTeam' has no write access",
			Severity:  codeowners.Error,
			CheckName: "Access",
		},
	}
	checker := checkers.Access{}
	validator := checker.NewValidator(codeowners.ValidatorOptions{
		Directory:              ".",
		CodeownersFileLocation: "CODEOWNERS",
		GithubToken:            "token",
	})
	got := []codeowners.CheckResult{}
	for _, inputLine := range input {
		got = append(got, validator.ValidateLine(inputLine.lineNo, inputLine.line)...)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Input: %v, Want: %v, Got: %v", input, want, got)
	}
}
