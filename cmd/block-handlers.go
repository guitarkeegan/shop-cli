package main

import (
	"context"
	"fmt"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
)

type CreateBlockOptions struct {
	Name   string
	Inputs string
}

func NewCreateBlockOptions() *CreateBlockOptions {
	return &CreateBlockOptions{}
}

func handleBlock(ctx context.Context) error {
	return ErrMissingSubcommand
}

func handleCreateBlock(opt *CreateBlockOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		settings, err := ren.ParseSettingStrings(opt.Inputs)
		if err != nil {
			return err
		}

		path, err := ren.CreateBlock(ctx, ren.WithBlockName(opt.Name), ren.WithBlockSettings(settings))
		if err != nil {
			return err
		}
		fmt.Println(colors.Success("Created " + path))
		return nil
	}
}
