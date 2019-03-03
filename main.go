package main

import (
	"context"
	"os"

	"github.com/int128/ghcp/infrastructure"
)

func main() {
	os.Exit(infrastructure.Run(context.Background(), os.Args))
}
