package main

import (
	"context"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/render"
	"github.com/guitarkeegan/shop-cli/internal/theme"
)

type blockCreateOptions struct {
	Name     string
	Settings string
}

func newBlockCmd() *clic.Clic {
	opts := &blockCreateOptions{}
	create := clic.NewFromFunc(createBlock(opts), "create")
	create.Flag(&opts.Name, "n|name", "Name of the block.")
	create.Operand(&opts.Settings, false, "settings", "A comma separated list of settings as type:id:label.")

	return clic.NewFromFunc(requireSubcommand, "block", create)
}

func createBlock(opts *blockCreateOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		settings, err := render.ParseSettings(opts.Settings)
		if err != nil {
			return err
		}

		content, err := render.Block(
			render.WithBlockName(opts.Name),
			render.WithBlockSettings(settings),
		)
		if err != nil {
			return err
		}

		path, err := theme.WriteBlock(opts.Name, content)
		if err != nil {
			return err
		}
		printCreated(path)
		return nil
	}
}
