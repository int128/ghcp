package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/go-github/v82/github"
	"github.com/shurcooL/githubv4"

	"github.com/int128/ghcp/pkg/git"
)

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
	slog.Debug("Querying the commit", "params", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", q)
	out := QueryCommitOutput{
		ChangedFiles: q.Repository.Object.Commit.ChangedFiles,
	}
	slog.Debug("Returning the commit", "commit", out)
	return &out, nil
}

// CreateCommit creates a commit and returns SHA of it.
func (c *GitHub) CreateCommit(ctx context.Context, n git.NewCommit) (git.CommitSHA, error) {
	slog.Debug("Creating a commit", "input", n)
	var parents []*github.Commit
	if n.ParentCommitSHA != "" {
		parents = append(parents, &github.Commit{SHA: github.Ptr(string(n.ParentCommitSHA))})
	}
	commit := github.Commit{
		Message: github.Ptr(string(n.Message)),
		Parents: parents,
		Tree:    &github.Tree{SHA: github.Ptr(string(n.TreeSHA))},
	}
	if n.Author != nil {
		commit.Author = &github.CommitAuthor{
			Name:  github.Ptr(n.Author.Name),
			Email: github.Ptr(n.Author.Email),
		}
	}
	if n.Committer != nil {
		commit.Committer = &github.CommitAuthor{
			Name:  github.Ptr(n.Committer.Name),
			Email: github.Ptr(n.Committer.Email),
		}
	}
	created, _, err := c.Client.CreateCommit(ctx, n.Repository.Owner, n.Repository.Name, &commit, nil)
	if err != nil {
		return "", fmt.Errorf("GitHub API error: %w", err)
	}
	return git.CommitSHA(created.GetSHA()), nil
}

// CreateTree creates a tree and returns SHA of it.
func (c *GitHub) CreateTree(ctx context.Context, n git.NewTree) (git.TreeSHA, error) {
	slog.Debug("Creating a tree", "input", n)
	entries := make([]*github.TreeEntry, len(n.Files))
	for i, file := range n.Files {
		entries[i] = &github.TreeEntry{
			Type: github.Ptr("blob"),
			Path: github.Ptr(file.Filename),
			Mode: github.Ptr(file.Mode()),
			SHA:  github.Ptr(string(file.BlobSHA)),
		}
	}
	tree, _, err := c.Client.CreateTree(ctx, n.Repository.Owner, n.Repository.Name, string(n.BaseTreeSHA), entries)
	if err != nil {
		return "", fmt.Errorf("GitHub API error: %w", err)
	}
	return git.TreeSHA(tree.GetSHA()), nil
}

// CreateBlob creates a blob and returns SHA of it.
func (c *GitHub) CreateBlob(ctx context.Context, n git.NewBlob) (git.BlobSHA, error) {
	slog.Debug("Creating a blob", "size", len(n.Content), "repository", n.Repository)
	blob, _, err := c.Client.CreateBlob(ctx, n.Repository.Owner, n.Repository.Name, &github.Blob{
		Encoding: github.Ptr("base64"),
		Content:  github.Ptr(n.Content),
	})
	if err != nil {
		return "", fmt.Errorf("GitHub API error: %w", err)
	}
	return git.BlobSHA(blob.GetSHA()), nil
}
