package commit

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/adaptors/mock_adaptors"
	"github.com/int128/ghcp/git"
	"github.com/int128/ghcp/usecases"
)

func TestCommit_Do(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("BasicOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file1").
			Return("base64content1", nil)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file2").
			Return("base64content2", nil)

		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			CreateBlob(ctx, git.NewBlob{
				Repository: repositoryID,
				Content:    "base64content1",
			}).
			Return(git.BlobSHA("blobSHA1"), nil)
		gitHub.EXPECT().
			CreateBlob(ctx, git.NewBlob{
				Repository: repositoryID,
				Content:    "base64content2",
			}).
			Return(git.BlobSHA("blobSHA2"), nil)
		gitHub.EXPECT().
			CreateTree(ctx, git.NewTree{
				Repository:  repositoryID,
				BaseTreeSHA: "masterTreeSHA",
				Files: []git.File{
					{
						Filename: "file1",
						BlobSHA:  "blobSHA1",
					}, {
						Filename:   "file2",
						BlobSHA:    "blobSHA2",
						Executable: true,
					},
				},
			}).
			Return(git.TreeSHA("treeSHA"), nil)
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

		useCase := Commit{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		out, err := useCase.Do(ctx, usecases.CommitIn{
			Files: []adaptors.File{
				{Path: "file1"},
				{Path: "file2", Executable: true},
			},
			Repository:      repositoryID,
			CommitMessage:   "message",
			ParentCommitSHA: "masterCommitSHA",
			ParentTreeSHA:   "masterTreeSHA",
		})
		if err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
		want := &usecases.CommitOut{
			CommitSHA:    "commitSHA",
			ChangedFiles: 1,
		}
		if diff := deep.Equal(want, out); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("NoFileMode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := mock_adaptors.NewMockFileSystem(ctrl)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file1").
			Return("base64content1", nil)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file2").
			Return("base64content2", nil)

		gitHub := mock_adaptors.NewMockGitHub(ctrl)
		gitHub.EXPECT().
			CreateBlob(ctx, git.NewBlob{
				Repository: repositoryID,
				Content:    "base64content1",
			}).
			Return(git.BlobSHA("blobSHA1"), nil)
		gitHub.EXPECT().
			CreateBlob(ctx, git.NewBlob{
				Repository: repositoryID,
				Content:    "base64content2",
			}).
			Return(git.BlobSHA("blobSHA2"), nil)
		gitHub.EXPECT().
			CreateTree(ctx, git.NewTree{
				Repository:  repositoryID,
				BaseTreeSHA: "masterTreeSHA",
				Files: []git.File{
					{
						Filename: "file1",
						BlobSHA:  "blobSHA1",
					}, {
						Filename: "file2",
						BlobSHA:  "blobSHA2",
						// no Executable
					},
				},
			}).
			Return(git.TreeSHA("treeSHA"), nil)
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

		useCase := Commit{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		out, err := useCase.Do(ctx, usecases.CommitIn{
			Files: []adaptors.File{
				{Path: "file1"},
				{Path: "file2", Executable: true},
			},
			Repository:      repositoryID,
			CommitMessage:   "message",
			ParentCommitSHA: "masterCommitSHA",
			ParentTreeSHA:   "masterTreeSHA",
			NoFileMode:      true,
		})
		if err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
		want := &usecases.CommitOut{
			CommitSHA:    "commitSHA",
			ChangedFiles: 1,
		}
		if diff := deep.Equal(want, out); diff != nil {
			t.Error(diff)
		}
	})
}
