package main

import (
	"context"

	"github.com/guitarkeegan/shop-cli/internal/fulltui"
)

// runInteractive starts the full-screen builder when shop runs with no
// subcommand.
func runInteractive(ctx context.Context) error {
	path, err := fulltui.Run()
	if err != nil {
		return err
	}
	if path != "" {
		printCreated(path)
	}
	return nil
}
