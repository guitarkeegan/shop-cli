// Command shop generates Shopify theme sections, blocks, and snippets.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/daved/clic"
	templates "github.com/guitarkeegan/shop-cli/assets"
	"github.com/guitarkeegan/shop-cli/internal/colors"
	"github.com/guitarkeegan/shop-cli/internal/features"
	"github.com/guitarkeegan/shop-cli/internal/render"
	"github.com/guitarkeegan/shop-cli/internal/theme"
)

// ErrMissingSubcommand is returned by parent commands run without a subcommand.
var ErrMissingSubcommand = errors.New("missing subcommand")

func main() {
	var showVersion bool
	root := clic.NewFromFunc(runRoot(&showVersion), "shop",
		newSectionCmd(features.New(templates.Templates)),
		newBlockCmd(),
		newSnippetCmd(),
		newExamplesCmd(),
	)
	root.Flag(&showVersion, "v|version", "Print the version and exit.")

	var showHelp bool
	root.Recursively(func(c *clic.Clic) {
		c.Flag(&showHelp, "h|help", "Show help for this command.")
	})

	cmd, err := root.Parse(os.Args[1:])
	if showHelp {
		fmt.Println(cmd.Usage())
		return
	}
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

// requireSubcommand handles parent commands like "shop section".
func requireSubcommand(context.Context) error {
	return ErrMissingSubcommand
}

// isUsageError reports whether err came from bad input, so usage help
// is worth showing alongside it.
func isUsageError(err error) bool {
	for _, target := range []error{
		ErrMissingSubcommand,
		theme.ErrMissingName,
		theme.ErrInvalidName,
		render.ErrInvalidSetting,
		render.ErrInvalidBlock,
	} {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func printCreated(path string) {
	fmt.Println(colors.Success("Created " + path))
}
