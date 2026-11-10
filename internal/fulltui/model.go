package fulltui

import (
	"context"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
)

// Model holds the state of the program. Forms are typically cleaned up,
// and need to be reinitialized before reuse.
type model struct {
	state                       state
	fileModel                   *huh.Form
	settingsModel               *huh.Form
	continueSettingsModel       *huh.Form
	blockSettingsModel          *huh.Form
	continueBlockSettingsModel  *huh.Form
	additionalBlocksPromptModel *huh.Form
	snippetDocModel             *huh.Form
	continueSnippetDocModel     *huh.Form
	previewFileModel            *huh.Form
	settings                    []shop2.Setting
	blockSettings               []shop2.SectionBlock
	// addingToBlock is true when blockSettingsModel only asks for a setting
	// to add to the last block, rather than for a new block.
	addingToBlock bool
	snippetDoc    snippetDoc
	width         int
	height        int
	preview       string
	viewportModel viewport.Model
	created       fileCreatedMsg
}

// initialModel starts the program on the file type selector.
func initialModel() model {
	fts := newFileSelector()
	fts.Init()

	return model{
		state:         selectingFile,
		fileModel:     fts,
		settingsModel: newSettingsForm(),
		settings:      make([]shop2.Setting, 0, 20),
		blockSettings: make([]shop2.SectionBlock, 0, 20),
		snippetDoc:    snippetDoc{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// restart returns a fresh model that keeps the current window size.
func (m model) restart() (model, tea.Cmd) {
	fresh := initialModel()
	fresh.width, fresh.height = m.width, m.height
	return fresh, fresh.fileModel.Init()
}

// newViewport sizes the preview pane to the right half of the window.
func (m model) newViewport() viewport.Model {
	return viewport.New((m.width/2)-4, m.height-6)
}

func (m model) fileName() string {
	return m.fileModel.GetString(nameKey)
}

func (m model) sectionOpts() []ren.SectionServiceOption {
	return []ren.SectionServiceOption{
		ren.WithSectionBlocks(m.blockSettings),
		ren.WithSectionName(m.fileName()),
		ren.WithSectionSettings(m.settings),
	}
}

func (m model) blockOpts() []ren.RenderBlockServiceOption {
	return []ren.RenderBlockServiceOption{
		ren.WithBlockName(m.fileName()),
		ren.WithBlockSettings(m.settings),
	}
}

func (m model) snippetOpts() []ren.SnippetOption {
	params := make([]shop2.Param, len(m.snippetDoc.Param))
	for i, p := range m.snippetDoc.Param {
		params[i] = shop2.Param{
			Type:        shop2.ParamTypeFromString(p.Type),
			Name:        p.Name,
			Required:    p.Required,
			Description: p.Description,
		}
	}
	return []ren.SnippetOption{
		ren.WithSnippetFileName(m.fileName()),
		ren.WithSnippetDescription(m.snippetDoc.Description),
		ren.WithSnippetAddExample(m.snippetDoc.Example),
		ren.WithSnippetParams(params),
	}
}

func (m model) sectionPreviewCmd() tea.Cmd {
	return showFilePreviewCmd(context.Background(), m.sectionOpts()...)
}

func (m model) blockPreviewCmd() tea.Cmd {
	return showBlockPreviewCmd(context.Background(), m.blockOpts()...)
}

func (m model) snippetPreviewCmd() tea.Cmd {
	return showSnippetPreviewCmd(context.Background(), m.snippetOpts()...)
}

// settingFromForm reads the type/id/label fields shared by the settings forms.
func settingFromForm(f *huh.Form) shop2.Setting {
	return shop2.Setting{
		ID:    f.GetString(idKey),
		Type:  f.GetString(inputTypeKey),
		Label: f.GetString(labelKey),
	}
}
