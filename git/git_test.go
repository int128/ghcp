package git_test

import (
	"testing"

	"github.com/int128/ghcp/git"
)

func TestBranchName_QualifiedName(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		qn := git.BranchName("master").QualifiedName()
		want := "refs/heads/master"
		if qn != want {
			t.Errorf("QualifiedName wants %s but %s", want, qn)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		qn := git.BranchName("").QualifiedName()
		if qn != "" {
			t.Errorf("QualifiedName wants empty but %s", qn)
		}
	})
}
