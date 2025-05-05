package forkcommit

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/wire"

	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/usecases/commit"
)

var Set = wire.NewSet(
	wire.Struct(new(ForkCommit), "*"),
	wire.Bind(new(Interface), new(*ForkCommit)),
)

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
	GitHub github.Interface
}

func (u *ForkCommit) Do(ctx context.Context, in Input) error {
	if !in.ParentRepository.IsValid() {
		return errors.New("you must set GitHub repository")
	}
	if in.TargetBranchName == "" {
		return errors.New("you must set target branch name")
	}
	if in.CommitMessage == "" {
		return errors.New("you must set commit message")
	}
	if len(in.Paths) == 0 {
		return errors.New("you must set one or more paths")
	}

	fork, err := u.GitHub.CreateFork(ctx, in.ParentRepository)
	if err != nil {
		return fmt.Errorf("could not fork the repository: %w", err)
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
		return fmt.Errorf("could not fork and commit: %w", err)
	}
	return nil
}
