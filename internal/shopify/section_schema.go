package shopify

import (
	"encoding/json"
	"fmt"
)

type SectionBlock struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Limit    *int      `json:"limit,omitempty"`
	Settings []Setting `json:"settings,omitempty"`
}

func (sb SectionBlock) String() string {
	json, err := json.MarshalIndent(sb, "", " ")
	if err != nil {
		return fmt.Sprintf("error marshalling string: %s\n", err)
	}
	return string(json)
}

type PresetSectionBlockConfig struct {
	Type     string    `json:"type"`
	Settings []Setting `json:"settings,omitempty"`
	Name     *string   `json:"name,omitempty"`
}

type Preset struct {
	Name     string                     `json:"name"`
	Category *string                    `json:"category,omitempty"`
	Settings []Setting                  `json:"settings,omitempty"`
	Blocks   []PresetSectionBlockConfig `json:"blocks,omitempty"`
}

// PageType is a template a section can be enabled or disabled on.
type PageType string

// GroupType is a section group a section can be enabled or disabled in.
type GroupType string

const (
	PageType404                    PageType = "404"
	PageTypeArticle                PageType = "article"
	PageTypeBlog                   PageType = "blog"
	PageTypeIndex                  PageType = "index"
	PageTypeCollection             PageType = "collection"
	PageTypeProduct                PageType = "product"
	PageTypeCart                   PageType = "cart"
	PageTypeSearch                 PageType = "search"
	PageTypeListCollections        PageType = "list-collections"
	PageTypePage                   PageType = "page"
	PageTypePassword               PageType = "password"
	PageTypeGiftCard               PageType = "gift_card"
	PageTypeCustomersLogin         PageType = "customers/login"
	PageTypeCustomersAccount       PageType = "customers/account"
	PageTypeCustomersOrder         PageType = "customers/order"
	PageTypeCustomersRegister      PageType = "customers/register"
	PageTypeCustomersResetPassword PageType = "customers/reset_password"
	PageTypeCustomersActivateAcct  PageType = "customers/activate_account"
	PageTypeAll                    PageType = "*"
)

const (
	GroupTypeHeader GroupType = "header"
	GroupTypeFooter GroupType = "footer"
	GroupTypeAside  GroupType = "aside"
	GroupTypeCustom GroupType = "custom"
)

type EnabledOn struct {
	Templates []PageType  `json:"templates,omitempty"`
	Groups    []GroupType `json:"groups,omitempty"`
}

type DisabledOn struct {
	Templates []PageType  `json:"templates,omitempty"`
	Groups    []GroupType `json:"groups,omitempty"`
}

type SectionSchema struct {
	Name       string                       `json:"name"`
	Tag        *string                      `json:"tag,omitempty"`
	Class      *string                      `json:"class,omitempty"`
	Limit      *int                         `json:"limit,omitempty"`
	Settings   []Setting                    `json:"settings,omitempty"`
	MaxBlocks  *int                         `json:"max_blocks,omitempty"`
	Blocks     []SectionBlock               `json:"blocks,omitempty"`
	Presets    []Preset                     `json:"presets,omitempty"`
	Locales    map[string]map[string]string `json:"locales,omitempty"`
	EnabledOn  *EnabledOn                   `json:"enabled_on,omitempty"`
	DisabledOn *DisabledOn                  `json:"disabled_on,omitempty"`
}

type Section struct {
	Content string
	Schema  SectionSchema
}
