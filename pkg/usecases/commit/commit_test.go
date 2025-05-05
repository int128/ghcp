package commit

import (
	"context"
	"testing"

	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/fs_mock"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/github_mock"
	"github.com/int128/ghcp/mocks/github.com/int128/ghcp/pkg/usecases/gitobject_mock"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/usecases/gitobject"
	"github.com/stretchr/testify/mock"
)

var parentRepositoryID = git.RepositoryID{Owner: "upstream", Name: "repo"}
var targetRepositoryID = git.RepositoryID{Owner: "owner", Name: "repo"}

var targetRepositoryNodeID = github.InternalRepositoryNodeID("OwnerRepo")
var targetBranchNodeID = github.InternalBranchNodeID("OwnerRepoTargetBranch")

var thePathFilter = mock.Anything
var theFiles = []fs.File{
	{Path: "file1"},
	{Path: "file2", Executable: true},
}

func newFileSystemMock(t *testing.T) *fs_mock.MockInterface {
	fileSystem := fs_mock.NewMockInterface(t)
	fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)
	return fileSystem
}

func newCreateGitObjectMock(ctx context.Context, t *testing.T, parentCommitSHA git.CommitSHA, parentTreeSHA git.TreeSHA, noFileMode bool, changedFiles int) *gitobject_mock.MockInterface {
	createGitObject := gitobject_mock.NewMockInterface(t)
	createGitObject.EXPECT().
		Do(ctx, gitobject.Input{
			Files:           theFiles,
			Repository:      targetRepositoryID,
			CommitMessage:   "message",
			ParentCommitSHA: parentCommitSHA,
			ParentTreeSHA:   parentTreeSHA,
			NoFileMode:      noFileMode,
		}).
		Return(&gitobject.Output{
			CommitSHA:    "commitSHA",
			ChangedFiles: changedFiles,
		}, nil)
	return createGitObject
}

