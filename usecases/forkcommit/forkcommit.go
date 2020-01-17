package forkcommit

import (
	"context"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors/github"
	"github.com/int128/ghcp/adaptors/logger"
	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/usecases/commit"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(ForkCommit), "*"),
	wire.Bind(new(Interface), new(*ForkCommit)),
)

//go:generate mockgen -destination mock_forkcommit/mock_forkcommit.go github.com/int128/ghcp/usecases/forkcommit Interface

type Interface interface {
	Do(ctx context.Context, in Input) error
}

type Input struct {
	ParentRepository git.RepositoryID
	ParentBranchName git.BranchName // if empty, the default branch of the parent repository
	TargetBranchName git.BranchName
	CommitMessage    git.CommitMessage
	Paths            []string
	NoFileMode       bool
	DryRun           bool
}

type ForkCommit struct {
	Commit commit.Interface
	Logger logger.Interface
	GitHub github.Interface
}

func (u *ForkCommit) Do(ctx context.Context, in Input) error {
	if !in.ParentRepository.IsValid() {
		return xerrors.New("you must set GitHub repository")
	}
	if in.TargetBranchName == "" {
		return xerrors.New("you must set target branch name")
	}
	if in.CommitMessage == "" {
		return xerrors.New("you must set commit message")
	}
	if len(in.Paths) == 0 {
		return xerrors.New("you must set one or more paths")
	}

	fork, err := u.GitHub.CreateFork(ctx, in.ParentRepository)
	if err != nil {
		return xerrors.Errorf("could not fork the repository: %w", err)
	}
	if err := u.Commit.Do(ctx, commit.Input{
		ParentRepository: in.ParentRepository,
		ParentBranch: commit.ParentBranch{
			FastForward: in.ParentBranchName == "",
			FromRef:     git.RefName(in.ParentBranchName),
		},
		TargetRepository: *fork,
		TargetBranchName: in.TargetBranchName,
		CommitMessage:    in.CommitMessage,
		Paths:            in.Paths,
		NoFileMode:       in.NoFileMode,
		DryRun:           in.DryRun,
	}); err != nil {
		return xerrors.Errorf("could not fork and commit: %w", err)
	}
	return nil
}
