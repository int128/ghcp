package infrastructure

import (
	"context"

	"github.com/google/go-github/v24/github"
)

//go:generate mockgen -package mock_infrastructure -destination ../mock_infrastructure/mock_infrastructure.go github.com/int128/ghcp/infrastructure/interfaces Cmd,GitHubClient,GitHubClientConfig

type Cmd interface {
	Run(ctx context.Context, args []string) int
}

type GitHubClient interface {
	// v4 API
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error

	// v3 API
	CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
	UpdateRef(ctx context.Context, owner string, repo string, ref *github.Reference, force bool) (*github.Reference, *github.Response, error)
	CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []github.TreeEntry) (*github.Tree, *github.Response, error)
	CreateBlob(ctx context.Context, owner string, repo string, blob *github.Blob) (*github.Blob, *github.Response, error)
}

type GitHubClientConfig interface {
	SetToken(token string)
}
