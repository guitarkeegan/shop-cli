package fulltui

import (
	"context"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	ren "github.com/guitarkeegan/shop-cli/internal/render-service"
	shop2 "github.com/guitarkeegan/shop-cli/internal/shopify-v2"
)

var includingSettingsUpdate = func() func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {

	var isFirstUpdate bool = true

	return func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
		// ensure the form is updating
		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.settingsModel.Update(msg)

		if f, ok := updated.(*huh.Form); ok {
			m.settingsModel = f
		}

		switch m.settingsModel.State {
		case huh.StateCompleted:
			if isFirstUpdate {
				wantsInputs := m.settingsModel.GetBool(responseKey)
				if !wantsInputs {
					m.state = includingBlocks
					m.settingsModel = nil
					m.blockSettingsModel = newIncludingBlocksForm()
					// Trigger preview update before transitioning to blocks
					previewCmd := showFilePreviewCmd(
						context.Background(),
						ren.WithSectionBlocks(m.blockSettings),
						ren.WithSectionName(m.fileModel.GetString(nameKey)),
						ren.WithSectionSettings(m.settings),
					)
					return m, tea.Batch(m.blockSettingsModel.Init(), previewCmd)
				}
				isFirstUpdate = false
			}
			m.settings = append(m.settings, shop2.Setting{
				ID:    m.settingsModel.GetString(idKey),
				Type:  m.settingsModel.GetString(inputTypeKey),
				Label: m.settingsModel.GetString(labelKey),
			})
			m.continueSettingsModel = newContinueSettingsForm()
			m.state = continueSettings
			m.includeInputsForm = nil
			// Trigger preview update with the new setting
			previewCmd := showFilePreviewCmd(
				context.Background(),
				ren.WithSectionBlocks(m.blockSettings),
				ren.WithSectionName(m.fileModel.GetString(nameKey)),
				ren.WithSectionSettings(m.settings),
			)
			return m, tea.Batch(m.continueSettingsModel.Init(), previewCmd)
		case huh.StateAborted:
			m.state = quitting
			return m, tea.Quit
		}
		return m, cmd
	}
}()

// Update asking the user if they would like to include more blocks in their section file.
var additionalBlocksPromptUpdate = func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.additionalBlocksPromptModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.additionalBlocksPromptModel = f
	}

	switch m.additionalBlocksPromptModel.State {
	case huh.StateCompleted:
		if m.additionalBlocksPromptModel.GetBool(continueKey) {
			m.state = includingBlocks
			m.additionalBlocksPromptModel = nil
			m.blockSettingsModel = newIncludingBlocksForm(func() bool { return true })
			return m, tea.Batch(m.blockSettingsModel.Init(), isFirstBlockCmd(false), isFirstUpdateCmd(false), isAdditionalUpdateCmd(false))
		}
		m.state = previewingFile
		m.additionalBlocksPromptModel = nil
		// setup viewport with proper dimensions for right pane
		viewportWidth := (m.width / 2) - 4 // Account for borders/padding
		viewportHeight := m.height - 6     // Account for borders/padding
		m.viewportModel = viewport.New(viewportWidth, viewportHeight)
		m.previewFileModel = newPreviewFileForm()

		previewCmd := showFilePreviewCmd(
			context.Background(),
			ren.WithSectionBlocks(m.blockSettings),
			ren.WithSectionName(m.fileModel.GetString(nameKey)),
			ren.WithSectionSettings(m.settings),
		)
		return m, tea.Batch(m.previewFileModel.Init(), previewCmd)
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}

	return m, cmd
}

