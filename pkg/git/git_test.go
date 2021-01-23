package git

import (
	"testing"
)

func TestBranchName_QualifiedName(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		qn := BranchName("master").QualifiedName()
		want := RefQualifiedName{Prefix: "refs/heads/", Name: "master"}
		if qn != want {
			t.Errorf("QualifiedName wants %+v but %+v", want, qn)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		qn := BranchName("").QualifiedName()
		if qn != (RefQualifiedName{}) {
			t.Errorf("QualifiedName wants zero but %+v", qn)
		}
	})
}

func TestRefQualifiedName_IsValid(t *testing.T) {
	for _, c := range []struct {
		RefQualifiedName RefQualifiedName
		Valid            bool
	}{
		{RefQualifiedName{}, false},
		{RefQualifiedName{Prefix: "refs/heads/"}, false},
		{RefQualifiedName{Name: "master"}, false},
		{RefQualifiedName{Prefix: "refs/heads/", Name: "master"}, true},
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
		RefQualifiedName RefQualifiedName
		String           string
	}{
		{RefQualifiedName{}, ""},
		{RefQualifiedName{Prefix: "refs/heads/"}, ""},
		{RefQualifiedName{Name: "master"}, ""},
		{RefQualifiedName{Prefix: "refs/heads/", Name: "master"}, "refs/heads/master"},
	} {
		t.Run(c.RefQualifiedName.String(), func(t *testing.T) {
			s := c.RefQualifiedName.String()
			if s != c.String {
				t.Errorf("String wants %v but %v", c.String, s)
			}
		})
	}
}
