package forkcommit

import (
	"context"

	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/domain/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/int128/ghcp/usecases/commit"
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
	TargetBranchName git.BranchName
	CommitStrategy   commitstrategy.CommitStrategy
	CommitMessage    git.CommitMessage
	Author           *git.CommitAuthor // optional
	Committer        *git.CommitAuthor // optional
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
		TargetRepository: *fork,
		TargetBranchName: in.TargetBranchName,
		ParentRepository: in.ParentRepository,
		CommitStrategy:   in.CommitStrategy,
		CommitMessage:    in.CommitMessage,
		Author:           in.Author,
		Committer:        in.Committer,
		Paths:            in.Paths,
		NoFileMode:       in.NoFileMode,
		DryRun:           in.DryRun,
	}); err != nil {
		return xerrors.Errorf("could not fork and commit: %w", err)
	}
	return nil
}
