package checkers

import (
	"fmt"
	"strings"

	"github.com/fmenezes/codeowners"
)

const accessCheckerName string = "Access"

func init() {
	codeowners.RegisterChecker(accessCheckerName, Access{})
}

// Access represents checker to validate if an owner has access to repo
type Access struct{}

// NewValidator returns validating capabilities for this checker
func (c Access) NewValidator(options codeowners.ValidatorOptions) codeowners.Validator {
	return accessValidator{
		options:    options,
		accessMemo: make(map[string]bool),
	}
}

type accessValidator struct {
	options    codeowners.ValidatorOptions
	accessMemo map[string]bool
}

// ValidateLine runs this NoOwner's check against each line
func (v accessValidator) ValidateLine(lineNo int, line string) []codeowners.CheckResult {
	results := []codeowners.CheckResult{}

	_, owners := codeowners.ParseLine(line)

	if len(owners) == 0 {
		return nil
	}

	for _, owner := range owners {
		if !ownerValid(owner) {
			continue
		}
		writeAccess, found := v.accessMemo[owner]
		if !found {
			writeAccess, _ = ownerHasWriteAccess(v.options, owner)
			v.accessMemo[owner] = writeAccess
		}
		if !writeAccess {
			result := codeowners.CheckResult{
				Position: codeowners.Position{
					FilePath:    v.options.CodeownersFileLocation,
					StartLine:   lineNo,
					EndLine:     lineNo,
					StartColumn: strings.Index(line, owner) + 1,
				},
				Message:   fmt.Sprintf("Owner '%s' has no write access", owner),
				Severity:  codeowners.Error,
				CheckName: accessCheckerName,
			}
			result.Position.EndColumn = result.Position.StartColumn + len(owner)

			results = append(results, result)
		}
	}

	if len(results) > 0 {
		return results
	}

	return nil
}
