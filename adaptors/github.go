package adaptors

import (
	"context"
	"fmt"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"go.uber.org/dig"
)

func NewGitHub(i GitHub) adaptors.GitHub {
	return &i
}

// GitHub provides GitHub API access.
type GitHub struct {
	dig.In
	Client infrastructure.GitHubClient
}

// QueryRepository returns the repository.
func (c *GitHub) QueryRepository(ctx context.Context, in adaptors.QueryRepositoryIn) (*adaptors.QueryRepositoryOut, error) {
	var q struct {
		Viewer struct {
			Login string
		}
		Repository struct {
			Name  string
			Owner struct{ Login string }

			// default branch (name, commit SHA and tree SHA)
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
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner": githubv4.String(in.Repository.Owner),
		"repo":  githubv4.String(in.Repository.Name),
	}
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, errors.Wrapf(err, "GitHub API error")
	}
	return &adaptors.QueryRepositoryOut{
		CurrentUserName:        q.Viewer.Login,
		Repository:             git.RepositoryID{Owner: q.Repository.Owner.Login, Name: q.Repository.Name},
		DefaultBranchName:      git.BranchName(q.Repository.DefaultBranchRef.Name),
		DefaultBranchCommitSHA: git.CommitSHA(q.Repository.DefaultBranchRef.Target.Commit.Oid),
		DefaultBranchTreeSHA:   git.TreeSHA(q.Repository.DefaultBranchRef.Target.Commit.Tree.Oid),
	}, nil
}

// CreateBranch creates a branch and returns nil or an error.
func (c *GitHub) CreateBranch(ctx context.Context, n git.NewBranch) error {
	_, _, err := c.Client.CreateRef(ctx, n.Repository.Owner, n.Repository.Name, &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", n.BranchName)),
		Object: &github.GitObject{SHA: github.String(string(n.CommitSHA))},
	})
	if err != nil {
		return errors.Wrapf(err, "GitHub API error")
	}
	return nil
}

// UpdateBranch updates the branch and returns nil or an error.
func (c *GitHub) UpdateBranch(ctx context.Context, n git.NewBranch, force bool) error {
	_, _, err := c.Client.UpdateRef(ctx, n.Repository.Owner, n.Repository.Name, &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", n.BranchName)),
		Object: &github.GitObject{SHA: github.String(string(n.CommitSHA))},
	}, force)
	if err != nil {
		return errors.Wrapf(err, "GitHub API error")
	}
	return nil
}

// CreateCommit creates a commit and returns SHA of it.
func (c *GitHub) CreateCommit(ctx context.Context, n git.NewCommit) (git.CommitSHA, error) {
	commit, _, err := c.Client.CreateCommit(ctx, n.Repository.Owner, n.Repository.Name, &github.Commit{
		Message: github.String(string(n.Message)),
		Parents: []github.Commit{{SHA: github.String(string(n.ParentCommitSHA))}},
		Tree:    &github.Tree{SHA: github.String(string(n.TreeSHA))},
	})
	if err != nil {
		return "", errors.Wrapf(err, "GitHub API error")
	}
	return git.CommitSHA(commit.GetSHA()), nil
}

// CreateTree creates a tree and returns SHA of it.
func (c *GitHub) CreateTree(ctx context.Context, n git.NewTree) (git.TreeSHA, error) {
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
		return "", errors.Wrapf(err, "GitHub API error")
	}
	return git.TreeSHA(tree.GetSHA()), nil
}

// CreateBlob creates a blob and returns SHA of it.
func (c *GitHub) CreateBlob(ctx context.Context, n git.NewBlob) (git.BlobSHA, error) {
	blob, _, err := c.Client.CreateBlob(ctx, n.Repository.Owner, n.Repository.Name, &github.Blob{
		Encoding: github.String("base64"),
		Content:  github.String(n.Content),
	})
	if err != nil {
		return "", errors.Wrapf(err, "GitHub API error")
	}
	return git.BlobSHA(blob.GetSHA()), nil
}
