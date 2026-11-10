// Package render turns Shopify theme file definitions into Liquid source.
// It has no side effects; see package theme for writing files.
package render

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

// writeAssigns writes a {% liquid %} tag that assigns each setting to a
// local variable, e.g. "assign title = section.settings.title". Settings
// without an id, such as headers, are skipped.
func writeAssigns(b *strings.Builder, scope string, settings []shopify.Setting) {
	b.WriteString("{% liquid\n")
	for _, s := range settings {
		if s.ID == "" {
			continue
		}
		fmt.Fprintf(b, "assign %s = %s.settings.%s\n", s.ID, scope, s.ID)
	}
	b.WriteString("%}\n")
}

// writeSchema writes the {% schema %} tag containing v as indented JSON.
func writeSchema(b *strings.Builder, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	b.WriteString("{% schema %}\n")
	b.Write(data)
	b.WriteString("\n{% endschema %}\n")
	return nil
}

// defaultPresets gives a schema a preset named after it, so the theme
// editor lists it, unless presets were set explicitly.
func defaultPresets(name string, presets []shopify.Preset) []shopify.Preset {
	if len(presets) == 0 && name != "" {
		return []shopify.Preset{{Name: name}}
	}
	return presets
}
