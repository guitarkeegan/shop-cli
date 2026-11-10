package main

import (
	"context"
	"fmt"
	"os"

	"github.com/daved/clic"
	"github.com/guitarkeegan/shop-cli/internal/render"
	"github.com/guitarkeegan/shop-cli/internal/templatepicker"
	"github.com/guitarkeegan/shop-cli/internal/theme"
)

type snippetCreateOptions struct {
	Name        string
	Params      string
	Template    bool
	Description string
	Example     string
}

func newSnippetCmd() *clic.Clic {
	opts := &snippetCreateOptions{}
	create := clic.NewFromFunc(createSnippet(opts), "create")
	create.Flag(&opts.Name, "n|name", "Name of the snippet.")
	create.Flag(&opts.Description, "d|description", "Description of the snippet.")
	create.Flag(&opts.Example, "e|example", "Example of rendering the snippet.")
	create.Flag(&opts.Template, "t|template", "Pick a .liquid file in this directory to use as the snippet body.")
	create.Operand(&opts.Params, false, "params", "A comma separated list of doc params.\nBasic: title,copy,image\nWith types: title:string,copy:string,image:object\nOptional: title?:string,copy?,image?:object")

	return clic.NewFromFunc(requireSubcommand, "snippet", create)
}

func createSnippet(opts *snippetCreateOptions) clic.HandlerFunc {
	return func(ctx context.Context) error {
		var body string
		if opts.Template {
			var err error
			if body, err = pickTemplate(); err != nil {
				return err
			}
		}

		content := render.Snippet(
			render.WithSnippetDescription(opts.Description),
			render.WithSnippetParams(render.ParseSnippetParams(opts.Params)),
			render.WithSnippetExample(opts.Example),
			render.WithSnippetBody(body),
		)

		path, err := theme.WriteSnippet(opts.Name, content)
		if err != nil {
			return err
		}
		printCreated(path)
		return nil
	}
}

// pickTemplate lets the user choose a .liquid file and returns its
// contents, or "" if they quit without choosing.
func pickTemplate() (string, error) {
	path, err := templatepicker.Run()
	if err != nil || path == "" {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read template: %w", err)
	}
	return string(data), nil
}
