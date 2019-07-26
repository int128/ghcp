package github

import (
	"net/http"
	"net/url"

	"github.com/google/go-github/v24/github"
	"github.com/google/wire"
	"github.com/int128/ghcp/infrastructure"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(Client)),
	wire.Bind(new(infrastructure.GitHubClient), new(*Client)),
	wire.Bind(new(infrastructure.GitHubClientInit), new(*Client)),
)

// Client provides GitHub access.
// Caller must call Init() before an API invocation.
type Client struct {
	*githubv4.Client
	*github.GitService
	*github.RepositoriesService
}

// Init initializes this client with the options.
func (c *Client) Init(o infrastructure.GitHubClientInitOptions) error {
	v4, v3, err := c.newClients(o)
	if err != nil {
		return xerrors.Errorf("error while initializing GitHub client: %w", err)
	}
	c.Client = v4
	c.GitService = v3.Git
	c.RepositoriesService = v3.Repositories
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
			return nil, nil, xerrors.Errorf("error while creating a GitHub v3 client: %w", err)
		}
		// https://developer.github.com/enterprise/2.16/v4/guides/forming-calls/
		v4URL, err := buildV4URL(v3.BaseURL)
		if err != nil {
			return nil, nil, xerrors.Errorf("error while creating a GitHub v4 client: %w", err)
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
		return "", xerrors.Errorf("error while building v4 URL: %w", err)
	}
	return v4.String(), nil
}
