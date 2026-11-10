package shopifyv2

type RadioOption func(*Setting)

func WithRadioDefault(val string) RadioOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

// options is required. If no default is set, the first option is selected.
func NewRadioSetting(id, label string, options []Option, opts ...RadioOption) Setting {
	s := Setting{ID: id, Type: "radio", Label: label, Options: options}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
