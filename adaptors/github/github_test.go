package github

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v24/github"
	testingLogger "github.com/int128/ghcp/adaptors/logger/testing"
	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/infrastructure/github/mock_github"
	"golang.org/x/xerrors"
)

func TestGitHub_CreateCommit(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("SingleParent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gitHubClient := mock_github.NewMockInterface(ctrl)
		gitHubClient.EXPECT().
			CreateCommit(ctx, "owner", "repo", &github.Commit{
				Message: github.String("message"),
				Parents: []github.Commit{{SHA: github.String("parentCommitSHA")}},
				Tree:    &github.Tree{SHA: github.String("treeSHA")},
			}).
			Return(&github.Commit{
				SHA: github.String("commitSHA"),
			}, nil, nil)
		gitHub := GitHub{
			Client: gitHubClient,
			Logger: testingLogger.New(t),
		}
		commitSHA, err := gitHub.CreateCommit(ctx, git.NewCommit{
			Repository:      repositoryID,
			Message:         "message",
			ParentCommitSHA: "parentCommitSHA",
			TreeSHA:         "treeSHA",
		})
		if err != nil {
			t.Fatalf("CreateCommit returned error: %+v", err)
		}
		if commitSHA != "commitSHA" {
			t.Errorf("commitSHA wants commitSHA but %s", commitSHA)
		}
	})

	t.Run("NoParent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gitHubClient := mock_github.NewMockInterface(ctrl)
		gitHubClient.EXPECT().
			CreateCommit(ctx, "owner", "repo", &github.Commit{
				Message: github.String("message"),
				Tree:    &github.Tree{SHA: github.String("treeSHA")},
			}).
			Return(&github.Commit{
				SHA: github.String("commitSHA"),
			}, nil, nil)
		gitHub := GitHub{
			Client: gitHubClient,
			Logger: testingLogger.New(t),
		}
		commitSHA, err := gitHub.CreateCommit(ctx, git.NewCommit{
			Repository: repositoryID,
			Message:    "message",
			TreeSHA:    "treeSHA",
		})
		if err != nil {
			t.Fatalf("CreateCommit returned error: %+v", err)
		}
		if commitSHA != "commitSHA" {
			t.Errorf("commitSHA wants commitSHA but %s", commitSHA)
		}
	})
}

func TestGitHub_GetReleaseByTagOrNil(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("Exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var resp github.Response
		resp.Response = &http.Response{StatusCode: 200}
		gitHubClient := mock_github.NewMockInterface(ctrl)
		gitHubClient.EXPECT().
			GetReleaseByTag(ctx, "owner", "repo", "v1.0.0").
			Return(&github.RepositoryRelease{
				ID:      github.Int64(1234567890),
				Name:    github.String("ReleaseName"),
				TagName: github.String("v1.0.0"),
			}, &resp, nil)
		gitHub := GitHub{
			Client: gitHubClient,
			Logger: testingLogger.New(t),
		}
		got, err := gitHub.GetReleaseByTagOrNil(ctx, repositoryID, "v1.0.0")
		if err != nil {
			t.Fatalf("GetReleaseByTagOrNil returned error: %+v", err)
		}
		want := &git.Release{
			ID: git.ReleaseID{
				Repository: repositoryID,
				InternalID: 1234567890,
			},
			TagName: "v1.0.0",
			Name:    "ReleaseName",
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("NotExist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		var resp github.Response
		resp.Response = &http.Response{StatusCode: 404}
		gitHubClient := mock_github.NewMockInterface(ctrl)
		gitHubClient.EXPECT().
			GetReleaseByTag(ctx, "owner", "repo", "v1.0.0").
			Return(nil, &resp, xerrors.New("not found"))
		gitHub := GitHub{
			Client: gitHubClient,
			Logger: testingLogger.New(t),
		}
		got, err := gitHub.GetReleaseByTagOrNil(ctx, repositoryID, "v1.0.0")
		if err != nil {
			t.Fatalf("GetReleaseByTagOrNil returned error: %+v", err)
		}
		if got != nil {
			t.Errorf("wants nil but got %+v", got)
		}
	})
}
