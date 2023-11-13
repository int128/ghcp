// Package branch provides use-cases for creating or updating a branch.
package commit

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/wire"

	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/int128/ghcp/pkg/usecases/gitobject"
)

var Set = wire.NewSet(
	wire.Struct(new(Commit), "*"),
	wire.Bind(new(Interface), new(*Commit)),
)

//go:generate mockgen -destination mock_commit/mock_commit.go github.com/int128/ghcp/pkg/usecases/commit Interface

type Interface interface {
	Do(ctx context.Context, in Input) error
}

type Input struct {
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // if empty, target is the default branch
	ParentRepository git.RepositoryID
	CommitStrategy   commitstrategy.CommitStrategy
	CommitMessage    git.CommitMessage
	Author           *git.CommitAuthor // optional
	Committer        *git.CommitAuthor // optional
	Paths            []string          // if Paths and DeletedPaths are empty or nil, create an empty commit
	DeletedPaths     []string          // if Paths and DeletedPaths are empty or nil, create an empty commit
	NoFileMode       bool
	DryRun           bool

	ForceUpdate bool //TODO: support force-update as well
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
		return errors.New("you must set GitHub repository")
	}
	if in.CommitMessage == "" {
		return errors.New("you must set commit message")
	}

	files, err := u.FileSystem.FindFiles(in.Paths, &pathFilter{Logger: u.Logger})
	if err != nil {
		return fmt.Errorf("could not find files: %w", err)
	}
	if len(in.Paths) > 0 && len(files) == 0 {
		return errors.New("no file exists in given paths")
	}

	if in.TargetBranchName == "" {
		q, err := u.GitHub.QueryDefaultBranch(ctx, github.QueryDefaultBranchInput{
			HeadRepository: in.TargetRepository,
			BaseRepository: in.ParentRepository, // mandatory but not used
		})
		if err != nil {
			return fmt.Errorf("could not determine the default branch: %w", err)
		}
		in.TargetBranchName = q.HeadDefaultBranchName
	}

	q, err := u.GitHub.QueryForCommit(ctx, github.QueryForCommitInput{
		ParentRepository: in.ParentRepository,
		ParentRef:        in.CommitStrategy.RebaseUpstream(), // valid only if rebase
		TargetRepository: in.TargetRepository,
		TargetBranchName: in.TargetBranchName,
	})
	if err != nil {
		return fmt.Errorf("could not find the repository: %w", err)
	}
	u.Logger.Infof("Author and committer: %s", q.CurrentUserName)
	if q.TargetBranchExists() {
		if err := u.updateExistingBranch(ctx, in, files, q); err != nil {
			return fmt.Errorf("could not update the existing branch (%s): %w", in.TargetBranchName, err)
		}
		return nil
	}
	if err := u.createNewBranch(ctx, in, files, q); err != nil {
		return fmt.Errorf("could not create a branch (%s) based on the default branch: %w", in.TargetBranchName, err)
	}
	return nil
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

