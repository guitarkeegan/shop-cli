package fulltui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
)

func continueBlockSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueBlockSettingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueBlockSettingsModel = f
	}

	switch m.continueBlockSettingsModel.State {
	case huh.StateCompleted:
		if m.continueBlockSettingsModel.GetBool(continueKey) {
			m.state = includingBlocks
			m.continueBlockSettingsModel = nil
			m.blockSettingsModel = newIncludingBlocksForm(func() bool { return true }, func() bool { return true })
			return m, tea.Batch(m.blockSettingsModel.Init(), isFirstUpdateCmd(false))
		}
		m.state = additionalBlocksPrompt
		m.additionalBlocksPromptModel = newAdditionalBlocksForm()
		// Trigger preview update when transitioning to additional blocks prompt
		previewCmd := showFilePreviewCmd(
			context.Background(),
			ren.WithSectionBlocks(m.blockSettings),
			ren.WithSectionName(m.fileModel.GetString(nameKey)),
			ren.WithSectionSettings(m.settings),
		)
		return m, tea.Batch(m.additionalBlocksPromptModel.Init(), previewCmd)
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

var includingBlockFileSettingsUpdate = func() func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var isFirstUpdate bool = true

	return func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.settingsModel.Update(msg)

		if f, ok := updated.(*huh.Form); ok {
			m.settingsModel = f
		}

		switch m.settingsModel.State {
		case huh.StateCompleted:
			if isFirstUpdate {
				wantsSettings := m.settingsModel.GetBool(responseKey)
				if !wantsSettings {
					m.state = quitting
					return m, tea.Quit
				}
				isFirstUpdate = false
			}
			m.settings = append(m.settings, shop2.Setting{
				ID:    m.settingsModel.GetString(idKey),
				Type:  m.settingsModel.GetString(inputTypeKey),
				Label: m.settingsModel.GetString(labelKey),
			})
			m.continueSettingsModel = newContinueSettingsForm()
			m.state = continueBlockFileSettings
			return m, m.continueSettingsModel.Init()
		case huh.StateAborted:
			m.state = quitting
			return m, tea.Quit
		}
		return m, cmd
	}
}()

func continueBlockFileSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueSettingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueSettingsModel = f
	}

	switch m.continueSettingsModel.State {
	case huh.StateCompleted:
		if m.continueSettingsModel.GetBool(continueKey) {
			m.state = includingBlockFileSettings
			m.continueSettingsModel = nil
			m.settingsModel = newSettingsForm(func() bool { return true })
			return m, m.settingsModel.Init()
		}
		m.state = quitting
		return m, tea.Quit
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}
