package fulltui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
)

// Model holds the state of the program. Forms are typically cleaned up,
// and need to be reinitialized before reuse.
type model struct {
	state                       state
	fileModel                   *huh.Form
	includeInputsForm           *huh.Form
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
	snippetDoc                  snippetDoc
	width                       int
	height                      int
	preview                     string
	viewportModel               viewport.Model
}

// Setup the initial view. Limits on block settings, and snippet doc params are also setup.
func initialModel() model {

	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.ShowDescription = false

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
