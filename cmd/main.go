package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/daved/clic"
	templates "github.com/guitarkeegan/shop-cli/assets"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	featureservice "github.com/guitarkeegan/shop-cli/internal/feature-service"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
)

var dbg = func() func(format string, as ...any) {
	if os.Getenv("DEBUG") == "" {
		return func(format string, as ...any) {}
	}
	return func(format string, as ...any) {
		fmt.Printf(format+"\n", as...)
	}
}()

func main() {
	featS := featureservice.New(templates.Templates)

	SectionCreateOptions := NewSectionCreateOptions(WithFeatureService(featS))
	create := clic.NewFromFunc(handleCreateSection(SectionCreateOptions), "create")
	create.Flag(&SectionCreateOptions.Name, "n|name", "Pass the name of the section.")
	create.Flag(&SectionCreateOptions.useTemplates, "t|template", "Start the interactive template picker")
	create.Flag(&SectionCreateOptions.Blocks, "b|blocks", "Add block with settings info (repeatable). Ex.\nshop section create -b \"blockName,settingType:settingID:settingLabel,settingType:settingID:settingLabel\".\nBlock type will be derived from the blockName.")
	create.Operand(&SectionCreateOptions.Settings, false, "inputs", "A comma separated list of settings. [type:id:label] Ex. \"text:title:Title,image_picker:image:Background Image\"")

	createBlockOptions := NewCreateBlockOptions()
	createBlock := clic.NewFromFunc(handleCreateBlock(createBlockOptions), "create")
	createBlock.Flag(&createBlockOptions.Name, "n|name", "Name of the block and file")
	createBlock.Operand(&createBlockOptions.Inputs, false, "inputs", "A comma separated list of settings. [type:id:label].")

	SnippetCreateOptions := NewSnippetCreateOptions()
	createSnippet := clic.NewFromFunc(handleCreateSnippet(SnippetCreateOptions), "create")
	createSnippet.Flag(&SnippetCreateOptions.Name, "n|name", "Name of the snippet file")
	createSnippet.Flag(&SnippetCreateOptions.Description, "d|description", "Give a description of the snippet")
	createSnippet.Flag(&SnippetCreateOptions.Example, "e|example", "Give an example implementation")
	createSnippet.Flag(&SnippetCreateOptions.Template, "t|template", "Start the interactive template picker")
	createSnippet.Operand(&SnippetCreateOptions.DocParams, false, "params", "A comma separated list of doc params.\nBasic Example title,copy,image\nExample with types: title:string,copy:string,image:object\nExample with optional params: title?:string,copy?,image?:object")

	examplesOptions := NewExampleOptions()
	useExample := clic.NewFromFunc(handleExamples(examplesOptions), "pull")
	useExample.Flag(&examplesOptions.Force, "f|force", "force download of examples, override hashes")

	section := clic.NewFromFunc(sectionHandler(), "section", create)
	block := clic.NewFromFunc(handleBlock, "block", createBlock)
	snippet := clic.NewFromFunc(snippetHandler(), "snippet", createSnippet)
	examples := clic.NewFromFunc(examplesHandler(), "examples", useExample)

	root := clic.NewFromFunc(rootHandler, "shop", section, block, snippet, examples)

	cmd, err := root.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, colors.Error(err.Error()))
		fmt.Fprintln(os.Stderr, cmd.Usage())
		os.Exit(2)
	}

	if err := cmd.Handle(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, colors.Error(err.Error()))
		if isUsageError(err) {
			fmt.Fprintln(os.Stderr, cmd.Usage())
		}
		os.Exit(1)
	}
}

// isUsageError reports whether err came from bad input, so usage help
// is worth showing alongside it.
func isUsageError(err error) bool {
	for _, target := range []error{
		ErrMissingSubcommand,
		ren.ErrMissingName,
		ren.ErrInvalidName,
		ren.ErrInvalidSetting,
		ren.ErrInvalidBlock,
	} {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
