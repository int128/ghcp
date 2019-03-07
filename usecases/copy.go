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
	filenames, err := u.FileSystem.FindFiles(in.Paths)
	if err != nil {
		return errors.Wrapf(err, "error while finding files")
	}

	out, err := u.GitHub.QueryRepository(ctx, adaptors.QueryRepositoryIn{
		Repository: git.RepositoryID{Owner: in.Repository.Owner, Name: in.Repository.Name},
	})
	if err != nil {
		return errors.Wrapf(err, "error while getting the repository")
	}
	u.Logger.Infof("Logged in as %s", out.CurrentUserName)

	gitFiles := make([]git.File, len(filenames))
	for i, filename := range filenames {
		content, err := u.FileSystem.ReadAsBase64EncodedContent(filename)
		if err != nil {
			return errors.Wrapf(err, "error while reading file %s", filename)
		}
		blobSHA, err := u.GitHub.CreateBlob(ctx, git.NewBlob{
			Repository: out.Repository,
			Content:    content,
		})
		if err != nil {
			return errors.Wrapf(err, "error while creating a blob for %s", filename)
		}
		gitFiles[i] = git.File{
			Filename: filename,
			BlobSHA:  blobSHA,
			//TODO: Executable
		}
		u.Logger.Infof("Uploaded %s as blob %s", filename, blobSHA)
	}

	treeSHA, err := u.GitHub.CreateTree(ctx, git.NewTree{
		Repository:  out.Repository,
		BaseTreeSHA: out.DefaultBranchTreeSHA,
		Files:       gitFiles,
	})
	if err != nil {
		return errors.Wrapf(err, "error while creating a tree")
	}
	u.Logger.Infof("Created tree %s", treeSHA)

	commitSHA, err := u.GitHub.CreateCommit(ctx, git.NewCommit{
		Repository:      out.Repository,
		Message:         in.CommitMessage,
		ParentCommitSHA: out.DefaultBranchCommitSHA,
		TreeSHA:         treeSHA,
	})
	if err != nil {
		return errors.Wrapf(err, "error while creating a commit")
	}
	u.Logger.Infof("Created commit %s", commitSHA)

	commit, err := u.GitHub.QueryCommit(ctx, adaptors.QueryCommitIn{
		Repository: out.Repository,
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
		u.Logger.Infof("Do not update %s branch due to dry-run", out.DefaultBranchName)
		return nil
	}

	if err := u.GitHub.UpdateBranch(ctx, git.NewBranch{
		Repository: out.Repository,
		BranchName: out.DefaultBranchName,
		CommitSHA:  commitSHA,
	}, false); err != nil {
		return errors.Wrapf(err, "error while updating %s branch", out.DefaultBranchName)
	}
	u.Logger.Infof("Updated %s branch", out.DefaultBranchName)

	return nil
}
