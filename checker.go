package codeowners

import (
	"bufio"
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

// CheckerOptions provide input arguments for checkers to use
type CheckerOptions struct {
	Directory              string
	CodeownersFileLocation string
}

// Checker provides tools for validating CODEOWNER file contents
type Checker interface {
	NewValidator(options CheckerOptions) CheckerValidator
}

// CheckerValidator provides tools for validating CODEOWNER file contents
type CheckerValidator interface {
	ValidateLine(lineNo int, line string) []CheckResult
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

// Position provides structured way to evaluate where a given validation result is located in the CODEOWNERs file
type Position struct {
	FilePath    string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

// String formats the position data
func (p Position) String() string {
	output := fmt.Sprintf("%d", p.StartLine)
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
	Directory string
	Checkers  []string
}

// Check evaluates the file contents against the checkers and return the results back.
func Check(options CheckOptions) ([]CheckResult, error) {

	fileLocation, result := findCodeownersFile(options.Directory)
	if result != nil {
		return []CheckResult{*result}, nil
	}

	file, err := os.Open(filepath.Join(options.Directory, fileLocation))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := []CheckResult{}
	scanner := bufio.NewScanner(file)
	lineNo := 0

	startedCheckers := make(map[string]CheckerValidator)
	for _, checker := range options.Checkers {
		c, ok := availableCheckers[checker]
		if !ok {
			return nil, fmt.Errorf("'%s' not found", checker)
		}
		startedCheckers[checker] = c.NewValidator(CheckerOptions{
			Directory:              options.Directory,
			CodeownersFileLocation: fileLocation,
		})
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		for _, c := range startedCheckers {
			lineResults := c.ValidateLine(lineNo, line)
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
				codeownersLocation = fileLocation
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
