package shopifyv2

type LinkListOption func(*Setting)

func WithLinkListDefault(val string) LinkListOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewLinkListSetting(id, label string, opts ...LinkListOption) Setting {
	s := Setting{ID: id, Type: "link_list", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
