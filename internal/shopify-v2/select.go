package shopifyv2

type SelectOption func(*Setting)

func WithSelectDefault(val string) SelectOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

// Renders as segmented control for <=5 options without groups,
// otherwise renders as a dropdown.
// options is required. If no default is set, the first option is selected.
func NewSelectSetting(id, label string, options []Option, opts ...SelectOption) Setting {
	s := Setting{ID: id, Type: "select", Label: label, Options: options}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

func GroupedOption(value, label, group string) Option {
	return Option{Value: value, Label: label, Group: ptr(group)}
}
