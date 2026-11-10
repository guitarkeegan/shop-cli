package main

import (
	"context"
	"fmt"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
)

func snippetHandler() clic.HandlerFunc {
	return func(ctx context.Context) error {
		return ErrMissingSubcommand
	}
}

func handleCreateSnippet(opt *SnippetCreateOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		path, err := ren.CreateSnippet(
			ctx,
			ren.WithSnippetFileName(opt.Name),
			ren.WithSnippetDescription(opt.Description),
			ren.WithSnippetParamsAsString(opt.DocParams),
			ren.WithSnippetTemplate(opt.Template),
			ren.WithSnippetAddExample(opt.Example),
		)
		if err != nil {
			return err
		}
		fmt.Println(colors.Success("Created " + path))
		return nil
	}
}
