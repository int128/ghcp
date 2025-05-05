package pullrequest

import (
	"context"
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/github_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
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
		}
		t.Run("when the pull request does not exist", func(t *testing.T) {
			gitHub := github_mock.NewMockInterface(t)
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
				}).
				Return(&github.CreatePullRequestOutput{
					URL: "https://github.com/octocat/Spoon-Knife/pull/19445",
				}, nil)
			useCase := PullRequest{
				GitHub: gitHub,
			}
			if err := useCase.Do(ctx, in); err != nil {
				t.Errorf("err wants nil but %+v", err)
			}
		})
		t.Run("when the pull request already exists", func(t *testing.T) {
			gitHub := github_mock.NewMockInterface(t)
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
			}
			if err := useCase.Do(ctx, in); err != nil {
				t.Errorf("err wants nil but %+v", err)
			}
		})
		t.Run("when the head branch does not exist", func(t *testing.T) {
			gitHub := github_mock.NewMockInterface(t)
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
		}
		gitHub := github_mock.NewMockInterface(t)
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
			}).
			Return(&github.CreatePullRequestOutput{
				URL: "https://github.com/octocat/Spoon-Knife/pull/19445",
			}, nil)
		useCase := PullRequest{
			GitHub: gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})

	t.Run("when a reviewer is set", func(t *testing.T) {
		in := Input{
			BaseRepository: baseRepositoryID,
			BaseBranchName: "develop",
			HeadRepository: headRepositoryID,
			HeadBranchName: "feature",
			Title:          "the-title",
			Reviewer:       "the-reviewer",
		}
		gitHub := github_mock.NewMockInterface(t)
		gitHub.EXPECT().
			QueryForPullRequest(ctx, github.QueryForPullRequestInput{
				BaseRepository: baseRepositoryID,
				BaseBranchName: "develop",
				HeadRepository: headRepositoryID,
				HeadBranchName: "feature",
				ReviewerUser:   "the-reviewer",
			}).
			Return(&github.QueryForPullRequestOutput{
				CurrentUserName:     "you",
				HeadBranchCommitSHA: "HeadCommitSHA",
				ReviewerUserNodeID:  "TheReviewerID",
			}, nil)
		gitHub.EXPECT().
			CreatePullRequest(ctx, github.CreatePullRequestInput{
				BaseRepository: baseRepositoryID,
				BaseBranchName: "develop",
				HeadRepository: headRepositoryID,
				HeadBranchName: "feature",
				Title:          "the-title",
			}).
			Return(&github.CreatePullRequestOutput{
				URL:               "https://github.com/octocat/Spoon-Knife/pull/19445",
				PullRequestNodeID: "ThePullRequestID",
			}, nil)
		gitHub.EXPECT().
			RequestPullRequestReview(ctx, github.RequestPullRequestReviewInput{
				PullRequest: "ThePullRequestID",
				User:        "TheReviewerID",
			}).
			Return(nil)
		useCase := PullRequest{
			GitHub: gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})

	t.Run("when optional values are set", func(t *testing.T) {
		in := Input{
			BaseRepository: baseRepositoryID,
			BaseBranchName: "develop",
			HeadRepository: headRepositoryID,
			HeadBranchName: "feature",
			Title:          "the-title",
			Body:           "the-body",
			Draft:          true,
		}
		gitHub := github_mock.NewMockInterface(t)
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
				ReviewerUserNodeID:  "TheReviewerID",
			}, nil)
		gitHub.EXPECT().
			CreatePullRequest(ctx, github.CreatePullRequestInput{
				BaseRepository: baseRepositoryID,
				BaseBranchName: "develop",
				HeadRepository: headRepositoryID,
				HeadBranchName: "feature",
				Title:          "the-title",
				Body:           "the-body",
				Draft:          true,
			}).
			Return(&github.CreatePullRequestOutput{
				URL:               "https://github.com/octocat/Spoon-Knife/pull/19445",
				PullRequestNodeID: "ThePullRequestID",
			}, nil)
		useCase := PullRequest{
			GitHub: gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})
}
