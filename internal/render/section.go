package render

import (
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

// SectionOption configures the schema of a rendered section.
type SectionOption func(*shopify.SectionSchema)

func WithSectionName(name string) SectionOption {
	return func(s *shopify.SectionSchema) { s.Name = name }
}

func WithSectionSettings(settings []shopify.Setting) SectionOption {
	return func(s *shopify.SectionSchema) { s.Settings = settings }
}

func WithSectionBlocks(blocks []shopify.SectionBlock) SectionOption {
	return func(s *shopify.SectionSchema) { s.Blocks = blocks }
}

// WithDefaults appends the settings most sections need: a theme header,
// a color scheme picker, and top and bottom padding.
func WithDefaults() SectionOption {
	return func(s *shopify.SectionSchema) {
		s.Settings = append(s.Settings,
			shopify.NewHeader("t:sections.general.theme.header"),
			shopify.NewDefaultColorSchemeSetting(),
		)
		s.Settings = append(s.Settings, shopify.NewSectionPaddingSettings()...)
	}
}

// Section renders a section file: variable assignments for each setting, a
// loop over any blocks, and the schema.
func Section(opts ...SectionOption) (string, error) {
	var schema shopify.SectionSchema
	for _, opt := range opts {
		opt(&schema)
	}
	schema.Presets = defaultPresets(schema.Name, schema.Presets)

	var b strings.Builder
	writeAssigns(&b, "section", schema.Settings)
	writeBlockLoop(&b, schema.Blocks)
	if err := writeSchema(&b, schema); err != nil {
		return "", err
	}
	return b.String(), nil
}

// writeBlockLoop writes a loop over section.blocks with a case for each
// block type that assigns its settings.
func writeBlockLoop(b *strings.Builder, blocks []shopify.SectionBlock) {
	if len(blocks) == 0 {
		return
	}

	b.WriteString("{% for block in section.blocks %}\n")
	b.WriteString("  {% case block.type %}\n")
	for _, block := range blocks {
		b.WriteString("    {% when '" + block.Type + "' %}\n")
		writeAssigns(b, "block", block.Settings)
	}
	b.WriteString("  {% endcase %}\n")
	b.WriteString("{% endfor %}\n")
}
