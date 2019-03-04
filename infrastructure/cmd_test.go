package infrastructure_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/di"
	"github.com/int128/ghcp/di/mock_di"
	"github.com/int128/ghcp/infrastructure"
)

func TestRun(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("FullOptions", func(t *testing.T) {
		cmd := mock_adaptors.NewMockCmd(ctrl)
		cmd.EXPECT().
			Run(ctx, adaptors.CmdOptions{
				RepositoryOwner: "owner",
				RepositoryName:  "repo",
				CommitMessage:   "commit message",
				Paths:           []string{"file1", "file2"},
			}).
			Return(nil)

		args := []string{
			"ghcp",
			"-token", "YOUR_TOKEN",
			"-u", "owner",
			"-r", "repo",
			"-m", "commit message",
			"file1",
			"file2",
		}
		exitCode := infrastructure.Run(ctx, newContainer(t, ctrl, cmd), args)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})
	t.Run("NoGitHubToken", func(t *testing.T) {
		args := []string{"ghcp"}
		container := mock_di.NewMockContainer(ctrl)
		exitCode := infrastructure.Run(ctx, container, args)
		if exitCode != 1 {
			t.Errorf("exitCode wants 1 but %d", exitCode)
		}
	})
}

func newContainer(t *testing.T, ctrl *gomock.Controller, cmd adaptors.Cmd) di.Container {
	t.Helper()
	container := mock_di.NewMockContainer(ctrl)
	container.EXPECT().
		Run(gomock.Any(), gomock.Any()).
		DoAndReturn(func(d di.ExtraDependencies, f func(adaptors.Cmd) error) error {
			if d.GitHubV3 == nil {
				t.Errorf("GitHubV3 wants non-nil but nil")
			}
			if d.GitHubV4 == nil {
				t.Errorf("GitHubV4 wants non-nil but nil")
			}
			return f(cmd)
		})
	return container
}
