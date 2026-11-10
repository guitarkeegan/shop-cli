package main

import (
	"context"
	"fmt"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	"github.com/guitarkeegan/shop-cli/internal/features"
	"github.com/guitarkeegan/shop-cli/internal/render"
	"github.com/guitarkeegan/shop-cli/internal/theme"
)

type sectionCreateOptions struct {
	Name         string
	Settings     string
	Blocks       []string
	useTemplates bool
}

func newSectionCmd(feats *features.FeatureService) *clic.Clic {
	opts := &sectionCreateOptions{}
	create := clic.NewFromFunc(createSection(opts, feats), "create")
	create.Flag(&opts.Name, "n|name", "Name of the section.")
	create.Flag(&opts.useTemplates, "t|template", "Pick a prebuilt feature to add to the theme.")
	create.Flag(&opts.Blocks, "b|blocks", "Add a block (repeatable). Ex.\nshop section create -b \"Slide,image_picker:image:Image,text:caption:Caption\"\nThe block type is derived from its name.")
	create.Operand(&opts.Settings, false, "settings", "A comma separated list of settings as type:id:label. Ex.\n\"text:title:Title,image_picker:image:Background Image\"")

	return clic.NewFromFunc(requireSubcommand, "section", create)
}

func createSection(opts *sectionCreateOptions, feats *features.FeatureService) clic.HandlerFunc {
	return func(ctx context.Context) error {
		if opts.useTemplates {
			if err := feats.Run(); err != nil {
				return err
			}
			fmt.Println(colors.Success("Feature created!"))
			return nil
		}

		settings, err := render.ParseSettings(opts.Settings)
		if err != nil {
			return err
		}
		blocks, err := render.ParseBlocks(opts.Blocks)
		if err != nil {
			return err
		}

		content, err := render.Section(
			render.WithSectionName(opts.Name),
			render.WithSectionSettings(settings),
			render.WithSectionBlocks(blocks),
			render.WithDefaults(),
		)
		if err != nil {
			return err
		}

		path, err := theme.WriteSection(opts.Name, content)
		if err != nil {
			return err
		}
		printCreated(path)
		return nil
	}
}
