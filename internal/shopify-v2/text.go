package shopifyv2

type TextOption func(*Setting)

func WithTextDefault(val string) TextOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func WithPlaceholder(val string) TextOption {
	return func(s *Setting) { s.Placeholder = ptr(val) }
}

func NewTextSetting(id, label string, opts ...TextOption) Setting {
	s := Setting{ID: id, Type: "text", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
