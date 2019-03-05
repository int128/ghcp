package infrastructure_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/infrastructure"
	"github.com/int128/ghcp/infrastructure/mock_infrastructure"
)

func TestCmd_Run(t *testing.T) {
	ctx := context.TODO()

	t.Run("FullOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmdAdaptor := mock_adaptors.NewMockCmd(ctrl)
		cmdAdaptor.EXPECT().
			Run(ctx, adaptors.CmdOptions{
				RepositoryOwner: "owner",
				RepositoryName:  "repo",
				CommitMessage:   "commit message",
				Paths:           []string{"file1", "file2"},
			}).
			Return(nil)

		clientConfig := mock_infrastructure.NewMockGitHubClientConfig(ctrl)
		clientConfig.EXPECT().
			SetToken("YOUR_TOKEN")

		cmd := infrastructure.Cmd{
			Cmd:                cmdAdaptor,
			GitHubClientConfig: clientConfig,
		}
		args := []string{
			"ghcp",
			"-token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit message",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("NoGitHubToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		cmdAdaptor := mock_adaptors.NewMockCmd(ctrl)
		clientConfig := mock_infrastructure.NewMockGitHubClientConfig(ctrl)
		cmd := infrastructure.Cmd{
			Cmd:                cmdAdaptor,
			GitHubClientConfig: clientConfig,
		}
		args := []string{"ghcp"}
		exitCode := cmd.Run(ctx, args)
		if exitCode != 1 {
			t.Errorf("exitCode wants 1 but %d", exitCode)
		}
	})
}
