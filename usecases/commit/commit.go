// Package commit provides an internal use-case for creating a commit.
package commit

import (
	"context"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(Commit), "*"),
	wire.Bind(new(usecases.Commit), new(*Commit)),
)

// Commit creates blob(s), a tree and a commit.
type Commit struct {
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
	GitHub     adaptors.GitHub
}

func (u *Commit) Do(ctx context.Context, in usecases.CommitIn) (*usecases.CommitOut, error) {
	files := make([]git.File, len(in.Files))
	for i, file := range in.Files {
		content, err := u.FileSystem.ReadAsBase64EncodedContent(file.Path)
		if err != nil {
			return nil, xerrors.Errorf("error while reading file %s: %w", file.Path, err)
		}
		blobSHA, err := u.GitHub.CreateBlob(ctx, git.NewBlob{
			Repository: in.Repository,
			Content:    content,
		})
		if err != nil {
			return nil, xerrors.Errorf("error while creating a blob for %s: %w", file.Path, err)
		}
		gitFile := git.File{
			Filename:   file.Path,
			BlobSHA:    blobSHA,
			Executable: !in.NoFileMode && file.Executable,
		}
		files[i] = gitFile
		u.Logger.Infof("Uploaded %s as blob %s", file.Path, blobSHA)
	}

	treeSHA, err := u.GitHub.CreateTree(ctx, git.NewTree{
		Repository:  in.Repository,
		BaseTreeSHA: in.ParentTreeSHA,
		Files:       files,
	})
	if err != nil {
		return nil, xerrors.Errorf("error while creating a tree: %w", err)
	}
	u.Logger.Infof("Created tree %s", treeSHA)

	commitSHA, err := u.GitHub.CreateCommit(ctx, git.NewCommit{
		Repository:      in.Repository,
		Message:         in.CommitMessage,
		ParentCommitSHA: in.ParentCommitSHA,
		TreeSHA:         treeSHA,
	})
	if err != nil {
		return nil, xerrors.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Created commit %s", commitSHA)

	commit, err := u.GitHub.QueryCommit(ctx, adaptors.QueryCommitIn{
		Repository: in.Repository,
		CommitSHA:  commitSHA,
	})
	if err != nil {
		return nil, xerrors.Errorf("error while getting the commit %s: %w", commitSHA, err)
	}

	return &usecases.CommitOut{
		CommitSHA:    commitSHA,
		ChangedFiles: commit.ChangedFiles,
	}, nil
}
