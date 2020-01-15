package github

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/google/go-github/v24/github"
	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/git"
	githubInfrastructure "github.com/int128/ghcp/infrastructure/github"
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
	QueryForCommitToBranch(ctx context.Context, in QueryForCommitToBranchIn) (*QueryForCommitToBranchOut, error)
	CreateBranch(ctx context.Context, branch git.NewBranch) error
	UpdateBranch(ctx context.Context, branch git.NewBranch, force bool) error
	CreateCommit(ctx context.Context, commit git.NewCommit) (git.CommitSHA, error)
	QueryCommit(ctx context.Context, in QueryCommitIn) (*QueryCommitOut, error)
	CreateTree(ctx context.Context, tree git.NewTree) (git.TreeSHA, error)
	CreateBlob(ctx context.Context, blob git.NewBlob) (git.BlobSHA, error)
}

type QueryForCommitToBranchIn struct {
	ParentRepository git.RepositoryID
	ParentRef        git.RefName // optional
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // optional
}

type QueryForCommitToBranchOut struct {
	CurrentUserName              string
	ParentDefaultBranchCommitSHA git.CommitSHA
	ParentDefaultBranchTreeSHA   git.TreeSHA
	ParentRefCommitSHA           git.CommitSHA // empty if the parent ref does not exist
	ParentRefTreeSHA             git.TreeSHA   // empty if the parent ref does not exist
	TargetRepository             git.RepositoryID
	TargetDefaultBranchName      git.BranchName
	TargetBranchCommitSHA        git.CommitSHA // empty if the branch does not exist
	TargetBranchTreeSHA          git.TreeSHA   // empty if the branch does not exist
}

type QueryCommitIn struct {
	Repository git.RepositoryID
	CommitSHA  git.CommitSHA
}

type QueryCommitOut struct {
	ChangedFiles int
}

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

// QueryForCommitToBranch returns the repository for updating the branch.
func (c *GitHub) QueryForCommitToBranch(ctx context.Context, in QueryForCommitToBranchIn) (*QueryForCommitToBranchOut, error) {
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
			Name  string
			Owner struct{ Login string }

			// default branch
			DefaultBranchRef struct {
				Name string
			}

			// branch (optional)
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
	out := QueryForCommitToBranchOut{
		CurrentUserName:              q.Viewer.Login,
		ParentDefaultBranchCommitSHA: git.CommitSHA(q.ParentRepository.DefaultBranchRef.Target.Commit.Oid),
		ParentDefaultBranchTreeSHA:   git.TreeSHA(q.ParentRepository.DefaultBranchRef.Target.Commit.Tree.Oid),
		ParentRefCommitSHA:           git.CommitSHA(q.ParentRepository.ParentRef.Target.Commit.Oid),
		ParentRefTreeSHA:             git.TreeSHA(q.ParentRepository.ParentRef.Target.Commit.Tree.Oid),
		TargetRepository:             git.RepositoryID{Owner: q.TargetRepository.Owner.Login, Name: q.TargetRepository.Name},
		TargetDefaultBranchName:      git.BranchName(q.TargetRepository.DefaultBranchRef.Name),
		TargetBranchCommitSHA:        git.CommitSHA(q.TargetRepository.Ref.Target.Commit.Oid),
		TargetBranchTreeSHA:          git.TreeSHA(q.TargetRepository.Ref.Target.Commit.Tree.Oid),
	}
	c.Logger.Debugf("Returning the repository: %+v", out)
	return &out, nil
}

// CreateBranch creates a branch and returns nil or an error.
func (c *GitHub) CreateBranch(ctx context.Context, n git.NewBranch) error {
	c.Logger.Debugf("Creating a branch %+v", n)
	_, _, err := c.Client.CreateRef(ctx, n.Repository.Owner, n.Repository.Name, &github.Reference{
		Ref:    github.String(n.BranchName.QualifiedName().String()),
		Object: &github.GitObject{SHA: github.String(string(n.CommitSHA))},
	})
	if err != nil {
		return xerrors.Errorf("GitHub API error: %w", err)
	}
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

// QueryCommit returns the commit.
func (c *GitHub) QueryCommit(ctx context.Context, in QueryCommitIn) (*QueryCommitOut, error) {
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
	out := QueryCommitOut{
		ChangedFiles: q.Repository.Object.Commit.ChangedFiles,
	}
	c.Logger.Debugf("Returning the commit: %+v", out)
	return &out, nil
}

// CreateCommit creates a commit and returns SHA of it.
func (c *GitHub) CreateCommit(ctx context.Context, n git.NewCommit) (git.CommitSHA, error) {
	c.Logger.Debugf("Creating a commit %+v", n)
	var parents []github.Commit
	if n.ParentCommitSHA != "" {
		parents = append(parents, github.Commit{SHA: github.String(string(n.ParentCommitSHA))})
	}
	commit, _, err := c.Client.CreateCommit(ctx, n.Repository.Owner, n.Repository.Name, &github.Commit{
		Message: github.String(string(n.Message)),
		Parents: parents,
		Tree:    &github.Tree{SHA: github.String(string(n.TreeSHA))},
	})
	if err != nil {
		return "", xerrors.Errorf("GitHub API error: %w", err)
	}
	return git.CommitSHA(commit.GetSHA()), nil
}

// CreateTree creates a tree and returns SHA of it.
func (c *GitHub) CreateTree(ctx context.Context, n git.NewTree) (git.TreeSHA, error) {
	c.Logger.Debugf("Creating a tree %+v", n)
	entries := make([]github.TreeEntry, len(n.Files))
	for i, file := range n.Files {
		entries[i] = github.TreeEntry{
			Type: github.String("blob"),
			Path: github.String(file.Filename),
			Mode: github.String(file.Mode()),
			SHA:  github.String(string(file.BlobSHA)),
		}
	}
	tree, _, err := c.Client.CreateTree(ctx, n.Repository.Owner, n.Repository.Name, string(n.BaseTreeSHA), entries)
	if err != nil {
		return "", xerrors.Errorf("GitHub API error: %w", err)
	}
	return git.TreeSHA(tree.GetSHA()), nil
}

// CreateBlob creates a blob and returns SHA of it.
func (c *GitHub) CreateBlob(ctx context.Context, n git.NewBlob) (git.BlobSHA, error) {
	c.Logger.Debugf("Creating a blob of %d byte(s) on the repository %+v", len(n.Content), n.Repository)
	blob, _, err := c.Client.CreateBlob(ctx, n.Repository.Owner, n.Repository.Name, &github.Blob{
		Encoding: github.String("base64"),
		Content:  github.String(n.Content),
	})
	if err != nil {
		return "", xerrors.Errorf("GitHub API error: %w", err)
	}
	return git.BlobSHA(blob.GetSHA()), nil
}
