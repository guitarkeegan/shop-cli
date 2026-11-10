package main

import (
	"context"
	"fmt"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
)

func sectionHandler() clic.HandlerFunc {
	return func(ctx context.Context) error {
		return ErrMissingSubcommand
	}
}

func handleCreateSection(opt *SectionCreateOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		if opt.useTemplates {
			if err := opt.Templates.Run(); err != nil {
				return err
			}
			fmt.Println(colors.Success("Feature created!"))
			return nil
		}

		settings, err := ren.ParseSettingStrings(opt.Settings)
		if err != nil {
			return err
		}
		blocks, err := ren.ParseBlockStrings(opt.Blocks)
		if err != nil {
			return err
		}

		path, err := ren.CreateSection(
			ctx,
			ren.WithSectionName(opt.Name),
			ren.WithSectionSettings(settings),
			ren.WithSectionBlocks(blocks),
			ren.WithDefaults(),
		)
		if err != nil {
			return err
		}
		fmt.Println(colors.Success("Created " + path))
		return nil
	}
}
