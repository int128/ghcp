package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/google/wire"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var Set = wire.NewSet(
	wire.Value(NewFunc(New)),
)

type NewFunc func(Option) (Interface, error)

//go:generate mockgen -destination mock_client/mock_client.go github.com/int128/ghcp/pkg/github/client Interface

type Interface interface {
	QueryService
	GitService
	RepositoriesService
}

type QueryService interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	Mutate(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error
}

type GitService interface {
	CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error)
	CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error)
}

type RepositoriesService interface {
	CreateFork(ctx context.Context, owner, repo string, opt *github.RepositoryCreateForkOptions) (*github.Repository, *github.Response, error)
	GetReleaseByTag(ctx context.Context, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error)
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
	UploadReleaseAsset(ctx context.Context, owner, repo string, id int64, opt *github.UploadOptions, file *os.File) (*github.ReleaseAsset, *github.Response, error)
}

type Option struct {
	// A token for GitHub API.
	Token string

	// GitHub API v3 URL (for GitHub Enterprise).
	// e.g. https://github.example.com/api/v3/
	URLv3 string
}

type clientSet struct {
	QueryService
	GitService
	RepositoriesService
}

func New(o Option) (Interface, error) {
	v4, v3, err := newClients(o)
	if err != nil {
		return nil, fmt.Errorf("error while initializing GitHub client: %w", err)
	}
	return &clientSet{
		QueryService:        v4,
		GitService:          v3.Git,
		RepositoriesService: v3.Repositories,
	}, nil
}

func newClients(o Option) (*githubv4.Client, *github.Client, error) {
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
			return nil, nil, fmt.Errorf("error while creating a GitHub v3 client: %w", err)
		}
		// https://developer.github.com/enterprise/2.16/v4/guides/forming-calls/
		v4URL, err := buildV4URL(v3.BaseURL)
		if err != nil {
			return nil, nil, fmt.Errorf("error while creating a GitHub v4 client: %w", err)
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
		return "", fmt.Errorf("error while building v4 URL: %w", err)
	}
	return v4.String(), nil
}
