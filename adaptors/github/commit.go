package github

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/domain/git"
	githubInfrastructure "github.com/int128/ghcp/infrastructure/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/xerrors"
)

// GitHub provides GitHub API access.
type GitHub struct {
	Client githubInfrastructure.Interface
	Logger logger.Interface
}

// CreateFork creates a fork of the repository.
// This returns ID of the fork.
func (c *GitHub) CreateFork(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error) {
	fork, _, err := c.Client.CreateFork(ctx, id.Owner, id.Name, nil)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); !ok {
			return nil, xerrors.Errorf("GitHub API error: %w", err)
		}
		c.Logger.Debugf("Fork in progress: %+v", err)
	}
	forkRepository := git.RepositoryID{
		Owner: fork.GetOwner().GetLogin(),
		Name:  fork.GetName(),
	}
	if err := c.waitUntilGitDataIsAvailable(ctx, forkRepository); err != nil {
		return nil, xerrors.Errorf("git data is not available on %s: %w", forkRepository, err)
	}
	return &forkRepository, nil
}

func (c *GitHub) waitUntilGitDataIsAvailable(ctx context.Context, id git.RepositoryID) error {
	operation := func() error {
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
		v := map[string]interface{}{
			"owner": githubv4.String(id.Owner),
			"repo":  githubv4.String(id.Name),
		}
		c.Logger.Debugf("Querying the repository with %+v", v)
		if err := c.Client.Query(ctx, &q, v); err != nil {
			return xerrors.Errorf("GitHub API error: %w", err)
		}
		c.Logger.Debugf("Got the result: %+v", q)
		return nil
	}
	if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
		return xerrors.Errorf("retry over: %w", err)
	}
	return nil
}

type QueryForCommitInput struct {
	ParentRepository git.RepositoryID
	ParentRef        git.RefName // optional
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // optional
}

type QueryForCommitOutput struct {
	CurrentUserName              string
	ParentDefaultBranchCommitSHA git.CommitSHA
	ParentDefaultBranchTreeSHA   git.TreeSHA
	ParentRefCommitSHA           git.CommitSHA // empty if the parent ref does not exist
	ParentRefTreeSHA             git.TreeSHA   // empty if the parent ref does not exist
	TargetRepositoryNodeID       InternalRepositoryNodeID
	TargetBranchCommitSHA        git.CommitSHA // empty if the branch does not exist
	TargetBranchTreeSHA          git.TreeSHA   // empty if the branch does not exist
}

func (q *QueryForCommitOutput) TargetBranchExists() bool {
	return q.TargetBranchCommitSHA != ""
}

// QueryForCommit returns the repository for updating the branch.
func (c *GitHub) QueryForCommit(ctx context.Context, in QueryForCommitInput) (*QueryForCommitOutput, error) {
	var q struct {
		Viewer struct {
			Login string
		}

		ParentRepository struct {
			// default branch
			DefaultBranchRef struct {
				Name   string
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			}

			// parent ref (optional)
			ParentRef struct {
				Prefix string
				Name   string
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			} `graphql:"parentRef: ref(qualifiedName: $parentRef)"`
		} `graphql:"parentRepository: repository(owner: $parentOwner, name: $parentRepo)"`

		TargetRepository struct {
			ID  githubv4.ID
			Ref struct {
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			} `graphql:"ref(qualifiedName: $targetRef)"`
		} `graphql:"targetRepository: repository(owner: $targetOwner, name: $targetRepo)"`
	}
	v := map[string]interface{}{
		"parentOwner": githubv4.String(in.ParentRepository.Owner),
		"parentRepo":  githubv4.String(in.ParentRepository.Name),
		"parentRef":   githubv4.String(in.ParentRef),
		"targetOwner": githubv4.String(in.TargetRepository.Owner),
		"targetRepo":  githubv4.String(in.TargetRepository.Name),
		"targetRef":   githubv4.String(in.TargetBranchName.QualifiedName().String()),
	}
	c.Logger.Debugf("Querying the repository with %+v", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", q)
	out := QueryForCommitOutput{
		CurrentUserName:              q.Viewer.Login,
		ParentDefaultBranchCommitSHA: git.CommitSHA(q.ParentRepository.DefaultBranchRef.Target.Commit.Oid),
		ParentDefaultBranchTreeSHA:   git.TreeSHA(q.ParentRepository.DefaultBranchRef.Target.Commit.Tree.Oid),
		ParentRefCommitSHA:           git.CommitSHA(q.ParentRepository.ParentRef.Target.Commit.Oid),
		ParentRefTreeSHA:             git.TreeSHA(q.ParentRepository.ParentRef.Target.Commit.Tree.Oid),
		TargetRepositoryNodeID:       q.TargetRepository.ID,
		TargetBranchCommitSHA:        git.CommitSHA(q.TargetRepository.Ref.Target.Commit.Oid),
		TargetBranchTreeSHA:          git.TreeSHA(q.TargetRepository.Ref.Target.Commit.Tree.Oid),
	}
	c.Logger.Debugf("Returning the repository: %+v", out)
	return &out, nil
}

type CreateBranchInput struct {
	RepositoryNodeID InternalRepositoryNodeID
	BranchName       git.BranchName
	CommitSHA        git.CommitSHA
}

// CreateBranch creates a branch and returns nil or an error.
func (c *GitHub) CreateBranch(ctx context.Context, in CreateBranchInput) error {
	c.Logger.Debugf("Creating a branch %+v", in.BranchName)
	v := githubv4.CreateRefInput{
		RepositoryID: in.RepositoryNodeID,
		Name:         githubv4.String(in.BranchName.QualifiedName().String()),
		Oid:          githubv4.GitObjectID(in.CommitSHA),
	}
	var m struct {
		CreateRef struct {
			Ref struct {
				Name string
			}
		} `graphql:"createRef(input: $input)"`
	}
	if err := c.Client.Mutate(ctx, &m, v, nil); err != nil {
		return xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", m)
	return nil
}

// UpdateBranch updates the branch and returns nil or an error.
func (c *GitHub) UpdateBranch(ctx context.Context, n git.NewBranch, force bool) error {
	c.Logger.Debugf("Updating the branch %+v, force: %v", n, force)
	_, _, err := c.Client.UpdateRef(ctx, n.Repository.Owner, n.Repository.Name, &github.Reference{
		Ref:    github.String(n.BranchName.QualifiedName().String()),
		Object: &github.GitObject{SHA: github.String(string(n.CommitSHA))},
	}, force)
	if err != nil {
		return xerrors.Errorf("GitHub API error: %w", err)
	}
	return nil
}

type QueryCommitInput struct {
	Repository git.RepositoryID
	CommitSHA  git.CommitSHA
}

type QueryCommitOutput struct {
	ChangedFiles int
}

// QueryCommit returns the commit.
func (c *GitHub) QueryCommit(ctx context.Context, in QueryCommitInput) (*QueryCommitOutput, error) {
	var q struct {
		Repository struct {
			Object struct {
				Commit struct {
					ChangedFiles int
				} `graphql:"... on Commit"`
			} `graphql:"object(oid: $commitSHA)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner":     githubv4.String(in.Repository.Owner),
		"repo":      githubv4.String(in.Repository.Name),
		"commitSHA": githubv4.GitObjectID(in.CommitSHA),
	}
	c.Logger.Debugf("Querying the commit with %+v", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", q)
	out := QueryCommitOutput{
		ChangedFiles: q.Repository.Object.Commit.ChangedFiles,
	}
	c.Logger.Debugf("Returning the commit: %+v", out)
	return &out, nil
}
