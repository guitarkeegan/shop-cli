package shopifyv2

type VideoURLOption func(*Setting)

// accept specifies which video providers are allowed (e.g. "youtube", "vimeo").
func NewVideoURLSetting(id, label string, accept []string, opts ...VideoURLOption) Setting {
	s := Setting{ID: id, Type: "video_url", Label: label, Accept: accept}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
