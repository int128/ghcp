// Package usecases provides use cases of this application.
package usecases

import (
	"context"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewCopyUseCase(i CopyUseCase) usecases.CopyUseCase {
	return &i
}

// CopyUseCase performs copying files to the repository.
type CopyUseCase struct {
	dig.In
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
	GitHub     adaptors.GitHub
}

func (u *CopyUseCase) Do(ctx context.Context, in usecases.CopyUseCaseIn) error {
	files, err := u.FileSystem.FindFiles(in.Paths)
	if err != nil {
		return errors.Wrapf(err, "error while finding files")
	}

	out, err := u.GitHub.QueryRepository(ctx, adaptors.QueryRepositoryIn{
		Repository: git.RepositoryID{Owner: in.Repository.Owner, Name: in.Repository.Name},
		BranchName: in.BranchName,
	})
	if err != nil {
		return errors.Wrapf(err, "error while getting the repository")
	}
	u.Logger.Infof("Logged in as %s", out.CurrentUserName)

	gitFiles := make([]git.File, len(files))
	for i, file := range files {
		content, err := u.FileSystem.ReadAsBase64EncodedContent(file.Path)
		if err != nil {
			return errors.Wrapf(err, "error while reading file %s", file.Path)
		}
		blobSHA, err := u.GitHub.CreateBlob(ctx, git.NewBlob{
			Repository: in.Repository,
			Content:    content,
		})
		if err != nil {
			return errors.Wrapf(err, "error while creating a blob for %s", file.Path)
		}
		gitFile := git.File{
			Filename:   file.Path,
			BlobSHA:    blobSHA,
			Executable: !in.NoFileMode && file.Executable,
		}
		gitFiles[i] = gitFile
		u.Logger.Infof("Uploaded %s as blob %s", file.Path, blobSHA)
	}

	if in.BranchName == "" {
		// copy to the default branch
		if err := u.copyToExistingBranch(ctx, copyToExistingBranchIn{
			CopyUseCaseIn:   in,
			Files:           gitFiles,
			Repository:      out.Repository,
			BranchName:      out.DefaultBranchName,
			ParentCommitSHA: out.DefaultBranchCommitSHA,
			ParentTreeSHA:   out.DefaultBranchTreeSHA,
		}); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	if out.BranchCommitSHA == "" || out.BranchTreeSHA == "" {
		return errors.Errorf("branch %s does not exist", in.BranchName)
	}
	if err := u.copyToExistingBranch(ctx, copyToExistingBranchIn{
		CopyUseCaseIn:   in,
		Files:           gitFiles,
		Repository:      out.Repository,
		BranchName:      in.BranchName,
		ParentCommitSHA: out.BranchCommitSHA,
		ParentTreeSHA:   out.BranchTreeSHA,
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type copyToExistingBranchIn struct {
	usecases.CopyUseCaseIn

	Files           []git.File
	Repository      git.RepositoryID
	BranchName      git.BranchName
	ParentCommitSHA git.CommitSHA
	ParentTreeSHA   git.TreeSHA
}

func (u *CopyUseCase) copyToExistingBranch(ctx context.Context, in copyToExistingBranchIn) error {
	treeSHA, err := u.GitHub.CreateTree(ctx, git.NewTree{
		Repository:  in.Repository,
		BaseTreeSHA: in.ParentTreeSHA,
		Files:       in.Files,
	})
	if err != nil {
		return errors.Wrapf(err, "error while creating a tree")
	}
	u.Logger.Infof("Created tree %s", treeSHA)

	commitSHA, err := u.GitHub.CreateCommit(ctx, git.NewCommit{
		Repository:      in.Repository,
		Message:         in.CommitMessage,
		ParentCommitSHA: in.ParentCommitSHA,
		TreeSHA:         treeSHA,
	})
	if err != nil {
		return errors.Wrapf(err, "error while creating a commit")
	}
	u.Logger.Infof("Created commit %s", commitSHA)

	commit, err := u.GitHub.QueryCommit(ctx, adaptors.QueryCommitIn{
		Repository: in.Repository,
		CommitSHA:  commitSHA,
	})
	if err != nil {
		return errors.Wrapf(err, "error while getting the commit %s", commitSHA)
	}
	u.Logger.Infof("Commit: %d changed file(s)", commit.ChangedFiles)
	if commit.ChangedFiles == 0 {
		u.Logger.Infof("Nothing to commit")
		return nil
	}

	if in.DryRun {
		u.Logger.Infof("Do not update %s branch due to dry-run", in.BranchName)
		return nil
	}

	if err := u.GitHub.UpdateBranch(ctx, git.NewBranch{
		Repository: in.Repository,
		BranchName: in.BranchName,
		CommitSHA:  commitSHA,
	}, false); err != nil {
		return errors.Wrapf(err, "error while updating %s branch", in.BranchName)
	}
	u.Logger.Infof("Updated %s branch", in.BranchName)

	return nil
}
