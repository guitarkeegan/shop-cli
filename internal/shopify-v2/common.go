package shopifyv2

type CommonOption func(*Setting)

// WithInfo sets informational text displayed below the setting.
// Supports the "t:" prefix for translations.
func WithInfo(val string) CommonOption {
	return func(s *Setting) { s.Info = ptr(val) }
}

func (c CommonOption) AsText() TextOption                     { return TextOption(c) }
func (c CommonOption) AsTextarea() TextareaOption             { return TextareaOption(c) }
func (c CommonOption) AsNumber() NumberOption                 { return NumberOption(c) }
func (c CommonOption) AsRange() RangeOption                   { return RangeOption(c) }
func (c CommonOption) AsCheckbox() CheckboxOption             { return CheckboxOption(c) }
func (c CommonOption) AsRadio() RadioOption                   { return RadioOption(c) }
func (c CommonOption) AsSelect() SelectOption                 { return SelectOption(c) }
func (c CommonOption) AsColor() ColorOption                   { return ColorOption(c) }
func (c CommonOption) AsColorBackground() ColorBgOption       { return ColorBgOption(c) }
func (c CommonOption) AsColorScheme() ColorSchemeOption        { return ColorSchemeOption(c) }
func (c CommonOption) AsColorSchemeGroup() ColorSchemeGrpOption { return ColorSchemeGrpOption(c) }
func (c CommonOption) AsFontPicker() FontPickerOption         { return FontPickerOption(c) }
func (c CommonOption) AsImagePicker() ImagePickerOption       { return ImagePickerOption(c) }
func (c CommonOption) AsRichText() RichTextOption             { return RichTextOption(c) }
func (c CommonOption) AsInlineRichText() InlineRichTextOption { return InlineRichTextOption(c) }
func (c CommonOption) AsHTML() HTMLOption                     { return HTMLOption(c) }
func (c CommonOption) AsLiquid() LiquidOption                 { return LiquidOption(c) }
func (c CommonOption) AsURL() URLOption                       { return URLOption(c) }
func (c CommonOption) AsVideoURL() VideoURLOption             { return VideoURLOption(c) }
func (c CommonOption) AsVideo() VideoOption                   { return VideoOption(c) }
func (c CommonOption) AsLinkList() LinkListOption             { return LinkListOption(c) }
func (c CommonOption) AsCollection() CollectionOption         { return CollectionOption(c) }
func (c CommonOption) AsCollectionList() CollectionListOption { return CollectionListOption(c) }
func (c CommonOption) AsProduct() ProductOption               { return ProductOption(c) }
func (c CommonOption) AsProductList() ProductListOption       { return ProductListOption(c) }
func (c CommonOption) AsArticle() ArticleOption               { return ArticleOption(c) }
func (c CommonOption) AsArticleList() ArticleListOption       { return ArticleListOption(c) }
func (c CommonOption) AsBlog() BlogOption                     { return BlogOption(c) }
func (c CommonOption) AsPage() PageOption                     { return PageOption(c) }
func (c CommonOption) AsTextAlignment() TextAlignmentOption   { return TextAlignmentOption(c) }
func (c CommonOption) AsMetaObject() MetaObjectOption         { return MetaObjectOption(c) }
func (c CommonOption) AsMetaObjectList() MetaObjectListOption { return MetaObjectListOption(c) }
