package shopifyv2

import (
	"encoding/json"
	"fmt"
)

// ThemeBlockReference represents a block type reference in a theme block's blocks array.
// Unlike section blocks, theme blocks cannot define blocks inline — they can only
// reference other blocks by type (e.g. "@app", "@theme", or a specific block type name).
type ThemeBlockReference struct {
	Type string `json:"type"`
}

// ThemeBlockSchema represents the schema for a Shopify theme block file.
// Theme blocks are standalone block files that differ from section-level blocks:
//   - They cannot define local/inline blocks (only reference other blocks by type)
//   - They support nesting up to 8 levels deep (excluding the section level)
//   - They use {% content_for 'blocks' %} to render nested child blocks
type ThemeBlockSchema struct {
	Name     string                `json:"name"`
	Tag      *string               `json:"tag,omitempty"`
	Class    *string               `json:"class,omitempty"`
	Settings []Setting             `json:"settings,omitempty"`
	Blocks   []ThemeBlockReference `json:"blocks,omitempty"`
	Presets  []Preset              `json:"presets,omitempty"`
}

func (tbs ThemeBlockSchema) String() string {
	b, err := json.MarshalIndent(tbs, "", " ")
	if err != nil {
		return fmt.Sprintf("error marshalling string: %s\n", err)
	}
	return string(b)
}

// ThemeBlock represents a rendered Shopify theme block file with its Liquid content and schema.
type ThemeBlock struct {
	Content string
	Schema  ThemeBlockSchema
}
