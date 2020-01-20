package commitstrategy

import "testing"

func TestRebaseOn(t *testing.T) {
	s := RebaseOn("develop")

	if s.RebaseUpstream() != "develop" {
		t.Errorf("RebaseUpstream wants %s but got %s", "develop", s.RebaseUpstream())
	}
	if s.IsFastForward() {
		t.Errorf("IsFastForward wants false but got true")
	}
	if !s.IsRebase() {
		t.Errorf("IsRebase wants true but got false")
	}
	if s.NoParent() {
		t.Errorf("NoParent wants false but got true")
	}
}
