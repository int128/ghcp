package github

import (
	"net/http"
	"net/url"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/infrastructure"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewClient returns a Client and GitHubClientInit.
// Caller must call GitHubClientInit.Init() before an API invocation.
func NewClient() (infrastructure.GitHubClient, infrastructure.GitHubClientInit) {
	var c Client
	return &c, &c
}

type Client struct {
	*githubv4.Client
	*github.GitService
}

// Init initializes this client with the options.
func (c *Client) Init(o infrastructure.GitHubClientInitOptions) error {
	v4, v3, err := c.newClients(o)
	if err != nil {
		return errors.Wrapf(err, "error while initializing GitHub client")
	}
	c.Client = v4
	c.GitService = v3.Git
	return nil
}

func (c *Client) newClients(o infrastructure.GitHubClientInitOptions) (*githubv4.Client, *github.Client, error) {
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Base:   http.DefaultTransport,
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: o.Token}),
		},
	}
	if o.URLv3 != "" {
		// https://developer.github.com/enterprise/2.16/v3/
		v3, err := github.NewEnterpriseClient(o.URLv3, o.URLv3, hc)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "error while creating a GitHub v3 client")
		}
		// https://developer.github.com/enterprise/2.16/v4/guides/forming-calls/
		v4URL, err := buildV4URL(v3.BaseURL)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "error while creating a GitHub v4 client")
		}
		v4 := githubv4.NewEnterpriseClient(v4URL, hc)
		return v4, v3, nil
	}
	v4 := githubv4.NewClient(hc)
	v3 := github.NewClient(hc)
	return v4, v3, nil
}

func buildV4URL(v3 *url.URL) (string, error) {
	v4, err := v3.Parse("../graphql")
	if err != nil {
		return "", errors.Wrapf(err, "error while building v4 URL")
	}
	return v4.String(), nil
}
