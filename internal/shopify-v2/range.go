package shopifyv2

// TODO: concider adding semantic helpers like: NewPaddingRangeSetting, or NewFontRangeSetting

import "strconv"

type RangeOption func(*Setting)

func WithRangeDefault(val int) RangeOption {
	return func(s *Setting) { s.Default = ptr(strconv.Itoa(val)) }
}

func WithUnit(val string) RangeOption {
	return func(s *Setting) { s.Unit = ptr(val) }
}

// min, max, and step are required per Shopify docs.
func NewRangeSetting(id, label string, min, max, step int, opts ...RangeOption) Setting {
	s := Setting{
		ID:    id,
		Type:  "range",
		Label: label,
		Min:   ptr(min),
		Max:   ptr(max),
		Step:  ptr(step),
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

func newSectionPaddingTop() Setting {
	return NewRangeSetting("padding_top", "Top padding", 0, 120, 4,
		WithUnit("px"),
		WithRangeDefault(36),
	)
}

func newSectionPaddingBottom() Setting {
	return NewRangeSetting("padding_bottom", "Bottom padding", 0, 120, 4,
		WithUnit("px"),
		WithRangeDefault(36),
	)
}

// make default padding inputs
func NewSectionPaddingSettings() []Setting {
	return []Setting{
		newSectionPaddingTop(),
		newSectionPaddingBottom(),
	}
}
