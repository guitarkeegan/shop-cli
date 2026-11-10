package main

import (
	"context"
	"errors"

	"github.com/daved/clic"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
)

var (
	missingBlocksDirErr = errors.New("missing blocks directory. Are you at the top level of your project?")
)

type CreateBlockOptions struct {
	Name   string
	Inputs string
}

func NewCreateBlockOptions() *CreateBlockOptions {
	return &CreateBlockOptions{
		Name: "",
	}
}

func handleBlock(ctx context.Context) error {
	return missingFlagErr
}

func handleCreateBlock(opt *CreateBlockOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		var blockOpts []ren.RenderBlockServiceOption

		blockOpts = append(blockOpts, ren.WithBlockName(opt.Name))

		if opt.Inputs != "" {
			ss, err := ren.ParseSettingStrings(opt.Inputs)
			if err != nil {
				return err
			}
			blockOpts = append(blockOpts, ren.WithBlockSettings(ss))
		}

		return ren.CreateBlock(ctx, blockOpts...)(ctx)
	}
}
