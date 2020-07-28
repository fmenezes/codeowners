package codeowners

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
func Check(directory string, checkers ...string) ([]CheckResult, error) {

	fileLocation, result := findCodeownersFile(directory)
	if result != nil {
		return []CheckResult{*result}, nil
	}

	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := []CheckResult{}
	decoder := NewDecoder(file)
	for decoder.More() {
		token, lineNo := decoder.Token()
		for _, checker := range checkers {
			c, ok := availableCheckers[checker]
			if !ok {
				return nil, fmt.Errorf("'%s' not found", checker)
			}
			lineResults := c.CheckLine(lineNo, token.Path(), token.Owners()...)
			if lineResults != nil {
				results = append(results, lineResults...)
			}
		}
	}

	if len(results) > 0 {
		return results, nil
	}

	return nil, nil
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	return !os.IsNotExist(err) && !info.IsDir()
}

func findCodeownersFile(dir string) (string, *CheckResult) {
	codeownersLocation := ""

	filesFound := []string{}
	for _, fileLocation := range DefaultLocations {
		currentFile := filepath.Join(dir, fileLocation)
		if fileExists(currentFile) {
			filesFound = append(filesFound, fileLocation)
			if len(codeownersLocation) == 0 {
				codeownersLocation = currentFile
			}
		}
	}

	if len(filesFound) == 0 {
		return "", &CheckResult{Message: "No CODEOWNERS file found", Severity: Error, CheckName: "NoCodeowners"}
	}

	if len(filesFound) > 1 {
		return "", &CheckResult{Message: fmt.Sprintf("Multiple CODEOWNERS files found (%s)", strings.Join(filesFound, ", ")), Severity: Warning, CheckName: "MultipleCodeowners"}
	}

	return codeownersLocation, nil
}
