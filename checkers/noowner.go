package checkers

import "github.com/fmenezes/codeowners"

const noOwnerCheckerName string = "NoOwner"

func init() {
	codeowners.RegisterChecker(noOwnerCheckerName, NoOwner{})
}

// NoOwner represents checker to decide validate presence of owners in each of CODEOWNERS lines
type NoOwner struct{}

// NewValidator returns validating capabilities for this checker
func (c NoOwner) NewValidator(options codeowners.ValidatorOptions) codeowners.Validator {
	return noOwnerValidator{
		options: options,
	}
}

type noOwnerValidator struct {
	options codeowners.ValidatorOptions
}

// ValidateLine runs this NoOwner's check against each line
func (v noOwnerValidator) ValidateLine(lineNo int, line string) []codeowners.CheckResult {
	var results []codeowners.CheckResult

	_, owners := codeowners.ParseLine(line)

	if len(owners) == 0 {
		results = []codeowners.CheckResult{
			{
				Position: codeowners.Position{
					FilePath:    v.options.CodeownersFileLocation,
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
