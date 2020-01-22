// Package gitobject provides the internal use-case for a set of blob, tree and commit.
package gitobject

import (
	"context"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors/fs"
	"github.com/int128/ghcp/adaptors/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/domain/git"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(CreateGitObject), "*"),
	wire.Bind(new(Interface), new(*CreateGitObject)),
)

//go:generate mockgen -destination mock_gitobject/mock_gitobject.go github.com/int128/ghcp/usecases/gitobject Interface

type Interface interface {
	Do(ctx context.Context, in Input) (*Output, error)
}

type Input struct {
	Files           []fs.File
	Repository      git.RepositoryID
	CommitMessage   git.CommitMessage
	ParentCommitSHA git.CommitSHA // no parent if empty
	ParentTreeSHA   git.TreeSHA   // no parent if empty
	NoFileMode      bool
}

type Output struct {
	CommitSHA    git.CommitSHA
	ChangedFiles int
}

// CreateGitObject creates blob(s), a tree and a commit.
type CreateGitObject struct {
	FileSystem fs.Interface
	Logger     logger.Interface
	GitHub     github.Interface
}

func (u *CreateGitObject) Do(ctx context.Context, in Input) (*Output, error) {
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

	commit, err := u.GitHub.QueryCommit(ctx, github.QueryCommitInput{
		Repository: in.Repository,
		CommitSHA:  commitSHA,
	})
	if err != nil {
		return nil, xerrors.Errorf("error while getting the commit %s: %w", commitSHA, err)
	}

	return &Output{
		CommitSHA:    commitSHA,
		ChangedFiles: commit.ChangedFiles,
	}, nil
}