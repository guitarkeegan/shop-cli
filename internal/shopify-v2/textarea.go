package shopifyv2

type TextareaOption func(*Setting)

func WithTextareaDefault(val string) TextareaOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func WithTextareaPlaceholder(val string) TextareaOption {
	return func(s *Setting) { s.Placeholder = ptr(val) }
}

func NewTextareaSetting(id, label string, opts ...TextareaOption) Setting {
	s := Setting{ID: id, Type: "textarea", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
