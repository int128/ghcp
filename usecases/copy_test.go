package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases/interfaces"
)

func newMockFileSystem(ctrl *gomock.Controller) adaptors.FileSystem {
	fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
	fileSystem.EXPECT().
		FindFiles([]string{"path"}).
		Return([]string{"file1", "file2"}, nil)
	fileSystem.EXPECT().
		ReadAsBase64EncodedContent("file1").
		Return("base64content1", nil)
	fileSystem.EXPECT().
		ReadAsBase64EncodedContent("file2").
		Return("base64content2", nil)
	return fileSystem
}

func gitHubCreateBlobTree(gitHub *mock_adaptors.MockGitHub, ctx context.Context, repositoryID git.RepositoryID, baseTreeSHA git.TreeSHA) {
	gitHub.EXPECT().
		CreateBlob(ctx, git.NewBlob{
			Repository: repositoryID,
			Content:    "base64content1",
		}).
		Return(git.BlobSHA("blobSHA"), nil)
	gitHub.EXPECT().
		CreateBlob(ctx, git.NewBlob{
			Repository: repositoryID,
			Content:    "base64content2",
		}).
		Return(git.BlobSHA("blobSHA"), nil)
	gitHub.EXPECT().
		CreateTree(ctx, git.NewTree{
			Repository:  repositoryID,
			BaseTreeSHA: baseTreeSHA,
			Files: []git.File{
				{
					Filename: "file1",
					BlobSHA:  "blobSHA",
				}, {
					Filename: "file2",
					BlobSHA:  "blobSHA",
				},
			},
		}).
		Return(git.TreeSHA("treeSHA"), nil)
}

