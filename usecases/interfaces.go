// Package usecases provides use cases of this application.
package usecases

import (
	"context"

	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
)

//go:generate mockgen -destination mock_usecases/mock_usecases.go github.com/int128/ghcp/usecases CommitToBranch,CreateBlobTreeCommit

type CommitToBranch interface {
	Do(ctx context.Context, in CommitToBranchIn) error
}

type CommitToBranchIn struct {
	Repository     git.RepositoryID
	BranchName     git.BranchName // default branch if empty
	ParentOfBranch ParentOfBranch
	CommitMessage  git.CommitMessage
	Paths          []string
	NoFileMode     bool
	DryRun         bool
}

// ParentOfBranch represents a parent ref of the branch to create or update.
// Exact one of the members must be valid.
type ParentOfBranch struct {
	NoParent    bool        // push a branch without any parent
	FastForward bool        // push the branch by fast-forward
	FromRef     git.RefName // push a branch based on the ref
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
