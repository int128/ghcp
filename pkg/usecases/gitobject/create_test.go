package gitobject

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"

	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/fs/mock_fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/github/mock_github"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
)

func TestCreateBlobTreeCommit_Do(t *testing.T) {
	ctx := context.TODO()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("BasicOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := mock_fs.NewMockInterface(ctrl)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file1").
			Return("base64content1", nil)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file2").
			Return("base64content2", nil)

		gitHub := mock_github.NewMockInterface(ctrl)
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
			QueryCommit(ctx, github.QueryCommitInput{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&github.QueryCommitOutput{
				ChangedFiles: 1,
			}, nil)

		useCase := CreateGitObject{
			FileSystem: fileSystem,
			Logger:     testingLogger.New(t),
			GitHub:     gitHub,
		}
		got, err := useCase.Do(ctx, Input{
			Files: []fs.File{
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
		want := &Output{
			CommitSHA:    "commitSHA",
			ChangedFiles: 1,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("NoFileMode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := mock_fs.NewMockInterface(ctrl)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file1").
			Return("base64content1", nil)
		fileSystem.EXPECT().
			ReadAsBase64EncodedContent("file2").
			Return("base64content2", nil)

		gitHub := mock_github.NewMockInterface(ctrl)
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
			QueryCommit(ctx, github.QueryCommitInput{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&github.QueryCommitOutput{
				ChangedFiles: 1,
			}, nil)

		useCase := CreateGitObject{
			FileSystem: fileSystem,
			Logger:     testingLogger.New(t),
			GitHub:     gitHub,
		}
		got, err := useCase.Do(ctx, Input{
			Files: []fs.File{
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
		want := &Output{
			CommitSHA:    "commitSHA",
			ChangedFiles: 1,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("NoFile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := mock_fs.NewMockInterface(ctrl)

		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "masterTreeSHA",
				ParentCommitSHA: "masterCommitSHA",
				Message:         "message",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, github.QueryCommitInput{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&github.QueryCommitOutput{
				ChangedFiles: 1,
			}, nil)

		useCase := CreateGitObject{
			FileSystem: fileSystem,
			Logger:     testingLogger.New(t),
			GitHub:     gitHub,
		}
		got, err := useCase.Do(ctx, Input{
			Files:           nil,
			Repository:      repositoryID,
			CommitMessage:   "message",
			ParentCommitSHA: "masterCommitSHA",
			ParentTreeSHA:   "masterTreeSHA",
		})
		if err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
		want := &Output{
			CommitSHA:    "commitSHA",
			ChangedFiles: 1,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CommitterAndAuthor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := mock_fs.NewMockInterface(ctrl)

		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "masterTreeSHA",
				ParentCommitSHA: "masterCommitSHA",
				Message:         "message",
				Committer:       &git.CommitAuthor{Name: "SomeCommitter", Email: "committer@example.com"},
				Author:          &git.CommitAuthor{Name: "SomeAuthor", Email: "author@example.com"},
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			QueryCommit(ctx, github.QueryCommitInput{
				Repository: repositoryID,
				CommitSHA:  "commitSHA",
			}).
			Return(&github.QueryCommitOutput{
				ChangedFiles: 1,
			}, nil)

		useCase := CreateGitObject{
			FileSystem: fileSystem,
			Logger:     testingLogger.New(t),
			GitHub:     gitHub,
		}
		got, err := useCase.Do(ctx, Input{
			Files:           nil,
			Repository:      repositoryID,
			CommitMessage:   "message",
			Committer:       &git.CommitAuthor{Name: "SomeCommitter", Email: "committer@example.com"},
			Author:          &git.CommitAuthor{Name: "SomeAuthor", Email: "author@example.com"},
			ParentCommitSHA: "masterCommitSHA",
			ParentTreeSHA:   "masterTreeSHA",
		})
		if err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
		want := &Output{
			CommitSHA:    "commitSHA",
			ChangedFiles: 1,
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}
