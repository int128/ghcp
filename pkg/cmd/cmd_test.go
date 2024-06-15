package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/int128/ghcp/pkg/env/mock_env"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/logger"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
	"github.com/int128/ghcp/pkg/usecases/commit"
	"github.com/int128/ghcp/pkg/usecases/commit/mock_commit"
)

const cmdName = "ghcp"
const version = "TEST"

func TestCmd_Run(t *testing.T) {
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
			Do(gomock.Any(), input)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{Debug: true}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
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
			Do(gomock.Any(), input)
		mockEnv := newEnv(ctrl, map[string]string{envGitHubAPI: ""})
		mockEnv.EXPECT().
			Chdir("dir")
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
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
			Do(gomock.Any(), input)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
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
			NewGitHub:         newGitHub(t, client.Option{}),
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
			Do(gomock.Any(), input)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN", URLv3: "https://github.example.com/api/v3/"}),
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
			Do(gomock.Any(), input)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN", URLv3: "https://github.example.com/api/v3/"}),
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
}

func newLogger(t *testing.T, want logger.Option) logger.NewFunc {
	return func(got logger.Option) logger.Interface {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		return testingLogger.New(t)
	}
}

func newGitHub(t *testing.T, want client.Option) client.NewFunc {
	return func(got client.Option) (client.Interface, error) {
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
	return func(l logger.Interface, g client.Interface) *InternalRunner {
		base.Logger = l
		return &base
	}
}
