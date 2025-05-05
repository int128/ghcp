package release

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
)

var Set = wire.NewSet(
	wire.Struct(new(Release), "*"),
	wire.Bind(new(Interface), new(*Release)),
)

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
	GitHub     github.Interface
}

func (u *Release) Do(ctx context.Context, in Input) error {
	if !in.Repository.IsValid() {
		return errors.New("you must set GitHub repository")
	}
	if in.TagName == "" {
		return errors.New("you must set the tag name")
	}
	if len(in.Paths) == 0 {
		return errors.New("you must set one or more paths")
	}

	files, err := u.FileSystem.FindFiles(in.Paths, nil)
	if err != nil {
		return fmt.Errorf("could not find files: %w", err)
	}
	if len(files) == 0 {
		return errors.New("no file exists in given paths")
	}

	release, err := u.GitHub.GetReleaseByTagOrNil(ctx, in.Repository, in.TagName)
	if err != nil {
		return fmt.Errorf("could not get the release: %w", err)
	}
	if release == nil {
		slog.Info("No release on the tag", "tag", in.TagName)
		if in.DryRun {
			slog.Info("Do not create a release due to dry-run")
			return nil
		}
		release, err = u.GitHub.CreateRelease(ctx, git.Release{
			ID:              git.ReleaseID{Repository: in.Repository},
			Name:            in.TagName.Name(),
			TagName:         in.TagName,
			TargetCommitish: in.TargetBranchOrCommitSHA,
		})
		if err != nil {
			return fmt.Errorf("could not create a release: %w", err)
		}
		slog.Info("Created a release", "release", release.Name)
	} else {
		slog.Info("Found the release on the tag", "tag", in.TagName)
	}

	if in.DryRun {
		slog.Info("Do not upload files to the release due to dry-run", "release", release.Name)
		return nil
	}
	slog.Info("Uploading", "files", len(files))
	for _, file := range files {
		if err := u.GitHub.CreateReleaseAsset(ctx, git.ReleaseAsset{
			Release:  release.ID,
			Name:     filepath.Base(file.Path),
			RealPath: file.Path,
		}); err != nil {
			return fmt.Errorf("could not create a release asset: %w", err)
		}
		slog.Info("Uploaded", "file", file.Path)
	}
	return nil
}
