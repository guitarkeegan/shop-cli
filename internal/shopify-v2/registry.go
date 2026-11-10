package shopifyv2

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

var ValidSettingTypes = NewSet(allSettingTypes...)

// this is a bit silly
func IsValidSettingType(t string) bool {
	return ValidSettingTypes.Has(t)
}

// NewSettingFromType creates a Setting with sensible defaults from a type string, id, and label.
// Returns the setting and true if the type is valid, or a zero Setting and false otherwise.
func NewSettingFromType(typeName, id, label string) (Setting, error) {
	if !IsValidSettingType(typeName) {
		return Setting{}, fmt.Errorf("setting type in invalid. setting: %s\n", typeName)
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
