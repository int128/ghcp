package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/go-github/v74/github"
	"github.com/int128/ghcp/pkg/git"
)

// GetReleaseByTagOrNil returns the release associated to the tag.
// It returns nil if it does not exist.
func (c *GitHub) GetReleaseByTagOrNil(ctx context.Context, repo git.RepositoryID, tag git.TagName) (*git.Release, error) {
	slog.Debug("Getting the release associated to the tag", "tag", tag, "repository", repo)
	release, resp, err := c.Client.GetReleaseByTag(ctx, repo.Owner, repo.Name, tag.Name())
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		slog.Debug("GitHub API returned 404", "tag", tag, "repository", repo, "error", err)
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	return &git.Release{
		ID: git.ReleaseID{
			Repository: repo,
			InternalID: release.GetID(),
		},
		TagName: git.TagName(release.GetTagName()),
		Name:    release.GetName(),
	}, nil
}

// CreateRelease creates a release.
func (c *GitHub) CreateRelease(ctx context.Context, r git.Release) (*git.Release, error) {
	slog.Debug("Creating a release", "release", r)
	release, _, err := c.Client.CreateRelease(ctx, r.ID.Repository.Owner, r.ID.Repository.Name, &github.RepositoryRelease{
		Name:            github.Ptr(r.Name),
		TagName:         github.Ptr(r.TagName.Name()),
		TargetCommitish: github.Ptr(r.TargetCommitish),
	})
	if err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	return &git.Release{
		ID: git.ReleaseID{
			Repository: r.ID.Repository,
			InternalID: release.GetID(),
		},
		TagName:         git.TagName(release.GetTagName()),
		TargetCommitish: release.GetTargetCommitish(),
		Name:            release.GetName(),
	}, nil
}

// CreateRelease creates a release asset.
func (c *GitHub) CreateReleaseAsset(ctx context.Context, a git.ReleaseAsset) error {
	slog.Debug("Creating a release asset", "asset", a)
	f, err := os.Open(a.RealPath)
	if err != nil {
		return fmt.Errorf("could not open the file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Failed to close the file", "error", err)
		}
	}()
	_, _, err = c.Client.UploadReleaseAsset(ctx, a.Release.Repository.Owner, a.Release.Repository.Name, a.Release.InternalID, &github.UploadOptions{
		Name: a.Name,
	}, f)
	if err != nil {
		return fmt.Errorf("GitHub API error: %w", err)
	}
	return nil
}
