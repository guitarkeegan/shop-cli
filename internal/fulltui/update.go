package fulltui

import tea "github.com/charmbracelet/bubbletea"

// hasPreview reports whether the state shows the live preview pane.
func (s state) hasPreview() bool {
	switch s {
	case includingSettings, continueSettings, includingBlocks, continueBlockSettings,
		additionalBlocksPrompt, includingBlockFileSettings, continueBlockFileSettings,
		includingSnippetDoc, continueSnippetParams, previewingFile:
		return true
	}
	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.state = quitting
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.state.hasPreview() {
			m.viewportModel.Width, m.viewportModel.Height = m.previewSize()
		}
	case previewMsg:
		if msg.Err != nil {
			m.preview = "Error: " + msg.Err.Error()
		} else {
			m.preview = msg.Preview
		}
		if m.state.hasPreview() {
			m.viewportModel.SetContent(m.preview)
		}
	case fileCreatedMsg:
		m.created = msg
		m.state = quitting
		return m, tea.Quit
	}

	// Update viewport for scrolling in file creation states. The preview
	// state updates its own viewport.
	if m.state.hasPreview() && m.state != previewingFile {
		var viewportCmd tea.Cmd
		m.viewportModel, viewportCmd = m.viewportModel.Update(msg)
		cmd = tea.Batch(cmd, viewportCmd)
	}

	switch m.state {

	case selectingFile:
		return selectingFileUpdate(m, msg)

	case includingSettings:
		return includingSettingsUpdate(m, msg)

	case continueSettings:
		return continueSettingsUpdate(m, msg)

	case includingBlocks:
		return includingBlocksSettingUpdate(m, msg)

	case continueBlockSettings:
		return continueBlockSettingsUpdate(m, msg)

	case additionalBlocksPrompt:
		return additionalBlocksPromptUpdate(m, msg)

	case includingBlockFileSettings:
		return includingBlockFileSettingsUpdate(m, msg)

	case continueBlockFileSettings:
		return continueBlockFileSettingsUpdate(m, msg)

	case includingSnippetDoc:
		return includingSnippetDocUpdate(m, msg)

	case continueSnippetParams:
		return continueSnippetDocUpdate(m, msg)

	case previewingFile:
		return previewFileFormUpdate(m, msg)

	case creatingFile:
		return m, nil

	case quitting:
		return m, tea.Quit
	}

	return m, cmd
}
