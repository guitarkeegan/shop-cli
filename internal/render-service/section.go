package renderservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
	"github.com/guitarkeegan/shop-cli/internal/utils"
)

var (
	ErrNoSettingString = errors.New("settings string required")
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

func WithSectionBlocksString(blocks []string) SectionServiceOption {
	return func(ss *SectionService) {
		// parse each block
		bs := make([]shop2.SectionBlock, 0, len(blocks))
		for _, blockStr := range blocks {
			parsedBlock, err := parseBlockString(blockStr)
			if err != nil {
				log.Fatalf("unparsable block string: %v\n", err)
			}
			bs = append(bs, parsedBlock)
		}
		ss.Section.Schema.Blocks = bs
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

// parseBlockString parses a block string in the format:
// "blockName,settingType:settingID:settingLabel,settingType:settingID:settingLabel"
// Returns a SectionBlock with the parsed name, type, and settings
func parseBlockString(blockStr string) (shop2.SectionBlock, error) {
	// Split by comma to separate block name from settings
	parts := strings.Split(blockStr, ",")
	if len(parts) == 0 {
		return shop2.SectionBlock{}, nil
	}

	// First part is the block name
	blockName := strings.TrimSpace(parts[0])
	blockType := utils.Handlize(blockName)

	// Parse settings from remaining parts
	settings := make([]shop2.Setting, 0, len(parts)-1)
	for i := 1; i < len(parts); i++ {
		settingStr := strings.TrimSpace(parts[i])
		if settingStr == "" {
			continue
		}

		setting, err := ParseSettingString(settingStr)
		if err != nil {
			return shop2.SectionBlock{}, err
		}
		settings = append(settings, setting)
	}

	return shop2.SectionBlock{
		Name:     blockName,
		Type:     blockType,
		Settings: settings,
	}, nil
}

func ParseSettingStrings(settingStrs string) ([]shop2.Setting, error) {

	if settingStrs == "" {
		return []shop2.Setting{}, ErrNoSettingString
	}

	var err error
	settings := strings.Split(settingStrs, ",")
	parsedSettings := make([]shop2.Setting, len(settings))
	for i, setting := range settings {
		parsedSettings[i], err = ParseSettingString(setting)
		if err != nil {
			return nil, err
		}
	}
	return parsedSettings, nil
}

// ParseSettingString parses a setting string in the format "type:id:label"
// Returns a Setting with the specified type, id, and label
func ParseSettingString(settingStr string) (shop2.Setting, error) {
	// Split by colon to get type, id, and label
	parts := strings.Split(settingStr, ":")
	if len(parts) < 3 {
		dbg("invalid setting format (expected type:id:label): %s", settingStr)
		return shop2.Setting{}, fmt.Errorf("invalid setting format (expected type:id:label): %s\n", settingStr)
	}

	settingType := strings.TrimSpace(parts[0])
	settingID := strings.TrimSpace(parts[1])
	settingLabel := strings.TrimSpace(parts[2])

	// Create the appropriate input type based on the settingType
	return shop2.NewSettingFromType(settingType, settingID, settingLabel)
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

var (
	missingNameErr       = errors.New("name is required")
	missingSectionDirErr = errors.New("missing sections directory. Are you at the top level of your project?")
	templatePathErr      = errors.New("sectionHandler: reading template")
)

var CreateSection = func(ctx context.Context, opts ...SectionServiceOption) func(ctx context.Context) error {

	return func(ctx context.Context) error {

		s := &SectionService{}

		// errors := make(map[string]error)
		for _, op := range opts {
			op(s)
		}

		if s.Section.Schema.Name == "" {
			return missingNameErr
		}

		// does sections dir exist?
		if !utils.FolderExistsInPwd("sections") {
			return missingSectionDirErr
		}

		// does file name already exist?
		path := filepath.Join("sections", utils.Handlize(s.Section.Schema.Name)+".liquid")
		if utils.FileExists(path) {
			path = utils.EnsureUniqueFilename(path)
		}

		if err := os.WriteFile(path, []byte(s.Render()), 0o644); err != nil {
			return fmt.Errorf("handleCreateBlock: write file: %w", err)
		}

		return nil
	}
}

// Similar to CreateSection, but just returns a string. Use this if within a TUI.
var CreateSectionPreview = func(ctx context.Context, opts ...SectionServiceOption) (string, error) {

	s := &SectionService{}

	// errors := make(map[string]error)
	for _, op := range opts {
		op(s)
	}

	if s.Section.Schema.Name == "" {
		return "", missingNameErr
	}

	return s.Render(), nil
}
