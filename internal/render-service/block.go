package renderservice

import (
	"context"
	"encoding/json"
	"strings"

	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
	"github.com/guitarkeegan/shop-cli/internal/utils"
)

type RenderBlockService struct {
	Block shop2.ThemeBlock
}

type RenderBlockServiceOption func(*RenderBlockService)

func NewRenderBlockService() *RenderBlockService {
	return &RenderBlockService{}
}

func (rbs RenderBlockService) Render() string {
	blk := strings.Builder{}

	if len(rbs.Block.Schema.Presets) == 0 && rbs.Block.Schema.Name != "" {
		rbs.Block.Schema.Presets = []shop2.Preset{{Name: rbs.Block.Schema.Name}}
	}

	blk.WriteString("{% liquid\n")
	for _, setting := range rbs.Block.Schema.Settings {
		blk.WriteString("assign ")
		blk.WriteString(setting.ID)
		blk.WriteString(" = block.settings.")
		blk.WriteString(setting.ID)
		blk.WriteByte('\n')
	}
	blk.WriteString("%}\n")

	if len(rbs.Block.Schema.Blocks) > 0 {
		blk.WriteString("{% content_for 'blocks' %}\n")
	}

	bites, err := json.MarshalIndent(rbs.Block.Schema, "", "  ")
	if err != nil {
		dbg("%s", err)
		return "error on marshal"
	}

	blk.WriteString("{% schema %}\n")
	blk.WriteString(string(bites))
	blk.WriteString("\n{% endschema %}")

	return blk.String()
}

func WithBlockName(name string) RenderBlockServiceOption {
	return func(rbs *RenderBlockService) {
		rbs.Block.Schema.Name = name
	}
}

func WithBlockSettings(settings []shop2.Setting) RenderBlockServiceOption {
	return func(rbs *RenderBlockService) {
		rbs.Block.Schema.Settings = settings
	}
}

func WithBlockPresets(presets []shop2.Preset) RenderBlockServiceOption {
	return func(rbs *RenderBlockService) {
		rbs.Block.Schema.Presets = presets
	}
}

func WithBlockBlocks(blocks []shop2.ThemeBlockReference) RenderBlockServiceOption {
	return func(rbs *RenderBlockService) {
		rbs.Block.Schema.Blocks = blocks
	}
}

// CreateBlock renders a theme block and writes it to blocks/<handle>.liquid.
// It returns the path of the created file.
func CreateBlock(ctx context.Context, opts ...RenderBlockServiceOption) (string, error) {
	b := NewRenderBlockService()
	for _, op := range opts {
		op(b)
	}

	if b.Block.Schema.Name == "" {
		return "", ErrMissingName
	}
	if !utils.FolderExistsInPwd("blocks") {
		return "", ErrMissingBlocksDir
	}

	return writeLiquid("blocks", b.Block.Schema.Name, b.Render())
}

// CreateBlockPreview renders a theme block without writing it. Use this within a TUI.
func CreateBlockPreview(ctx context.Context, opts ...RenderBlockServiceOption) (string, error) {
	b := NewRenderBlockService()
	for _, op := range opts {
		op(b)
	}

	if b.Block.Schema.Name == "" {
		return "", ErrMissingName
	}

	return b.Render(), nil
}
