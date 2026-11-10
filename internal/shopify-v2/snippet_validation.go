package shopifyv2

import (
	"fmt"
	"strings"
)

// validSnippetInputTypes maps setting type strings that are valid as snippet parameter types.
var validSnippetInputTypes = NewSet(
	"checkbox",
	"number",
	"radio",
	"range",
	"select",
	"text",
	"textarea",
	"article",
	"blog",
	"collection",
	"collection_list",
	"color",
	"color_background",
	"color_scheme",
	"color_scheme_group",
	"font_picker",
	"html",
	"image_picker",
	"image",
	"inline_richtext",
	"link_list",
	"liquid",
	"metaobject",
	"metaobject_list",
	"page",
	"product",
	"product_list",
	"richtext",
	"text_alignment",
	"url",
	"video",
	"video_url",
)

// IsValidSnippetParamType checks if a string is a valid snippet parameter type.
func IsValidSnippetParamType(s string) bool {
	return validSnippetInputTypes.Has(s)
}

// ParsedSnippetParam represents a parsed/validated snippet param from CLI input.
type ParsedSnippetParam struct {
	Value       string
	Optional    bool
	IsValidType bool
}

func (p ParsedSnippetParam) String() string {
	return fmt.Sprintf("Value: %s, Optional: %t, IsValidType: %t\n", p.Value, p.Optional, p.IsValidType)
}

// ExtractSnippetParams parses a comma-separated string of snippet params.
// Each param is validated and marked as optional if it ends with '?'.
func ExtractSnippetParams(s string) []ParsedSnippetParam {
	paramSeq := strings.SplitSeq(s, ",")
	params := make([]ParsedSnippetParam, 0, 20)

	for param := range paramSeq {
		trimmed := strings.TrimSpace(param)
		if len(trimmed) == 0 {
			continue
		}
		switch {
		case IsValidSnippetParamType(trimmed):
			params = append(params, ParsedSnippetParam{
				Value:       trimmed,
				Optional:    false,
				IsValidType: true,
			})
		case len(trimmed) > 1 && IsValidSnippetParamType(trimmed[:len(trimmed)-1]) && trimmed[len(trimmed)-1] == '?':
			params = append(params, ParsedSnippetParam{
				Value:       trimmed,
				Optional:    true,
				IsValidType: true,
			})
		default:
			params = append(params, ParsedSnippetParam{
				Value:       trimmed,
				Optional:    IsOptionalParam(trimmed),
				IsValidType: false,
			})
		}
	}

	return params
}

// IsOptionalParam checks if a param string ends with '?'.
func IsOptionalParam(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s)-1] == '?'
}
