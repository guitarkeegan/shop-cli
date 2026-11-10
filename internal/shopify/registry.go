package shopify

import "fmt"

// allSettingTypes lists every valid Shopify setting type string.
var allSettingTypes = []string{
	"article",
	"article_list",
	"blog",
	"checkbox",
	"collection",
	"collection_list",
	"color",
	"color_background",
	"color_scheme",
	"color_scheme_group",
	"font_picker",
	"header",
	"html",
	"image_picker",
	"inline_richtext",
	"link_list",
	"liquid",
	"metaobject",
	"metaobject_list",
	"number",
	"page",
	"paragraph",
	"product",
	"product_list",
	"radio",
	"range",
	"richtext",
	"select",
	"text",
	"text_alignment",
	"textarea",
	"url",
	"video",
	"video_url",
}

var validSettingTypes = func() map[string]bool {
	m := make(map[string]bool, len(allSettingTypes))
	for _, t := range allSettingTypes {
		m[t] = true
	}
	return m
}()

// IsValidSettingType reports whether t is a Shopify setting type.
func IsValidSettingType(t string) bool {
	return validSettingTypes[t]
}

// NewSettingFromType creates a Setting with sensible defaults from a type string, id, and label.
// It returns an error for unknown types.
func NewSettingFromType(typeName, id, label string) (Setting, error) {
	if !IsValidSettingType(typeName) {
		return Setting{}, fmt.Errorf("unknown setting type %q", typeName)
	}

	s := Setting{
		ID:    id,
		Type:  typeName,
		Label: label,
	}

	switch typeName {
	case "range":
		s.Min = ptr(0)
		s.Max = ptr(100)
		s.Step = ptr(1)
	case "header":
		s.ID = ""
		s.Label = ""
		s.Content = ptr(label)
	case "paragraph":
		s.ID = ""
		s.Label = ""
		s.Content = ptr(label)
	}

	return s, nil
}
