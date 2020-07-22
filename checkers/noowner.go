package checkers

import "github.com/fmenezes/codeowners"

const noOwnerCheckerName string = "NoOwner"

func init() {
	codeowners.RegisterChecker(noOwnerCheckerName, NoOwner{})
}

// NoOwner represents checker to decide validate presence of owners in each of CODEOWNERS lines
type NoOwner struct{}

// CheckLine runs this NoOwner's check against each line
func (c NoOwner) CheckLine(lineNo int, pattern string, owners ...string) []codeowners.CheckResult {
	results := []codeowners.CheckResult{}

	if len(owners) == 0 {
		results = append(results, codeowners.CheckResult{
			LineNo:    lineNo,
			Message:   "No owners specified",
			Severity:  codeowners.Error,
			CheckName: noOwnerCheckerName,
		})

	}
	return results
}