var includingBlocksSettingUpdate = func() func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {

	var isFirstBlock = true
	var isFirstUpdate = true
	var isAdditionalUpdate = true

	return func(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msg := msg.(type) {
		case isFirstUpdateMsg:
			isFirstUpdate = bool(msg)
		case isFirstBlockMsg:
			isFirstBlock = bool(msg)
		case isAdditionalUpdateMsg:
			isAdditionalUpdate = bool(msg)
		}

		var cmd tea.Cmd
		var updated tea.Model
		updated, cmd = m.blockSettingsModel.Update(msg)

		if f, ok := updated.(*huh.Form); ok {
			m.blockSettingsModel = f
		}

		switch m.blockSettingsModel.State {
		case huh.StateCompleted:
			if isFirstUpdate {
				wantsBlockSettings := m.blockSettingsModel.GetBool(responseKey)
				if !wantsBlockSettings {
					m.state = previewingFile
					m.blockSettingsModel = nil
					m.previewFileModel = newPreviewFileForm()

					previewCmd := showFilePreviewCmd(
						context.Background(),
						ren.WithSectionBlocks(m.blockSettings),
						ren.WithSectionName(m.fileModel.GetString(nameKey)),
						ren.WithSectionSettings(m.settings),
					)
					return m, tea.Batch(m.previewFileModel.Init(), previewCmd)
				}
				newSetting := shop2.Setting{
					ID:    m.blockSettingsModel.GetString(idKey),
					Type:  m.blockSettingsModel.GetString(inputTypeKey),
					Label: m.blockSettingsModel.GetString(labelKey),
				}
				newBlock := shop2.SectionBlock{
					Name:     m.blockSettingsModel.GetString(nameKey),
					Type:     m.blockSettingsModel.GetString(typeKey),
					Settings: []shop2.Setting{},
				}
				newBlock.Settings = append(newBlock.Settings, newSetting)
				m.blockSettings = append(m.blockSettings, newBlock)
				isFirstUpdate = false
			} else if !isFirstBlock && !isAdditionalUpdate {
				newSetting := shop2.Setting{
					ID:    m.blockSettingsModel.GetString(idKey),
					Type:  m.blockSettingsModel.GetString(inputTypeKey),
					Label: m.blockSettingsModel.GetString(labelKey),
				}
				newBlock := shop2.SectionBlock{
					Name:     m.blockSettingsModel.GetString(nameKey),
					Type:     m.blockSettingsModel.GetString(typeKey),
					Settings: []shop2.Setting{},
				}
				newBlock.Settings = append(newBlock.Settings, newSetting)
				m.blockSettings = append(m.blockSettings, newBlock)
				isAdditionalUpdate = true
			} else if isAdditionalUpdate {
				newSetting := shop2.Setting{
					ID:    m.blockSettingsModel.GetString(idKey),
					Type:  m.blockSettingsModel.GetString(inputTypeKey),
					Label: m.blockSettingsModel.GetString(labelKey),
				}
				m.blockSettings[len(m.blockSettings)-1].Settings = append(m.blockSettings[len(m.blockSettings)-1].Settings, newSetting)
			} else {
				dbg("else, is first block, not first update...")
				newSetting := shop2.Setting{
					ID:    m.blockSettingsModel.GetString(idKey),
					Type:  m.blockSettingsModel.GetString(inputTypeKey),
					Label: m.blockSettingsModel.GetString(labelKey),
				}
				m.blockSettings[len(m.blockSettings)-1].Settings = append(m.blockSettings[len(m.blockSettings)-1].Settings, newSetting)
			}
			m.state = continueBlockSettings
			m.blockSettingsModel = nil
			m.continueBlockSettingsModel = newContinueSettingsForm()
			// Trigger preview update with the new block setting
			previewCmd := showFilePreviewCmd(
				context.Background(),
				ren.WithSectionBlocks(m.blockSettings),
				ren.WithSectionName(m.fileModel.GetString(nameKey)),
				ren.WithSectionSettings(m.settings),
			)
			return m, tea.Batch(m.continueBlockSettingsModel.Init(), previewCmd)
		case huh.StateAborted:
			m.state = quitting
			return m, tea.Quit
		}
		return m, cmd
	}
}()

func continueSettingsUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updated tea.Model
	updated, cmd = m.continueSettingsModel.Update(msg)

	if f, ok := updated.(*huh.Form); ok {
		m.continueSettingsModel = f
	}

	switch m.continueSettingsModel.State {
	case huh.StateCompleted:
		if m.continueSettingsModel.GetBool(continueKey) {
			m.state = includingSettings
			m.continueSettingsModel = nil
			m.settingsModel = newSettingsForm(func() bool { return true })
			return m, m.settingsModel.Init()
		}
		m.state = includingBlocks
		m.continueSettingsModel = nil
		m.blockSettingsModel = newIncludingBlocksForm()
		// Trigger preview update when transitioning to blocks
		previewCmd := showFilePreviewCmd(
			context.Background(),
			ren.WithSectionBlocks(m.blockSettings),
			ren.WithSectionName(m.fileModel.GetString(nameKey)),
			ren.WithSectionSettings(m.settings),
		)
		return m, tea.Batch(m.blockSettingsModel.Init(), previewCmd)
	case huh.StateAborted:
		m.state = quitting
		return m, tea.Quit
	}
	return m, cmd
}
