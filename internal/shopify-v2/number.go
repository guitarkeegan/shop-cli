package shopifyv2

import "strconv"

type NumberOption func(*Setting)

func WithNumberDefault(val int) NumberOption {
	return func(s *Setting) { s.Default = ptr(strconv.Itoa(val)) }
}

func WithNumberPlaceholder(val string) NumberOption {
	return func(s *Setting) { s.Placeholder = ptr(val) }
}

func NewNumberSetting(id, label string, opts ...NumberOption) Setting {
	s := Setting{ID: id, Type: "number", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
