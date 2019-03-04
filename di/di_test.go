package di_test

import (
	"testing"

	"github.com/int128/ghcp/adaptors/interfaces"
	"github.com/int128/ghcp/di"
)

func TestContainer_Run(t *testing.T) {
	c := di.New()
	if err := c.Run(di.ExtraDependencies{}, func(adaptors.Cmd) {}); err != nil {
		t.Fatalf("could not resolve dependencies: %+v", err)
	}
}
