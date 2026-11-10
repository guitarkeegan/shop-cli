package shopifyv2

import "strconv"

type CheckboxOption func(*Setting)

func WithCheckboxDefault(val bool) CheckboxOption {
	return func(s *Setting) { s.Default = ptr(strconv.FormatBool(val)) }
}

// Default is false if not specified.
func NewCheckboxSetting(id, label string, opts ...CheckboxOption) Setting {
	s := Setting{ID: id, Type: "checkbox", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
