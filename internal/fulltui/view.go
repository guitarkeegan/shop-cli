package fulltui

import "strings"

func (m model) View() string {
	switch m.state {
	case selectingFile:
		return withPageTitle("Welcome to Shop CLI!", m.fileModel.View())
	case includingSettings:
		return m.splitScreen(
			func() string { return withPageTitle("Section -> Settings", m.settingsModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case continueSettings:
		return m.splitScreen(
			func() string { return withPageTitle("Section -> Settings", m.continueSettingsModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case includingBlocks:
		return m.splitScreen(
			func() string { return withPageTitle("Section -> Blocks -> Settings", m.blockSettingsModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case continueBlockSettings:
		return m.splitScreen(
			func() string {
				return withPageTitle("Section -> Blocks -> Settings", m.continueBlockSettingsModel.View())
			},
			func() string { return m.viewportModel.View() },
		)
	case additionalBlocksPrompt:
		return m.splitScreen(
			func() string { return withPageTitle("Section -> Blocks", m.additionalBlocksPromptModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case includingBlockFileSettings:
		return withPageTitle("Block Settings", m.settingsModel.View())
	case continueBlockFileSettings:
		return withPageTitle("Block Settings", m.continueSettingsModel.View())
	case includingSnippetDoc:
		return m.splitScreen(
			func() string { return withPageTitle("Snippet Params", m.snippetDocModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case continueSnippetParams:
		return m.splitScreen(
			func() string { return withPageTitle("Snippet Params", m.continueSnippetDocModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case previewingFile:
		return m.splitScreen(m.previewFileModel.View, func() string { return m.viewportModel.View() })
	case creatingFile:
		return withPageTitle("Creating File...", "Let's do this!")
	case quitting:
		settingsBuilder := strings.Builder{}

		for _, sSettings := range m.settings {
			_, _ = settingsBuilder.WriteString(sSettings.String() + "\n")
		}
		for _, bSetting := range m.blockSettings {
			_, _ = settingsBuilder.WriteString(bSetting.String() + "\n")
		}
		if len(m.snippetDoc.Param) > 1 {
			_, _ = settingsBuilder.WriteString(m.snippetDoc.String() + "\n")
		}

		dbg(withPageTitle("bye for now!", settingsBuilder.String()))

		return withPageTitle("bye for now!", settingsBuilder.String())

	default:
		return "no view found..."
	}
}
