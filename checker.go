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

	validators := make(map[string]Validator)
	for _, checker := range options.Checkers {
		c, ok := availableCheckers[checker]
		if !ok {
			return nil, fmt.Errorf("'%s' not found", checker)
		}
		validators[checker] = c.NewValidator(ValidatorOptions{
			Directory:              options.Directory,
			CodeownersFileLocation: fileLocation,
			GithubToken:            options.GithubToken,
			GithubTokenType:        options.GithubTokenType,
		})
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		for _, c := range validators {
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
		return "", &CheckResult{Position: Position{FilePath: "CODEOWNERS"}, Message: "No CODEOWNERS file found", Severity: Error, CheckName: "NoCodeowners"}
	}

	if len(filesFound) > 1 {
		return "", &CheckResult{Position: Position{FilePath: codeownersLocation}, Message: fmt.Sprintf("Multiple CODEOWNERS files found (%s)", strings.Join(filesFound, ", ")), Severity: Warning, CheckName: "MultipleCodeowners"}
	}

	return codeownersLocation, nil
}
