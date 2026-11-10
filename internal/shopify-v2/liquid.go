package shopifyv2

type LiquidOption func(*Setting)

func WithLiquidDefault(val string) LiquidOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewLiquidSetting(id, label string, opts ...LiquidOption) Setting {
	s := Setting{ID: id, Type: "liquid", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
