package main

import (
	"context"

	"github.com/guitarkeegan/shop-cli/internal/fulltui"
)

// runRoot handles "shop" with no subcommand: it prints the version when
// asked, and otherwise starts the full-screen builder.
func runRoot(showVersion *bool) func(context.Context) error {
	return func(ctx context.Context) error {
		if *showVersion {
			printVersion()
			return nil
		}
		return runInteractive()
	}
}

func runInteractive() error {
	path, err := fulltui.Run()
	if err != nil {
		return err
	}
	if path != "" {
		printCreated(path)
	}
	return nil
}
