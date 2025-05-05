package cmd

import (
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/forkcommit_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/usecases/forkcommit"
	"github.com/stretchr/testify/mock"
)

func TestCmd_Run_forkcommit(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		commitUseCase := forkcommit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, forkcommit.Input{
				ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				TargetBranchName: "topic",
				CommitStrategy:   commitstrategy.FastForward,
				CommitMessage:    "commit-message",
				Paths:            []string{"file1", "file2"},
			}).
			Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
			NewInternalRunner: newInternalRunner(InternalRunner{ForkCommitUseCase: commitUseCase}),
		}
		args := []string{
			cmdName,
			forkCommitCmdName,
			"--token", "YOUR_TOKEN",
			"-r", "owner/repo",
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
		commitUseCase := forkcommit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, forkcommit.Input{
				TargetBranchName: "topic",
				ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitStrategy:   commitstrategy.RebaseOn("develop"),
				CommitMessage:    "commit-message",
				Paths:            []string{"file1", "file2"},
			}).
			Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
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
}
