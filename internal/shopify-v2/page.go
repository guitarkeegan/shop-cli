package shopifyv2

type PageOption func(*Setting)

// Note: page settings do not support the default attribute.
func NewPageSetting(id, label string, opts ...PageOption) Setting {
	s := Setting{ID: id, Type: "page", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
