// +build !unit

package checkers

import (
	"os/exec"
	"strings"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func (a *accessAPI) extractRepoURL() error {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = a.directory
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	a.repoURL = strings.ReplaceAll(string(out), "\n", "")
	return nil
}

func (a *accessAPI) initiateClient() {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.token, TokenType: a.tokenType},
	)
	oauthClient := oauth2.NewClient(a.ctx, tokenSource)
	a.client = github.NewClient(oauthClient)
}

func (a *accessAPI) terminateClient() {}
