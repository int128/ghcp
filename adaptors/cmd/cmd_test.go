package cmd

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/int128/ghcp/adaptors/env/mock_env"
	"github.com/int128/ghcp/adaptors/logger"
	testingLogger "github.com/int128/ghcp/adaptors/logger/testing"
	"github.com/int128/ghcp/domain/git"
	"github.com/int128/ghcp/domain/git/commitstrategy"
	"github.com/int128/ghcp/infrastructure/github"
	"github.com/int128/ghcp/usecases/commit"
	"github.com/int128/ghcp/usecases/commit/mock_commit"
	"github.com/int128/ghcp/usecases/forkcommit"
	"github.com/int128/ghcp/usecases/forkcommit/mock_forkcommit"
	"github.com/int128/ghcp/usecases/pullrequest"
	"github.com/int128/ghcp/usecases/pullrequest/mock_pullrequest"
	"github.com/int128/ghcp/usecases/release"
	"github.com/int128/ghcp/usecases/release/mock_release"
)

func TestCmd_Run(t *testing.T) {
	const cmdName = "ghcp"
	const version = "TEST"
	ctx := context.TODO()

	t.Run("commit", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, commit.Input{
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitStrategy:   commitstrategy.FastForward,
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
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "gh-pages",
					CommitStrategy:   commitstrategy.FastForward,
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
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitStrategy:   commitstrategy.RebaseOn("develop"),
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
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitStrategy:   commitstrategy.NoParent,
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
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitStrategy:   commitstrategy.FastForward,
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
					TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitStrategy:   commitstrategy.FastForward,
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

	t.Run("fork", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_forkcommit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, forkcommit.Input{
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					TargetBranchName: "topic",
					CommitStrategy:   commitstrategy.FastForward,
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
				forkCommitCmdName,
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
			commitUseCase := mock_forkcommit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, forkcommit.Input{
					TargetBranchName: "topic",
					ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					CommitStrategy:   commitstrategy.RebaseOn("develop"),
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
				forkCommitCmdName,
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

	t.Run("pull-request", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			useCase := mock_pullrequest.NewMockInterface(ctrl)
			useCase.EXPECT().
				Do(ctx, pullrequest.Input{
					HeadRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					HeadBranchName: "feature",
					BaseRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					Title:          "commit-message",
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{PullRequestUseCase: useCase}),
			}
			args := []string{
				cmdName,
				pullRequestCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-b", "feature",
				"--title", "commit-message",
			}
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
		t.Run("--base", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			useCase := mock_pullrequest.NewMockInterface(ctrl)
			useCase.EXPECT().
				Do(ctx, pullrequest.Input{
					HeadRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
					HeadBranchName: "feature",
					BaseRepository: git.RepositoryID{Owner: "upstream-owner", Name: "upstream-repo"},
					BaseBranchName: "develop",
					Title:          "commit-message",
					Body:           "body",
				})
			r := Runner{
				NewLogger:         newLogger(t, logger.Option{}),
				NewGitHub:         newGitHub(t, github.Option{Token: "YOUR_TOKEN"}),
				Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
				NewInternalRunner: newInternalRunner(InternalRunner{PullRequestUseCase: useCase}),
			}
			args := []string{
				cmdName,
				pullRequestCmdName,
				"--token", "YOUR_TOKEN",
				"-u", "owner",
				"-r", "repo",
				"-b", "feature",
				"--base-owner", "upstream-owner",
				"--base-repo", "upstream-repo",
				"--base", "develop",
				"--title", "commit-message",
				"--body", "body",
			}
			exitCode := r.Run(args, version)
			if exitCode != exitCodeOK {
				t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
			}
		})
	})

	t.Run("release", func(t *testing.T) {
		t.Run("BasicOptions", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			releaseUseCase := mock_release.NewMockInterface(ctrl)
			releaseUseCase.EXPECT().
				Do(ctx, release.Input{
					Repository:              git.RepositoryID{Owner: "owner", Name: "repo"},
					TagName:                 "v1.0.0",
					TargetBranchOrCommitSHA: "COMMIT_SHA",
					Paths:                   []string{"file1", "file2"},
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
				"--target", "COMMIT_SHA",
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
		input := commit.Input{
			TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
			ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
			CommitStrategy:   commitstrategy.FastForward,
			CommitMessage:    "commit-message",
			Paths:            []string{"file1", "file2"},
		}

		t.Run("--debug", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			commitUseCase := mock_commit.NewMockInterface(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, input)
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
				Do(ctx, input)
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
				Do(ctx, input)
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
				Do(ctx, input)
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
				Do(ctx, input)
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
