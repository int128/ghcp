package main

import (
	"context"
	"log"
	"os"

	"github.com/int128/ghcp/di"
	"github.com/int128/ghcp/infrastructure/interfaces"
)

func main() {
	if err := di.Invoke(func(cmd infrastructure.Cmd) {
		os.Exit(cmd.Run(context.Background(), os.Args))
	}); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
