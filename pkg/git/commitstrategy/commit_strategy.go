package commitstrategy

import (
	"fmt"

	"github.com/int128/ghcp/pkg/git"
)

// CommitStrategy represents a method to create a commit object.
type CommitStrategy interface {
	IsFastForward() bool
	IsRebase() bool
	RebaseUpstream() git.RefName
	NoParent() bool
	String() string
}

type commitStrategy struct {
	name        string
	fastForward bool
	rebase      bool
	noParent    bool
}

func (s *commitStrategy) IsFastForward() bool         { return s.fastForward }
func (s *commitStrategy) IsRebase() bool              { return s.rebase }
func (s *commitStrategy) RebaseUpstream() git.RefName { return "" }
func (s *commitStrategy) NoParent() bool              { return s.noParent }
func (s *commitStrategy) String() string              { return s.name }

// FastForward represents the fast-forward.
var FastForward CommitStrategy = &commitStrategy{name: "fast-forward", fastForward: true}

// NoParent represents the method to create a commit without any parent
var NoParent CommitStrategy = &commitStrategy{name: "no-parent", noParent: true}

// RebaseOn represents the rebase on the upstream.
func RebaseOn(upstreamRef git.RefName) CommitStrategy {
	return &rebase{
		commitStrategy: commitStrategy{name: fmt.Sprintf("rebase-on-%v", upstreamRef), rebase: true},
		upstreamRef:    upstreamRef,
	}
}

type rebase struct {
	commitStrategy
	upstreamRef git.RefName
}

func (f *rebase) RebaseUpstream() git.RefName { return f.upstreamRef }
