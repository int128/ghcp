package forkcommit

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github/mock_github"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
	"github.com/int128/ghcp/usecases/commit"
	"github.com/int128/ghcp/usecases/commit/mock_commit"
)

func TestForkCommit_Do(t *testing.T) {
	ctx := context.TODO()
	parentRepositoryID := git.RepositoryID{Owner: "upstream", Name: "repo"}
	forkedRepositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("FromDefaultBranch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		in := Input{
			ParentRepository: parentRepositoryID,
			TargetBranchName: "topic",
			CommitStrategy:   commitstrategy.FastForward,
			CommitMessage:    "message",
			Paths:            []string{"path"},
		}

		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			CreateFork(ctx, parentRepositoryID).
			Return(&forkedRepositoryID, nil)

		commitUseCase := mock_commit.NewMockInterface(ctrl)
		commitUseCase.EXPECT().
			Do(ctx, commit.Input{
				TargetRepository: forkedRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.FastForward,
				CommitMessage:    "message",
				Paths:            []string{"path"},
			})

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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		in := Input{
			TargetBranchName: "topic",
			ParentRepository: parentRepositoryID,
			CommitStrategy:   commitstrategy.RebaseOn("develop"),
			CommitMessage:    "message",
			Paths:            []string{"path"},
		}

		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			CreateFork(ctx, parentRepositoryID).
			Return(&forkedRepositoryID, nil)

		commitUseCase := mock_commit.NewMockInterface(ctrl)
		commitUseCase.EXPECT().
			Do(ctx, commit.Input{
				TargetRepository: forkedRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.RebaseOn("develop"),
				CommitMessage:    "message",
				Paths:            []string{"path"},
			})

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
