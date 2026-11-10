package shopifyv2

type HTMLOption func(*Setting)

func WithHTMLDefault(val string) HTMLOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewHTMLSetting(id, label string, opts ...HTMLOption) Setting {
	s := Setting{ID: id, Type: "html", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
