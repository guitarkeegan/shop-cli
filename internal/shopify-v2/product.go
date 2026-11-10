package shopifyv2

type ProductOption func(*Setting)

// Note: product settings do not support the default attribute.
func NewProductSetting(id, label string, opts ...ProductOption) Setting {
	s := Setting{ID: id, Type: "product", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type ProductListOption func(*Setting)

// Defaults to 50 if not specified.
func WithProductListLimit(val int) ProductListOption {
	return func(s *Setting) { s.Limit = ptr(val) }
}

func NewProductListSetting(id, label string, opts ...ProductListOption) Setting {
	s := Setting{ID: id, Type: "product_list", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
