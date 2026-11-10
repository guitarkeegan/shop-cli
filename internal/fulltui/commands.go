package fulltui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guitarkeegan/shop-cli/internal/render"
	"github.com/guitarkeegan/shop-cli/internal/theme"
)

type previewMsg struct {
	Preview string
	Err     error
}

// fileCreatedMsg reports the result of writing a file. The program quits on
// receiving it, and Run reports the path or error once the screen is restored.
type fileCreatedMsg struct {
	path string
	err  error
}

// fileType returns the kind of file being built: "section", "block" or "snippet".
func (m model) fileType() string {
	return m.fileModel.GetString(typeKey)
}

// render renders the file being built from the answers so far.
func (m model) render() (string, error) {
	switch m.fileType() {
	case "section":
		return render.Section(
			render.WithSectionName(m.fileName()),
			render.WithSectionSettings(m.settings),
			render.WithSectionBlocks(m.blockSettings),
		)
	case "block":
		return render.Block(
			render.WithBlockName(m.fileName()),
			render.WithBlockSettings(m.settings),
		)
	case "snippet":
		return render.Snippet(
			render.WithSnippetDescription(m.snippetDoc.Description),
			render.WithSnippetExample(m.snippetDoc.Example),
			render.WithSnippetParams(m.snippetDoc.shopifyParams()),
		), nil
	}
	return "", nil
}

// previewCmd renders the file for the preview pane.
func (m model) previewCmd() tea.Cmd {
	return func() tea.Msg {
		preview, err := m.render()
		return previewMsg{Preview: preview, Err: err}
	}
}

// createCmd renders the file and writes it into the theme.
func (m model) createCmd() tea.Cmd {
	return func() tea.Msg {
		content, err := m.render()
		if err != nil {
			return fileCreatedMsg{err: err}
		}

		write := map[string]func(name, content string) (string, error){
			"section": theme.WriteSection,
			"block":   theme.WriteBlock,
			"snippet": theme.WriteSnippet,
		}[m.fileType()]

		path, err := write(m.fileName(), content)
		return fileCreatedMsg{path: path, err: err}
	}
}
