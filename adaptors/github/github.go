package github

import (
	"context"

	"github.com/google/wire"
	"github.com/int128/ghcp/domain/git"
	"github.com/shurcooL/githubv4"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(GitHub), "*"),
	wire.Bind(new(Interface), new(*GitHub)),
)

//go:generate mockgen -destination mock_github/mock_github.go github.com/int128/ghcp/adaptors/github Interface

type Interface interface {
	CreateFork(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error)

	QueryForCommit(ctx context.Context, in QueryForCommitInput) (*QueryForCommitOutput, error)
	CreateBranch(ctx context.Context, branch git.NewBranch) error
	UpdateBranch(ctx context.Context, branch git.NewBranch, force bool) error
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

type QueryDefaultBranchInput struct {
	BaseRepository git.RepositoryID
	HeadRepository git.RepositoryID
}

type QueryDefaultBranchOutput struct {
	BaseDefaultBranchName git.BranchName
	HeadDefaultBranchName git.BranchName
}

// QueryDefaultBranch returns the default branch names.
// You can set both repositories or either repository.
func (c *GitHub) QueryDefaultBranch(ctx context.Context, in QueryDefaultBranchInput) (*QueryDefaultBranchOutput, error) {
	if !in.BaseRepository.IsValid() && !in.HeadRepository.IsValid() {
		return nil, xerrors.New("BaseRepository and HeadRepository are zero")
	}
	var q struct {
		BaseRepository struct {
			DefaultBranchRef struct {
				Name string
			}
		} `graphql:"baseRepository: repository(owner: $baseOwner, name: $baseRepo)"`
		HeadRepository struct {
			DefaultBranchRef struct {
				Name string
			}
		} `graphql:"headRepository: repository(owner: $headOwner, name: $headRepo)"`
	}
	v := map[string]interface{}{
		"baseOwner": githubv4.String(in.BaseRepository.Owner),
		"baseRepo":  githubv4.String(in.BaseRepository.Name),
		"headOwner": githubv4.String(in.HeadRepository.Owner),
		"headRepo":  githubv4.String(in.HeadRepository.Name),
	}
	c.Logger.Debugf("Querying the default branch name with %+v", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", q)
	return &QueryDefaultBranchOutput{
		BaseDefaultBranchName: git.BranchName(q.BaseRepository.DefaultBranchRef.Name),
		HeadDefaultBranchName: git.BranchName(q.HeadRepository.DefaultBranchRef.Name),
	}, nil
}
