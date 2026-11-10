package shopifyv2

type TextAlignmentOption func(*Setting)

// WithTextAlignmentDefault sets the default alignment ("left", "center", "right").
func WithTextAlignmentDefault(val string) TextAlignmentOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewTextAlignmentSetting(id, label string, opts ...TextAlignmentOption) Setting {
	s := Setting{ID: id, Type: "text_alignment", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
