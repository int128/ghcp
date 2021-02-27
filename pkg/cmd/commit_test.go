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

func TestCmd_Run_commit(t *testing.T) {
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
				Paths:            []string{"file1", "file2"},
			})
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
			NewInternalRunner: newInternalRunner(InternalRunner{CommitUseCase: commitUseCase}),
		}
		args := []string{
			cmdName,
			commitCmdName,
			"--token", "YOUR_TOKEN",
			"-r", "owner/repo",
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
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
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

	t.Run("only --author-name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
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
			"--author-name", "Some Author",
			"file1",
		}
		exitCode := r.Run(args, version)
		if exitCode != exitCodeError {
			t.Errorf("exitCode wants %d but %d", exitCodeError, exitCode)
		}
	})

	t.Run("only --committer-email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
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
			"--committer-email", "committer@example.com",
			"file1",
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
			"--dry-run",
			"file1",
			"file2",
		}
		exitCode := r.Run(args, version)
		if exitCode != exitCodeOK {
			t.Errorf("exitCode wants %d but %d", exitCodeOK, exitCode)
		}
	})
}
