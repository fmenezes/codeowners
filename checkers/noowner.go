package checkers

import "github.com/fmenezes/codeowners"

const noOwnerCheckerName string = "NoOwner"

func init() {
	codeowners.RegisterChecker(noOwnerCheckerName, NoOwner{})
}

// NoOwner represents checker to decide validate presence of owners in each of CODEOWNERS lines
type NoOwner struct{}

// CheckLine runs this NoOwner's check against each line
func (c NoOwner) CheckLine(lineNo int, line string) []codeowners.CheckResult {
	var results []codeowners.CheckResult

	_, owners := codeowners.ParseLine(line)

	if len(owners) == 0 {
		results = []codeowners.CheckResult{
			{
				Position: codeowners.Position{
					StartLine:   lineNo,
					EndLine:     lineNo,
					StartColumn: 0,
					EndColumn:   0,
				},
				Message:   "No owners specified",
				Severity:  codeowners.Error,
				CheckName: noOwnerCheckerName,
			},
		}
	}

	return results
}
