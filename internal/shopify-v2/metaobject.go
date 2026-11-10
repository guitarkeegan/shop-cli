package shopifyv2

type MetaObjectOption func(*Setting)

func NewMetaObjectSetting(id, label string, opts ...MetaObjectOption) Setting {
	s := Setting{ID: id, Type: "metaobject", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type MetaObjectListOption func(*Setting)

// Defaults to 50 if not specified.
func WithMetaObjectListLimit(val int) MetaObjectListOption {
	return func(s *Setting) { s.Limit = ptr(val) }
}

func NewMetaObjectListSetting(id, label string, opts ...MetaObjectListOption) Setting {
	s := Setting{ID: id, Type: "metaobject_list", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
