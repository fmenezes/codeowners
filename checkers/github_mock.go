// +build unit

package checkers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v32/github"
)

func (a *accessAPI) extractRepoURL() error {
	if a.directory == "bad" {
		return errors.New("Mocked error")
	}
	a.repoURL = "git@github.com:owner/repo.git"
	return nil
}

var server *httptest.Server

func mockedServer(w http.ResponseWriter, r *http.Request) {
	h := w.Header()

	isCollaboratorUrlRegex := regexp.MustCompile(`/repos/([^/]+)/([^/]+)/collaborators/([^/]+)$`)
	parts := isCollaboratorUrlRegex.FindStringSubmatch(r.URL.String())
	if len(parts) > 0 {
		switch parts[3] {
		case "noOwner":
			w.WriteHeader(404)
			return
		case "noOwnerWithError":
			h.Add("Content-Type", "text/plain")
			w.WriteHeader(504)
			c, _ := ioutil.ReadFile(filepath.Join("..", "test", "error.txt"))
			w.Write(c)
			return
		default:
			w.WriteHeader(204)
			return
		}
	}

	permissionUrlRegex := regexp.MustCompile(`/repos/([^/]+)/([^/]+)/collaborators/([^/]+)/permission$`)
	parts = permissionUrlRegex.FindStringSubmatch(r.URL.String())
	if len(parts) > 0 {
		switch parts[3] {
		case "ownerWithAccess":
			h.Add("Content-Type", "application/json")
			w.WriteHeader(200)
			c, _ := ioutil.ReadFile(filepath.Join("..", "test", "fixtures", "permission", "pass.json"))
			w.Write(c)
			return
		case "ownerWithError":
			h.Add("Content-Type", "text/plain")
			w.WriteHeader(504)
			c, _ := ioutil.ReadFile(filepath.Join("..", "test", "error.txt"))
			w.Write(c)
			return
		default:
			h.Add("Content-Type", "application/json")
			w.WriteHeader(200)
			c, _ := ioutil.ReadFile(filepath.Join("..", "test", "fixtures", "permission", "deny.json"))
			w.Write(c)
			return
		}
	}

	teamUrlRegex := regexp.MustCompile(`/repos/([^/]+)/([^/]+)/teams$`)
	parts = teamUrlRegex.FindStringSubmatch(r.URL.String())
	if len(parts) > 0 {
		h.Add("Content-Type", "application/json")
		w.WriteHeader(200)
		c, _ := ioutil.ReadFile(filepath.Join("..", "test", "fixtures", "team", "pass.json"))
		w.Write(c)
	}

	searchUrlRegex := regexp.MustCompile(`/search/users\?q=(.+)$`)
	parts = searchUrlRegex.FindStringSubmatch(r.URL.String())
	if len(parts) > 0 {
		h.Add("Content-Type", "application/json")
		w.WriteHeader(200)
		f := "found.json"
		if strings.Index(parts[1], "notfound") >= 0 {
			f = "notFound.json"
		}
		c, _ := ioutil.ReadFile(filepath.Join("..", "test", "fixtures", "search", f))
		w.Write(c)
	}
}

func (a *accessAPI) initiateClient() {
	server = httptest.NewServer(http.HandlerFunc(mockedServer))
	client, _ := github.NewEnterpriseClient(server.URL, server.URL, nil)
	a.client = client
}

func (a *accessAPI) terminateClient() {
	server.Close()
}
