package main

import (
	"os"

	"github.com/rarimo/geo-forms-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
