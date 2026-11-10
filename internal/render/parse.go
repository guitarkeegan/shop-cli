package render

import (
	"errors"
	"fmt"
	"strings"

	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

// Parse errors. Callers can match these with errors.Is to show usage help.
var (
	ErrInvalidSetting = errors.New("invalid setting")
	ErrInvalidBlock   = errors.New("invalid block")
)

// ParseSettings parses a comma separated list of "type:id:label" settings.
// An empty string is valid and returns no settings.
func ParseSettings(list string) ([]shopify.Setting, error) {
	var settings []shopify.Setting
	for _, item := range strings.Split(list, ",") {
		if strings.TrimSpace(item) == "" {
			continue
		}
		setting, err := ParseSetting(item)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

// ParseSetting parses a setting in the format "type:id:label".
// The label may itself contain colons.
func ParseSetting(s string) (shopify.Setting, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 3 {
		return shopify.Setting{}, fmt.Errorf("%w %q: expected type:id:label", ErrInvalidSetting, strings.TrimSpace(s))
	}

	setting, err := shopify.NewSettingFromType(
		strings.TrimSpace(parts[0]),
		strings.TrimSpace(parts[1]),
		strings.TrimSpace(parts[2]),
	)
	if err != nil {
		return shopify.Setting{}, fmt.Errorf("%w: %w", ErrInvalidSetting, err)
	}
	return setting, nil
}

// ParseBlocks parses each block definition in the format
// "Block Name,type:id:label,type:id:label". The block type is derived from
// the name.
func ParseBlocks(defs []string) ([]shopify.SectionBlock, error) {
	blocks := make([]shopify.SectionBlock, 0, len(defs))
	for _, def := range defs {
		block, err := parseBlock(def)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func parseBlock(def string) (shopify.SectionBlock, error) {
	name, settingList, _ := strings.Cut(def, ",")
	name = strings.TrimSpace(name)

	blockType := shopify.Handle(name)
	if blockType == "" {
		return shopify.SectionBlock{}, fmt.Errorf("%w %q: expected \"name,type:id:label,...\"", ErrInvalidBlock, def)
	}

	settings, err := ParseSettings(settingList)
	if err != nil {
		return shopify.SectionBlock{}, fmt.Errorf("block %q: %w", name, err)
	}

	return shopify.SectionBlock{Name: name, Type: blockType, Settings: settings}, nil
}

// ParseSnippetParams parses params like "title:string,image?:object,copy".
// A trailing "?" on the name marks a param optional; the type defaults to
// string.
func ParseSnippetParams(list string) []shopify.Param {
	var params []shopify.Param
	for item := range strings.SplitSeq(list, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		name, typ, _ := strings.Cut(item, ":")
		name = strings.TrimSpace(name)
		optional := strings.HasSuffix(name, "?")
		params = append(params, shopify.Param{
			Type:     shopify.ParamTypeFromString(strings.TrimSpace(typ)),
			Name:     strings.TrimSuffix(name, "?"),
			Required: !optional,
		})
	}
	return params
}
