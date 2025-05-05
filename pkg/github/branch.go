package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/int128/ghcp/pkg/git"
	"github.com/shurcooL/githubv4"
)

type QueryForCommitInput struct {
	ParentRepository git.RepositoryID
	ParentRef        git.RefName // optional
	TargetRepository git.RepositoryID
	TargetBranchName git.BranchName // optional
}

type QueryForCommitOutput struct {
	CurrentUserName              string
	ParentDefaultBranchCommitSHA git.CommitSHA
	ParentDefaultBranchTreeSHA   git.TreeSHA
	ParentRefCommitSHA           git.CommitSHA // empty if the parent ref does not exist
	ParentRefTreeSHA             git.TreeSHA   // empty if the parent ref does not exist
	TargetRepositoryNodeID       InternalRepositoryNodeID
	TargetBranchNodeID           InternalBranchNodeID
	TargetBranchCommitSHA        git.CommitSHA // empty if the branch does not exist
	TargetBranchTreeSHA          git.TreeSHA   // empty if the branch does not exist
}

func (q *QueryForCommitOutput) TargetBranchExists() bool {
	return q.TargetBranchCommitSHA != ""
}

// QueryForCommit returns the repository for creating or updating the branch.
func (c *GitHub) QueryForCommit(ctx context.Context, in QueryForCommitInput) (*QueryForCommitOutput, error) {
	var q struct {
		Viewer struct {
			Login string
		}

		ParentRepository struct {
			// default branch
			DefaultBranchRef struct {
				Name   string
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			}

			// parent ref (optional)
			ParentRef struct {
				Prefix string
				Name   string
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			} `graphql:"parentRef: ref(qualifiedName: $parentRef)"`
		} `graphql:"parentRepository: repository(owner: $parentOwner, name: $parentRepo)"`

		TargetRepository struct {
			ID  githubv4.ID
			Ref struct {
				ID     githubv4.ID
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			} `graphql:"ref(qualifiedName: $targetRef)"`
		} `graphql:"targetRepository: repository(owner: $targetOwner, name: $targetRepo)"`
	}
	v := map[string]interface{}{
		"parentOwner": githubv4.String(in.ParentRepository.Owner),
		"parentRepo":  githubv4.String(in.ParentRepository.Name),
		"parentRef":   githubv4.String(in.ParentRef),
		"targetOwner": githubv4.String(in.TargetRepository.Owner),
		"targetRepo":  githubv4.String(in.TargetRepository.Name),
		"targetRef":   githubv4.String(in.TargetBranchName.QualifiedName().String()),
	}
	slog.Debug("Querying the repository with", "params", v)
	if err := c.Client.Query(ctx, &q, v); err != nil {
		return nil, fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", q)
	out := QueryForCommitOutput{
		CurrentUserName:              q.Viewer.Login,
		ParentDefaultBranchCommitSHA: git.CommitSHA(q.ParentRepository.DefaultBranchRef.Target.Commit.Oid),
		ParentDefaultBranchTreeSHA:   git.TreeSHA(q.ParentRepository.DefaultBranchRef.Target.Commit.Tree.Oid),
		ParentRefCommitSHA:           git.CommitSHA(q.ParentRepository.ParentRef.Target.Commit.Oid),
		ParentRefTreeSHA:             git.TreeSHA(q.ParentRepository.ParentRef.Target.Commit.Tree.Oid),
		TargetRepositoryNodeID:       q.TargetRepository.ID,
		TargetBranchNodeID:           q.TargetRepository.Ref.ID,
		TargetBranchCommitSHA:        git.CommitSHA(q.TargetRepository.Ref.Target.Commit.Oid),
		TargetBranchTreeSHA:          git.TreeSHA(q.TargetRepository.Ref.Target.Commit.Tree.Oid),
	}
	slog.Debug("Returning the repository", "repository", out)
	return &out, nil
}

type CreateBranchInput struct {
	RepositoryNodeID InternalRepositoryNodeID
	BranchName       git.BranchName
	CommitSHA        git.CommitSHA
}

// CreateBranch creates a branch and returns nil or an error.
func (c *GitHub) CreateBranch(ctx context.Context, in CreateBranchInput) error {
	// https://docs.github.com/en/graphql/reference/mutations#createref
	v := githubv4.CreateRefInput{
		RepositoryID: in.RepositoryNodeID,
		Name:         githubv4.String(in.BranchName.QualifiedName().String()),
		Oid:          githubv4.GitObjectID(in.CommitSHA),
	}
	slog.Debug("Mutation createRef", "params", v)
	var m struct {
		CreateRef struct {
			Ref struct {
				Name string
			}
		} `graphql:"createRef(input: $input)"`
	}
	if err := c.Client.Mutate(ctx, &m, v, nil); err != nil {
		return fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", m)
	return nil
}

type UpdateBranchInput struct {
	BranchRefNodeID InternalBranchNodeID
	CommitSHA       git.CommitSHA
	Force           bool
}

// UpdateBranch updates the branch and returns nil or an error.
func (c *GitHub) UpdateBranch(ctx context.Context, in UpdateBranchInput) error {
	// https://docs.github.com/en/graphql/reference/mutations#updateref
	v := githubv4.UpdateRefInput{
		RefID: in.BranchRefNodeID,
		Oid:   githubv4.GitObjectID(in.CommitSHA),
		Force: githubv4.NewBoolean(githubv4.Boolean(in.Force)),
	}
	slog.Debug("Mutation updateRef", "params", v)
	var m struct {
		UpdateRef struct {
			Ref struct {
				Name string
			}
		} `graphql:"updateRef(input: $input)"`
	}
	if err := c.Client.Mutate(ctx, &m, v, nil); err != nil {
		return fmt.Errorf("GitHub API error: %w", err)
	}
	slog.Debug("Got the response", "response", m)
	return nil
}
