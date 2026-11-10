package shopifyv2

type RichTextOption func(*Setting)

// WithRichTextDefault sets the default HTML content for a richtext setting.
func WithRichTextDefault(val string) RichTextOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewRichTextSetting(id, label string, opts ...RichTextOption) Setting {
	s := Setting{ID: id, Type: "richtext", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type InlineRichTextOption func(*Setting)

// WithInlineRichTextDefault sets the default HTML content for an inline richtext setting.
func WithInlineRichTextDefault(val string) InlineRichTextOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewInlineRichTextSetting(id, label string, opts ...InlineRichTextOption) Setting {
	s := Setting{ID: id, Type: "inline_richtext", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
