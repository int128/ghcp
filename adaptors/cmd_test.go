package adaptors

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure/interfaces"
	"github.com/int128/ghcp/infrastructure/mock_infrastructure"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestCmd_Run(t *testing.T) {
	ctx := context.TODO()

	t.Run("FullOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)

		clientInit := mock_infrastructure.NewMockGitHubClientInit(ctrl)
		clientInit.EXPECT().
			Init(infrastructure.GitHubClientInitOptions{
				Token: "YOUR_TOKEN",
			})

		push := mock_usecases.NewMockPush(ctrl)
		push.EXPECT().
			Do(ctx, usecases.PushIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			Push:             push,
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			GitHubClientInit: clientInit,
		}
		args := []string{
			"ghcp",
			"-token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("EnvGitHubToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Get(envGitHubToken).
			Return("YOUR_TOKEN")

		clientInit := mock_infrastructure.NewMockGitHubClientInit(ctrl)
		clientInit.EXPECT().
			Init(infrastructure.GitHubClientInitOptions{
				Token: "YOUR_TOKEN",
			})

		push := mock_usecases.NewMockPush(ctrl)
		push.EXPECT().
			Do(ctx, usecases.PushIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			Push:             push,
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			GitHubClientInit: clientInit,
		}
		args := []string{
			"ghcp",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
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

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Get(envGitHubToken).
			Return("")

		clientInit := mock_infrastructure.NewMockGitHubClientInit(ctrl)
		push := mock_usecases.NewMockPush(ctrl)

		cmd := Cmd{
			Push:             push,
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			GitHubClientInit: clientInit,
		}
		args := []string{
			"ghcp",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != 1 {
			t.Errorf("exitCode wants 1 but %d", exitCode)
		}
	})
}
