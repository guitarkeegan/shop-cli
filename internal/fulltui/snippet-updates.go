package fulltui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
)

var includingSnippetDocUpdate = func() func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var isFirstUpdate = true

	return func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.snippetDocModel.Update(msg)

		if f, ok := updated.(*huh.Form); ok {
			m.snippetDocModel = f
		}

		switch m.snippetDocModel.State {
		case huh.StateCompleted:
			if isFirstUpdate {
				isFirstUpdate = false
				if !m.snippetDocModel.GetBool(responseKey) {
					m.state = previewingFile
					m.snippetDocModel = nil
					m.previewFileModel = newPreviewFileForm()
					// Trigger preview with current state (no params)
					previewCmd := showSnippetPreviewCmd(
						context.Background(),
						ren.WithSnippetFileName(m.fileModel.GetString(nameKey)),
					)
					return m, tea.Batch(m.previewFileModel.Init(), previewCmd)
				}
				newParam := param{
					Name:        m.snippetDocModel.GetString(nameKey),
					Type:        m.snippetDocModel.GetString(typeKey),
					Description: m.snippetDocModel.GetString(textKey),
					Required:    m.snippetDocModel.GetBool(optionalKey),
				}
				newSnippetDoc := snippetDoc{
					Description: m.snippetDocModel.GetString(descriptionKey),
					Example:     m.snippetDocModel.GetString(exampleKey),
				}
				newSnippetDoc.Param = append(newSnippetDoc.Param, newParam)
				m.snippetDoc = newSnippetDoc
			} else {
				newParam := param{
					Name:        m.snippetDocModel.GetString(nameKey),
					Type:        m.snippetDocModel.GetString(typeKey),
					Description: m.snippetDocModel.GetString(textKey),
					Required:    m.snippetDocModel.GetBool(optionalKey),
				}
				m.snippetDoc.Param = append(m.snippetDoc.Param, newParam)
			}
			m.state = continueSnippetParams
			m.snippetDocModel = nil
			m.continueSnippetDocModel = newContinueSnippetParams()
			// Trigger preview update with the new param
			params := make([]shop2.Param, len(m.snippetDoc.Param))
			for i, param := range m.snippetDoc.Param {
				params[i] = shop2.Param{
					Type:        shop2.ParamTypeFromString(param.Type),
					Name:        param.Name,
					Required:    param.Required,
					Description: param.Description,
				}
			}
			previewCmd := showSnippetPreviewCmd(
				context.Background(),
				ren.WithSnippetFileName(m.fileModel.GetString(nameKey)),
				ren.WithSnippetDescription(m.snippetDoc.Description),
				ren.WithSnippetAddExample(m.snippetDoc.Example),
				ren.WithSnippetParams(params),
			)
			return m, tea.Batch(m.continueSnippetDocModel.Init(), previewCmd)
		case huh.StateAborted:
			m.state = quitting
			return m, tea.Quit
		}
		return m, cmd
	}
}()

func continueSnippetDocUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueSnippetDocModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueSnippetDocModel = f
	}

	switch m.continueSnippetDocModel.State {
	case huh.StateCompleted:
		if m.continueSnippetDocModel.GetBool(continueKey) {
			m.state = includingSnippetDoc
			m.continueSnippetDocModel = nil
			m.snippetDocModel = newSnippetDocForm(func() bool { return true }, func() bool { return true })
			return m, m.snippetDocModel.Init()
		}
		m.state = previewingFile
		m.continueSnippetDocModel = nil
		m.previewFileModel = newPreviewFileForm()
		// make the shopify params
		params := make([]shop2.Param, len(m.snippetDoc.Param))
		for i, param := range m.snippetDoc.Param {
			params[i] = shop2.Param{
				Type:        shop2.ParamTypeFromString(param.Type),
				Name:        param.Name,
				Required:    param.Required,
				Description: param.Description,
			}
		}
		previewCmd := showSnippetPreviewCmd(
			context.Background(),
			ren.WithSnippetFileName(m.fileModel.GetString(nameKey)),
			ren.WithSnippetDescription(m.snippetDoc.Description),
			ren.WithSnippetAddExample(m.snippetDoc.Example),
			ren.WithSnippetParams(params),
		)
		return m, tea.Batch(m.previewFileModel.Init(), previewCmd)
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}