func TestCopyUseCase_Do(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("DefaultBranch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := newMockFileSystem(ctrl)
		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			QueryRepository(ctx, adaptors.QueryRepositoryIn{
				Repository: repositoryID,
			}).
			Return(&adaptors.QueryRepositoryOut{
				CurrentUserName:        "current",
				Repository:             repositoryID,
				DefaultBranchName:      "master",
				DefaultBranchCommitSHA: "masterCommitSHA",
				DefaultBranchTreeSHA:   "masterTreeSHA",
			}, nil)
		gitHubCreateBlobTree(gitHub, ctx, repositoryID, "masterTreeSHA")
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "treeSHA",
				ParentCommitSHA: "masterCommitSHA",
				Message:         "message",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, adaptors.QueryCommitIn{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&adaptors.QueryCommitOut{
				ChangedFiles: 1,
			}, nil)
		gitHub.EXPECT().
			UpdateBranch(ctx, git.NewBranch{
				Repository: repositoryID,
				BranchName: "master",
				CommitSHA:  "commitSHA",
			}, false).
			Return(nil)

		useCase := CopyUseCase{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		err := useCase.Do(ctx, usecases.CopyUseCaseIn{
			Repository:    repositoryID,
			CommitMessage: "message",
			Paths:         []string{"path"},
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("DefaultBranch/NothingToCommit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := newMockFileSystem(ctrl)
		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			QueryRepository(ctx, adaptors.QueryRepositoryIn{
				Repository: repositoryID,
			}).
			Return(&adaptors.QueryRepositoryOut{
				CurrentUserName:        "current",
				Repository:             repositoryID,
				DefaultBranchName:      "master",
				DefaultBranchCommitSHA: "masterCommitSHA",
				DefaultBranchTreeSHA:   "masterTreeSHA",
			}, nil)
		gitHubCreateBlobTree(gitHub, ctx, repositoryID, "masterTreeSHA")
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "treeSHA",
				ParentCommitSHA: "masterCommitSHA",
				Message:         "message",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, adaptors.QueryCommitIn{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&adaptors.QueryCommitOut{
				ChangedFiles: 0,
			}, nil)

		useCase := CopyUseCase{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		err := useCase.Do(ctx, usecases.CopyUseCaseIn{
			Repository:    repositoryID,
			CommitMessage: "message",
			Paths:         []string{"path"},
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("DefaultBranch/DryRun", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := newMockFileSystem(ctrl)
		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			QueryRepository(ctx, adaptors.QueryRepositoryIn{
				Repository: repositoryID,
			}).
			Return(&adaptors.QueryRepositoryOut{
				CurrentUserName:        "current",
				Repository:             repositoryID,
				DefaultBranchName:      "master",
				DefaultBranchCommitSHA: "masterCommitSHA",
				DefaultBranchTreeSHA:   "masterTreeSHA",
			}, nil)
		gitHubCreateBlobTree(gitHub, ctx, repositoryID, "masterTreeSHA")
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "treeSHA",
				ParentCommitSHA: "masterCommitSHA",
				Message:         "message",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, adaptors.QueryCommitIn{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&adaptors.QueryCommitOut{
				ChangedFiles: 1,
			}, nil)

		useCase := CopyUseCase{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		err := useCase.Do(ctx, usecases.CopyUseCaseIn{
			Repository:    repositoryID,
			CommitMessage: "message",
			Paths:         []string{"path"},
			DryRun:        true,
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("GivenBranch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := newMockFileSystem(ctrl)
		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			QueryRepository(ctx, adaptors.QueryRepositoryIn{
				Repository: repositoryID,
				BranchName: "gh-pages",
			}).
			Return(&adaptors.QueryRepositoryOut{
				CurrentUserName:        "current",
				Repository:             repositoryID,
				DefaultBranchName:      "master",
				DefaultBranchCommitSHA: "masterCommitSHA",
				DefaultBranchTreeSHA:   "masterTreeSHA",
				BranchCommitSHA:        "ghCommitSHA",
				BranchTreeSHA:          "ghTreeSHA",
			}, nil)
		gitHubCreateBlobTree(gitHub, ctx, repositoryID, "ghTreeSHA")
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "treeSHA",
				ParentCommitSHA: "ghCommitSHA",
				Message:         "message",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, adaptors.QueryCommitIn{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&adaptors.QueryCommitOut{
				ChangedFiles: 1,
			}, nil)
		gitHub.EXPECT().
			UpdateBranch(ctx, git.NewBranch{
				Repository: repositoryID,
				BranchName: "gh-pages",
				CommitSHA:  "commitSHA",
			}, false).
			Return(nil)

		useCase := CopyUseCase{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		err := useCase.Do(ctx, usecases.CopyUseCaseIn{
			Repository:    repositoryID,
			BranchName:    "gh-pages",
			CommitMessage: "message",
			Paths:         []string{"path"},
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("GivenBranch/DryRun", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := newMockFileSystem(ctrl)
		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			QueryRepository(ctx, adaptors.QueryRepositoryIn{
				Repository: repositoryID,
				BranchName: "gh-pages",
			}).
			Return(&adaptors.QueryRepositoryOut{
				CurrentUserName:        "current",
				Repository:             repositoryID,
				DefaultBranchName:      "master",
				DefaultBranchCommitSHA: "masterCommitSHA",
				DefaultBranchTreeSHA:   "masterTreeSHA",
				BranchCommitSHA:        "ghCommitSHA",
				BranchTreeSHA:          "ghTreeSHA",
			}, nil)
		gitHubCreateBlobTree(gitHub, ctx, repositoryID, "ghTreeSHA")
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "treeSHA",
				ParentCommitSHA: "ghCommitSHA",
				Message:         "message",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, adaptors.QueryCommitIn{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&adaptors.QueryCommitOut{
				ChangedFiles: 1,
			}, nil)

		useCase := CopyUseCase{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		err := useCase.Do(ctx, usecases.CopyUseCaseIn{
			Repository:    repositoryID,
			BranchName:    "gh-pages",
			CommitMessage: "message",
			Paths:         []string{"path"},
			DryRun:        true,
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})
}
