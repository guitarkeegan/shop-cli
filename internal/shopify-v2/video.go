package shopifyv2

type VideoOption func(*Setting)

func NewVideoSetting(id, label string, opts ...VideoOption) Setting {
	s := Setting{ID: id, Type: "video", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
