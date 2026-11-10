package renderservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
	"github.com/guitarkeegan/shop-cli/internal/utils"
)

var (
	missingBlocksDirErr = errors.New("missing blocks directory. Are you at the top level of your project?")
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

var CreateBlock = func(ctx context.Context, opts ...RenderBlockServiceOption) func(ctx context.Context) error {
	return func(ctx context.Context) error {

		b := NewRenderBlockService()
		for _, op := range opts {
			op(b)
		}

		if b.Block.Schema.Name == "" {
			return missingNameErr
		}

		if !utils.FolderExistsInPwd("blocks") {
			return missingBlocksDirErr
		}

		path := filepath.Join("blocks", utils.Handlize(b.Block.Schema.Name)+".liquid")
		if utils.FileExists(path) {
			path = utils.EnsureUniqueFilename(path)
		}

		if err := os.WriteFile(path, []byte(b.Render()), 0o644); err != nil {
			return fmt.Errorf("handleCreateBlock: write file: %w", err)
		}

		return nil
	}
}

var CreateBlockPreview = func(ctx context.Context, opts ...RenderBlockServiceOption) (string, error) {

	b := NewRenderBlockService()
	for _, op := range opts {
		op(b)
	}

	if b.Block.Schema.Name == "" {
		return "", missingNameErr
	}

	return b.Render(), nil
}
