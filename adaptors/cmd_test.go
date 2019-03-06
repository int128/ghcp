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
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
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
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
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

		cmd := Cmd{
			Push:             mock_usecases.NewMockPush(ctrl),
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: mock_infrastructure.NewMockGitHubClientInit(ctrl),
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

	t.Run("DryRun", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

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
				DryRun:        true,
			})

		cmd := Cmd{
			Push:             push,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: clientInit,
		}
		args := []string{
			"ghcp",
			"-token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"-dry-run",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("Debug", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

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

		loggerConfig := mock_adaptors.NewMockLoggerConfig(ctrl)
		loggerConfig.EXPECT().
			SetDebug(true)

		cmd := Cmd{
			Push:             push,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     loggerConfig,
			GitHubClientInit: clientInit,
		}
		args := []string{
			"ghcp",
			"-token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"-debug",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})
}