func (u *Commit) createNewBranch(ctx context.Context, in Input, files []fs.File, q *github.QueryForCommitOutput) error {
	gitObj := gitobject.Input{
		Files:         files,
		DeletedFiles:  in.DeletedPaths,
		Repository:    in.TargetRepository,
		CommitMessage: in.CommitMessage,
		Author:        in.Author,
		Committer:     in.Committer,
		NoFileMode:    in.NoFileMode,
	}
	switch {
	case in.CommitStrategy.IsFastForward():
		u.Logger.Infof("Creating a branch (%s) based on the default branch", in.TargetBranchName)
		gitObj.ParentCommitSHA = q.ParentDefaultBranchCommitSHA
		gitObj.ParentTreeSHA = q.ParentDefaultBranchTreeSHA
	case in.CommitStrategy.IsRebase():
		u.Logger.Infof("Creating a branch (%s) based on the ref (%s)", in.TargetBranchName, in.CommitStrategy.RebaseUpstream())
		gitObj.ParentCommitSHA = q.ParentRefCommitSHA
		gitObj.ParentTreeSHA = q.ParentRefTreeSHA
	case in.CommitStrategy.NoParent():
		u.Logger.Infof("Creating a branch (%s) with no parent", in.TargetBranchName)
	default:
		return fmt.Errorf("unknown commit strategy %+v", in.CommitStrategy)
	}

	u.Logger.Debugf("Creating a commit with the %d file(s)", len(gitObj.Files))
	commit, err := u.CreateGitObject.Do(ctx, gitObj)
	if err != nil {
		return fmt.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Created a commit with %d changed file(s)", commit.ChangedFiles)
	if len(files) > 0 && commit.ChangedFiles == 0 {
		u.Logger.Warnf("Nothing to commit because the branch has the same file(s)")
		return nil
	}
	if in.DryRun {
		u.Logger.Infof("Do not create %s branch due to dry-run", in.TargetBranchName)
		return nil
	}

	u.Logger.Debugf("Creating a branch (%s)", in.TargetBranchName)
	createBranchIn := github.CreateBranchInput{
		RepositoryNodeID: q.TargetRepositoryNodeID,
		BranchName:       in.TargetBranchName,
		CommitSHA:        commit.CommitSHA,
	}
	if err := u.GitHub.CreateBranch(ctx, createBranchIn); err != nil {
		return fmt.Errorf("error while creating %s branch: %w", in.TargetBranchName, err)
	}
	u.Logger.Infof("Created a branch (%s)", in.TargetBranchName)
	return nil
}

func (u *Commit) updateExistingBranch(ctx context.Context, in Input, files []fs.File, q *github.QueryForCommitOutput) error {
	gitObj := gitobject.Input{
		Files:         files,
		DeletedFiles:  in.DeletedPaths,
		Repository:    in.TargetRepository,
		CommitMessage: in.CommitMessage,
		Author:        in.Author,
		Committer:     in.Committer,
		NoFileMode:    in.NoFileMode,
	}
	switch {
	case in.CommitStrategy.IsFastForward():
		u.Logger.Infof("Updating the branch (%s) by fast-forward", in.TargetBranchName)
		gitObj.ParentCommitSHA = q.TargetBranchCommitSHA
		gitObj.ParentTreeSHA = q.TargetBranchTreeSHA
	case in.CommitStrategy.IsRebase():
		u.Logger.Infof("Rebasing the branch (%s) on the ref (%s)", in.TargetBranchName, in.CommitStrategy.RebaseUpstream())
		gitObj.ParentCommitSHA = q.ParentRefCommitSHA
		gitObj.ParentTreeSHA = q.ParentRefTreeSHA
	case in.CommitStrategy.NoParent():
		u.Logger.Infof("Updating the branch (%s) to a commit with no parent", in.TargetBranchName)
	default:
		return fmt.Errorf("unknown commit strategy %+v", in.CommitStrategy)
	}

	u.Logger.Debugf("Creating a commit with the %d file(s)", len(gitObj.Files))
	commit, err := u.CreateGitObject.Do(ctx, gitObj)
	if err != nil {
		return fmt.Errorf("error while creating a commit: %w", err)
	}
	u.Logger.Infof("Created a commit with %d changed file(s)", commit.ChangedFiles)
	if len(files) > 0 && commit.ChangedFiles == 0 {
		u.Logger.Warnf("Nothing to commit because %s branch has the same file(s)", in.TargetBranchName)
		return nil
	}
	if in.DryRun {
		u.Logger.Infof("Do not update %s branch due to dry-run", in.TargetBranchName)
		return nil
	}

	u.Logger.Debugf("Updating the branch (%s)", in.TargetBranchName)
	updateBranchIn := github.UpdateBranchInput{
		BranchRefNodeID: q.TargetBranchNodeID,
		CommitSHA:       commit.CommitSHA,
		Force:           in.ForceUpdate,
	}
	if err := u.GitHub.UpdateBranch(ctx, updateBranchIn); err != nil {
		return fmt.Errorf("error while updating %s branch: %w", in.TargetBranchName, err)
	}
	u.Logger.Infof("Updated the branch (%s)", in.TargetBranchName)
	return nil
}
