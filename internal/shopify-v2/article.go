package shopifyv2

type ArticleOption func(*Setting)

// Note: article settings do not support the default attribute.
func NewArticleSetting(id, label string, opts ...ArticleOption) Setting {
	s := Setting{ID: id, Type: "article", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

type ArticleListOption func(*Setting)

// Defaults to 50 if not specified.
func WithArticleListLimit(val int) ArticleListOption {
	return func(s *Setting) { s.Limit = ptr(val) }
}

func NewArticleListSetting(id, label string, opts ...ArticleListOption) Setting {
	s := Setting{ID: id, Type: "article_list", Label: label}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}
