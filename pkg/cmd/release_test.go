package cmd

import (
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/release_mock"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github/client"
	"github.com/int128/ghcp/pkg/usecases/release"
	"github.com/stretchr/testify/mock"
)

func TestCmd_Run_release(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		releaseUseCase := release_mock.NewMockInterface(t)
		releaseUseCase.EXPECT().
			Do(mock.Anything, release.Input{
				Repository:              git.RepositoryID{Owner: "owner", Name: "repo"},
				TagName:                 "v1.0.0",
				TargetBranchOrCommitSHA: "COMMIT_SHA",
				Paths:                   []string{"file1", "file2"},
			}).
			Return(nil)
		r := Runner{
			NewGitHub:         newGitHub(t, client.Option{Token: "YOUR_TOKEN"}),
			Env:               newEnv(t, map[string]string{envGitHubAPI: ""}),
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
