package di_test

import (
	"testing"

	"github.com/google/go-github/v24/github"
	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/di"
	"github.com/shurcooL/githubv4"
)

func TestNew(t *testing.T) {
	c, err := di.New()
	if err != nil {
		t.Fatalf("could not create a container: %+v", err)
	}
	if err := c.Provide(func() *github.Client { return nil }); err != nil {
		t.Error(err)
	}
	if err := c.Provide(func() *githubv4.Client { return nil }); err != nil {
		t.Error(err)
	}
	if err := c.Invoke(func(adaptors.Cmd) {}); err != nil {
		t.Fatalf("could not resolve dependencies: %+v", err)
	}
}
