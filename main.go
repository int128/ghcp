package main

import (
	"context"
	"log"
	"os"

	"github.com/int128/ghcp/adaptors"
	"github.com/int128/ghcp/di"
)

func main() {
	if err := di.Invoke(func(cmd adaptors.Cmd) {
		os.Exit(cmd.Run(context.Background(), os.Args))
	}); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
