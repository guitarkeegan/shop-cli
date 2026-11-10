package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/daved/clic"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
)

var (
	missingFlagErr       = errors.New("missing required flag")
	missingSectionDirErr = errors.New("missing sections directory. Are you at the top level of your project?")
)

func sectionHandler() clic.HandlerFunc {
	return func(ctx context.Context) error {
		return missingFlagErr
	}
}

func handleCreateSection(opt *SectionCreateOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {

		if opt.useTemplates {
			err := opt.Templates.Run()
			if err != nil {
				return err
			}
			fmt.Println("feature created!")
			return nil
		}

		ss, err := ren.ParseSettingStrings(opt.Settings)
		if err != nil {
			return err
		}

		return ren.CreateSection(
			ctx,
			ren.WithSectionBlocksString(opt.Blocks),
			ren.WithSectionTemplate(opt.useTemplates),
			ren.WithSectionName(opt.Name),
			ren.WithSectionSettings(ss),
			ren.WithDefaults())(ctx)
	}
}
