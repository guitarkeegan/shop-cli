package fulltui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/guitarkeegan/shop-cli/internal/shopify"
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
	settings                    []shopify.Setting
	blockSettings               []shopify.SectionBlock
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
		settingsModel: newSettingsForm(false),
		settings:      make([]shopify.Setting, 0, 20),
		blockSettings: make([]shopify.SectionBlock, 0, 20),
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

// settingFromForm reads the type/id/label fields shared by the settings forms.
func settingFromForm(f *huh.Form) shopify.Setting {
	return shopify.Setting{
		ID:    f.GetString(idKey),
		Type:  f.GetString(inputTypeKey),
		Label: f.GetString(labelKey),
	}
}
