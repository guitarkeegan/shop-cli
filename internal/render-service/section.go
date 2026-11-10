package renderservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
	"github.com/guitarkeegan/shop-cli/internal/utils"
)

type SectionService struct {
	Section     shop2.Section
	UseTemplate bool `json:"-"`
}

func NewSectionService() *SectionService {
	return &SectionService{}
}

// Considering calling this render rather than String()
// since it will serve maybe a slightly different purpose.
func (ss SectionService) Render() string {
	sec := strings.Builder{}

	if len(ss.Section.Schema.Presets) == 0 && ss.Section.Schema.Name != "" {
		ss.Section.Schema.Presets = []shop2.Preset{{Name: ss.Section.Schema.Name}}
	}

	sec.WriteString("{% liquid\n")
	// vars, there needs to be a nauty list here
	for _, setting := range ss.Section.Schema.Settings {
		if setting.Type == "header" {
			continue
		}
		sec.WriteString("assign ")
		sec.WriteString(setting.ID)
		sec.WriteString(" = section.settings.")
		sec.WriteString(setting.ID)
		sec.WriteByte('\n')
	}
	sec.WriteString("%}\n")
	sec.WriteString(ss.RenderBlockLoop())
	bites, err := json.MarshalIndent(ss.Section.Schema, "", "  ")
	if err != nil {
		dbg("%s", err)
		return "error on marshal"
	}

	sec.WriteString("{% schema %}\n")
	sec.WriteString(string(bites))
	sec.WriteString("\n{% endschema %}")

	return sec.String()
}

func (ss SectionService) RenderBlockLoop() string {
	if len(ss.Section.Schema.Blocks) == 0 {
		return ""
	}
	blockLoop := strings.Builder{}

	blockLoop.WriteString("{% for block in section.blocks %}\n")
	blockLoop.WriteString("  {% case block.type %}\n")

	for _, block := range ss.Section.Schema.Blocks {
		blockLoop.WriteString("    {% when '")
		blockLoop.WriteString(block.Type)
		blockLoop.WriteString("' %}\n")
		// make liquid
		blockLoop.WriteString("{% liquid\n")
		for _, blockSetting := range block.Settings {
			blockLoop.WriteString("assign ")
			blockLoop.WriteString(blockSetting.ID)
			blockLoop.WriteString(" = block.settings.")
			blockLoop.WriteString(blockSetting.ID)
			blockLoop.WriteString("\n")
		}
		blockLoop.WriteString("%}\n")
	}

	blockLoop.WriteString("  {% endcase %}\n")
	blockLoop.WriteString("{% endfor %}\n")

	return blockLoop.String()
}

type SectionServiceOption func(*SectionService)

func WithSectionName(name string) SectionServiceOption {
	return func(ss *SectionService) {
		ss.Section.Schema.Name = name
	}
}

func WithSectionSettings(settings []shop2.Setting) SectionServiceOption {
	return func(ss *SectionService) {
		ss.Section.Schema.Settings = settings
	}
}

func WithSectionTemplate(b bool) SectionServiceOption {
	return func(ss *SectionService) {
		ss.UseTemplate = b
	}
}

func WithSectionBlocks(blocks []shop2.SectionBlock) SectionServiceOption {
	return func(ss *SectionService) {
		ss.Section.Schema.Blocks = blocks
	}
}

func WithSectionPresets(presets []shop2.Preset) SectionServiceOption {
	return func(ss *SectionService) {
		ss.Section.Schema.Presets = presets
	}
}

func WithDefaults() SectionServiceOption {
	return func(ss *SectionService) {
		paddingInputs := shop2.NewSectionPaddingSettings()
		colorScheme := shop2.NewDefaultColorSchemeSetting()
		header := shop2.NewHeader("t:sections.general.theme.header")
		ss.Section.Schema.Settings = append(ss.Section.Schema.Settings, header, colorScheme)
		ss.Section.Schema.Settings = append(ss.Section.Schema.Settings, paddingInputs...)
	}
}

