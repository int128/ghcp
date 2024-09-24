package github

import (
	"context"
	"testing"

	"github.com/google/go-github/v63/github"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/github/client_mock"
	"github.com/int128/ghcp/pkg/git"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
)

func TestGitHub_CreateCommit(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("SingleParent", func(t *testing.T) {
		gitHubClient := client_mock.NewMockInterface(t)
		gitHubClient.EXPECT().
			CreateCommit(ctx, "owner", "repo", &github.Commit{
				Message: github.String("message"),
				Parents: []*github.Commit{{SHA: github.String("parentCommitSHA")}},
				Tree:    &github.Tree{SHA: github.String("treeSHA")},
			}, (*github.CreateCommitOptions)(nil)).
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
		gitHubClient := client_mock.NewMockInterface(t)
		gitHubClient.EXPECT().
			CreateCommit(ctx, "owner", "repo", &github.Commit{
				Message: github.String("message"),
				Tree:    &github.Tree{SHA: github.String("treeSHA")},
			}, (*github.CreateCommitOptions)(nil)).
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
