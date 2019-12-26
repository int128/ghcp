package cmd

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure"
	"github.com/int128/ghcp/infrastructure/mock_infrastructure"
	"github.com/int128/ghcp/usecases"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestCmd_Run(t *testing.T) {
	const cmdName = "ghcp"
	const version = "TEST"
	ctx := context.TODO()

	t.Run("Commit", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--branch", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "gh-pages",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"-b", "gh-pages",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FromRef: "develop"},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"-b", "topic",
				"--parent", "develop",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--no-parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{NoParent: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"-b", "topic",
				"--no-parent",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--parent and --no-parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cmd := Cmd{
				Commit:           mock_usecases.NewMockCommit(ctrl),
				Env:              mock_adaptors.NewMockEnv(ctrl),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: mock_infrastructure.NewMockGitHubClientInit(ctrl),
			}
			args := []string{
				cmdName,
				commitCmdName,
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
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeError {
				t.Errorf("exitCode wants %d but %d", exitCodeError, exitCode)
			}
		})

		t.Run("--no-file-mode", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
					NoFileMode:       true,
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"--no-file-mode",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--dry-run", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
					DryRun:           true,
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"--dry-run",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})

	t.Run("CommitToFork", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommitToFork(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitToForkIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				CommitToFork:     commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitToForkCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-b", "topic",
				"-m", "commit-message",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommitToFork(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitToForkIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranchName: "develop",
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				CommitToFork:     commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitToForkCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"-b", "topic",
				"--parent", "develop",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})

	t.Run("GlobalOptions", func(t *testing.T) {
		t.Run("--debug", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			loggerConfig := mock_adaptors.NewMockLoggerConfig(ctrl)
			loggerConfig.EXPECT().
				SetDebug(true)
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     loggerConfig,
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"--debug",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--directory", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			env := newEnv(ctrl, map[string]string{envGitHubAPI: ""})
			env.EXPECT().
				Chdir("dir")
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              env,
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"-C", "dir",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("env/GITHUB_TOKEN", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:           commitUseCase,
				Env:              newEnv(ctrl, map[string]string{envGitHubToken: "YOUR_TOKEN", envGitHubAPI: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: newGitHubClientInit(ctrl, infrastructure.GitHubClientInitOptions{Token: "YOUR_TOKEN"}),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("NoGitHubToken", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cmd := Cmd{
				Commit:           mock_usecases.NewMockCommit(ctrl),
				Env:              newEnv(ctrl, map[string]string{envGitHubToken: ""}),
				Logger:           mock_adaptors.NewLogger(t),
				LoggerConfig:     mock_adaptors.NewMockLoggerConfig(ctrl),
				GitHubClientInit: mock_infrastructure.NewMockGitHubClientInit(ctrl),
			}
			args := []string{
				cmdName,
				commitCmdName,
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeError {
				t.Errorf("exitCode wants %d but %d", exitCodeError, exitCode)
			}
		})

		t.Run("--api", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:       commitUseCase,
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
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"--api", "https://github.example.com/api/v3/",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("env/GITHUB_API", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     usecases.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			cmd := Cmd{
				Commit:       commitUseCase,
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
				commitCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-m", "commit-message",
				"file1",
				"file2",
			}
			exitCode := cmd.Run(ctx, args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
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
