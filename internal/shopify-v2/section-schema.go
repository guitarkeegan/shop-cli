package shopifyv2

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

type PageType string
type GroupType string

const (
	fourOFour              PageType = "404"
	article                PageType = "article"
	blog                   PageType = "blog"
	index                  PageType = "index"
	collection             PageType = "collection"
	product                PageType = "product"
	cart                   PageType = "cart"
	search                 PageType = "search"
	listCollections        PageType = "list-collections"
	page                   PageType = "page"
	password               PageType = "password"
	giftCard               PageType = "gift_card"
	customersLogin         PageType = "customers/login"
	customersAccount       PageType = "customers/account"
	customersOrder         PageType = "customers/order"
	customersRegister      PageType = "customers/register"
	customersResetPassword PageType = "customers/reset_password"
	customersActivateAcct  PageType = "customers/activate_account"
	wildcard               PageType = "*"
)

const (
	header GroupType = "header"
	footer GroupType = "footer"
	aside  GroupType = "aside"
	custom GroupType = "custom"
)

type PType struct {
	FourOFour              PageType
	Article                PageType
	Blog                   PageType
	Index                  PageType
	Collection             PageType
	Product                PageType
	Cart                   PageType
	Search                 PageType
	ListCollections        PageType
	Page                   PageType
	Password               PageType
	GiftCard               PageType
	CustomersLogin         PageType
	CustomersAccount       PageType
	CustomersOrder         PageType
	CustomersRegister      PageType
	CustomersResetPassword PageType
	CustomersActivateAcct  PageType
	Wildcard               PageType
}

type GType struct {
	Header GroupType
	Footer GroupType
	Aside  GroupType
	Custom GroupType
}

var PageTypes = PType{
	FourOFour:              fourOFour,
	Article:                article,
	Blog:                   blog,
	Index:                  index,
	Collection:             collection,
	Product:                product,
	Cart:                   cart,
	Search:                 search,
	ListCollections:        listCollections,
	Page:                   page,
	Password:               password,
	GiftCard:               giftCard,
	CustomersLogin:         customersLogin,
	CustomersAccount:       customersAccount,
	CustomersOrder:         customersOrder,
	CustomersRegister:      customersRegister,
	CustomersResetPassword: customersResetPassword,
	CustomersActivateAcct:  customersActivateAcct,
	Wildcard:               wildcard,
}

var GroupTypes = GType{
	Header: header,
	Footer: footer,
	Aside:  aside,
	Custom: custom,
}

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
