package github

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v74/github"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/github/client_mock"
	"github.com/int128/ghcp/pkg/git"
)

func TestGitHub_GetReleaseByTagOrNil(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("Exists", func(t *testing.T) {
		var resp github.Response
		resp.Response = &http.Response{StatusCode: 200}
		gitHubClient := client_mock.NewMockInterface(t)
		gitHubClient.EXPECT().
			GetReleaseByTag(ctx, "owner", "repo", "v1.0.0").
			Return(&github.RepositoryRelease{
				ID:      github.Ptr(int64(1234567890)),
				Name:    github.Ptr("ReleaseName"),
				TagName: github.Ptr("v1.0.0"),
			}, &resp, nil)
		gitHub := GitHub{
			Client: gitHubClient,
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
		var resp github.Response
		resp.Response = &http.Response{StatusCode: 404}
		gitHubClient := client_mock.NewMockInterface(t)
		gitHubClient.EXPECT().
			GetReleaseByTag(ctx, "owner", "repo", "v1.0.0").
			Return(nil, &resp, errors.New("not found"))
		gitHub := GitHub{
			Client: gitHubClient,
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
