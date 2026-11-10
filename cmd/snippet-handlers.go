package main

import (
	"context"
	"errors"

	"github.com/daved/clic"
	render "github.com/guitarkeegan/shop-cli/internal/render-service"
)

var (
	missingSnippetDirErr = errors.New("No 'snippets' dir found. Are you at the top level of the project?")
	missingNameErr       = errors.New("Snippet must have a name")
)

func snippetHandler() clic.HandlerFunc {
	return func(ctx context.Context) error {
		return missingFlagErr
	}
}

func handleCreateSnippet(opt *SnippetCreateOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		dbg("opt: name: %s", opt.Name)
		handler := render.CreateSnippet(
			ctx,
			render.WithSnippetDescription(opt.Description),
			render.WithSnippetFileName(opt.Name),
			render.WithSnippetParamsAsString(opt.DocParams),
			render.WithSnippetTemplate(opt.Template),
			render.WithSnippetAddExample(opt.Example),
		)
		return handler(ctx)
	}
}
