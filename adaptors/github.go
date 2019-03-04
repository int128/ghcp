package adaptors

import (
	"context"
	"fmt"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
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
	if err := c.v4.Query(ctx, &q, v); err != nil {
		return nil, errors.Wrapf(err, "GitHub API error")
	}
	return &adaptors.GetRepositoryOut{
		CurrentUserName:        q.Viewer.Login,
		Repository:             git.RepositoryID{Owner: q.Repository.Owner.Login, Name: q.Repository.Name},
		DefaultBranchName:      git.BranchName(q.Repository.DefaultBranchRef.Name),
		DefaultBranchCommitSHA: git.CommitSHA(q.Repository.DefaultBranchRef.Target.Commit.Oid),
		DefaultBranchTreeSHA:   git.TreeSHA(q.Repository.DefaultBranchRef.Target.Commit.Tree.Oid),
	}, nil
}

// CreateBranch creates a branch and returns nil or an error.
func (c *GitHub) CreateBranch(ctx context.Context, n adaptors.NewBranch) error {
	_, _, err := c.v3.Git.CreateRef(ctx, n.Repository.Owner, n.Repository.Name, &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", n.BranchName)),
		Object: &github.GitObject{SHA: github.String(string(n.CommitSHA))},
	})
	if err != nil {
		return errors.Wrapf(err, "GitHub API error")
	}
	return nil
}

// UpdateBranch updates the branch and returns nil or an error.
func (c *GitHub) UpdateBranch(ctx context.Context, n adaptors.NewBranch, force bool) error {
	_, _, err := c.v3.Git.UpdateRef(ctx, n.Repository.Owner, n.Repository.Name, &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", n.BranchName)),
		Object: &github.GitObject{SHA: github.String(string(n.CommitSHA))},
	}, force)
	if err != nil {
		return errors.Wrapf(err, "GitHub API error")
	}
	return nil
}

// CreateCommit creates a commit and returns SHA of it.
func (c *GitHub) CreateCommit(ctx context.Context, n adaptors.NewCommit) (git.CommitSHA, error) {
	commit, _, err := c.v3.Git.CreateCommit(ctx, n.Repository.Owner, n.Repository.Name, &github.Commit{
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
func (c *GitHub) CreateTree(ctx context.Context, n adaptors.NewTree) (git.TreeSHA, error) {
	entries := make([]github.TreeEntry, len(n.Files))
	for i, file := range n.Files {
		entries[i] = github.TreeEntry{
			Type: github.String("blob"),
			Path: github.String(file.Filename),
			Mode: github.String(file.Mode()),
			SHA:  github.String(string(file.BlobSHA)),
		}
	}
	tree, _, err := c.v3.Git.CreateTree(ctx, n.Repository.Owner, n.Repository.Name, string(n.BaseTreeSHA), entries)
	if err != nil {
		return "", errors.Wrapf(err, "GitHub API error")
	}
	return git.TreeSHA(tree.GetSHA()), nil
}

// CreateBlob creates a blob and returns SHA of it.
func (c *GitHub) CreateBlob(ctx context.Context, n adaptors.NewBlob) (git.BlobSHA, error) {
	blob, _, err := c.v3.Git.CreateBlob(ctx, n.Repository.Owner, n.Repository.Name, &github.Blob{
		Encoding: github.String("base64"),
		Content:  github.String(n.Content),
	})
	if err != nil {
		return "", errors.Wrapf(err, "GitHub API error")
	}
	return git.BlobSHA(blob.GetSHA()), nil
}
