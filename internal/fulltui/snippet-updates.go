package fulltui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func includingSnippetDocUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.snippetDocModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.snippetDocModel = f
	}

	switch m.snippetDocModel.State {
	case huh.StateCompleted:
		f := m.snippetDocModel
		// Only the first form asks "Add doc params?" and for the description.
		if len(m.snippetDoc.Param) == 0 {
			if !f.GetBool(responseKey) {
				m.snippetDocModel = nil
				return m.toPreview(m.snippetPreviewCmd())
			}
			m.snippetDoc.Description = f.GetString(descriptionKey)
			m.snippetDoc.Example = f.GetString(exampleKey)
		}
		m.snippetDoc.Param = append(m.snippetDoc.Param, param{
			Name:        f.GetString(nameKey),
			Type:        f.GetString(typeKey),
			Description: f.GetString(textKey),
			Required:    f.GetBool(optionalKey),
		})

		m.state = continueSnippetParams
		m.snippetDocModel = nil
		m.continueSnippetDocModel = newContinueSnippetParams()
		return m, tea.Batch(m.continueSnippetDocModel.Init(), m.snippetPreviewCmd())
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

func continueSnippetDocUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueSnippetDocModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueSnippetDocModel = f
	}

	switch m.continueSnippetDocModel.State {
	case huh.StateCompleted:
		wantsMore := m.continueSnippetDocModel.GetBool(continueKey)
		m.continueSnippetDocModel = nil
		if wantsMore {
			m.state = includingSnippetDoc
			m.snippetDocModel = newSnippetDocForm(func() bool { return true }, func() bool { return true })
			return m, m.snippetDocModel.Init()
		}
		return m.toPreview(m.snippetPreviewCmd())
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}
