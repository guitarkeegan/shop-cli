package shopifyv2

type URLOption func(*Setting)

func WithURLDefault(val string) URLOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewURLSetting(id, label string, opts ...URLOption) Setting {
	s := Setting{ID: id, Type: "url", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
