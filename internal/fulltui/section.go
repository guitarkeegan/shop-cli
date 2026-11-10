package fulltui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/guitarkeegan/shop-cli/internal/shopify"
)

// toBlocks moves from section settings to the section blocks questions.
func (m model) toBlocks() (model, tea.Cmd) {
	m.state = includingBlocks
	m.addingToBlock = false
	m.blockSettingsModel = newIncludingBlocksForm(false, false)
	return m, tea.Batch(m.blockSettingsModel.Init(), m.previewCmd())
}

func includingSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m.toBlocks()
		}
		m.settings = append(m.settings, settingFromForm(m.settingsModel))
		m.continueSettingsModel = newContinueSettingsForm()
		m.state = continueSettings
		return m, tea.Batch(m.continueSettingsModel.Init(), m.previewCmd())
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

func continueSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.state = includingSettings
			m.settingsModel = newSettingsForm(true)
			return m, m.settingsModel.Init()
		}
		return m.toBlocks()
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

func includingBlocksSettingUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.blockSettingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.blockSettingsModel = f
	}

	switch m.blockSettingsModel.State {
	case huh.StateCompleted:
		f := m.blockSettingsModel
		// Only the first form asks "Include section blocks?".
		if len(m.blockSettings) == 0 && !f.GetBool(responseKey) {
			m.blockSettingsModel = nil
			return m.toPreview()
		}

		setting := settingFromForm(f)
		if m.addingToBlock && len(m.blockSettings) > 0 {
			last := &m.blockSettings[len(m.blockSettings)-1]
			last.Settings = append(last.Settings, setting)
		} else {
			m.blockSettings = append(m.blockSettings, shopify.SectionBlock{
				Name:     f.GetString(nameKey),
				Type:     f.GetString(typeKey),
				Settings: []shopify.Setting{setting},
			})
		}

		m.state = continueBlockSettings
		m.blockSettingsModel = nil
		m.continueBlockSettingsModel = newContinueSettingsForm()
		return m, tea.Batch(m.continueBlockSettingsModel.Init(), m.previewCmd())
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

// Asks whether to add another setting to the current block.
func continueBlockSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueBlockSettingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueBlockSettingsModel = f
	}

	switch m.continueBlockSettingsModel.State {
	case huh.StateCompleted:
		wantsMore := m.continueBlockSettingsModel.GetBool(continueKey)
		m.continueBlockSettingsModel = nil
		if wantsMore {
			m.state = includingBlocks
			m.addingToBlock = true
			m.blockSettingsModel = newIncludingBlocksForm(true, true)
			return m, m.blockSettingsModel.Init()
		}
		m.state = additionalBlocksPrompt
		m.additionalBlocksPromptModel = newAdditionalBlocksForm()
		return m, tea.Batch(m.additionalBlocksPromptModel.Init(), m.previewCmd())
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}

// Asks whether to add another block to the section.
func additionalBlocksPromptUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.additionalBlocksPromptModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.additionalBlocksPromptModel = f
	}

	switch m.additionalBlocksPromptModel.State {
	case huh.StateCompleted:
		wantsMore := m.additionalBlocksPromptModel.GetBool(continueKey)
		m.additionalBlocksPromptModel = nil
		if wantsMore {
			m.state = includingBlocks
			m.addingToBlock = false
			m.blockSettingsModel = newIncludingBlocksForm(true, false)
			return m, m.blockSettingsModel.Init()
		}
		return m.toPreview()
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}

	return m, cmd
}
