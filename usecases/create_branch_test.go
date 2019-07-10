package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases/interfaces"
	"github.com/int128/ghcp/usecases/mock_usecases"
)

func TestCreateBranch_Do(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	type testCase struct {
		DryRun            bool
		NoFileMode        bool
		ChangedFiles      int
		CreateBranchTimes int
	}

	run := func(t *testing.T, c testCase) {
		t.Run("FromDefaultBranch", func(t *testing.T) {
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
				QueryForCreateBranch(ctx, adaptors.QueryForCreateBranchIn{
					Repository:    repositoryID,
					NewBranchName: "topic",
				}).
				Return(&adaptors.QueryForCreateBranchOut{
					CurrentUserName: "current",
					Repository:      repositoryID,
					DefaultBranchRefName: git.RefQualifiedName{
						Prefix: "refs/heads/",
						Name:   "master",
					},
					DefaultBranchCommitSHA: "masterCommitSHA",
					DefaultBranchTreeSHA:   "masterTreeSHA",
				}, nil)
			gitHub.EXPECT().
				CreateBranch(ctx, git.NewBranch{
					Repository: repositoryID,
					BranchName: "topic",
					CommitSHA:  "commitSHA",
				}).
				Return(nil).
				Times(c.CreateBranchTimes)

			useCase := CreateBranch{
				Commit:     commitUseCase,
				FileSystem: fileSystem,
				Logger:     mock_adaptors.NewLogger(t),
				GitHub:     gitHub,
			}
			err := useCase.Do(ctx, usecases.CreateBranchIn{
				Repository:        repositoryID,
				NewBranchName:     "topic",
				ParentOfNewBranch: usecases.ParentOfNewBranch{FromDefaultBranch: true},
				CommitMessage:     "message",
				Paths:             []string{"path"},
				NoFileMode:        c.NoFileMode,
				DryRun:            c.DryRun,
			})
			if err != nil {
				t.Errorf("Do returned error: %+v", err)
			}
		})

		t.Run("FromGivenBranch", func(t *testing.T) {
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
					ParentCommitSHA: "developCommitSHA",
					ParentTreeSHA:   "developTreeSHA",
					NoFileMode:      c.NoFileMode,
				}).
				Return(&usecases.CommitOut{
					CommitSHA:    "commitSHA",
					ChangedFiles: c.ChangedFiles,
				}, nil)

			gitHub := mock_adaptors.NewMockGitHub(ctrl)
			gitHub.EXPECT().
				QueryForCreateBranch(ctx, adaptors.QueryForCreateBranchIn{
					Repository:    repositoryID,
					NewBranchName: "topic",
					ParentRef:     "develop",
				}).
				Return(&adaptors.QueryForCreateBranchOut{
					CurrentUserName: "current",
					Repository:      repositoryID,
					DefaultBranchRefName: git.RefQualifiedName{
						Prefix: "refs/heads/",
						Name:   "master",
					},
					DefaultBranchCommitSHA: "masterCommitSHA",
					DefaultBranchTreeSHA:   "masterTreeSHA",
					ParentRefName: git.RefQualifiedName{
						Prefix: "refs/heads/",
						Name:   "develop",
					},
					ParentRefCommitSHA: "developCommitSHA",
					ParentRefTreeSHA:   "developTreeSHA",
				}, nil)
			gitHub.EXPECT().
				CreateBranch(ctx, git.NewBranch{
					Repository: repositoryID,
					BranchName: "topic",
					CommitSHA:  "commitSHA",
				}).
				Return(nil).
				Times(c.CreateBranchTimes)

			useCase := CreateBranch{
				Commit:     commitUseCase,
				FileSystem: fileSystem,
				Logger:     mock_adaptors.NewLogger(t),
				GitHub:     gitHub,
			}
			err := useCase.Do(ctx, usecases.CreateBranchIn{
				Repository:        repositoryID,
				NewBranchName:     "topic",
				ParentOfNewBranch: usecases.ParentOfNewBranch{FromRef: "develop"},
				CommitMessage:     "message",
				Paths:             []string{"path"},
				NoFileMode:        c.NoFileMode,
				DryRun:            c.DryRun,
			})
			if err != nil {
				t.Errorf("Do returned error: %+v", err)
			}
		})

		t.Run("NoParent", func(t *testing.T) {
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
					Repository:    repositoryID,
					CommitMessage: "message",
					NoFileMode:    c.NoFileMode,
				}).
				Return(&usecases.CommitOut{
					CommitSHA:    "commitSHA",
					ChangedFiles: c.ChangedFiles,
				}, nil)

			gitHub := mock_adaptors.NewMockGitHub(ctrl)
			gitHub.EXPECT().
				QueryForCreateBranch(ctx, adaptors.QueryForCreateBranchIn{
					Repository:    repositoryID,
					NewBranchName: "gh-pages",
				}).
				Return(&adaptors.QueryForCreateBranchOut{
					CurrentUserName: "current",
					Repository:      repositoryID,
					DefaultBranchRefName: git.RefQualifiedName{
						Prefix: "refs/heads/",
						Name:   "master",
					},
					DefaultBranchCommitSHA: "masterCommitSHA",
					DefaultBranchTreeSHA:   "masterTreeSHA",
					ParentRefName: git.RefQualifiedName{
						Prefix: "refs/heads/",
						Name:   "gh-pages",
					},
					ParentRefCommitSHA: "ghCommitSHA",
					ParentRefTreeSHA:   "ghTreeSHA",
				}, nil)
			gitHub.EXPECT().
				CreateBranch(ctx, git.NewBranch{
					Repository: repositoryID,
					BranchName: "gh-pages",
					CommitSHA:  "commitSHA",
				}).
				Return(nil).
				Times(c.CreateBranchTimes)

			useCase := CreateBranch{
				Commit:     commitUseCase,
				FileSystem: fileSystem,
				Logger:     mock_adaptors.NewLogger(t),
				GitHub:     gitHub,
			}
			err := useCase.Do(ctx, usecases.CreateBranchIn{
				Repository:        repositoryID,
				NewBranchName:     "gh-pages",
				ParentOfNewBranch: usecases.ParentOfNewBranch{NoParent: true},
				CommitMessage:     "message",
				Paths:             []string{"path"},
				NoFileMode:        c.NoFileMode,
				DryRun:            c.DryRun,
			})
			if err != nil {
				t.Errorf("Do returned error: %+v", err)
			}
		})
	}

	for name, c := range map[string]testCase{
		"Success": {
			ChangedFiles:      1,
			CreateBranchTimes: 1,
		},
		"NothingToCommit": {
			ChangedFiles:      0,
			CreateBranchTimes: 0,
		},
		"DryRun": {
			DryRun:            true,
			CreateBranchTimes: 0,
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
