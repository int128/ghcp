package commit

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/ghcp/pkg/fs"
	"github.com/int128/ghcp/pkg/fs/mock_fs"
	"github.com/int128/ghcp/pkg/git"
	"github.com/int128/ghcp/pkg/git/commitstrategy"
	"github.com/int128/ghcp/pkg/github"
	"github.com/int128/ghcp/pkg/github/mock_github"
	testingLogger "github.com/int128/ghcp/pkg/logger/testing"
	"github.com/int128/ghcp/usecases/gitobject"
	"github.com/int128/ghcp/usecases/gitobject/mock_gitobject"
)

var parentRepositoryID = git.RepositoryID{Owner: "upstream", Name: "repo"}
var targetRepositoryID = git.RepositoryID{Owner: "owner", Name: "repo"}

var targetRepositoryNodeID = github.InternalRepositoryNodeID("OwnerRepo")
var targetBranchNodeID = github.InternalBranchNodeID("OwnerRepoTargetBranch")

var thePathFilter = gomock.AssignableToTypeOf(&pathFilter{})
var theFiles = []fs.File{
	{Path: "file1"},
	{Path: "file2", Executable: true},
}

func newFileSystemMock(ctrl *gomock.Controller) *mock_fs.MockInterface {
	fileSystem := mock_fs.NewMockInterface(ctrl)
	fileSystem.EXPECT().FindFiles([]string{"path"}, thePathFilter).Return(theFiles, nil)
	return fileSystem
}

func newCreateGitObjectMock(ctx context.Context, ctrl *gomock.Controller, parentCommitSHA git.CommitSHA, parentTreeSHA git.TreeSHA, noFileMode bool, changedFiles int) *mock_gitobject.MockInterface {
	createGitObject := mock_gitobject.NewMockInterface(ctrl)
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			gitHub := mock_github.NewMockInterface(ctrl)
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
			gitHub.EXPECT().
				CreateBranch(ctx, github.CreateBranchInput{
					RepositoryNodeID: targetRepositoryNodeID,
					BranchName:       "topic",
					CommitSHA:        "commitSHA",
				}).
				Return(nil).
				Times(c.branchOperationTimes)

			useCase := Commit{
				CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "masterCommitSHA", "masterTreeSHA", c.noFileMode, c.changedFiles),
				FileSystem:      newFileSystemMock(ctrl),
				Logger:          testingLogger.New(t),
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
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				gitHub := mock_github.NewMockInterface(ctrl)
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
				gitHub.EXPECT().
					CreateBranch(ctx, github.CreateBranchInput{
						RepositoryNodeID: targetRepositoryNodeID,
						BranchName:       "topic",
						CommitSHA:        "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "masterCommitSHA", "masterTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(ctrl),
					Logger:          testingLogger.New(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("when the branch exists, it should update it", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				gitHub := mock_github.NewMockInterface(ctrl)
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
				gitHub.EXPECT().
					UpdateBranch(ctx, github.UpdateBranchInput{
						BranchRefNodeID: targetBranchNodeID,
						CommitSHA:       "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "topicCommitSHA", "topicTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(ctrl),
					Logger:          testingLogger.New(t),
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
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				gitHub := mock_github.NewMockInterface(ctrl)
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
				gitHub.EXPECT().
					CreateBranch(ctx, github.CreateBranchInput{
						RepositoryNodeID: targetRepositoryNodeID,
						BranchName:       "topic",
						CommitSHA:        "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "", "", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(ctrl),
					Logger:          testingLogger.New(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("when the branch exists, it should update it", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				gitHub := mock_github.NewMockInterface(ctrl)
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
				gitHub.EXPECT().
					UpdateBranch(ctx, github.UpdateBranchInput{
						BranchRefNodeID: targetBranchNodeID,
						CommitSHA:       "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "", "", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(ctrl),
					Logger:          testingLogger.New(t),
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
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				gitHub := mock_github.NewMockInterface(ctrl)
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
				gitHub.EXPECT().
					CreateBranch(ctx, github.CreateBranchInput{
						RepositoryNodeID: targetRepositoryNodeID,
						BranchName:       "topic",
						CommitSHA:        "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "developCommitSHA", "developTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(ctrl),
					Logger:          testingLogger.New(t),
					GitHub:          gitHub,
				}
				if err := useCase.Do(ctx, in); err != nil {
					t.Errorf("err wants nil but %+v", err)
				}
			})

			t.Run("when the branch exists, it should update it", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				gitHub := mock_github.NewMockInterface(ctrl)
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
				gitHub.EXPECT().
					UpdateBranch(ctx, github.UpdateBranchInput{
						BranchRefNodeID: targetBranchNodeID,
						CommitSHA:       "commitSHA",
					}).
					Return(nil).
					Times(c.branchOperationTimes)

				useCase := Commit{
					CreateGitObject: newCreateGitObjectMock(ctx, ctrl, "developCommitSHA", "developTreeSHA", c.noFileMode, c.changedFiles),
					FileSystem:      newFileSystemMock(ctrl),
					Logger:          testingLogger.New(t),
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
			changedFiles:         0,
			branchOperationTimes: 0,
		},
		"DryRun": {
			dryRun:               true,
			branchOperationTimes: 0,
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
			f := &pathFilter{Logger: testingLogger.New(t)}
			skip := f.SkipDir(c.path)
			if skip != c.skip {
				t.Errorf("skip wants %v but %v", c.skip, skip)
			}
		})
	}
}

func Test_pathFilter_ExcludeFile(t *testing.T) {
	f := &pathFilter{Logger: testingLogger.New(t)}
	exclude := f.ExcludeFile("foo")
	if exclude {
		t.Errorf("exclude wants %v but %v", false, exclude)
	}
}
