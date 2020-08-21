package checkers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/fmenezes/codeowners"
	"github.com/google/go-github/github"
)

type accessApi struct {
	directory   string
	tokenType   string
	token       string
	repoURL     string
	repoOwner   string
	repoName    string
	accessLevel string

	client *github.Client
	ctx    context.Context
}

func (a *accessApi) extractRepoData() {
	r := regexp.MustCompile(`github.com[\:\/]([A-Za-z0-9-]+)\/([A-Za-z0-9-]+)\.git`)
	data := r.FindStringSubmatch(a.repoURL)
	a.repoOwner = data[1]
	a.repoName = data[2]
}

func (a accessApi) fetchTeamAccess(org, team string) (string, error) {
	teams, _, err := a.client.Repositories.ListTeams(a.ctx, a.repoOwner, a.repoName, nil)
	if err != nil {
		return "", err
	}
	for _, t := range teams {
		if t.GetOrganization().GetLogin() == org && t.GetSlug() == team {
			return t.GetPermission(), nil
		}
	}
	return "none", nil
}

func (a accessApi) hasWriteAccess() bool {
	switch a.accessLevel {
	case "admin", "push", "maintain", "write", "email": // allowing emails to pass
		return true
	}
	return false
}

func (a accessApi) fetchUserAccess(user string) (string, error) {
	isCollaborator, _, err := a.client.Repositories.IsCollaborator(a.ctx, a.repoOwner, a.repoName, user)
	if err != nil {
		return "", err
	}
	if !isCollaborator {
		return "none", nil
	}
	permissionLevel, _, err := a.client.Repositories.GetPermissionLevel(a.ctx, a.repoOwner, a.repoName, user)
	if err != nil {
		return "", err
	}
	return permissionLevel.GetPermission(), nil
}

func (a accessApi) findUserFromEmail(email string) (string, error) {
	res, _, err := a.client.Search.Users(a.ctx, fmt.Sprintf("%s in:email type:user", email), nil)
	if err != nil {
		return "", err
	}
	if res.GetTotal() > 0 {
		return res.Users[0].GetLogin(), nil
	}
	return "", nil
}

func (a *accessApi) fetchAccess(user string) error {
	var err error
	login := ""
	if string(user[0]) == "@" {
		login = strings.TrimPrefix(user, "@")
	} else {
		login, err = a.findUserFromEmail(user)
		if err != nil {
			return err
		}
	}

	if len(login) == 0 {
		a.accessLevel = "email"
		return nil
	}

	if strings.Index(login, "/") >= 0 {
		parts := strings.Split(login, "/")
		access, err := a.fetchTeamAccess(parts[0], parts[1])
		if err != nil {
			return err
		}
		a.accessLevel = access
		return nil
	}

	access, err := a.fetchUserAccess(login)
	if err != nil {
		return err
	}
	a.accessLevel = access
	return nil
}

func ownerHasWriteAccess(options codeowners.ValidatorOptions, user string) (bool, error) {
	a := accessApi{
		ctx:       context.Background(),
		directory: options.Directory,
		token:     options.GithubToken,
		tokenType: options.GithubTokenType,
	}
	a.initiateClient()
	defer a.terminateClient()

	err := a.extractRepoURL()
	if err != nil {
		return false, err
	}
	a.extractRepoData()

	err = a.fetchAccess(user)
	if err != nil {
		return false, err
	}
	return a.hasWriteAccess(), nil
}
