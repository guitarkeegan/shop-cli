package shopify

import "strconv"

// Defaults added to every section generated from the CLI.

// NewHeader creates a header pseudo-setting for the theme editor sidebar.
// Headers display static text and do not have an id or label.
func NewHeader(content string) Setting {
	return Setting{Type: "header", Content: ptr(content)}
}

// NewDefaultColorSchemeSetting is the standard color scheme picker.
func NewDefaultColorSchemeSetting() Setting {
	return Setting{
		ID:      "color_scheme",
		Type:    "color_scheme",
		Label:   "t:sections.all.colors.label",
		Default: ptr("Primary"),
	}
}

// NewSectionPaddingSettings returns top and bottom padding range sliders.
func NewSectionPaddingSettings() []Setting {
	return []Setting{
		newPaddingSetting("padding_top", "Top padding"),
		newPaddingSetting("padding_bottom", "Bottom padding"),
	}
}

func newPaddingSetting(id, label string) Setting {
	return Setting{
		ID:      id,
		Type:    "range",
		Label:   label,
		Min:     ptr(0),
		Max:     ptr(120),
		Step:    ptr(4),
		Unit:    ptr("px"),
		Default: ptr(strconv.Itoa(36)),
	}
}
