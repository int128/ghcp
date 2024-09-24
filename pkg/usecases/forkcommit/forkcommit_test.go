package forkcommit

import (
	"context"
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/github_mock"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/commit_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
	"github.com/int128/ghcp/pkg/usecases/commit"
)

func TestForkCommit_Do(t *testing.T) {
	ctx := context.TODO()
	parentRepositoryID := git.RepositoryID{Owner: "upstream", Name: "repo"}
	forkedRepositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("FromDefaultBranch", func(t *testing.T) {
		in := Input{
			ParentRepository: parentRepositoryID,
			TargetBranchName: "topic",
			CommitStrategy:   commitstrategy.FastForward,
			CommitMessage:    "message",
			Paths:            []string{"path"},
		}

		gitHub := github_mock.NewMockInterface(t)
		gitHub.EXPECT().
			CreateFork(ctx, parentRepositoryID).
			Return(&forkedRepositoryID, nil)

		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(ctx, commit.Input{
				TargetRepository: forkedRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.FastForward,
				CommitMessage:    "message",
				Paths:            []string{"path"},
			}).
			Return(nil)

		u := ForkCommit{
			Commit: commitUseCase,
			GitHub: gitHub,
			Logger: testingLogger.New(t),
		}
		if err := u.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})

	t.Run("FromBranch", func(t *testing.T) {
		in := Input{
			TargetBranchName: "topic",
			ParentRepository: parentRepositoryID,
			CommitStrategy:   commitstrategy.RebaseOn("develop"),
			CommitMessage:    "message",
			Paths:            []string{"path"},
		}

		gitHub := github_mock.NewMockInterface(t)
		gitHub.EXPECT().
			CreateFork(ctx, parentRepositoryID).
			Return(&forkedRepositoryID, nil)

		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(ctx, commit.Input{
				TargetRepository: forkedRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.RebaseOn("develop"),
				CommitMessage:    "message",
				Paths:            []string{"path"},
			}).
			Return(nil)

		u := ForkCommit{
			Commit: commitUseCase,
			GitHub: gitHub,
			Logger: testingLogger.New(t),
		}
		if err := u.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})
}
