package release

import (
	"context"
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/fs_mock"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/github_mock"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/stretchr/testify/mock"
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
		in := Input{
			Repository:              targetRepositoryID,
			TagName:                 targetTagName,
			TargetBranchOrCommitSHA: "TARGET_COMMIT",
			Paths:                   []string{"path"},
		}
		fileSystem := fs_mock.NewMockInterface(t)
		fileSystem.EXPECT().FindFiles([]string{"path"}, mock.Anything).Return(theFiles, nil)
		gitHub := github_mock.NewMockInterface(t)
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
			}).
			Return(nil)
		gitHub.EXPECT().
			CreateReleaseAsset(ctx, git.ReleaseAsset{
				Release: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				Name:     "file2",
				RealPath: "dir2/file2",
			}).
			Return(nil)

		useCase := Release{
			FileSystem: fileSystem,
			GitHub:     gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})

	t.Run("ReleaseAlreadyExists", func(t *testing.T) {
		in := Input{
			Repository: targetRepositoryID,
			TagName:    targetTagName,
			Paths:      []string{"path"},
		}
		fileSystem := fs_mock.NewMockInterface(t)
		fileSystem.EXPECT().FindFiles([]string{"path"}, mock.Anything).Return(theFiles, nil)
		gitHub := github_mock.NewMockInterface(t)
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
			}).
			Return(nil)
		gitHub.EXPECT().
			CreateReleaseAsset(ctx, git.ReleaseAsset{
				Release: git.ReleaseID{
					Repository: targetRepositoryID,
					InternalID: 1234567890,
				},
				Name:     "file2",
				RealPath: "dir2/file2",
			}).
			Return(nil)

		useCase := Release{
			FileSystem: fileSystem,
			GitHub:     gitHub,
		}
		if err := useCase.Do(ctx, in); err != nil {
			t.Errorf("err wants nil but %+v", err)
		}
	})
}
