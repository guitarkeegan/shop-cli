package main

import (
	"context"
	"fmt"

	"github.com/guitarkeegan/shop-cli/internal/colors"
	"github.com/guitarkeegan/shop-cli/internal/fulltui"
)

func rootHandler(ctx context.Context) error {
	path, err := fulltui.Run()
	if err != nil {
		return err
	}
	if path != "" {
		fmt.Println(colors.Success("Created " + path))
	}
	return nil
}
