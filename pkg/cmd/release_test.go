package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/logger"
	"github.com/int128/ghcp/pkg/usecases/release"
	"github.com/int128/ghcp/pkg/usecases/release/mock_release"
)

func TestCmd_Run_release(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		releaseUseCase := mock_release.NewMockInterface(ctrl)
		releaseUseCase.EXPECT().
			Do(gomock.Any(), release.Input{
				Repository:              git.RepositoryID{Owner: "owner", Name: "repo"},
				TagName:                 "v1.0.0",
				TargetBranchOrCommitSHA: "COMMIT_SHA",
				Paths:                   []string{"file1", "file2"},
			})
		r := Runner{
			NewLogger:         newLogger(t, logger.Option{}),
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(ctrl, map[string]string{envGitHubAPI: ""}),
			NewInternalRunner: newInternalRunner(InternalRunner{ReleaseUseCase: releaseUseCase}),
		}
		args := []string{
			cmdName,
			releaseCmdName,
			"--token", "YOUR_TOKEN",
			"-r", "owner/repo",
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
}
