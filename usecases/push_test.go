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

func TestPush_Do(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	t.Run("Files", func(t *testing.T) {
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
				BaseTreeSHA: "masterTreeSHA",
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
		gitHub.EXPECT().
			CreateCommit(ctx, git.NewCommit{
				Repository:      repositoryID,
				TreeSHA:         "treeSHA",
				ParentCommitSHA: "masterCommitSHA",
			}).
			Return(git.CommitSHA("commitSHA"), nil)
		gitHub.EXPECT().
			UpdateBranch(ctx, git.NewBranch{
				Repository: repositoryID,
				BranchName: "master",
				CommitSHA:  "commitSHA",
			}, false).
			Return(nil)

		push := Push{
			FileSystem: fileSystem,
			Logger:     mock_adaptors.NewLogger(t),
			GitHub:     gitHub,
		}
		err := push.Do(ctx, usecases.PushIn{
			Repository: repositoryID,
			Paths:      []string{"path"},
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})
}
