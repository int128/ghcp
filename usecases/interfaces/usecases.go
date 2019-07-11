package usecases

import (
	"context"

	adaptors2 "github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
)

//go:generate mockgen -package mock_usecases -destination ../mock_usecases/mock_usecases.go github.com/int128/ghcp/usecases/interfaces UpdateBranch,CreateBranch,Commit

type UpdateBranch interface {
	Do(ctx context.Context, in UpdateBranchIn) error
}

type UpdateBranchIn struct {
	Repository    git.RepositoryID
	BranchName    git.BranchName // default branch if empty
	CommitMessage git.CommitMessage
	Paths         []string
	NoFileMode    bool
	DryRun        bool
}

type CreateBranch interface {
	Do(ctx context.Context, in CreateBranchIn) error
}

type CreateBranchIn struct {
	Repository        git.RepositoryID
	NewBranchName     git.BranchName
	ParentOfNewBranch ParentOfNewBranch
	CommitMessage     git.CommitMessage
	Paths             []string
	NoFileMode        bool
	DryRun            bool
}

// ParentOfNewBranch represents a parent of a branch to create.
// Exact one of the members must be valid.
type ParentOfNewBranch struct {
	NoParent          bool
	FromDefaultBranch bool
	FromRef           git.RefName
}

type Commit interface {
	Do(ctx context.Context, in CommitIn) (*CommitOut, error)
}

type CommitIn struct {
	Files           []adaptors2.File
	Repository      git.RepositoryID
	CommitMessage   git.CommitMessage
	ParentCommitSHA git.CommitSHA // no parent if empty
	ParentTreeSHA   git.TreeSHA   // no parent if empty
	NoFileMode      bool
}

type CommitOut struct {
	CommitSHA    git.CommitSHA
	ChangedFiles int
}
