package di_test

import (
	"testing"

	"github.com/int128/ghcp/di"
	"github.com/int128/ghcp/infrastructure/interfaces"
)

func TestContainer_Run(t *testing.T) {
	if err := di.Invoke(func(infrastructure.Cmd) {}); err != nil {
		t.Fatalf("could not resolve dependencies: %+v", err)
	}
}
