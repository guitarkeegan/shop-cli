package main

import (
	"context"

	"github.com/guitarkeegan/shop-cli/internal/fulltui"
)

func rootHandler(ctx context.Context) error {
	if err := fulltui.Run(); err != nil {
		return err
	}
	return nil
}
