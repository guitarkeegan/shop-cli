package fulltui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func selectingFileUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.fileModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.fileModel = f
	}

	switch m.fileModel.State {
	case huh.StateCompleted:
		m.viewportModel = m.newViewport()
		switch m.fileType() {
		case "section":
			m.state = includingSettings
			return m, tea.Batch(m.settingsModel.Init(), m.previewCmd())
		case "block":
			m.state = includingBlockFileSettings
			return m, tea.Batch(m.settingsModel.Init(), m.previewCmd())
		case "snippet":
			m.state = includingSnippetDoc
			m.snippetDocModel = newSnippetDocForm(false)
			return m, tea.Batch(m.snippetDocModel.Init(), m.previewCmd())
		}
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

// toPreview moves to the final confirm step, showing the rendered file.
func (m model) toPreview() (model, tea.Cmd) {
	m.state = previewingFile
	m.previewFileModel = newPreviewFileForm()
	return m, tea.Batch(m.previewFileModel.Init(), m.previewCmd())
}

func previewFileFormUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var updated tea.Model

	// Update the viewport to handle scrolling
	m.viewportModel, cmd = m.viewportModel.Update(msg)
	cmds = append(cmds, cmd)

	// Update the form
	updated, cmd = m.previewFileModel.Update(msg)
	cmds = append(cmds, cmd)

	if f, ok := updated.(*huh.Form); ok {
		m.previewFileModel = f
	}

	switch m.previewFileModel.State {
	case huh.StateCompleted:
		if !m.previewFileModel.GetBool(responseKey) {
			return m.restart()
		}
		m.state = creatingFile
		m.previewFileModel = nil
		return m, m.createCmd()
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, tea.Batch(cmds...)
}
