// Code generated by mockery v2.52.2. DO NOT EDIT.

package client_mock

import (
	context "context"

	github "github.com/google/go-github/v70/github"
	mock "github.com/stretchr/testify/mock"
)

// MockPullRequestsService is an autogenerated mock type for the PullRequestsService type
type MockPullRequestsService struct {
	mock.Mock
}

type MockPullRequestsService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPullRequestsService) EXPECT() *MockPullRequestsService_Expecter {
	return &MockPullRequestsService_Expecter{mock: &_m.Mock}
}

// RequestReviewers provides a mock function with given fields: ctx, owner, repo, number, reviewers
func (_m *MockPullRequestsService) RequestReviewers(ctx context.Context, owner string, repo string, number int, reviewers github.ReviewersRequest) (*github.PullRequest, *github.Response, error) {
	ret := _m.Called(ctx, owner, repo, number, reviewers)

	if len(ret) == 0 {
		panic("no return value specified for RequestReviewers")
	}

	var r0 *github.PullRequest
	var r1 *github.Response
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, github.ReviewersRequest) (*github.PullRequest, *github.Response, error)); ok {
		return rf(ctx, owner, repo, number, reviewers)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int, github.ReviewersRequest) *github.PullRequest); ok {
		r0 = rf(ctx, owner, repo, number, reviewers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*github.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int, github.ReviewersRequest) *github.Response); ok {
		r1 = rf(ctx, owner, repo, number, reviewers)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*github.Response)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, int, github.ReviewersRequest) error); ok {
		r2 = rf(ctx, owner, repo, number, reviewers)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockPullRequestsService_RequestReviewers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestReviewers'
type MockPullRequestsService_RequestReviewers_Call struct {
	*mock.Call
}

// RequestReviewers is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - repo string
//   - number int
//   - reviewers github.ReviewersRequest
func (_e *MockPullRequestsService_Expecter) RequestReviewers(ctx interface{}, owner interface{}, repo interface{}, number interface{}, reviewers interface{}) *MockPullRequestsService_RequestReviewers_Call {
	return &MockPullRequestsService_RequestReviewers_Call{Call: _e.mock.On("RequestReviewers", ctx, owner, repo, number, reviewers)}
}

func (_c *MockPullRequestsService_RequestReviewers_Call) Run(run func(ctx context.Context, owner string, repo string, number int, reviewers github.ReviewersRequest)) *MockPullRequestsService_RequestReviewers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int), args[4].(github.ReviewersRequest))
	})
	return _c
}

func (_c *MockPullRequestsService_RequestReviewers_Call) Return(_a0 *github.PullRequest, _a1 *github.Response, _a2 error) *MockPullRequestsService_RequestReviewers_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockPullRequestsService_RequestReviewers_Call) RunAndReturn(run func(context.Context, string, string, int, github.ReviewersRequest) (*github.PullRequest, *github.Response, error)) *MockPullRequestsService_RequestReviewers_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPullRequestsService creates a new instance of MockPullRequestsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPullRequestsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPullRequestsService {
	mock := &MockPullRequestsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
