package di_test

import (
	"testing"

	adaptors2 "github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/di"
)

func TestContainer_Run(t *testing.T) {
	if err := di.Invoke(func(adaptors2.Cmd) {}); err != nil {
		t.Fatalf("could not resolve dependencies: %+v", err)
	}
}
