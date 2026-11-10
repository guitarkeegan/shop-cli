package render

import (
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

// BlockOption configures the schema of a rendered theme block.
type BlockOption func(*shopify.ThemeBlockSchema)

func WithBlockName(name string) BlockOption {
	return func(s *shopify.ThemeBlockSchema) { s.Name = name }
}

func WithBlockSettings(settings []shopify.Setting) BlockOption {
	return func(s *shopify.ThemeBlockSchema) { s.Settings = settings }
}

// Block renders a theme block file: variable assignments for each setting,
// a content_for tag if it accepts nested blocks, and the schema.
func Block(opts ...BlockOption) (string, error) {
	var schema shopify.ThemeBlockSchema
	for _, opt := range opts {
		opt(&schema)
	}
	schema.Presets = defaultPresets(schema.Name, schema.Presets)

	var b strings.Builder
	writeAssigns(&b, "block", schema.Settings)
	if len(schema.Blocks) > 0 {
		b.WriteString("{% content_for 'blocks' %}\n")
	}
	if err := writeSchema(&b, schema); err != nil {
		return "", err
	}
	return b.String(), nil
}
