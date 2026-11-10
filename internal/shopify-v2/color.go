package shopifyv2

type ColorOption func(*Setting)

// WithColorDefault sets the default hex color value (e.g. "#000000").
func WithColorDefault(val string) ColorOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewColorSetting(id, label string, opts ...ColorOption) Setting {
	s := Setting{ID: id, Type: "color", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type ColorBgOption func(*Setting)

// WithColorBgDefault sets the default CSS background value.
func WithColorBgDefault(val string) ColorBgOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewColorBackgroundSetting(id, label string, opts ...ColorBgOption) Setting {
	s := Setting{ID: id, Type: "color_background", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type ColorSchemeOption func(*Setting)

// WithColorSchemeDefault sets the default color scheme ID.
func WithColorSchemeDefault(val string) ColorSchemeOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewColorSchemeSetting(id, label string, opts ...ColorSchemeOption) Setting {
	s := Setting{ID: id, Type: "color_scheme", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

func NewDefaultColorSchemeSetting() Setting {
	return NewColorSchemeSetting(
		"color_scheme",
		"t:sections.all.colors.label",
		WithColorSchemeDefault("Primary"),
	)
}

type ColorSchemeGrpOption func(*Setting)

// WithColorSchemeGrpDefault sets the default color scheme group ID.
func WithColorSchemeGrpDefault(val string) ColorSchemeGrpOption {
	return func(s *Setting) { s.Default = ptr(val) }
}

func NewColorSchemeGroupSetting(id, label string, opts ...ColorSchemeGrpOption) Setting {
	s := Setting{ID: id, Type: "color_scheme_group", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
