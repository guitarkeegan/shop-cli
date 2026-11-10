package shopify

import (
	"encoding/json"
	"fmt"
)

type Setting struct {
	ID      string  `json:"id,omitempty"`
	Type    string  `json:"type"`
	Label   string  `json:"label,omitempty"`
	Default *string `json:"default,omitempty"`
	Info    *string `json:"info,omitempty"`
	Content *string `json:"content,omitempty"` // header/paragraph sidebar content

	Options     []Option `json:"options,omitempty"`     // radio/select
	Min         *int     `json:"min,omitempty"`         // range
	Max         *int     `json:"max,omitempty"`         // range
	Step        *int     `json:"step,omitempty"`        // range
	Unit        *string  `json:"unit,omitempty"`        // range
	Placeholder *string  `json:"placeholder,omitempty"` // text/textarea/number
	Accept      []string `json:"accept,omitempty"`      // video_url
	Limit       *int     `json:"limit,omitempty"`       // collection_list/product_list/article_list/metaobject_list
}

func (s Setting) String() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Sprintf("Setting{error: %v}", err)
	}
	return string(b)
}

type Option struct {
	Value string  `json:"value"`
	Label string  `json:"label"`
	Group *string `json:"group,omitempty"` // select option grouping
}

func ptr[T any](v T) *T {
	return &v
}
