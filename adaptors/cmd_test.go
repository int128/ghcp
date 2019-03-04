package adaptors

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestCmd_Run(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("NoArgs", func(t *testing.T) {
		push := mock_usecases.NewMockPush(ctrl)
		c := Cmd{Push: push}
		err := c.Run(ctx, adaptors.CmdOptions{})
		if err == nil {
			t.Errorf("err wants non-nil but nil")
		}
	})

	t.Run("OK", func(t *testing.T) {
		push := mock_usecases.NewMockPush(ctrl)
		push.EXPECT().
			Do(ctx, usecases.PushIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file"},
			})

		c := Cmd{Push: push}
		err := c.Run(ctx, adaptors.CmdOptions{
			RepositoryOwner: "owner",
			RepositoryName:  "repo",
			CommitMessage:   "commit-message",
			Paths:           []string{"file"},
		})
		if err != nil {
			t.Errorf("Run returned error: %+v", err)
		}
	})
}
