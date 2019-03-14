package git_test

import (
	"testing"

	"github.com/int128/ghcp/git"
)

func TestBranchName_QualifiedName(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		qn := git.BranchName("master").QualifiedName()
		want := git.RefQualifiedName{Prefix: "refs/heads/", Name: "master"}
		if qn != want {
			t.Errorf("QualifiedName wants %+v but %+v", want, qn)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		qn := git.BranchName("").QualifiedName()
		if qn != (git.RefQualifiedName{}) {
			t.Errorf("QualifiedName wants zero but %+v", qn)
		}
	})
}

func TestRefQualifiedName_IsValid(t *testing.T) {
	for _, c := range []struct {
		RefQualifiedName git.RefQualifiedName
		Valid            bool
	}{
		{git.RefQualifiedName{}, false},
		{git.RefQualifiedName{Prefix: "refs/heads/"}, false},
		{git.RefQualifiedName{Name: "master"}, false},
		{git.RefQualifiedName{Prefix: "refs/heads/", Name: "master"}, true},
	} {
		t.Run(c.RefQualifiedName.String(), func(t *testing.T) {
			valid := c.RefQualifiedName.IsValid()
			if valid != c.Valid {
				t.Errorf("IsValid wants %v but %v", c.Valid, valid)
			}
		})
	}
}

func TestRefQualifiedName_String(t *testing.T) {
	for _, c := range []struct {
		RefQualifiedName git.RefQualifiedName
		String           string
	}{
		{git.RefQualifiedName{}, ""},
		{git.RefQualifiedName{Prefix: "refs/heads/"}, ""},
		{git.RefQualifiedName{Name: "master"}, ""},
		{git.RefQualifiedName{Prefix: "refs/heads/", Name: "master"}, "refs/heads/master"},
	} {
		t.Run(c.RefQualifiedName.String(), func(t *testing.T) {
			s := c.RefQualifiedName.String()
			if s != c.String {
				t.Errorf("String wants %v but %v", c.String, s)
			}
		})
	}
}
