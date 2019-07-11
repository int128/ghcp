package branch

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestUpdateBranch_Do(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	type testCase struct {
		DryRun            bool
		NoFileMode        bool
		ChangedFiles      int
		UpdateBranchTimes int
	}

	run := func(t *testing.T, c testCase) {
		t.Run("DefaultBranch", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
			fileSystem.EXPECT().
				FindFiles([]string{"path"}).
				Return([]adaptors.File{
					{Path: "file1"},
					{Path: "file2", Executable: true},
				}, nil)

			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					Files: []adaptors.File{
						{Path: "file1"},
						{Path: "file2", Executable: true},
					},
					Repository:      repositoryID,
					CommitMessage:   "message",
					ParentCommitSHA: "masterCommitSHA",
					ParentTreeSHA:   "masterTreeSHA",
					NoFileMode:      c.NoFileMode,
				}).
				Return(&usecases.CommitOut{
					CommitSHA:    "commitSHA",
					ChangedFiles: c.ChangedFiles,
				}, nil)

			gitHub := mock_adaptors.NewMockGitHub(ctrl)
			gitHub.EXPECT().
				QueryForUpdateBranch(ctx, adaptors.QueryForUpdateBranchIn{
					Repository: repositoryID,
				}).
				Return(&adaptors.QueryForUpdateBranchOut{
					CurrentUserName:        "current",
					Repository:             repositoryID,
					DefaultBranchName:      "master",
					DefaultBranchCommitSHA: "masterCommitSHA",
					DefaultBranchTreeSHA:   "masterTreeSHA",
				}, nil)
			gitHub.EXPECT().
				UpdateBranch(ctx, git.NewBranch{
					Repository: repositoryID,
					BranchName: "master",
					CommitSHA:  "commitSHA",
				}, false).
				Return(nil).
				Times(c.UpdateBranchTimes)

			useCase := UpdateBranch{
				Commit:     commitUseCase,
				FileSystem: fileSystem,
				Logger:     mock_adaptors.NewLogger(t),
				GitHub:     gitHub,
			}
			err := useCase.Do(ctx, usecases.UpdateBranchIn{
				Repository:    repositoryID,
				CommitMessage: "message",
				Paths:         []string{"path"},
				NoFileMode:    c.NoFileMode,
				DryRun:        c.DryRun,
			})
			if err != nil {
				t.Errorf("Do returned error: %+v", err)
			}
		})

		t.Run("GivenBranch", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
			fileSystem.EXPECT().
				FindFiles([]string{"path"}).
				Return([]adaptors.File{
					{Path: "file1"},
					{Path: "file2", Executable: true},
				}, nil)

			commitUseCase := mock_usecases.NewMockCommit(ctrl)
			commitUseCase.EXPECT().
				Do(ctx, usecases.CommitIn{
					Files: []adaptors.File{
						{Path: "file1"},
						{Path: "file2", Executable: true},
					},
					Repository:      repositoryID,
					CommitMessage:   "message",
					ParentCommitSHA: "ghCommitSHA",
					ParentTreeSHA:   "ghTreeSHA",
					NoFileMode:      c.NoFileMode,
				}).
				Return(&usecases.CommitOut{
					CommitSHA:    "commitSHA",
					ChangedFiles: c.ChangedFiles,
				}, nil)

			gitHub := mock_adaptors.NewMockGitHub(ctrl)
			gitHub.EXPECT().
				QueryForUpdateBranch(ctx, adaptors.QueryForUpdateBranchIn{
					Repository: repositoryID,
					BranchName: "gh-pages",
				}).
				Return(&adaptors.QueryForUpdateBranchOut{
					CurrentUserName:        "current",
					Repository:             repositoryID,
					DefaultBranchName:      "master",
					DefaultBranchCommitSHA: "masterCommitSHA",
					DefaultBranchTreeSHA:   "masterTreeSHA",
					BranchCommitSHA:        "ghCommitSHA",
					BranchTreeSHA:          "ghTreeSHA",
				}, nil)
			gitHub.EXPECT().
				UpdateBranch(ctx, git.NewBranch{
					Repository: repositoryID,
					BranchName: "gh-pages",
					CommitSHA:  "commitSHA",
				}, false).
				Return(nil).
				Times(c.UpdateBranchTimes)

			useCase := UpdateBranch{
				Commit:     commitUseCase,
				FileSystem: fileSystem,
				Logger:     mock_adaptors.NewLogger(t),
				GitHub:     gitHub,
			}
			err := useCase.Do(ctx, usecases.UpdateBranchIn{
				Repository:    repositoryID,
				BranchName:    "gh-pages",
				CommitMessage: "message",
				Paths:         []string{"path"},
				NoFileMode:    c.NoFileMode,
				DryRun:        c.DryRun,
			})
			if err != nil {
				t.Errorf("Do returned error: %+v", err)
			}
		})
	}

	for name, c := range map[string]testCase{
		"Success": {
			ChangedFiles:      1,
			UpdateBranchTimes: 1,
		},
		"NothingToCommit": {
			ChangedFiles:      0,
			UpdateBranchTimes: 0,
		},
		"DryRun": {
			DryRun:            true,
			UpdateBranchTimes: 0,
		},
		"NoFileMode": {
			NoFileMode: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			run(t, c)
		})
	}
}
