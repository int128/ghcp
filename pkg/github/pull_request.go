package github

import (
	"context"
	"fmt"

	"github.com/int128/ghcp/pkg/git"
	"github.com/shurcooL/githubv4"
	"golang.org/x/xerrors"
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
	PullRequestURL       string      // URL of the pull request associated to the head branch, if exists
	ReviewerUserNodeID   githubv4.ID // optional
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
					Nodes []struct {
						URL string
					}
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
	c.Logger.Debugf("Querying the existing pull request with %+v", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", q)

	out := QueryForPullRequestOutput{
		CurrentUserName:      q.Viewer.Login,
		BaseRepositoryNodeID: q.BaseRepository.ID,
		HeadBranchCommitSHA:  git.CommitSHA(q.HeadRepository.Ref.Target.OID),
		ReviewerUserNodeID:   q.ReviewerUser.ID,
	}
	if len(q.HeadRepository.Ref.AssociatedPullRequests.Nodes) > 0 {
		out.PullRequestURL = q.HeadRepository.Ref.AssociatedPullRequests.Nodes[0].URL
	}
	c.Logger.Debugf("Returning the result: %+v", out)
	return &out, nil
}

type CreatePullRequestInput struct {
	BaseRepository       git.RepositoryID
	BaseBranchName       git.BranchName
	BaseRepositoryNodeID InternalRepositoryNodeID
	HeadRepository       git.RepositoryID
	HeadBranchName       git.BranchName
	Title                string
	Body                 string
}

type CreatePullRequestOutput struct {
	PullRequestNodeID githubv4.ID
	URL               string
}

func (c *GitHub) CreatePullRequest(ctx context.Context, in CreatePullRequestInput) (*CreatePullRequestOutput, error) {
	c.Logger.Debugf("Creating a pull request %+v", in)
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
		Body:         githubv4.NewString(githubv4.String(in.Body)),
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
		return nil, xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", m)
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
	c.Logger.Debugf("Requesting a review for the pull request %+v", in)
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
		return xerrors.Errorf("GitHub API error: %w", err)
	}
	c.Logger.Debugf("Got the result: %+v", m)
	return nil
}
