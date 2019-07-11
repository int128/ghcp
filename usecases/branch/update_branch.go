package branch

import (
	"context"

	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"golang.org/x/xerrors"
)

// UpdateBranch copies files to the default/given branch on the repository.
type UpdateBranch struct {
	Commit     usecases.Commit
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
	GitHub     adaptors.GitHub
}

func (u *UpdateBranch) Do(ctx context.Context, in usecases.UpdateBranchIn) error {
	if !in.Repository.IsValid() {
		return xerrors.New("you must set GitHub repository")
	}
	if in.CommitMessage == "" {
		return xerrors.New("you must set commit message")
	}
	if len(in.Paths) == 0 {
		return xerrors.New("you must set one or more paths")
	}

	files, err := u.FileSystem.FindFiles(in.Paths)
	if err != nil {
		return xerrors.Errorf("error while finding files: %w", err)
	}
	if len(files) == 0 {
		return xerrors.New("no file exists in given paths")
	}

	out, err := u.GitHub.QueryForUpdateBranch(ctx, adaptors.QueryForUpdateBranchIn{
		Repository: in.Repository,
		BranchName: in.BranchName,
	})
	if err != nil {
		return xerrors.Errorf("error while getting the repository: %w", err)
	}
	u.Logger.Infof("Author and committer: %s", out.CurrentUserName)

	if in.BranchName == "" {
		// copy to the default branch
		if err := u.doInternal(ctx, updateBranchInternalIn{
			CommitIn: usecases.CommitIn{
				Files:           files,
				Repository:      out.Repository,
				CommitMessage:   in.CommitMessage,
				ParentCommitSHA: out.DefaultBranchCommitSHA,
				ParentTreeSHA:   out.DefaultBranchTreeSHA,
				NoFileMode:      in.NoFileMode,
			},
			BranchName: out.DefaultBranchName,
			DryRun:     in.DryRun,
		}); err != nil {
			return xerrors.Errorf("error while updating the default branch: %w", err)
		}
		return nil
	}

	// copy to the given branch
	if out.BranchCommitSHA == "" || out.BranchTreeSHA == "" {
		return xerrors.Errorf("branch %s does not exist", in.BranchName)
	}
	if err := u.doInternal(ctx, updateBranchInternalIn{
		CommitIn: usecases.CommitIn{
			Files:           files,
			Repository:      out.Repository,
			CommitMessage:   in.CommitMessage,
			ParentCommitSHA: out.BranchCommitSHA,
			ParentTreeSHA:   out.BranchTreeSHA,
			NoFileMode:      in.NoFileMode,
		},
		BranchName: in.BranchName,
		DryRun:     in.DryRun,
	}); err != nil {
		return xerrors.Errorf("error while updating %s branch: %w", in.BranchName, err)
	}
	return nil
}

type updateBranchInternalIn struct {
	usecases.CommitIn

	BranchName git.BranchName
	DryRun     bool
}

func (u *UpdateBranch) doInternal(ctx context.Context, in updateBranchInternalIn) error {
	u.Logger.Infof("Copying %d file(s) to %s branch of %s", len(in.Files), in.BranchName, in.Repository)

	commit, err := u.Commit.Do(ctx, in.CommitIn)
	if err != nil {
		return xerrors.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Commit: %d changed file(s)", commit.ChangedFiles)
	if commit.ChangedFiles == 0 {
		u.Logger.Warnf("Nothing to commit because %s branch has the same file(s)", in.BranchName)
		return nil
	}
	if in.DryRun {
		u.Logger.Infof("Do not update %s branch due to dry-run", in.BranchName)
		return nil
	}

	if err := u.GitHub.UpdateBranch(ctx, git.NewBranch{
		Repository: in.Repository,
		BranchName: in.BranchName,
		CommitSHA:  commit.CommitSHA,
	}, false); err != nil {
		return xerrors.Errorf("error while updating %s branch: %w", in.BranchName, err)
	}
	u.Logger.Infof("Updated %s branch of %s", in.BranchName, in.Repository)

	return nil
}
