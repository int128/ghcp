package usecases

import (
	"context"

	"github.com/int128/ghcp/adaptors/interfaces"
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
	Repository    git.RepositoryID
	NewBranchName git.BranchName
	ParentRef     git.RefName // default branch if empty
	CommitMessage git.CommitMessage
	Paths         []string
	NoFileMode    bool
	DryRun        bool
}

type Commit interface {
	Do(ctx context.Context, in CommitIn) (*CommitOut, error)
}

type CommitIn struct {
	Files           []adaptors.File
	Repository      git.RepositoryID
	CommitMessage   git.CommitMessage
	ParentCommitSHA git.CommitSHA
	ParentTreeSHA   git.TreeSHA
	NoFileMode      bool
}

type CommitOut struct {
	CommitSHA    git.CommitSHA
	ChangedFiles int
}
