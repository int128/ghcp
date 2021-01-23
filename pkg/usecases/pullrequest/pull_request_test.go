package pullrequest

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/github/mock_github"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
)

func TestPullRequest_Do(t *testing.T) {
	ctx := context.TODO()
	baseRepositoryID := git.RepositoryID{Owner: "base", Name: "repo"}
	headRepositoryID := git.RepositoryID{Owner: "head", Name: "repo"}

	t.Run("when head and base branch name are given", func(t *testing.T) {
		in := Input{
			BaseRepository: baseRepositoryID,
			BaseBranchName: "develop",
			HeadRepository: headRepositoryID,
			HeadBranchName: "feature",
			Title:          "the-title",
			Body:           "the-body",
		}
		t.Run("when the pull request does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			gitHub := mock_github.NewMockInterface(ctrl)
			gitHub.EXPECT().
				QueryForPullRequest(ctx, github.QueryForPullRequestInput{
					BaseRepository: baseRepositoryID,
					BaseBranchName: "develop",
					HeadRepository: headRepositoryID,
					HeadBranchName: "feature",
				}).
				Return(&github.QueryForPullRequestOutput{
					CurrentUserName:     "you",
					HeadBranchCommitSHA: "HeadCommitSHA",
				}, nil)
			gitHub.EXPECT().
				CreatePullRequest(ctx, github.CreatePullRequestInput{
					BaseRepository: baseRepositoryID,
					BaseBranchName: "develop",
					HeadRepository: headRepositoryID,
					HeadBranchName: "feature",
					Title:          "the-title",
					Body:           "the-body",
				}).
				Return(&github.CreatePullRequestOutput{
					URL: "https://github.com/octocat/Spoon-Knife/pull/19445",
				}, nil)
			useCase := PullRequest{
				GitHub: gitHub,
				Logger: testingLogger.New(t),
			}
			if err := useCase.Do(ctx, in); err != nil {
				t.Errorf("err wants nil but %+v", err)
			}
		})
		t.Run("when the pull request already exists", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			gitHub := mock_github.NewMockInterface(ctrl)
			gitHub.EXPECT().
				QueryForPullRequest(ctx, github.QueryForPullRequestInput{
					BaseRepository: baseRepositoryID,
					BaseBranchName: "develop",
					HeadRepository: headRepositoryID,
					HeadBranchName: "feature",
				}).
				Return(&github.QueryForPullRequestOutput{
					CurrentUserName:     "you",
					HeadBranchCommitSHA: "HeadCommitSHA",
					PullRequestURL:      "https://github.com/octocat/Spoon-Knife/pull/19445",
				}, nil)
			useCase := PullRequest{
				GitHub: gitHub,
				Logger: testingLogger.New(t),
			}
			if err := useCase.Do(ctx, in); err != nil {
				t.Errorf("err wants nil but %+v", err)
			}
		})
		t.Run("when the head branch does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			gitHub := mock_github.NewMockInterface(ctrl)
			gitHub.EXPECT().
				QueryForPullRequest(ctx, github.QueryForPullRequestInput{
					BaseRepository: baseRepositoryID,
					BaseBranchName: "develop",
					HeadRepository: headRepositoryID,
					HeadBranchName: "feature",
				}).
				Return(&github.QueryForPullRequestOutput{
					CurrentUserName: "you",
				}, nil)
			useCase := PullRequest{
				GitHub: gitHub,
				Logger: testingLogger.New(t),
			}
			if err := useCase.Do(ctx, in); err == nil {
				t.Errorf("err wants non-nil but got nil")
			}
		})
	})

	t.Run("when the default base branch is given", func(t *testing.T) {
		in := Input{
			BaseRepository: baseRepositoryID,
			HeadRepository: headRepositoryID,
			BaseBranchName: "staging",
			Title:          "the-title",
			Body:           "the-body",
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			QueryDefaultBranch(ctx, github.QueryDefaultBranchInput{
				BaseRepository: baseRepositoryID,
				HeadRepository: headRepositoryID,
			}).
			Return(&github.QueryDefaultBranchOutput{
				BaseDefaultBranchName: "master",
				HeadDefaultBranchName: "develop",
			}, nil)
		gitHub.EXPECT().
			QueryForPullRequest(ctx, github.QueryForPullRequestInput{
				BaseRepository: baseRepositoryID,
				BaseBranchName: "staging",
				HeadRepository: headRepositoryID,
				HeadBranchName: "develop",
			}).
			Return(&github.QueryForPullRequestOutput{
				CurrentUserName:     "you",
				HeadBranchCommitSHA: "HeadCommitSHA",
			}, nil)
		gitHub.EXPECT().
			CreatePullRequest(ctx, github.CreatePullRequestInput{
				BaseRepository: baseRepositoryID,
				BaseBranchName: "staging",
				HeadRepository: headRepositoryID,
				HeadBranchName: "develop",
				Title:          "the-title",
				Body:           "the-body",
			}).
			Return(&github.CreatePullRequestOutput{
				URL: "https://github.com/octocat/Spoon-Knife/pull/19445",
			}, nil)
		useCase := PullRequest{
			GitHub: gitHub,
			Logger: testingLogger.New(t),
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})
}
