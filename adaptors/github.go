package adaptors

import (
	"context"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

func NewGitHub(v3 *github.Client, v4 *githubv4.Client) adaptors.GitHub {
	return &GitHub{v3, v4}
}

// GitHub provides GitHub API access.
type GitHub struct {
	v3 *github.Client
	v4 *githubv4.Client
}

func (c *GitHub) GetRepository(ctx context.Context, in adaptors.GetRepositoryIn) (*adaptors.GetRepositoryOut, error) {
	var q struct {
		Viewer struct {
			Login string
		}
		Repository struct {
			// default branch (name, commit SHA and tree SHA)
			DefaultBranchRef struct {
				Name string
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner": githubv4.String(in.Owner),
		"repo":  githubv4.String(in.Name),
	}
	if err := c.v4.Query(ctx, &q, v); err != nil {
		return nil, errors.Wrapf(err, "GitHub API error")
	}
	return &adaptors.GetRepositoryOut{
		CurrentUserName:   q.Viewer.Login,
		DefaultBranchName: q.Repository.DefaultBranchRef.Name,
	}, nil
}
