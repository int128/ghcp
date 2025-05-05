// Package gitobject provides the internal use-case for a set of blob, tree and commit.
package gitobject

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"

	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
)

var Set = wire.NewSet(
	wire.Struct(new(CreateGitObject), "*"),
	wire.Bind(new(Interface), new(*CreateGitObject)),
)

type Interface interface {
	Do(ctx context.Context, in Input) (*Output, error)
}

type Input struct {
	Files           []fs.File // nil or empty to create an empty commit
	Repository      git.RepositoryID
	CommitMessage   git.CommitMessage
	Author          *git.CommitAuthor // optional
	Committer       *git.CommitAuthor // optional
	ParentCommitSHA git.CommitSHA     // no parent if empty
	ParentTreeSHA   git.TreeSHA       // no parent if empty
	NoFileMode      bool
}

type Output struct {
	CommitSHA    git.CommitSHA
	ChangedFiles int
}

// CreateGitObject creates blob(s), a tree and a commit.
type CreateGitObject struct {
	FileSystem fs.Interface
	GitHub     github.Interface
}

func (u *CreateGitObject) Do(ctx context.Context, in Input) (*Output, error) {
	treeSHA, err := u.uploadFilesIfSet(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("error while creating a tree: %w", err)
	}

	commitSHA, err := u.GitHub.CreateCommit(ctx, git.NewCommit{
		Repository:      in.Repository,
		Message:         in.CommitMessage,
		Author:          in.Author,
		Committer:       in.Committer,
		ParentCommitSHA: in.ParentCommitSHA,
		TreeSHA:         treeSHA,
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating a commit: %w", err)
	}
	slog.Info("Created commit", "sha", commitSHA)

	commit, err := u.GitHub.QueryCommit(ctx, github.QueryCommitInput{
		Repository: in.Repository,
		CommitSHA:  commitSHA,
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting the commit %s: %w", commitSHA, err)
	}

	return &Output{
		CommitSHA:    commitSHA,
		ChangedFiles: commit.ChangedFiles,
	}, nil
}

func (u *CreateGitObject) uploadFilesIfSet(ctx context.Context, in Input) (git.TreeSHA, error) {
	if len(in.Files) == 0 {
		slog.Debug("Using the parent tree", "tree", in.ParentTreeSHA)
		return in.ParentTreeSHA, nil
	}

	files := make([]git.File, len(in.Files))
	for i, file := range in.Files {
		content, err := u.FileSystem.ReadAsBase64EncodedContent(file.Path)
		if err != nil {
			return "", fmt.Errorf("error while reading file %s: %w", file.Path, err)
		}
		blobSHA, err := u.GitHub.CreateBlob(ctx, git.NewBlob{
			Repository: in.Repository,
			Content:    content,
		})
		if err != nil {
			return "", fmt.Errorf("error while creating a blob for %s: %w", file.Path, err)
		}
		gitFile := git.File{
			Filename:   file.Path,
			BlobSHA:    blobSHA,
			Executable: !in.NoFileMode && file.Executable,
		}
		files[i] = gitFile
		slog.Info("Uploaded", "file", file.Path, "blob", blobSHA)
	}

	treeSHA, err := u.GitHub.CreateTree(ctx, git.NewTree{
		Repository:  in.Repository,
		BaseTreeSHA: in.ParentTreeSHA,
		Files:       files,
	})
	if err != nil {
		return "", fmt.Errorf("error while creating a tree: %w", err)
	}
	slog.Info("Created a tree", "tree", treeSHA)
	return treeSHA, nil
}
