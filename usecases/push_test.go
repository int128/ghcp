package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/usecases/interfaces"
)

func TestPush_Do(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitHub := mock_adaptors.NewMockGitHub(ctrl)
	gitHub.EXPECT().
		GetRepository(ctx, adaptors.GetRepositoryIn{
			Owner: "owner",
			Name:  "repo",
		}).
		Return(&adaptors.GetRepositoryOut{
			CurrentUserName:   "current",
			DefaultBranchName: "master",
		}, nil)

	push := Push{GitHub: gitHub}
	err := push.Do(ctx, usecases.PushIn{
		RepositoryOwner: "owner",
		RepositoryName:  "repo",
		Paths:           []string{"file"},
	})
	if err != nil {
		t.Errorf("Do returned error: %+v", err)
	}
}
