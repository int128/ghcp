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
	const cmdName = "ghcp"
	ctx := context.TODO()

	t.Run("BasicOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--branch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				BranchName:    "gh-pages",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"-b", "gh-pages",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--no-file-mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
				NoFileMode:    true,
			})

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"--no-file-mode",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--dry-run", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
				DryRun:        true,
			})

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"--dry-run",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--debug", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		loggerConfig := mock_adaptors.NewMockLoggerConfig(ctrl)
		loggerConfig.EXPECT().
			SetDebug(true)

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              mock_adaptors.NewMockEnv(ctrl),
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     loggerConfig,
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"--debug",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--directory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Chdir("dir")

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"-C", "dir",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("env/GITHUB_TOKEN", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Getenv(envGitHubToken).
			Return("YOUR_TOKEN")

		copyUseCase := mock_usecases.NewMockCopyUseCase(ctrl)
		copyUseCase.EXPECT().
			Do(ctx, usecases.CopyUseCaseIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			CopyUseCase:      copyUseCase,
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
		}
		args := []string{
			cmdName,
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("NoGitHubToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Getenv(envGitHubToken).
			Return("")

		cmd := Cmd{
			CopyUseCase:      mock_usecases.NewMockCopyUseCase(ctrl),
			Env:              env,
			Logger:           mock_adaptors.NewLogger(t),
			LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: mock_infrastructure.NewMockGitHubClientInit(ctrl),
		}
		args := []string{
			cmdName,
			"-u", "owner",
			"-r", "repo",
			"-m", "commit-message",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodePreconditionError {
			t.Errorf("exitCode wants %d but %d", exitCodePreconditionError, exitCode)
		}
	})

	t.Run("InvalidArguments", func(t *testing.T) {
		for _, c := range []struct {
			name string
			args []string
		}{
			{
				"NoOwner",
				[]string{
					cmdName,
					"--token", "YOUR_TOKEN",
					"-r", "repo",
					"-m", "commit-message",
					"file1",
					"file2",
				},
			}, {
				"NoRepo",
				[]string{
					cmdName,
					"--token", "YOUR_TOKEN",
					"-u", "owner",
					"-m", "commit-message",
					"file1",
					"file2",
				},
			}, {
				"NoMessage",
				[]string{
					cmdName,
					"--token", "YOUR_TOKEN",
					"-r", "repo",
					"-u", "owner",
					"file1",
					"file2",
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				cmd := Cmd{
					CopyUseCase:      mock_usecases.NewMockCopyUseCase(ctrl),
					Env:              mock_adaptors.NewMockEnv(ctrl),
					Logger:           mock_adaptors.NewLogger(t),
					LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
					GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
				}
				exitCode := cmd.Run(ctx, c.args)
				if exitCode != exitCodePreconditionError {
					t.Errorf("exitCode wants %d but %d", exitCodePreconditionError, exitCode)
				}
			})
		}
	})
}

func newGitHubClientInit(ctrl *gomock.Controller, o infrastructure.GitHubClientInitOptions) *mock_infrastructure.MockGitHubClientInit {
	clientInit := mock_infrastructure.NewMockGitHubClientInit(ctrl)
	clientInit.EXPECT().
		Init(o)
	return clientInit
}
