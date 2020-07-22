package codeowners

import (
	"fmt"
	"io"
)

var availableCheckers map[string]Checker

func init() {
	availableCheckers = make(map[string]Checker)
}

// AvailableCheckers returns list of registered checkers
func AvailableCheckers() []string {
	names := make([]string, len(availableCheckers))
	i := 0
	for checkerName := range availableCheckers {
		names[i] = checkerName
		i++
	}
	return names
}

// RegisterChecker adds checker to be used later when checking CODEOWNERS files
func RegisterChecker(name string, checker Checker) error {
	_, found := availableCheckers[name]
	if found {
		return fmt.Errorf("Checker %s already exists", name)
	}
	availableCheckers[name] = checker
	return nil
}

// Checker provides tools for validating CODEOWNER file contents
type Checker interface {
	CheckLine(lineNo int, filePattern string, owners ...string) []CheckResult
}

// SeverityLevel exposes all possible levels of severity check results
type SeverityLevel int

// All possible severiy levels
const (
	Error   SeverityLevel = iota // Error serverity level
	Warning                      // Warning serverity level
)

// String returns the string representation of this severity level
func (l SeverityLevel) String() string {
	return [...]string{"Error", "Warning"}[l]
}

// CheckResult provides structured way to evaluate results of a CODEOWNERS validation check
type CheckResult struct {
	LineNo    int
	Message   string
	Severity  SeverityLevel
	CheckName string
}

// Check evaluates the file contents against the checkers and return the results back.
func Check(r io.Reader, checkers ...string) ([]CheckResult, error) {
	results := []CheckResult{}
	decoder := NewDecoder(r)
	for decoder.More() {
		token, lineNo := decoder.Token()
		for _, checker := range checkers {
			c, ok := availableCheckers[checker]
			if !ok {
				return nil, fmt.Errorf("'%s' not found", checker)
			}
			lineResults := c.CheckLine(lineNo, token.Path(), token.Owners()...)
			results = append(results, lineResults...)
		}
	}

	return results, nil
}
