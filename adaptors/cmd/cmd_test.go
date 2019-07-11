package cmd

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

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			UpdateBranch:     updateBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				BranchName:    "gh-pages",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			UpdateBranch:     updateBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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

	t.Run("--new-branch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		createBranch := mock_usecases.NewMockCreateBranch(ctrl)
		createBranch.EXPECT().
			Do(ctx, usecases.CreateBranchIn{
				Repository:        git.RepositoryID{Owner: "owner", Name: "repo"},
				NewBranchName:     "topic",
				ParentOfNewBranch: usecases.ParentOfNewBranch{FromDefaultBranch: true},
				CommitMessage:     "commit-message",
				Paths:             []string{"file1", "file2"},
			})

		cmd := Cmd{
			CreateBranch:     createBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
			"-B", "topic",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--parent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		createBranch := mock_usecases.NewMockCreateBranch(ctrl)
		createBranch.EXPECT().
			Do(ctx, usecases.CreateBranchIn{
				Repository:        git.RepositoryID{Owner: "owner", Name: "repo"},
				NewBranchName:     "topic",
				ParentOfNewBranch: usecases.ParentOfNewBranch{FromRef: "develop"},
				CommitMessage:     "commit-message",
				Paths:             []string{"file1", "file2"},
			})

		cmd := Cmd{
			CreateBranch:     createBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
			"-B", "topic",
			"--parent", "develop",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--no-parent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		createBranch := mock_usecases.NewMockCreateBranch(ctrl)
		createBranch.EXPECT().
			Do(ctx, usecases.CreateBranchIn{
				Repository:        git.RepositoryID{Owner: "owner", Name: "repo"},
				NewBranchName:     "topic",
				ParentOfNewBranch: usecases.ParentOfNewBranch{NoParent: true},
				CommitMessage:     "commit-message",
				Paths:             []string{"file1", "file2"},
			})

		cmd := Cmd{
			CreateBranch:     createBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
			"-B", "topic",
			"--no-parent",
			"file1",
			"file2",
		}
		exitCode := cmd.Run(ctx, args)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("BadOptions", func(t *testing.T) {
		t.Run("--branch and --new-branch", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := Cmd{
				CreateBranch:     mock_usecases.NewMockCreateBranch(ctrl),
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
				"-b", "topic",
				"-B", "gh-pages",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args)
			if exitCode != exitCodePreconditionError {
				t.Errorf("exitCode wants %d but %d", exitCodePreconditionError, exitCode)
			}
		})

		t.Run("--branch and --parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := Cmd{
				CreateBranch:     mock_usecases.NewMockCreateBranch(ctrl),
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
				"-b", "topic",
				"--parent", "develop",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args)
			if exitCode != exitCodePreconditionError {
				t.Errorf("exitCode wants %d but %d", exitCodePreconditionError, exitCode)
			}
		})

		t.Run("--parent and --no-parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := Cmd{
				CreateBranch:     mock_usecases.NewMockCreateBranch(ctrl),
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
				"-b", "topic",
				"--parent", "develop",
				"--no-parent",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args)
			if exitCode != exitCodePreconditionError {
				t.Errorf("exitCode wants %d but %d", exitCodePreconditionError, exitCode)
			}
		})
	})

	t.Run("--no-file-mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
				NoFileMode:    true,
			})

		cmd := Cmd{
			UpdateBranch:     updateBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
				DryRun:        true,
			})

		cmd := Cmd{
			UpdateBranch:     updateBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		loggerConfig := mock_adaptors.NewMockLoggerConfig(ctrl)
		loggerConfig.EXPECT().
			SetDebug(true)

		cmd := Cmd{
			UpdateBranch:     updateBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		env := newEnv(ctrl, map[string]string{envGitHubAPI: ""})
		env.EXPECT().
			Chdir("dir")

		cmd := Cmd{
			UpdateBranch:     updateBranch,
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

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			UpdateBranch:     updateBranch,
			Env:              newEnv(ctrl, map[string]string{envGitHubToken: "YOUR_TOKEN", envGitHubAPI: ""}),
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

	t.Run("Error/NoGitHubToken", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := Cmd{
			UpdateBranch:     mock_usecases.NewMockUpdateBranch(ctrl),
			Env:              newEnv(ctrl, map[string]string{envGitHubToken: ""}),
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

	t.Run("--api", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			UpdateBranch: updateBranch,
			Env:          newEnv(ctrl, map[string]string{}),
			Logger:       mock_adaptors.NewLogger(t),
			LoggerConfig: mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{
				Token: "YOUR_TOKEN",
				URLv3: "https://github.example.com/api/v3/",
			}),
		}
		args := []string{
			cmdName,
			"--token", "YOUR_TOKEN",
			"--api", "https://github.example.com/api/v3/",
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

	t.Run("env/GITHUB_API", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		updateBranch := mock_usecases.NewMockUpdateBranch(ctrl)
		updateBranch.EXPECT().
			Do(ctx, usecases.UpdateBranchIn{
				Repository:    git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitMessage: "commit-message",
				Paths:         []string{"file1", "file2"},
			})

		cmd := Cmd{
			UpdateBranch: updateBranch,
			Env:          newEnv(ctrl, map[string]string{envGitHubAPI: "https://github.example.com/api/v3/"}),
			Logger:       mock_adaptors.NewLogger(t),
			LoggerConfig: mock_adaptors.NewMockLoggerConfig(ctrl),
			GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{
				Token: "YOUR_TOKEN",
				URLv3: "https://github.example.com/api/v3/",
			}),
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
}

func newGitHubClientInit(ctrl *gomock.Controller, o infrastructure.GitHubClientInitOptions) *mock_infrastructure.MockGitHubClientInit {
	clientInit := mock_infrastructure.NewMockGitHubClientInit(ctrl)
	clientInit.EXPECT().
		Init(o)
	return clientInit
}

func newEnv(ctrl *gomock.Controller, getenv map[string]string) *mock_adaptors.MockEnv {
	env := mock_adaptors.NewMockEnv(ctrl)
	for k, v := range getenv {
		env.EXPECT().Getenv(k).Return(v)
	}
	return env
}