func TestCommitToBranch_Do(t *testing.T) {
	ctx := context.TODO()

	type testCase struct {
		dryRun               bool
		noFileMode           bool
		changedFiles         int
		branchOperationTimes int
	}

	run := func(t *testing.T, c testCase) {
		t.Run("when the branch name is not set, it should resolve the default branch", func(t *testing.T) {
			in := Input{
				TargetRepository: targetRepositoryID,
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.FastForward,
				CommitMessage:    "message",
				Paths:            []string{"path"},
				NoFileMode:       c.noFileMode,
				DryRun:           c.dryRun,
			}
			gitHub := github_mock.NewMockInterface(t)
			gitHub.EXPECT().
				QueryDefaultBranch(ctx, github.QueryDefaultBranchInput{
					HeadRepository: targetRepositoryID,
					BaseRepository: parentRepositoryID,
				}).
				Return(&github.QueryDefaultBranchOutput{
					HeadDefaultBranchName: "topic",
				}, nil)
			gitHub.EXPECT().
				QueryForCommit(ctx, github.QueryForCommitInput{
					ParentRepository: parentRepositoryID,
					TargetRepository: targetRepositoryID,
					TargetBranchName: "topic",
				}).
				Return(&github.QueryForCommitOutput{
					CurrentUserName:              "current",
					ParentDefaultBranchCommitSHA: "masterCommitSHA",
					ParentDefaultBranchTreeSHA:   "masterTreeSHA",
					TargetRepositoryNodeID:       targetRepositoryNodeID,
				}, nil)
			if c.branchOperationTimes > 0 {
				gitHub.EXPECT().
					CreateBranch(ctx, github.CreateBranchInput{
						RepositoryNodeID: targetRepositoryNodeID,
						BranchName:       "topic",
						CommitSHA:        "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)
			}

			useCase := Commit{
				CreateGitObject: newCreateGitObjectMock(ctx, t, "masterCommitSHA", "masterTreeSHA", c.noFileMode, c.changedFiles),
				FileSystem:      newFileSystemMock(t),
				GitHub:          gitHub,
			}
			if err := useCase.Do(ctx, in); err != nil {
				t.Errorf("err wants nil but %+v", err)
			}
		})

		t.Run("by fast-forward", func(t *testing.T) {
			in := Input{
				TargetRepository: targetRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.FastForward,
				CommitMessage:    "message",
				Paths:            []string{"path"},
				NoFileMode:       c.noFileMode,
				DryRun:           c.dryRun,
			}

			t.Run("when a branch does not exist, it should create it", func(t *testing.T) {
				gitHub := github_mock.NewMockInterface(t)
				gitHub.EXPECT().
					QueryForCommit(ctx, github.QueryForCommitInput{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&github.QueryForCommitOutput{
						CurrentUserName:              "current",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetRepositoryNodeID:       targetRepositoryNodeID,
					}, nil)
				if c.branchOperationTimes > 0 {
					gitHub.EXPECT().
						CreateBranch(ctx, github.CreateBranchInput{
							RepositoryNodeID: targetRepositoryNodeID,
							BranchName:       "topic",
							CommitSHA:        "commitSHA",
						}).
						Return(nil).
						Times(c.branchOperationTimes)
				}

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, t, "masterCommitSHA", "masterTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("when the branch exists, it should update it", func(t *testing.T) {
				gitHub := github_mock.NewMockInterface(t)
				gitHub.EXPECT().
					QueryForCommit(ctx, github.QueryForCommitInput{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&github.QueryForCommitOutput{
						CurrentUserName:              "current",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchNodeID:           targetBranchNodeID,
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
					}, nil)
				if c.branchOperationTimes > 0 {
					gitHub.EXPECT().
						UpdateBranch(ctx, github.UpdateBranchInput{
							BranchRefNodeID: targetBranchNodeID,
							CommitSHA:       "commitSHA",
						}).
						Return(nil).
						Times(c.branchOperationTimes)
				}

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, t, "topicCommitSHA", "topicTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})
		})

		t.Run("with no parent", func(t *testing.T) {
			in := Input{
				TargetRepository: targetRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.NoParent,
				CommitMessage:    "message",
				Paths:            []string{"path"},
				NoFileMode:       c.noFileMode,
				DryRun:           c.dryRun,
			}

			t.Run("when a branch does not exist, it should create it", func(t *testing.T) {
				gitHub := github_mock.NewMockInterface(t)
				gitHub.EXPECT().
					QueryForCommit(ctx, github.QueryForCommitInput{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&github.QueryForCommitOutput{
						CurrentUserName:              "current",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetRepositoryNodeID:       targetRepositoryNodeID,
					}, nil)
				if c.branchOperationTimes > 0 {
					gitHub.EXPECT().
						CreateBranch(ctx, github.CreateBranchInput{
							RepositoryNodeID: targetRepositoryNodeID,
							BranchName:       "topic",
							CommitSHA:        "commitSHA",
						}).
						Return(nil).
						Times(c.branchOperationTimes)
				}

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, t, "", "", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("when the branch exists, it should update it", func(t *testing.T) {
				gitHub := github_mock.NewMockInterface(t)
				gitHub.EXPECT().
					QueryForCommit(ctx, github.QueryForCommitInput{
						ParentRepository: parentRepositoryID,
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&github.QueryForCommitOutput{
						CurrentUserName:              "current",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchNodeID:           targetBranchNodeID,
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
					}, nil)
				if c.branchOperationTimes > 0 {
					gitHub.EXPECT().
						UpdateBranch(ctx, github.UpdateBranchInput{
							BranchRefNodeID: targetBranchNodeID,
							CommitSHA:       "commitSHA",
						}).
						Return(nil).
						Times(c.branchOperationTimes)
				}

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, t, "", "", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})
		})

		t.Run("by rebase", func(t *testing.T) {
			in := Input{
				TargetRepository: targetRepositoryID,
				TargetBranchName: "topic",
				ParentRepository: parentRepositoryID,
				CommitStrategy:   commitstrategy.RebaseOn("develop"),
				CommitMessage:    "message",
				Paths:            []string{"path"},
				NoFileMode:       c.noFileMode,
				DryRun:           c.dryRun,
			}

			t.Run("when a branch does not exist, it should create it", func(t *testing.T) {
				gitHub := github_mock.NewMockInterface(t)
				gitHub.EXPECT().
					QueryForCommit(ctx, github.QueryForCommitInput{
						ParentRepository: parentRepositoryID,
						ParentRef:        "develop",
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&github.QueryForCommitOutput{
						CurrentUserName:              "current",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						ParentRefCommitSHA:           "developCommitSHA",
						ParentRefTreeSHA:             "developTreeSHA",
						TargetRepositoryNodeID:       targetRepositoryNodeID,
					}, nil)
				if c.branchOperationTimes > 0 {
					gitHub.EXPECT().
						CreateBranch(ctx, github.CreateBranchInput{
							RepositoryNodeID: targetRepositoryNodeID,
							BranchName:       "topic",
							CommitSHA:        "commitSHA",
						}).
						Return(nil).
						Times(c.branchOperationTimes)
				}

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, t, "developCommitSHA", "developTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("when the branch exists, it should update it", func(t *testing.T) {
				gitHub := github_mock.NewMockInterface(t)
				gitHub.EXPECT().
					QueryForCommit(ctx, github.QueryForCommitInput{
						ParentRepository: parentRepositoryID,
						ParentRef:        "develop",
						TargetRepository: targetRepositoryID,
						TargetBranchName: "topic",
					}).
					Return(&github.QueryForCommitOutput{
						CurrentUserName:              "current",
						ParentDefaultBranchCommitSHA: "masterCommitSHA",
						ParentDefaultBranchTreeSHA:   "masterTreeSHA",
						TargetBranchNodeID:           targetBranchNodeID,
						TargetBranchCommitSHA:        "topicCommitSHA",
						TargetBranchTreeSHA:          "topicTreeSHA",
						ParentRefCommitSHA:           "developCommitSHA",
						ParentRefTreeSHA:             "developTreeSHA",
					}, nil)
				if c.branchOperationTimes > 0 {
					gitHub.EXPECT().
						UpdateBranch(ctx, github.UpdateBranchInput{
							BranchRefNodeID: targetBranchNodeID,
							CommitSHA:       "commitSHA",
						}).
						Return(nil).
						Times(c.branchOperationTimes)
				}

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, t, "developCommitSHA", "developTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})
		})
	}

	for name, c := range map[string]testCase{
		"Success": {
			changedFiles:         1,
			branchOperationTimes: 1,
		},
		"NothingToCommit": {
			changedFiles: 0,
		},
		"DryRun": {
			dryRun: true,
		},
		"NoFileMode": {
			noFileMode: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			run(t, c)
		})
	}
}

func Test_pathFilter_SkipDir(t *testing.T) {
	for _, c := range []struct {
		path string
		skip bool
	}{
		{path: "."},
		{path: "foo"},
		{path: ".git", skip: true},
		{path: "foo/bar"},
		{path: "foo/.git", skip: true},
	} {
		t.Run(c.path, func(t *testing.T) {
			f := &pathFilter{}
			skip := f.SkipDir(c.path)
			if skip != c.skip {
				t.Errorf("skip wants %v but %v", c.skip, skip)
			}
		})
	}
}

func Test_pathFilter_ExcludeFile(t *testing.T) {
	f := &pathFilter{}
	exclude := f.ExcludeFile("foo")
	if exclude {
		t.Errorf("exclude wants %v but %v", false, exclude)
	}
}
