package shopifyv2

type FontPickerOption func(*Setting)

// WithFontPickerDefault sets the default font family (e.g. "helvetica_n4").
func WithFontPickerDefault(val string) FontPickerOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewFontPickerSetting(id, label string, opts ...FontPickerOption) Setting {
	s := Setting{ID: id, Type: "font_picker", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
