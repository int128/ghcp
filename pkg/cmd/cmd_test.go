package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/env_mock"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/commit_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/usecases/commit"
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
		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, input).Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
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
		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, input).Return(nil)
		mockEnv := newEnv(t, map[string]string{envGitHubAPI: ""})
		mockEnv.EXPECT().
			Chdir("dir").Return(nil)
		r := Runner{
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
		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, input).Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubToken: "YOUR_TOKEN", envGitHubAPI: ""}),
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
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{}),
			Env:               newEnv(t, map[string]string{envGitHubToken: ""}),
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
		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, input).Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN", URLv3: "https://github.example.com/api/v3/"}),
			Env:               newEnv(t, nil),
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
		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, input).Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN", URLv3: "https://github.example.com/api/v3/"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: "https://github.example.com/api/v3/"}),
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

func newGitHub(t *testing.T, want client.Option) client.NewFunc {
	return func(got client.Option) (client.Interface, error) {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
		return nil, nil
	}
}

func newEnv(t *testing.T, getenv map[string]string) *env_mock.MockInterface {
	env := env_mock.NewMockInterface(t)
	for k, v := range getenv {
		env.EXPECT().Getenv(k).Return(v)
	}
	return env
}

func newInternalRunner(base InternalRunner) NewInternalRunnerFunc {
	return func(g client.Interface) *InternalRunner {
		return &base
	}
}
