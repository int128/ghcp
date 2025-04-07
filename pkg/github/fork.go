package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/int128/ghcp/pkg/git"
	"github.com/sethvargo/go-retry"
	"github.com/shurcooL/githubv4"
)

// CreateFork creates a fork of the repository.
// This returns ID of the fork.
func (c *GitHub) CreateFork(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error) {
	fork, _, err := c.Client.CreateFork(ctx, id.Owner, id.Name, nil)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); !ok {
			return nil, fmt.Errorf("GitHub API error: %w", err)
		}
		c.Logger.Debugf("Fork in progress: %+v", err)
	}
	forkRepository := git.RepositoryID{
		Owner: fork.GetOwner().GetLogin(),
		Name:  fork.GetName(),
	}
	if err := c.waitUntilGitDataIsAvailable(ctx, forkRepository); err != nil {
		return nil, fmt.Errorf("git data is not available on %s: %w", forkRepository, err)
	}
	return &forkRepository, nil
}

func (c *GitHub) waitUntilGitDataIsAvailable(ctx context.Context, id git.RepositoryID) error {
	operation := func(ctx context.Context) error {
		var q struct {
			Repository struct {
				DefaultBranchRef struct {
					Target struct {
						Commit struct {
							Oid string
						} `graphql:"... on Commit"`
					}
				}
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}
		v := map[string]any{
			"owner": githubv4.String(id.Owner),
			"repo":  githubv4.String(id.Name),
		}
		c.Logger.Debugf("Querying the repository with %+v", v)
		if err := c.Client.Query(ctx, &q, v); err != nil {
			return retry.RetryableError(fmt.Errorf("GitHub API error: %w", err))
		}
		c.Logger.Debugf("Got the result: %+v", q)
		return nil
	}
	if err := retry.Exponential(ctx, 500*time.Millisecond, operation); err != nil {
		return fmt.Errorf("retry over: %w", err)
	}
	return nil
}
