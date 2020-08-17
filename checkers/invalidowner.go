package checkers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fmenezes/codeowners"
)

const invalidOwnerCheckerName string = "InvalidOwner"

func init() {
	codeowners.RegisterChecker(invalidOwnerCheckerName, InvalidOwner{})
}

// InvalidOwner represents checker to decide validate owners in each of CODEOWNERS lines
type InvalidOwner struct{}

func ownerValid(owner string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if emailRegex.MatchString(owner) {
		return true
	}

	if owner[0] != byte('@') { // should be @company/user or @user
		return false
	}

	parts := strings.Split(owner[1:], "/")
	if len(parts) > 2 { // should be user or company/user
		return false
	}

	var githubUsernameRegex = regexp.MustCompile("^[A-Za-z0-9](?:-?[A-Za-z0-9])*$")
	for _, username := range parts {
		if len(username) > 39 || !githubUsernameRegex.MatchString(username) {
			return false
		}
	}

	return true
}

// CheckLine runs this InvalidOwner's check against each line
func (c InvalidOwner) CheckLine(file string, lineNo int, line string) []codeowners.CheckResult {
	var results []codeowners.CheckResult

	_, owners := codeowners.ParseLine(line)

	for _, owner := range owners {
		if ownerValid(owner) {
			continue
		}
		result := codeowners.CheckResult{
			Position: codeowners.Position{
				FilePath:    file,
				StartLine:   lineNo,
				EndLine:     lineNo,
				StartColumn: strings.Index(line, owner) + 1,
			},
			Message:   fmt.Sprintf("Owner '%s' is invalid", owner),
			Severity:  codeowners.Error,
			CheckName: invalidOwnerCheckerName,
		}
		result.Position.EndColumn = result.Position.StartColumn + len(owner)

		if results == nil {
			results = []codeowners.CheckResult{result}
		} else {
			results = append(results, result)
		}
	}

	return results
}
