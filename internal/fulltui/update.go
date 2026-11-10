package fulltui

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			m.state = quitting
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update viewport dimensions when window resizes for all file creation states
		switch m.state {
		case includingSettings, continueSettings, includingBlocks, continueBlockSettings, additionalBlocksPrompt, includingSnippetDoc, continueSnippetParams, previewingFile:
			m.viewportModel.Width = (m.width / 2) - 4 // Account for borders/padding
			m.viewportModel.Height = m.height - 6     // Account for borders/padding
		}
	case previewMsg:
		if msg.Err != nil {
			m.preview = "Error: " + msg.Err.Error()
		} else {
			m.preview = msg.Preview
		}
		// Set the viewport content when preview is ready for all file creation states
		switch m.state {
		case includingSettings, continueSettings, includingBlocks, continueBlockSettings, additionalBlocksPrompt, includingSnippetDoc, continueSnippetParams, previewingFile:
			m.viewportModel.SetContent(m.preview)
		}
	case createSectionMsg, createSnippetMsg:
		return m, tea.Quit
	}

	// Update viewport for scrolling in file creation states
	switch m.state {
	case includingSettings, continueSettings, includingBlocks, continueBlockSettings, additionalBlocksPrompt, includingSnippetDoc, continueSnippetParams:
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
