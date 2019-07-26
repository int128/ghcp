// Package usecases provides use cases of this application.
package usecases

import (
	"context"

	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
)

//go:generate mockgen -destination mock_usecases/mock_usecases.go github.com/int128/ghcp/usecases Commit,CommitToFork,CreateBlobTreeCommit

type Commit interface {
	Do(ctx context.Context, in CommitIn) error
}

type CommitIn struct {
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

type CommitToFork interface {
	Do(ctx context.Context, in CommitToForkIn) error
}

type CommitToForkIn struct {
	ParentRepository git.RepositoryID
	ParentBranchName git.BranchName // if empty, the default branch of the parent repository
	TargetBranchName git.BranchName
	CommitMessage    git.CommitMessage
	Paths            []string
	NoFileMode       bool
	DryRun           bool
}

type CreateBlobTreeCommit interface {
	Do(ctx context.Context, in CreateBlobTreeCommitIn) (*CreateBlobTreeCommitOut, error)
}

type CreateBlobTreeCommitIn struct {
	Files           []adaptors.File
	Repository      git.RepositoryID
	CommitMessage   git.CommitMessage
	ParentCommitSHA git.CommitSHA // no parent if empty
	ParentTreeSHA   git.TreeSHA   // no parent if empty
	NoFileMode      bool
}

type CreateBlobTreeCommitOut struct {
	CommitSHA    git.CommitSHA
	ChangedFiles int
}
