package shopifyv2

type CollectionOption func(*Setting)

// Note: collection settings do not support the default attribute.
func NewCollectionSetting(id, label string, opts ...CollectionOption) Setting {
	s := Setting{ID: id, Type: "collection", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type CollectionListOption func(*Setting)

// Defaults to 50 if not specified.
func WithCollectionListLimit(val int) CollectionListOption {
	return func(s *Setting) { s.Limit = ptr(val) }
}

func NewCollectionListSetting(id, label string, opts ...CollectionListOption) Setting {
	s := Setting{ID: id, Type: "collection_list", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
