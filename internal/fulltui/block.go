package fulltui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Collects settings for a standalone theme block file.
func includingBlockFileSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.settingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.settingsModel = f
	}

	switch m.settingsModel.State {
	case huh.StateCompleted:
		// Only the first form asks "Include settings?".
		if len(m.settings) == 0 && !m.settingsModel.GetBool(responseKey) {
			m.settingsModel = nil
			return m.toPreview()
		}
		m.settings = append(m.settings, settingFromForm(m.settingsModel))
		m.continueSettingsModel = newContinueSettingsForm()
		m.state = continueBlockFileSettings
		return m, tea.Batch(m.continueSettingsModel.Init(), m.previewCmd())
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

func continueBlockFileSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueSettingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueSettingsModel = f
	}

	switch m.continueSettingsModel.State {
	case huh.StateCompleted:
		wantsMore := m.continueSettingsModel.GetBool(continueKey)
		m.continueSettingsModel = nil
		if wantsMore {
			m.state = includingBlockFileSettings
			m.settingsModel = newSettingsForm(true)
			return m, m.settingsModel.Init()
		}
		return m.toPreview()
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}
