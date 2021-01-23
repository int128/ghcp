package release

import (
	"context"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/logger"
	"golang.org/x/xerrors"
)

var Set = wire.NewSet(
	wire.Struct(new(Release), "*"),
	wire.Bind(new(Interface), new(*Release)),
)

//go:generate mockgen -destination mock_release/mock_release.go github.com/int128/ghcp/usecases/release Interface

type Interface interface {
	Do(ctx context.Context, in Input) error
}

type Input struct {
	Repository              git.RepositoryID
	TagName                 git.TagName
	TargetBranchOrCommitSHA string // optional
	Paths                   []string
	DryRun                  bool
}

// Release create a release with the files to the tag in the repository.
type Release struct {
	FileSystem fs.Interface
	Logger     logger.Interface
	GitHub     github.Interface
}

func (u *Release) Do(ctx context.Context, in Input) error {
	if !in.Repository.IsValid() {
		return xerrors.New("you must set GitHub repository")
	}
	if in.TagName == "" {
		return xerrors.New("you must set the tag name")
	}
	if len(in.Paths) == 0 {
		return xerrors.New("you must set one or more paths")
	}

	files, err := u.FileSystem.FindFiles(in.Paths, nil)
	if err != nil {
		return xerrors.Errorf("could not find files: %w", err)
	}
	if len(files) == 0 {
		return xerrors.New("no file exists in given paths")
	}

	release, err := u.GitHub.GetReleaseByTagOrNil(ctx, in.Repository, in.TagName)
	if err != nil {
		return xerrors.Errorf("could not get the release: %w", err)
	}
	if release == nil {
		u.Logger.Infof("No release on the tag %s", in.TagName)
		if in.DryRun {
			u.Logger.Infof("Do not create a release due to dry-run")
			return nil
		}
		release, err = u.GitHub.CreateRelease(ctx, git.Release{
			ID:              git.ReleaseID{Repository: in.Repository},
			Name:            in.TagName.Name(),
			TagName:         in.TagName,
			TargetCommitish: in.TargetBranchOrCommitSHA,
		})
		if err != nil {
			return xerrors.Errorf("could not create a release: %w", err)
		}
		u.Logger.Infof("Created a release %s", release.Name)
	} else {
		u.Logger.Infof("Found the release on the tag %s", in.TagName)
	}

	if in.DryRun {
		u.Logger.Infof("Do not upload files to the release %s due to dry-run", release.Name)
		return nil
	}
	u.Logger.Infof("Uploading %d file(s)", len(files))
	for _, file := range files {
		if err := u.GitHub.CreateReleaseAsset(ctx, git.ReleaseAsset{
			Release:  release.ID,
			Name:     filepath.Base(file.Path),
			RealPath: file.Path,
		}); err != nil {
			return xerrors.Errorf("could not create a release asset: %w", err)
		}
		u.Logger.Infof("Uploaded %s", file.Path)
	}
	return nil
}
