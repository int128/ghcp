package main

import (
	"os"

	"github.com/int128/ghcp/pkg/di"
)

var version = "HEAD"

func main() {
	os.Exit(di.NewCmd().Run(os.Args, version))
}
