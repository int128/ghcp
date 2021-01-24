package github

import (
	"context"

	"github.com/int128/ghcp/pkg/git"
	"github.com/shurcooL/githubv4"
	"golang.org/x/xerrors"
)

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
	if !in.BaseRepository.IsValid() || !in.HeadRepository.IsValid() {
		return nil, xerrors.New("you need to set both BaseRepository and HeadRepository")
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
