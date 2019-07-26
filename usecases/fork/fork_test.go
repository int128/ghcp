package fork

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestCommitToFork_Do(t *testing.T) {
	ctx := context.TODO()
	parentRepositoryID := git.RepositoryID{Owner: "upstream", Name: "repo"}
	forkedRepositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("FromDefaultBranch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		in := usecases.CommitToForkIn{
			ParentRepository: parentRepositoryID,
			TargetBranchName: "topic",
			CommitMessage:    "message",
			Paths:            []string{"path"},
		}

		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			CreateFork(ctx, parentRepositoryID).
			Return(&forkedRepositoryID, nil)

		commitUseCase := mock_usecases.NewMockCommit(ctrl)
		commitUseCase.EXPECT().
			Do(ctx, usecases.CommitIn{
				ParentRepository: parentRepositoryID,
				ParentBranch:     usecases.ParentBranch{FastForward: true},
				TargetRepository: forkedRepositoryID,
				TargetBranchName: "topic",
				CommitMessage:    "message",
				Paths:            []string{"path"},
			})

		u := CommitToFork{
			Commit: commitUseCase,
			GitHub: gitHub,
			Logger: mock_adaptors.NewLogger(t),
		}
		if err := u.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})

	t.Run("FromBranch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		in := usecases.CommitToForkIn{
			ParentRepository: parentRepositoryID,
			ParentBranchName: "develop",
			TargetBranchName: "topic",
			CommitMessage:    "message",
			Paths:            []string{"path"},
		}

		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			CreateFork(ctx, parentRepositoryID).
			Return(&forkedRepositoryID, nil)

		commitUseCase := mock_usecases.NewMockCommit(ctrl)
		commitUseCase.EXPECT().
			Do(ctx, usecases.CommitIn{
				ParentRepository: parentRepositoryID,
				ParentBranch:     usecases.ParentBranch{FromRef: "develop"},
				TargetRepository: forkedRepositoryID,
				TargetBranchName: "topic",
				CommitMessage:    "message",
				Paths:            []string{"path"},
			})

		u := CommitToFork{
			Commit: commitUseCase,
			GitHub: gitHub,
			Logger: mock_adaptors.NewLogger(t),
		}
		if err := u.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})
}
