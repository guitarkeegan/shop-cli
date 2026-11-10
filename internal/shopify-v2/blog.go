package shopifyv2

type BlogOption func(*Setting)

// Note: blog settings do not support the default attribute.
func NewBlogSetting(id, label string, opts ...BlogOption) Setting {
	s := Setting{ID: id, Type: "blog", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
