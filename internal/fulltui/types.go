package fulltui

import "github.com/guitarkeegan/shop-cli/internal/shopify"

// param is a snippet doc param as entered in the form.
type param struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

// snippetDoc collects the answers for a snippet's {% doc %} block.
type snippetDoc struct {
	Description string
	Example     string
	Param       []param
}

func (sd snippetDoc) shopifyParams() []shopify.Param {
	params := make([]shopify.Param, len(sd.Param))
	for i, p := range sd.Param {
		params[i] = shopify.Param{
			Type:        shopify.ParamTypeFromString(p.Type),
			Name:        p.Name,
			Required:    p.Required,
			Description: p.Description,
		}
	}
	return params
}
