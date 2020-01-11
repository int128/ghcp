// Package branch provides use-cases for creating or updating a branch.
package commit

import (
	"context"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors/fs"
	"github.com/int128/ghcp/adaptors/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases/btc"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(Commit), "*"),
	wire.Bind(new(Interface), new(*Commit)),
)

//go:generate mockgen -destination mock_commit/mock_commit.go github.com/int128/ghcp/usecases/commit Interface

type Interface interface {
	Do(ctx context.Context, in Input) error
}

type Input struct {
	ParentRepository git.RepositoryID
	ParentBranch     ParentBranch
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // default branch if empty
	CommitMessage    git.CommitMessage
	Paths            []string
	NoFileMode       bool
	DryRun           bool
}

// ParentBranch represents a parent ref of the branch to create or update.
// Exact one of the members must be valid.
type ParentBranch struct {
	NoParent    bool        // push a branch without any parent
	FastForward bool        // push the branch by fast-forward
	FromRef     git.RefName // push a branch based on the ref
}

// Commit commits files to the default/given branch on the repository.
type Commit struct {
	CreateBlobTreeCommit btc.Interface
	FileSystem           fs.Interface
	Logger               logger.Interface
	GitHub               github.Interface
}

func (u *Commit) Do(ctx context.Context, in Input) error {
	if !in.TargetRepository.IsValid() {
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

	out, err := u.GitHub.QueryForCommitToBranch(ctx, github.QueryForCommitToBranchIn{
		ParentRepository: in.ParentRepository,
		ParentRef:        in.ParentBranch.FromRef, // optional
		TargetRepository: in.TargetRepository,
		TargetBranchName: in.TargetBranchName, // optional
	})
	if err != nil {
		return xerrors.Errorf("could not find the repository: %w", err)
	}
	u.Logger.Infof("Author and committer: %s", out.CurrentUserName)

	if in.ParentBranch.FastForward {
		if in.TargetBranchName == "" {
			u.Logger.Debugf("Updating the default branch by fast-forward")
			if err := u.updateBranch(ctx, updateBranchIn{
				Input: btc.Input{
					Files:           files,
					Repository:      out.TargetRepository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.ParentDefaultBranchCommitSHA,
					ParentTreeSHA:   out.ParentDefaultBranchTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				BranchName: out.TargetDefaultBranchName,
				DryRun:     in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not update the default branch by fast-forward: %w", err)
			}
			return nil
		}
		if out.TargetBranchCommitSHA == "" {
			u.Logger.Debugf("Creating a branch (%s) based on the default branch", in.TargetBranchName)
			if err := u.createBranch(ctx, createBranchIn{
				Input: btc.Input{
					Files:           files,
					Repository:      out.TargetRepository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.ParentDefaultBranchCommitSHA,
					ParentTreeSHA:   out.ParentDefaultBranchTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				NewBranchName: in.TargetBranchName,
				DryRun:        in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not create a branch (%s) based on the default branch: %w", in.TargetBranchName, err)
			}
			return nil
		}
		u.Logger.Debugf("Updating the branch (%s) by fast-forward", in.TargetBranchName)
		if err := u.updateBranch(ctx, updateBranchIn{
			Input: btc.Input{
				Files:           files,
				Repository:      out.TargetRepository,
				CommitMessage:   in.CommitMessage,
				ParentCommitSHA: out.TargetBranchCommitSHA,
				ParentTreeSHA:   out.TargetBranchTreeSHA,
				NoFileMode:      in.NoFileMode,
			},
			BranchName: in.TargetBranchName,
			DryRun:     in.DryRun,
		}); err != nil {
			return xerrors.Errorf("could not update the branch (%s) by fast-forward: %w", in.TargetBranchName, err)
		}
		return nil
	}

	if in.ParentBranch.NoParent {
		if in.TargetBranchName == "" {
			//TODO: this requires force update
			u.Logger.Debugf("Updating the default branch to a commit without any parent")
			if err := u.updateBranch(ctx, updateBranchIn{
				Input: btc.Input{
					Files:         files,
					Repository:    out.TargetRepository,
					CommitMessage: in.CommitMessage,
					NoFileMode:    in.NoFileMode,
				},
				BranchName: out.TargetDefaultBranchName,
				DryRun:     in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not update the default branch to a commit without any parent: %w", err)
			}
			return nil
		}
		if out.TargetBranchCommitSHA == "" {
			u.Logger.Debugf("Creating a branch (%s) without any parent", in.TargetBranchName)
			if err := u.createBranch(ctx, createBranchIn{
				Input: btc.Input{
					Files:         files,
					Repository:    out.TargetRepository,
					CommitMessage: in.CommitMessage,
					NoFileMode:    in.NoFileMode,
				},
				NewBranchName: in.TargetBranchName,
				DryRun:        in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not create a branch (%s) without any parent: %w", in.TargetBranchName, err)
			}
			return nil
		}
		//TODO: this may require force update
		u.Logger.Debugf("Updating the branch (%s) to a commit without any parent", in.TargetBranchName)
		if err := u.updateBranch(ctx, updateBranchIn{
			Input: btc.Input{
				Files:         files,
				Repository:    out.TargetRepository,
				CommitMessage: in.CommitMessage,
				NoFileMode:    in.NoFileMode,
			},
			BranchName: in.TargetBranchName,
			DryRun:     in.DryRun,
		}); err != nil {
			return xerrors.Errorf("could not update the branch (%s) to a commit without any parent: %w", in.TargetBranchName, err)
		}
		return nil
	}

	if in.ParentBranch.FromRef != "" {
		if in.TargetBranchName == "" {
			//TODO: this requires force update
			u.Logger.Debugf("Updating the default branch to a commit based on the parent ref (%s)", in.ParentBranch.FromRef)
			if err := u.updateBranch(ctx, updateBranchIn{
				Input: btc.Input{
					Files:           files,
					Repository:      out.TargetRepository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.ParentRefCommitSHA,
					ParentTreeSHA:   out.ParentRefTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				BranchName: out.TargetDefaultBranchName,
				DryRun:     in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not update the default branch to a commit based on the parent ref %s: %w", in.ParentBranch.FromRef, err)
			}
			return nil
		}
		if out.TargetBranchCommitSHA == "" {
			u.Logger.Debugf("Creating a branch (%s) with a commit based on the parent ref (%s)", in.TargetBranchName, in.ParentBranch.FromRef)
			if err := u.createBranch(ctx, createBranchIn{
				Input: btc.Input{
					Files:           files,
					Repository:      out.TargetRepository,
					CommitMessage:   in.CommitMessage,
					ParentCommitSHA: out.ParentRefCommitSHA,
					ParentTreeSHA:   out.ParentRefTreeSHA,
					NoFileMode:      in.NoFileMode,
				},
				NewBranchName: in.TargetBranchName,
				DryRun:        in.DryRun,
			}); err != nil {
				return xerrors.Errorf("could not create a branch (%s) with a commit based on the parent ref (%s): %w", in.TargetBranchName, in.ParentBranch.FromRef, err)
			}
			return nil
		}
		//TODO: this requires force update
		u.Logger.Debugf("Updating the branch (%s) to a commit based on the parent ref (%s)", in.TargetBranchName, in.ParentBranch.FromRef)
		if err := u.updateBranch(ctx, updateBranchIn{
			Input: btc.Input{
				Files:           files,
				Repository:      out.TargetRepository,
				CommitMessage:   in.CommitMessage,
				ParentCommitSHA: out.ParentRefCommitSHA,
				ParentTreeSHA:   out.ParentRefTreeSHA,
				NoFileMode:      in.NoFileMode,
			},
			BranchName: in.TargetBranchName,
			DryRun:     in.DryRun,
		}); err != nil {
			return xerrors.Errorf("could not update the branch (%s) to a commit based on the parent ref (%s): %w", in.TargetBranchName, in.ParentBranch.FromRef, err)
		}
		return nil
	}

	return xerrors.New("exact one of ParentBranch members must be valid")
}

type pathFilter struct {
	Logger logger.Interface
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
	btc.Input

	NewBranchName git.BranchName
	DryRun        bool
}

func (u *Commit) createBranch(ctx context.Context, in createBranchIn) error {
	u.Logger.Debugf("Creating a commit with the %d file(s)", len(in.Files))
	commit, err := u.CreateBlobTreeCommit.Do(ctx, in.Input)
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
	btc.Input

	BranchName  git.BranchName
	DryRun      bool
	ForceUpdate bool
}

func (u *Commit) updateBranch(ctx context.Context, in updateBranchIn) error {
	u.Logger.Debugf("Creating a commit with the %d file(s)", len(in.Files))
	commit, err := u.CreateBlobTreeCommit.Do(ctx, in.Input)
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
