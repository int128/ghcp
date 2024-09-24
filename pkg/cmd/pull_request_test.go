package cmd

import (
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/pullrequest_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/int128/ghcp/pkg/usecases/pullrequest"
	"github.com/stretchr/testify/mock"
)

func TestCmd_Run_pull_request(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		useCase := pullrequest_mock.NewMockInterface(t)
		useCase.EXPECT().
			Do(mock.Anything, pullrequest.Input{
				HeadRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				HeadBranchName: "feature",
				BaseRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				Title:          "commit-message",
			}).
			Return(nil)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
			NewInternalRunner: newInternalRunner(InternalRunner{PullRequestUseCase: useCase}),
		}
		args := []string{
			cmdName,
			pullRequestCmdName,
			"--token", "YOUR_TOKEN",
			"-r", "owner/repo",
			"-b", "feature",
			"--title", "commit-message",
		}
		exitCode := r.Run(args, version)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("--base-repo", func(t *testing.T) {
		useCase := pullrequest_mock.NewMockInterface(t)
		useCase.EXPECT().
			Do(mock.Anything, pullrequest.Input{
				HeadRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				HeadBranchName: "feature",
				BaseRepository: git.RepositoryID{Owner: "upstream-owner", Name: "upstream-repo"},
				Title:          "commit-message",
			}).
			Return(nil)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
			NewInternalRunner: newInternalRunner(InternalRunner{PullRequestUseCase: useCase}),
		}
		args := []string{
			cmdName,
			pullRequestCmdName,
			"--token", "YOUR_TOKEN",
			"-r", "owner/repo",
			"--base-repo", "upstream-owner/upstream-repo",
			"-b", "feature",
			"--title", "commit-message",
		}
		exitCode := r.Run(args, version)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})

	t.Run("optional-flags", func(t *testing.T) {
		useCase := pullrequest_mock.NewMockInterface(t)
		useCase.EXPECT().
			Do(mock.Anything, pullrequest.Input{
				HeadRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				HeadBranchName: "feature",
				BaseRepository: git.RepositoryID{Owner: "upstream-owner", Name: "upstream-repo"},
				BaseBranchName: "develop",
				Title:          "commit-message",
				Body:           "body",
				Reviewer:       "the-reviewer",
				Draft:          true,
			}).
			Return(nil)
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
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
			"--draft",
			"--reviewer", "the-reviewer",
		}
		exitCode := r.Run(args, version)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})
}
