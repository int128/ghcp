package release

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/fs/mock_fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github/mock_github"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
)

func TestRelease_Do(t *testing.T) {
	ctx := context.TODO()
	targetRepositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	targetTagName := git.TagName("v1.0.0")
	theFiles := []fs.File{
		{Path: "file1"},
		{Path: "dir2/file2", Executable: true},
	}

	t.Run("CreateReleaseIfNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		in := Input{
			Repository:              targetRepositoryID,
			TagName:                 targetTagName,
			TargetBranchOrCommitSHA: "TARGET_COMMIT",
			Paths:                   []string{"path"},
		}
		fileSystem := mock_fs.NewMockInterface(ctrl)
		fileSystem.EXPECT().FindFiles([]string{"path"}, gomock.Any()).Return(theFiles, nil)
		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			GetReleaseByTagOrNil(ctx, targetRepositoryID, targetTagName).
			Return(nil, nil)
		gitHub.EXPECT().
			CreateRelease(ctx, git.Release{
				ID: git.ReleaseID{
					Repository: targetRepositoryID,
				},
				TagName:         targetTagName,
				Name:            targetTagName.Name(),
				TargetCommitish: "TARGET_COMMIT",
			}).
			Return(&git.Release{
				ID: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				TagName:         targetTagName,
				Name:            targetTagName.Name(),
				TargetCommitish: "TARGET_COMMIT",
			}, nil)
		gitHub.EXPECT().
			CreateReleaseAsset(ctx, git.ReleaseAsset{
				Release: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				Name:     "file1",
				RealPath: "file1",
			})
		gitHub.EXPECT().
			CreateReleaseAsset(ctx, git.ReleaseAsset{
				Release: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				Name:     "file2",
				RealPath: "dir2/file2",
			})

		useCase := Release{
			FileSystem: fileSystem,
			Logger:     testingLogger.New(t),
			GitHub:     gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})

	t.Run("ReleaseAlreadyExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		in := Input{
			Repository: targetRepositoryID,
			TagName:    targetTagName,
			Paths:      []string{"path"},
		}
		fileSystem := mock_fs.NewMockInterface(ctrl)
		fileSystem.EXPECT().FindFiles([]string{"path"}, gomock.Any()).Return(theFiles, nil)
		gitHub := mock_github.NewMockInterface(ctrl)
		gitHub.EXPECT().
			GetReleaseByTagOrNil(ctx, targetRepositoryID, targetTagName).
			Return(&git.Release{
				ID: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				TagName: targetTagName,
				Name:    targetTagName.Name(),
			}, nil)
		gitHub.EXPECT().
			CreateReleaseAsset(ctx, git.ReleaseAsset{
				Release: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				Name:     "file1",
				RealPath: "file1",
			})
		gitHub.EXPECT().
			CreateReleaseAsset(ctx, git.ReleaseAsset{
				Release: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				Name:     "file2",
				RealPath: "dir2/file2",
			})

		useCase := Release{
			FileSystem: fileSystem,
			Logger:     testingLogger.New(t),
			GitHub:     gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})
}
