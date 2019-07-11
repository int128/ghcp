package branch

import (
	"context"

	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"golang.org/x/xerrors"
)

// CreateBranch creates a new branch based on the default/given branch of the repository.
type CreateBranch struct {
	Commit     usecases.Commit
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
	GitHub     adaptors.GitHub
}

func (u *CreateBranch) Do(ctx context.Context, in usecases.CreateBranchIn) error {
	if !in.Repository.IsValid() {
		return xerrors.New("you must set GitHub repository")
	}
	if in.NewBranchName == "" {
		return xerrors.New("you must set new branch name")
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
		return xerrors.Errorf("no file exists in %v", in.Paths)
	}

	out, err := u.GitHub.QueryForCreateBranch(ctx, adaptors.QueryForCreateBranchIn{
		Repository:    in.Repository,
		ParentRef:     in.ParentOfNewBranch.FromRef,
		NewBranchName: in.NewBranchName,
	})
	if err != nil {
		return xerrors.Errorf("error while getting the repository: %w", err)
	}
	if out.NewBranchExists {
		return xerrors.Errorf("branch %s already exists", in.NewBranchName)
	}
	u.Logger.Infof("Author and committer: %s", out.CurrentUserName)

	if in.ParentOfNewBranch.FromDefaultBranch {
		if err := u.doInternal(ctx, createBranchInternalIn{
			CommitIn: usecases.CommitIn{
				Files:           files,
				Repository:      out.Repository,
				CommitMessage:   in.CommitMessage,
				ParentCommitSHA: out.DefaultBranchCommitSHA,
				ParentTreeSHA:   out.DefaultBranchTreeSHA,
				NoFileMode:      in.NoFileMode,
			},
			ParentRefName: out.DefaultBranchRefName,
			NewBranchName: in.NewBranchName,
			DryRun:        in.DryRun,
		}); err != nil {
			return xerrors.Errorf("error while creating a branch: %w", err)
		}
		return nil
	}
	if in.ParentOfNewBranch.FromRef != "" {
		if !out.ParentRefName.IsValid() {
			return xerrors.Errorf("parent ref %s does not exist", in.ParentOfNewBranch.FromRef)
		}
		if err := u.doInternal(ctx, createBranchInternalIn{
			CommitIn: usecases.CommitIn{
				Files:           files,
				Repository:      out.Repository,
				CommitMessage:   in.CommitMessage,
				ParentCommitSHA: out.ParentRefCommitSHA,
				ParentTreeSHA:   out.ParentRefTreeSHA,
				NoFileMode:      in.NoFileMode,
			},
			ParentRefName: out.ParentRefName,
			NewBranchName: in.NewBranchName,
			DryRun:        in.DryRun,
		}); err != nil {
			return xerrors.Errorf("error while creating a branch: %w", err)
		}
		return nil
	}
	if in.ParentOfNewBranch.NoParent {
		if err := u.doInternal(ctx, createBranchInternalIn{
			CommitIn: usecases.CommitIn{
				Files:         files,
				Repository:    out.Repository,
				CommitMessage: in.CommitMessage,
				NoFileMode:    in.NoFileMode,
			},
			NewBranchName: in.NewBranchName,
			DryRun:        in.DryRun,
		}); err != nil {
			return xerrors.Errorf("error while creating a branch: %w", err)
		}
		return nil
	}
	return xerrors.New("you must set one of ParentOfNewBranch")
}

type createBranchInternalIn struct {
	usecases.CommitIn

	ParentRefName git.RefQualifiedName
	NewBranchName git.BranchName
	DryRun        bool
}

func (u *CreateBranch) doInternal(ctx context.Context, in createBranchInternalIn) error {
	u.Logger.Infof("Copying %d file(s) to %s branch based on %s of %s", len(in.Files), in.NewBranchName, in.ParentRefName, in.Repository)

	commit, err := u.Commit.Do(ctx, in.CommitIn)
	if err != nil {
		return xerrors.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Commit: %d changed file(s)", commit.ChangedFiles)
	if commit.ChangedFiles == 0 {
		u.Logger.Warnf("Nothing to commit because %s has the same file(s)", in.ParentRefName)
		return nil
	}
	if in.DryRun {
		u.Logger.Infof("Do not create %s branch due to dry-run", in.NewBranchName)
		return nil
	}

	if err := u.GitHub.CreateBranch(ctx, git.NewBranch{
		Repository: in.Repository,
		BranchName: in.NewBranchName,
		CommitSHA:  commit.CommitSHA,
	}); err != nil {
		return xerrors.Errorf("error while creating %s branch: %w", in.NewBranchName, err)
	}
	u.Logger.Infof("Created %s branch on %s", in.NewBranchName, in.Repository)

	return nil
}
