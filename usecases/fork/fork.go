package fork

import (
	"context"

	"github.com/google/wire"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(CommitToFork), "*"),
	wire.Bind(new(usecases.CommitToFork), new(*CommitToFork)),
)

type CommitToFork struct {
	Commit usecases.Commit
	Logger adaptors.Logger
	GitHub adaptors.GitHub
}

func (u *CommitToFork) Do(ctx context.Context, in usecases.CommitToForkIn) error {
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
		return xerrors.Errorf("could not create a fork: %w", err)
	}
	if err := u.Commit.Do(ctx, usecases.CommitIn{
		ParentRepository: in.ParentRepository,
		ParentBranch: usecases.ParentBranch{
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
		return xerrors.Errorf("could not commit to the fork: %w", err)
	}
	return nil
}
