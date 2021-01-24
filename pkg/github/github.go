package github

import (
	"context"

	"github.com/google/wire"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/shurcooL/githubv4"
)

var Set = wire.NewSet(
	wire.Struct(new(GitHub), "*"),
	wire.Bind(new(Interface), new(*GitHub)),
)

//go:generate mockgen -destination mock_github/mock_github.go github.com/int128/ghcp/pkg/github Interface

type Interface interface {
	CreateFork(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error)

	QueryForCommit(ctx context.Context, in QueryForCommitInput) (*QueryForCommitOutput, error)
	CreateBranch(ctx context.Context, in CreateBranchInput) error
	UpdateBranch(ctx context.Context, in UpdateBranchInput) error
	CreateCommit(ctx context.Context, commit git.NewCommit) (git.CommitSHA, error)

	QueryCommit(ctx context.Context, in QueryCommitInput) (*QueryCommitOutput, error)
	CreateTree(ctx context.Context, tree git.NewTree) (git.TreeSHA, error)
	CreateBlob(ctx context.Context, blob git.NewBlob) (git.BlobSHA, error)

	GetReleaseByTagOrNil(ctx context.Context, repo git.RepositoryID, tag git.TagName) (*git.Release, error)
	CreateRelease(ctx context.Context, r git.Release) (*git.Release, error)
	CreateReleaseAsset(ctx context.Context, a git.ReleaseAsset) error

	QueryForPullRequest(ctx context.Context, in QueryForPullRequestInput) (*QueryForPullRequestOutput, error)
	CreatePullRequest(ctx context.Context, in CreatePullRequestInput) (*CreatePullRequestOutput, error)

	QueryDefaultBranch(ctx context.Context, in QueryDefaultBranchInput) (*QueryDefaultBranchOutput, error)
}

// GitHub provides GitHub API access.
type GitHub struct {
	Client client.Interface
	Logger logger.Interface
}

type InternalRepositoryNodeID githubv4.ID
type InternalBranchNodeID githubv4.ID