// ParseBlockStrings parses each -b flag value into a SectionBlock.
func ParseBlockStrings(blocks []string) ([]shop2.SectionBlock, error) {
	parsed := make([]shop2.SectionBlock, 0, len(blocks))
	for _, b := range blocks {
		block, err := parseBlockString(b)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, block)
	}
	return parsed, nil
}

// parseBlockString parses a block string in the format:
// "blockName,settingType:settingID:settingLabel,settingType:settingID:settingLabel"
// Returns a SectionBlock with the parsed name, type, and settings
func parseBlockString(blockStr string) (shop2.SectionBlock, error) {
	parts := strings.Split(blockStr, ",")

	blockName := strings.TrimSpace(parts[0])
	blockType := utils.Handlize(blockName)
	if blockType == "" {
		return shop2.SectionBlock{}, fmt.Errorf("%w %q: expected \"name,type:id:label,...\"", ErrInvalidBlock, blockStr)
	}

	settings := make([]shop2.Setting, 0, len(parts)-1)
	for _, part := range parts[1:] {
		settingStr := strings.TrimSpace(part)
		if settingStr == "" {
			continue
		}

		setting, err := ParseSettingString(settingStr)
		if err != nil {
			return shop2.SectionBlock{}, fmt.Errorf("block %q: %w", blockName, err)
		}
		settings = append(settings, setting)
	}

	return shop2.SectionBlock{
		Name:     blockName,
		Type:     blockType,
		Settings: settings,
	}, nil
}

// ParseSettingStrings parses a comma separated list of "type:id:label" settings.
// An empty string is valid and returns no settings.
func ParseSettingStrings(settingStrs string) ([]shop2.Setting, error) {
	var parsed []shop2.Setting
	for _, settingStr := range strings.Split(settingStrs, ",") {
		if strings.TrimSpace(settingStr) == "" {
			continue
		}
		setting, err := ParseSettingString(settingStr)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, setting)
	}
	return parsed, nil
}

// ParseSettingString parses a setting string in the format "type:id:label".
// The label may itself contain colons.
func ParseSettingString(settingStr string) (shop2.Setting, error) {
	parts := strings.SplitN(settingStr, ":", 3)
	if len(parts) < 3 {
		return shop2.Setting{}, fmt.Errorf("%w %q: expected type:id:label", ErrInvalidSetting, strings.TrimSpace(settingStr))
	}

	settingType := strings.TrimSpace(parts[0])
	settingID := strings.TrimSpace(parts[1])
	settingLabel := strings.TrimSpace(parts[2])

	setting, err := shop2.NewSettingFromType(settingType, settingID, settingLabel)
	if err != nil {
		return shop2.Setting{}, fmt.Errorf("%w: %w", ErrInvalidSetting, err)
	}
	return setting, nil
}

func (ss *SectionService) CreateLiquidAssignments() string {
	output := strings.Builder{}

	for _, setting := range ss.Section.Schema.Settings {
		if shop2.IsValidSettingType(setting.Type) {
			output.WriteString("assign ")
			output.WriteString(setting.ID)
			output.WriteString(" = section.settings.")
			output.WriteString(setting.ID)
			output.WriteByte('\n')
		}
	}

	return strings.TrimSpace(output.String())
}

// CreateSection renders a section and writes it to sections/<handle>.liquid.
// It returns the path of the created file.
func CreateSection(ctx context.Context, opts ...SectionServiceOption) (string, error) {
	s := NewSectionService()
	for _, op := range opts {
		op(s)
	}

	if s.Section.Schema.Name == "" {
		return "", ErrMissingName
	}
	if !utils.FolderExistsInPwd("sections") {
		return "", ErrMissingSectionsDir
	}

	return writeLiquid("sections", s.Section.Schema.Name, s.Render())
}

// CreateSectionPreview renders a section without writing it. Use this within a TUI.
func CreateSectionPreview(ctx context.Context, opts ...SectionServiceOption) (string, error) {
	s := NewSectionService()
	for _, op := range opts {
		op(s)
	}

	if s.Section.Schema.Name == "" {
		return "", ErrMissingName
	}

	return s.Render(), nil
}
