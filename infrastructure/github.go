package infrastructure

import (
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewGitHubClient returns a GitHubClient and GitHubClientInit.
func NewGitHubClient() (infrastructure.GitHubClient, infrastructure.GitHubClientInit) {
	var c GitHubClient
	c.Init(infrastructure.GitHubClientInitOptions{})
	return &c, &c
}

type GitHubClient struct {
	*githubv4.Client
	*github.GitService
}

// Init initializes this client with the options.
func (c *GitHubClient) Init(options infrastructure.GitHubClientInitOptions) {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   http.DefaultTransport,
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: options.Token}),
		},
	}
	v4 := githubv4.NewClient(hc)
	v3 := github.NewClient(hc)

	c.Client = v4
	c.GitService = v3.Git
}
