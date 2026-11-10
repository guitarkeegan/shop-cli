package shopifyv2

type ImagePickerOption func(*Setting)

func NewImagePickerSetting(id, label string, opts ...ImagePickerOption) Setting {
	s := Setting{ID: id, Type: "image_picker", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
