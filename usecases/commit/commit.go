// Package branch provides use-cases for creating or updating a branch.
package commit

import (
	"context"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors/fs"
	"github.com/int128/ghcp/adaptors/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/domain/git/commitstrategy"
	"github.com/int128/ghcp/usecases/gitobject"
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
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // if empty, target is the default branch
	ParentRepository git.RepositoryID
	CommitStrategy   commitstrategy.CommitStrategy
	CommitMessage    git.CommitMessage
	Paths            []string
	NoFileMode       bool
	DryRun           bool
}

// Commit commits files to the default/given branch on the repository.
type Commit struct {
	CreateGitObject gitobject.Interface
	FileSystem      fs.Interface
	Logger          logger.Interface
	GitHub          github.Interface
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

	if in.TargetBranchName == "" {
		return u.commitDefaultBranch(ctx, in, files)
	}
	return u.commitTargetBranch(ctx, in, files)
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

func (f *pathFilter) ExcludeFile(string) bool {
	return false
}

func (u *Commit) commitDefaultBranch(ctx context.Context, in Input, files []fs.File) error {
	q, err := u.GitHub.QueryForCommitToBranch(ctx, github.QueryForCommitToBranchIn{
		ParentRepository: in.ParentRepository,
		ParentRef:        in.CommitStrategy.RebaseUpstream(), // valid only if rebase
		TargetRepository: in.TargetRepository,
	})
	if err != nil {
		return xerrors.Errorf("could not find the repository: %w", err)
	}
	u.Logger.Infof("Author and committer: %s", q.CurrentUserName)

	bi := updateBranchInput{
		Input: gitobject.Input{
			Files:         files,
			Repository:    q.TargetRepository,
			CommitMessage: in.CommitMessage,
			NoFileMode:    in.NoFileMode,
		},
		BranchName: q.TargetDefaultBranchName,
		DryRun:     in.DryRun,
	}
	switch {
	case in.CommitStrategy.IsFastForward():
		u.Logger.Debugf("Updating the default branch by fast-forward")
		bi.ParentCommitSHA = q.ParentDefaultBranchCommitSHA
		bi.ParentTreeSHA = q.ParentDefaultBranchTreeSHA
	case in.CommitStrategy.IsRebase():
		u.Logger.Debugf("Rebasing the default branch on the ref (%s)", in.CommitStrategy.RebaseUpstream())
		bi.ParentCommitSHA = q.ParentRefCommitSHA
		bi.ParentTreeSHA = q.ParentRefTreeSHA
	case in.CommitStrategy.NoParent():
		u.Logger.Debugf("Updating the default branch to a commit with no parent")
	default:
		return xerrors.Errorf("unknown commit strategy %+v", in.CommitStrategy)
	}
	if err := u.updateBranch(ctx, bi); err != nil {
		return xerrors.Errorf("could not update the default branch: %w", err)
	}
	return nil
}

func (u *Commit) commitTargetBranch(ctx context.Context, in Input, files []fs.File) error {
	q, err := u.GitHub.QueryForCommitToBranch(ctx, github.QueryForCommitToBranchIn{
		ParentRepository: in.ParentRepository,
		ParentRef:        in.CommitStrategy.RebaseUpstream(), // valid only if rebase
		TargetRepository: in.TargetRepository,
		TargetBranchName: in.TargetBranchName,
	})
	if err != nil {
		return xerrors.Errorf("could not find the repository: %w", err)
	}
	u.Logger.Infof("Author and committer: %s", q.CurrentUserName)

	if q.TargetBranchExists() {
		bi := updateBranchInput{
			Input: gitobject.Input{
				Files:         files,
				Repository:    q.TargetRepository,
				CommitMessage: in.CommitMessage,
				NoFileMode:    in.NoFileMode,
			},
			BranchName: in.TargetBranchName,
			DryRun:     in.DryRun,
		}
		switch {
		case in.CommitStrategy.IsFastForward():
			u.Logger.Debugf("Updating the branch (%s) by fast-forward", in.TargetBranchName)
			bi.ParentCommitSHA = q.TargetBranchCommitSHA
			bi.ParentTreeSHA = q.TargetBranchTreeSHA
		case in.CommitStrategy.IsRebase():
			u.Logger.Debugf("Rebasing the branch (%s) on the ref (%s)", in.TargetBranchName, in.CommitStrategy.RebaseUpstream())
			bi.ParentCommitSHA = q.ParentRefCommitSHA
			bi.ParentTreeSHA = q.ParentRefTreeSHA
		case in.CommitStrategy.NoParent():
			u.Logger.Debugf("Updating the branch (%s) to a commit with no parent", in.TargetBranchName)
		default:
			return xerrors.Errorf("unknown commit strategy %+v", in.CommitStrategy)
		}
		if err := u.updateBranch(ctx, bi); err != nil {
			return xerrors.Errorf("could not update the branch (%s) by fast-forward: %w", in.TargetBranchName, err)
		}
		return nil
	}

	bi := createBranchInput{
		Input: gitobject.Input{
			Files:         files,
			Repository:    q.TargetRepository,
			CommitMessage: in.CommitMessage,
			NoFileMode:    in.NoFileMode,
		},
		NewBranchName: in.TargetBranchName,
		DryRun:        in.DryRun,
	}
	switch {
	case in.CommitStrategy.IsFastForward():
		u.Logger.Debugf("Creating a branch (%s) based on the default branch", in.TargetBranchName)
		bi.ParentCommitSHA = q.ParentDefaultBranchCommitSHA
		bi.ParentTreeSHA = q.ParentDefaultBranchTreeSHA
	case in.CommitStrategy.IsRebase():
		u.Logger.Debugf("Creating a branch (%s) based on the ref (%s)", in.TargetBranchName, in.CommitStrategy.RebaseUpstream())
		bi.ParentCommitSHA = q.ParentRefCommitSHA
		bi.ParentTreeSHA = q.ParentRefTreeSHA
	case in.CommitStrategy.NoParent():
		u.Logger.Debugf("Creating a branch (%s) with no parent", in.TargetBranchName)
	default:
		return xerrors.Errorf("unknown commit strategy %+v", in.CommitStrategy)
	}
	if err := u.createBranch(ctx, bi); err != nil {
		return xerrors.Errorf("could not create a branch (%s) based on the default branch: %w", in.TargetBranchName, err)
	}
	return nil
}

type createBranchInput struct {
	gitobject.Input

	NewBranchName git.BranchName
	DryRun        bool
}

func (u *Commit) createBranch(ctx context.Context, in createBranchInput) error {
	u.Logger.Debugf("Creating a commit with the %d file(s)", len(in.Files))
	commit, err := u.CreateGitObject.Do(ctx, in.Input)
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

type updateBranchInput struct {
	gitobject.Input

	BranchName  git.BranchName
	DryRun      bool
	ForceUpdate bool
}

func (u *Commit) updateBranch(ctx context.Context, in updateBranchInput) error {
	u.Logger.Debugf("Creating a commit with the %d file(s)", len(in.Files))
	commit, err := u.CreateGitObject.Do(ctx, in.Input)
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
