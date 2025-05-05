package cmd

import (
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/commit_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/usecases/commit"
	"github.com/stretchr/testify/mock"
)

func TestCmd_Run_empty_commit(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		commitUseCase := commit_mock.NewMockInterface(t)
		commitUseCase.EXPECT().
			Do(mock.Anything, commit.Input{
				TargetRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				ParentRepository: git.RepositoryID{Owner: "owner", Name: "repo"},
				CommitStrategy:   commitstrategy.FastForward,
				CommitMessage:    "commit-message",
			}).
			Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
			NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
		}
		args := []string{
			cmdName,
			emptyCommitCmdName,
			"--token", "YOUR_TOKEN",
			"-r", "owner/repo",
			"-m", "commit-message",
		}
		exitCode := r.Run(args, version)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})
}
