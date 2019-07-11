package main

import (
	"context"
	"os"

	"github.com/int128/ghcp/di"
)

func main() {
	os.Exit(di.NewCmd().Run(context.Background(), os.Args))
}
