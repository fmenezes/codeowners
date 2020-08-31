package codeowners

import (
	"fmt"
)

// ValidatorOptions provide input arguments for checkers to use
type ValidatorOptions struct {
	Directory              string
	CodeownersFileLocation string
	GithubTokenType        string
	GithubToken            string
}

// Checker provides tools for validating CODEOWNER file contents
type Checker interface {
	NewValidator(options ValidatorOptions) Validator
}

// Validator provides tools for validating CODEOWNER file contents
type Validator interface {
	ValidateLine(lineNo int, line string) []CheckResult
}

// SeverityLevel exposes all possible levels of severity check results
type SeverityLevel int

// All possible severiy levels
const (
	Error   SeverityLevel = iota // Error serverity level
	Warning                      // Warning serverity level
)

// Name returns the string representation of this severity level
func (l SeverityLevel) Name() string {
	return [...]string{"Error", "Warning"}[l]
}

// Position provides structured way to evaluate where a given validation result is located in the CODEOWNERs file
type Position struct {
	FilePath    string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

// Format converts the position data into a string
func (p Position) Format() string {
	output := fmt.Sprintf("%s %d", p.FilePath, p.StartLine)
	if p.StartColumn >= 1 {
		output = fmt.Sprintf("%s:%d", output, p.StartColumn)
	}
	if p.EndLine > p.StartLine {
		output = fmt.Sprintf("%s-%d:%d", output, p.EndLine, p.EndColumn)
	} else if p.StartColumn >= 1 && p.EndColumn > p.StartColumn {
		output = fmt.Sprintf("%s-%d", output, p.EndColumn)
	}

	return output
}

// CheckResult provides structured way to evaluate results of a CODEOWNERS validation check
type CheckResult struct {
	Position  Position
	Message   string
	Severity  SeverityLevel
	CheckName string
}

// CheckOptions provides parameters for running a list of checks
type CheckOptions struct {
	Directory       string
	Checkers        []string
	GithubTokenType string
	GithubToken     string
}
