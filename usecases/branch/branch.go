// Package branch provides use-cases for creating or updating a branch.
package branch

import (
	"context"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(CommitToBranch), "*"),
	wire.Bind(new(usecases.CommitToBranch), new(*CommitToBranch)),
)

// CommitToBranch commits files to the default/given branch on the repository.
type CommitToBranch struct {
	Commit     usecases.Commit
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
	GitHub     adaptors.GitHub
}

func (u *CommitToBranch) Do(ctx context.Context, in usecases.CommitToBranchIn) error {
	if !in.Repository.IsValid() {
		return xerrors.New("you must set GitHub repository")
	}
	if in.CommitMessage == "" {
		return xerrors.New("you must set commit message")
	}
	if len(in.Paths) == 0 {
		return xerrors.New("you must set one or more paths")
	}

	files, err := u.FileSystem.FindFiles(in.Paths, &pathFilter{Logger: u.Logger})
	if err != nil {
		return xerrors.Errorf("could not find files: %w", err)
	}
	if len(files) == 0 {
		return xerrors.New("no file exists in given paths")
	}

	out, err := u.GitHub.QueryForCommitToBranch(ctx, adaptors.QueryForCommitToBranchIn{
		Repository: in.Repository,
		BranchName: in.BranchName,             // optional
		ParentRef:  in.ParentOfBranch.FromRef, // optional
	})
	if err != nil {
		return xerrors.Errorf("could not find the repository: %w", err)
	}
	u.Logger.Infof("Author and committer: %s", out.CurrentUserName)

	if in.ParentOfBranch.FastForward {
		if in.BranchName == "" {
			u.Logger.Debugf("Updating the default branch by fast-forward")
			if err := u.updateBranch(ctx, updateBranchIn{
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
				return xerrors.Errorf("could not update the default branch by fast-forward: %w", err)
			}
			return nil
		}
		if out.BranchCommitSHA == "" {
			u.Logger.Debugf("Creating a branch (%s) based on the default branch", in.BranchName)
			if err := u.createBranch(ctx, createBranchIn{
				CommitIn: usecases.CommitIn{
					Files:           files,
					Repository:      out.Repository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.DefaultBranchCommitSHA,
					ParentTreeSHA:   out.DefaultBranchTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				NewBranchName: in.BranchName,
				DryRun:        in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not create a branch (%s) based on the default branch: %w", in.BranchName, err)
			}
			return nil
		}
		u.Logger.Debugf("Updating the branch (%s) by fast-forward", in.BranchName)
		if err := u.updateBranch(ctx, updateBranchIn{
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
			return xerrors.Errorf("could not update the branch (%s) by fast-forward: %w", in.BranchName, err)
		}
		return nil
	}

	if in.ParentOfBranch.NoParent {
		if in.BranchName == "" {
			//TODO: this requires force update
			u.Logger.Debugf("Updating the default branch to a commit without any parent")
			if err := u.updateBranch(ctx, updateBranchIn{
				CommitIn: usecases.CommitIn{
					Files:         files,
					Repository:    out.Repository,
					CommitMessage: in.CommitMessage,
					NoFileMode:    in.NoFileMode,
				},
				BranchName: out.DefaultBranchName,
				DryRun:     in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not update the default branch to a commit without any parent: %w", err)
			}
			return nil
		}
		if out.BranchCommitSHA == "" {
			u.Logger.Debugf("Creating a branch (%s) without any parent", in.BranchName)
			if err := u.createBranch(ctx, createBranchIn{
				CommitIn: usecases.CommitIn{
					Files:         files,
					Repository:    out.Repository,
					CommitMessage: in.CommitMessage,
					NoFileMode:    in.NoFileMode,
				},
				NewBranchName: in.BranchName,
				DryRun:        in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not create a branch (%s) without any parent: %w", in.BranchName, err)
			}
			return nil
		}
		//TODO: this may require force update
		u.Logger.Debugf("Updating the branch (%s) to a commit without any parent", in.BranchName)
		if err := u.updateBranch(ctx, updateBranchIn{
			CommitIn: usecases.CommitIn{
				Files:         files,
				Repository:    out.Repository,
				CommitMessage: in.CommitMessage,
				NoFileMode:    in.NoFileMode,
			},
			BranchName: in.BranchName,
			DryRun:     in.DryRun,
		}); err != nil {
			return xerrors.Errorf("could not update the branch (%s) to a commit without any parent: %w", in.BranchName, err)
		}
		return nil
	}

	if in.ParentOfBranch.FromRef != "" {
		if in.BranchName == "" {
			//TODO: this requires force update
			u.Logger.Debugf("Updating the default branch to a commit based on the parent ref (%s)", in.ParentOfBranch.FromRef)
			if err := u.updateBranch(ctx, updateBranchIn{
				CommitIn: usecases.CommitIn{
					Files:           files,
					Repository:      out.Repository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.ParentRefCommitSHA,
					ParentTreeSHA:   out.ParentRefTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				BranchName: out.DefaultBranchName,
				DryRun:     in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not update the default branch to a commit based on the parent ref %s: %w", in.ParentOfBranch.FromRef, err)
			}
			return nil
		}
		if out.BranchCommitSHA == "" {
			u.Logger.Debugf("Creating a branch (%s) with a commit based on the parent ref (%s)", in.BranchName, in.ParentOfBranch.FromRef)
			if err := u.createBranch(ctx, createBranchIn{
				CommitIn: usecases.CommitIn{
					Files:           files,
					Repository:      out.Repository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.ParentRefCommitSHA,
					ParentTreeSHA:   out.ParentRefTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				NewBranchName: in.BranchName,
				DryRun:        in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not create a branch (%s) with a commit based on the parent ref (%s): %w", in.BranchName, in.ParentOfBranch.FromRef, err)
			}
			return nil
		}
		//TODO: this requires force update
		u.Logger.Debugf("Updating the branch (%s) to a commit based on the parent ref (%s)", in.BranchName, in.ParentOfBranch.FromRef)
		if err := u.updateBranch(ctx, updateBranchIn{
			CommitIn: usecases.CommitIn{
				Files:           files,
				Repository:      out.Repository,
				CommitMessage:   in.CommitMessage,
				ParentCommitSHA: out.ParentRefCommitSHA,
				ParentTreeSHA:   out.ParentRefTreeSHA,
				NoFileMode:      in.NoFileMode,
			},
			BranchName: in.BranchName,
			DryRun:     in.DryRun,
		}); err != nil {
			return xerrors.Errorf("could not update the branch (%s) to a commit based on the parent ref (%s): %w", in.BranchName, in.ParentOfBranch.FromRef, err)
		}
		return nil
	}

	return xerrors.New("exact one of ParentOfBranch members must be valid")
}

type pathFilter struct {
	Logger adaptors.Logger
}

func (f *pathFilter) SkipDir(path string) bool {
	base := filepath.Base(path)
	if base == ".git" {
		f.Logger.Debugf("Exclude .git directory: %s", path)
		return true
	}
	return false
}

func (f *pathFilter) ExcludeFile(path string) bool {
	return false
}

type createBranchIn struct {
	usecases.CommitIn

	NewBranchName git.BranchName
	DryRun        bool
}

func (u *CommitToBranch) createBranch(ctx context.Context, in createBranchIn) error {
	u.Logger.Debugf("Creating a commit with the %d file(s)", len(in.Files))
	commit, err := u.Commit.Do(ctx, in.CommitIn)
	if err != nil {
		return xerrors.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Created a commit with %d changed file(s)", commit.ChangedFiles)
	if commit.ChangedFiles == 0 {
		u.Logger.Warnf("Nothing to commit because the branch has the same file(s)")
		return nil
	}
	if in.DryRun {
		u.Logger.Infof("Do not create %s branch due to dry-run", in.NewBranchName)
		return nil
	}

	u.Logger.Debugf("Creating a branch (%s)", in.NewBranchName)
	if err := u.GitHub.CreateBranch(ctx, git.NewBranch{
		Repository: in.Repository,
		BranchName: in.NewBranchName,
		CommitSHA:  commit.CommitSHA,
	}); err != nil {
		return xerrors.Errorf("error while creating %s branch: %w", in.NewBranchName, err)
	}
	u.Logger.Infof("Created a branch (%s)", in.NewBranchName)
	return nil
}

type updateBranchIn struct {
	usecases.CommitIn

	BranchName  git.BranchName
	DryRun      bool
	ForceUpdate bool
}

func (u *CommitToBranch) updateBranch(ctx context.Context, in updateBranchIn) error {
	u.Logger.Debugf("Creating a commit with the %d file(s)", len(in.Files))
	commit, err := u.Commit.Do(ctx, in.CommitIn)
	if err != nil {
		return xerrors.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Created a commit with %d changed file(s)", commit.ChangedFiles)
	if commit.ChangedFiles == 0 {
		u.Logger.Warnf("Nothing to commit because %s branch has the same file(s)", in.BranchName)
		return nil
	}
	if in.DryRun {
		u.Logger.Infof("Do not update %s branch due to dry-run", in.BranchName)
		return nil
	}

	u.Logger.Debugf("Updating the branch (%s)", in.BranchName)
	if err := u.GitHub.UpdateBranch(ctx, git.NewBranch{
		Repository: in.Repository,
		BranchName: in.BranchName,
		CommitSHA:  commit.CommitSHA,
	}, in.ForceUpdate); err != nil {
		return xerrors.Errorf("error while updating %s branch: %w", in.BranchName, err)
	}
	u.Logger.Infof("Updated the branch (%s)", in.BranchName)
	return nil
}
