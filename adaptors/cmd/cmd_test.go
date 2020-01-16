package cmd

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/int128/ghcp/adaptors/env/mock_env"
	"github.com/int128/ghcp/adaptors/logger"
	testingLogger "github.com/int128/ghcp/adaptors/logger/testing"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/infrastructure/github"
	"github.com/int128/ghcp/usecases/commit"
	"github.com/int128/ghcp/usecases/commit/mock_commit"
	"github.com/int128/ghcp/usecases/fork"
	"github.com/int128/ghcp/usecases/fork/mock_fork"
	"github.com/int128/ghcp/usecases/release"
	"github.com/int128/ghcp/usecases/release/mock_release"
)

func TestCmd_Run(t *testing.T) {
	const cmdName = "ghcp"
	const version = "TEST"
	ctx := context.TODO()

	t.Run("Commit", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--branch", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "gh-pages",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FromRef: "develop"},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--no-parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{NoParent: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--parent and --no-parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, nil),
				NewInternalRunner: newInternalRunner(InternalRunner{}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeError {
				t.Errorf("exitCode wants %d but %d", exitCodeError, exitCode)
			}
		})

		t.Run("--no-file-mode", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
					NoFileMode:       true,
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--dry-run", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
					DryRun:           true,
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})

	t.Run("CommitToFork", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_fork.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, fork.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{ForkCommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--parent", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_fork.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, fork.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranchName: "develop",
					TargetBranchName: "topic",
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{ForkCommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})

	t.Run("Release", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			releaseUseCase := mock_release.NewMockInterface(ctrl)
			releaseUseCase.EXPECT().
				Do(ctx, release.Input{
					Repository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TagName:    "v1.0.0",
					Paths:      []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{ReleaseUseCase: releaseUseCase}),
			}
			args := []string{
				cmdName,
				releaseCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-t", "v1.0.0",
				"file1",
				"file2",
			}
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})

	t.Run("GlobalOptions", func(t *testing.T) {
		t.Run("--debug", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{Debug: true}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("--directory", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			mockEnv := newEnv(ctrl, map[string]string{envGitHubAPI: ""})
			mockEnv.EXPECT().
				Chdir("dir")
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               mockEnv,
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("env/GITHUB_TOKEN", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubToken: "YOUR_TOKEN", envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("NoGitHubToken", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{}),
				Env:               newEnv(ctrl, map[string]string{envGitHubToken: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeError {
				t.Errorf("exitCode wants %d but %d", exitCodeError, exitCode)
			}
		})

		t.Run("--api", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN", URLv3: "https://github.example.com/api/v3/"}),
				Env:               newEnv(ctrl, nil),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})

		t.Run("env/GITHUB_API", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentBranch:     commit.ParentBranch{FastForward: true},
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitMessage:    "commit-message",
					Paths:            []string{"file1", "file2"},
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN", URLv3: "https://github.example.com/api/v3/"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: "https://github.example.com/api/v3/"}),
				NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
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
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})
}

func newLogger(t *testing.T, want logger.Option) logger.NewFunc {
	return func(got logger.Option) logger.Interface {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		return testingLogger.New(t)
	}
}

func newGitHub(t *testing.T, want github.Option) github.NewFunc {
	return func(got github.Option) (github.Interface, error) {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		return nil, nil
	}
}

func newEnv(ctrl *gomock.Controller, getenv map[string]string) *mock_env.MockInterface {
	env := mock_env.NewMockInterface(ctrl)
	for k, v := range getenv {
		env.EXPECT().Getenv(k).Return(v)
	}
	return env
}

func newInternalRunner(base InternalRunner) NewInternalRunnerFunc {
	return func(l logger.Interface, g github.Interface) *InternalRunner {
		base.Logger = l
		return &base
	}
}
