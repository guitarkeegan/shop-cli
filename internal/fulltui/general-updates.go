package fulltui

import (
	"context"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
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
		fileType := m.fileModel.GetString(typeKey)
		switch fileType {
		case "section":
			m.state = includingSettings
			// Initialize viewport for live preview
			viewportWidth := (m.width / 2) - 4
			viewportHeight := m.height - 6
			m.viewportModel = viewport.New(viewportWidth, viewportHeight)
			// Trigger initial preview with just the section name
			previewCmd := showFilePreviewCmd(
				context.Background(),
				ren.WithSectionBlocks(m.blockSettings),
				ren.WithSectionName(m.fileModel.GetString(nameKey)),
				ren.WithSectionSettings(m.settings),
			)
			return m, tea.Batch(m.settingsModel.Init(), previewCmd)
		case "block":
			m.state = includingBlockFileSettings
			return m, m.settingsModel.Init()
		case "snippet":
			m.state = includingSnippetDoc
			m.snippetDocModel = newSnippetDocForm()
			// Initialize viewport for live preview
			viewportWidth := (m.width / 2) - 4
			viewportHeight := m.height - 6
			m.viewportModel = viewport.New(viewportWidth, viewportHeight)
			// Trigger initial preview with just the snippet name
			previewCmd := showSnippetPreviewCmd(
				context.Background(),
				ren.WithSnippetFileName(m.fileModel.GetString(nameKey)),
			)
			return m, tea.Batch(m.snippetDocModel.Init(), previewCmd)
		}
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
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
		if m.previewFileModel.GetBool(responseKey) {
			switch m.fileModel.GetString(typeKey) {
			case "section":
				m.state = creatingFile
				m.previewFileModel = nil
				// make the file
				return m, createSectionCmd(
					context.Background(),
					ren.WithSectionBlocks(m.blockSettings),
					ren.WithSectionName(m.fileModel.GetString(nameKey)),
					ren.WithSectionSettings(m.settings),
					ren.WithSectionPresets([]shop2.Preset{
						{Name: m.fileModel.GetString(nameKey)},
					}),
				)
			case "snippet":
				m.state = creatingFile
				m.previewFileModel = nil
				params := make([]shop2.Param, len(m.snippetDoc.Param))
				for i, param := range m.snippetDoc.Param {
					params[i] = shop2.Param{
						Type:        shop2.ParamTypeFromString(param.Type),
						Name:        param.Name,
						Required:    param.Required,
						Description: param.Description,
					}
				}
				return m, tea.Batch(
					cmd,
					createSnippetCmd(
						context.Background(),
						ren.WithSnippetFileName(m.fileModel.GetString(nameKey)),
						ren.WithSnippetDescription(m.snippetDoc.Description),
						ren.WithSnippetAddExample(m.snippetDoc.Example),
						ren.WithSnippetParams(params),
					),
				)
			}
		}
		m.state = selectingFile
		m.fileModel = nil
		m.fileModel = newFileSelector()
		return m, m.fileModel.Init()
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, tea.Batch(cmds...)
}
