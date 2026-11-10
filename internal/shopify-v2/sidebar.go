package shopifyv2

// NewHeader creates a header pseudo-setting for the theme editor sidebar.
// Headers display static text and do not have an id or label.
func NewHeader(content string) Setting {
	return Setting{Type: "header", Content: ptr(content)}
}

// NewParagraph creates a paragraph pseudo-setting for the theme editor sidebar.
// Paragraphs display static text and do not have an id or label.
func NewParagraph(content string) Setting {
	return Setting{Type: "paragraph", Content: ptr(content)}
}
