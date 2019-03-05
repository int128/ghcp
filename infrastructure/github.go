package infrastructure

import (
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewGitHubClient returns a GitHubClient.
func NewGitHubClient() (infrastructure.GitHubClient, infrastructure.GitHubClientConfig) {
	var token oauth2.Token
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   http.DefaultTransport,
			Source: oauth2.StaticTokenSource(&token),
		},
	}
	v4 := githubv4.NewClient(hc)
	v3 := github.NewClient(hc)
	c := &GitHubClient{v4, v3.Git, &token}
	return c, c
}

type GitHubClient struct {
	*githubv4.Client
	*github.GitService

	token *oauth2.Token
}

func (c *GitHubClient) SetToken(token string) {
	c.token.AccessToken = token
}
