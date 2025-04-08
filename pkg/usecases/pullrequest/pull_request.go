package pullrequest

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/wire"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/logger"
)

var Set = wire.NewSet(
	wire.Bind(new(Interface), new(*PullRequest)),
	wire.Struct(new(PullRequest), "*"),
)

type Interface interface {
	Do(ctx context.Context, in Input) error
}

type Input struct {
	BaseRepository git.RepositoryID
	BaseBranchName git.BranchName // if empty, use the default branch of base
	HeadRepository git.RepositoryID
	HeadBranchName git.BranchName // if empty, use the default branch of head
	Title          string
	Body           string // optional
	Reviewer       string // optional
	Draft          bool
}

// PullRequest provides the use-case to create a pull request.
type PullRequest struct {
	GitHub github.Interface
	Logger logger.Interface
}

func (u *PullRequest) Do(ctx context.Context, in Input) error {
	if !in.BaseRepository.IsValid() {
		return errors.New("you must set the base repository")
	}
	if !in.HeadRepository.IsValid() {
		return errors.New("you must set the head repository")
	}

	if in.HeadBranchName == "" || in.BaseBranchName == "" {
		q, err := u.GitHub.QueryDefaultBranch(ctx, github.QueryDefaultBranchInput{
			BaseRepository: in.BaseRepository,
			HeadRepository: in.HeadRepository,
		})
		if err != nil {
			return fmt.Errorf("could not determine the default branch: %w", err)
		}
		if in.BaseBranchName == "" {
			in.BaseBranchName = q.BaseDefaultBranchName
		}
		if in.HeadBranchName == "" {
			in.HeadBranchName = q.HeadDefaultBranchName
		}
	}

	q, err := u.GitHub.QueryForPullRequest(ctx, github.QueryForPullRequestInput{
		BaseRepository: in.BaseRepository,
		BaseBranchName: in.BaseBranchName,
		HeadRepository: in.HeadRepository,
		HeadBranchName: in.HeadBranchName,
	})
	if err != nil {
		return fmt.Errorf("could not query for creating a pull request: %w", err)
	}
	u.Logger.Infof("Logged in as %s", q.CurrentUserName)
	if q.HeadBranchCommitSHA == "" {
		return fmt.Errorf("the head branch (%s) does not exist", in.HeadBranchName)
	}
	u.Logger.Debugf("Found the head branch (%s) with the commit %s", in.HeadBranchName, q.HeadBranchCommitSHA)
	if q.PullRequestURL != "" {
		u.Logger.Warnf("A pull request already exists: %s", q.PullRequestURL)
		return nil
	}

	createdPR, err := u.GitHub.CreatePullRequest(ctx, github.CreatePullRequestInput{
		BaseRepository:       in.BaseRepository,
		BaseBranchName:       in.BaseBranchName,
		BaseRepositoryNodeID: q.BaseRepositoryNodeID,
		HeadRepository:       in.HeadRepository,
		HeadBranchName:       in.HeadBranchName,
		Title:                in.Title,
		Body:                 in.Body,
		Draft:                in.Draft,
	})
	if err != nil {
		return fmt.Errorf("could not create a pull request: %w", err)
	}
	u.Logger.Infof("Created a pull request: %s", createdPR.URL)

	if in.Reviewer == "" {
		return nil
	}
	u.Logger.Infof("Requesting a review to %s", in.Reviewer)
	if err := u.GitHub.RequestPullRequestReview(ctx, github.RequestPullRequestReviewInput{
		Repository: in.BaseRepository,
		Number:     createdPR.PullRequestNumber,
		User:       in.Reviewer,
	}); err != nil {
		return fmt.Errorf("could not request a review for the pull request: %w", err)
	}
	return nil
}
