package cmd

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/int128/ghcp/pkg/usecases/commit"
	"github.com/int128/ghcp/pkg/usecases/commit/mock_commit"
)

func TestCmd_Run_empty_commit(t *testing.T) {
	ctx := context.TODO()

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
			})
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
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
