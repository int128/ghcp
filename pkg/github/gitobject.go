package github

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/xerrors"

	"github.com/int128/ghcp/pkg/git"
)

// CreateCommit creates a commit and returns SHA of it.
func (c *GitHub) CreateCommit(ctx context.Context, n git.NewCommit) (git.CommitSHA, error) {
	c.Logger.Debugf("Creating a commit %+v", n)
	var parents []*github.Commit
	if n.ParentCommitSHA != "" {
		parents = append(parents, &github.Commit{SHA: github.String(string(n.ParentCommitSHA))})
	}
	commit := github.Commit{
		Message: github.String(string(n.Message)),
		Parents: parents,
		Tree:    &github.Tree{SHA: github.String(string(n.TreeSHA))},
	}
	if n.Author != nil {
		commit.Author = &github.CommitAuthor{
			Name:  github.String(n.Author.Name),
			Email: github.String(n.Author.Email),
		}
	}
	if n.Committer != nil {
		commit.Committer = &github.CommitAuthor{
			Name:  github.String(n.Committer.Name),
			Email: github.String(n.Committer.Email),
		}
	}
	created, _, err := c.Client.CreateCommit(ctx, n.Repository.Owner, n.Repository.Name, &commit)
	if err != nil {
		return "", xerrors.Errorf("GitHub API error: %w", err)
	}
	return git.CommitSHA(created.GetSHA()), nil
}

// CreateTree creates a tree and returns SHA of it.
func (c *GitHub) CreateTree(ctx context.Context, n git.NewTree) (git.TreeSHA, error) {
	c.Logger.Debugf("Creating a tree %+v", n)
	entries := make([]*github.TreeEntry, len(n.Files))
	for i, file := range n.Files {
		entries[i] = &github.TreeEntry{
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
