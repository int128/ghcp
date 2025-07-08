package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/int128/ghcp/pkg/git"
	"github.com/shurcooL/githubv4"
)

type QueryForPullRequestInput struct {
	BaseRepository git.RepositoryID
	BaseBranchName git.BranchName
	HeadRepository git.RepositoryID
	HeadBranchName git.BranchName
	ReviewerUser   string // optional
}

type QueryForPullRequestOutput struct {
	CurrentUserName      string
	BaseRepositoryNodeID InternalRepositoryNodeID
	HeadBranchCommitSHA  git.CommitSHA
	ExistingPullRequests []ExistingPullRequest
	ReviewerUserNodeID   githubv4.ID // optional
}

type ExistingPullRequest struct {
	URL   string
	State githubv4.PullRequestState
}

// QueryForPullRequest performs the query for creating a pull request.
func (c *GitHub) QueryForPullRequest(ctx context.Context, in QueryForPullRequestInput) (*QueryForPullRequestOutput, error) {
	var q struct {
		Viewer struct {
			Login string
		}
		BaseRepository struct {
			ID githubv4.ID
		} `graphql:"baseRepository: repository(owner: $baseOwner, name: $baseRepo)"`
		HeadRepository struct {
			Ref struct {
				Target struct {
					OID string
				}
				AssociatedPullRequests struct {
					Nodes []ExistingPullRequest
				} `graphql:"associatedPullRequests(baseRefName: $baseRefName, first: 1)"`
			} `graphql:"ref(qualifiedName: $headRefName)"`
		} `graphql:"headRepository: repository(owner: $headOwner, name: $headRepo)"`
		ReviewerUser struct {
			ID githubv4.ID
		} `graphql:"reviewer: user(login: $reviewerUser) @include(if: $withReviewerUser)"`
	}
	v := map[string]interface{}{
		"baseOwner":        githubv4.String(in.BaseRepository.Owner),
		"baseRepo":         githubv4.String(in.BaseRepository.Name),
		"baseRefName":      githubv4.String(in.BaseBranchName),
		"headOwner":        githubv4.String(in.HeadRepository.Owner),
		"headRepo":         githubv4.String(in.HeadRepository.Name),
		"headRefName":      githubv4.String(in.HeadBranchName.QualifiedName().String()),
		"reviewerUser":     githubv4.String(in.ReviewerUser),
		"withReviewerUser": githubv4.Boolean(in.ReviewerUser != ""),
	}
	slog.Debug("Querying the existing pull request", "params", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", q)

	out := QueryForPullRequestOutput{
		CurrentUserName:      q.Viewer.Login,
		BaseRepositoryNodeID: q.BaseRepository.ID,
		HeadBranchCommitSHA:  git.CommitSHA(q.HeadRepository.Ref.Target.OID),
		ExistingPullRequests: q.HeadRepository.Ref.AssociatedPullRequests.Nodes,
		ReviewerUserNodeID:   q.ReviewerUser.ID,
	}
	slog.Debug("Returning the result", "result", out)
	return &out, nil
}

type CreatePullRequestInput struct {
	BaseRepository       git.RepositoryID
	BaseBranchName       git.BranchName
	BaseRepositoryNodeID InternalRepositoryNodeID
	HeadRepository       git.RepositoryID
	HeadBranchName       git.BranchName
	Title                string
	Body                 string // optional
	Draft                bool
}

type CreatePullRequestOutput struct {
	PullRequestNodeID githubv4.ID
	URL               string
}

func (c *GitHub) CreatePullRequest(ctx context.Context, in CreatePullRequestInput) (*CreatePullRequestOutput, error) {
	slog.Debug("Creating a pull request", "input", in)
	headRefName := string(in.HeadBranchName)
	if in.BaseRepository != in.HeadRepository {
		// For cross-repository pull requests.
		// https://developer.github.com/v4/input_object/createpullrequestinput/
		headRefName = fmt.Sprintf("%s:%s", in.HeadRepository.Owner, in.HeadBranchName)
	}
	v := githubv4.CreatePullRequestInput{
		RepositoryID: in.BaseRepositoryNodeID,
		BaseRefName:  githubv4.String(in.BaseBranchName),
		HeadRefName:  githubv4.String(headRefName),
		Title:        githubv4.String(in.Title),
	}
	if in.Body != "" {
		v.Body = githubv4.NewString(githubv4.String(in.Body))
	}
	if in.Draft {
		v.Draft = githubv4.NewBoolean(true)
	}
	var m struct {
		CreatePullRequest struct {
			PullRequest struct {
				ID  githubv4.ID
				URL string
			}
		} `graphql:"createPullRequest(input: $input)"`
	}
	if err := c.Client.Mutate(ctx, &m, v, nil); err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", m)
	return &CreatePullRequestOutput{
		PullRequestNodeID: m.CreatePullRequest.PullRequest.ID,
		URL:               m.CreatePullRequest.PullRequest.URL,
	}, nil
}

type RequestPullRequestReviewInput struct {
	PullRequest githubv4.ID
	User        githubv4.ID
}

func (c *GitHub) RequestPullRequestReview(ctx context.Context, in RequestPullRequestReviewInput) error {
	slog.Debug("Requesting a review for the pull request", "pullRequest", in.PullRequest, "user", in.User)
	v := githubv4.RequestReviewsInput{
		PullRequestID: in.PullRequest,
		UserIDs:       &[]githubv4.ID{in.User},
	}
	var m struct {
		RequestReviews struct {
			Actor struct {
				Login string
			}
		} `graphql:"requestReviews(input: $input)"`
	}
	if err := c.Client.Mutate(ctx, &m, v, nil); err != nil {
		return fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", m)
	return nil
}
