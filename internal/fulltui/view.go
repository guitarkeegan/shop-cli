package fulltui

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
		return m.splitScreen(
			func() string { return withPageTitle("Block -> Settings", m.settingsModel.View()) },
			func() string { return m.viewportModel.View() },
		)
	case continueBlockFileSettings:
		return m.splitScreen(
			func() string { return withPageTitle("Block -> Settings", m.continueSettingsModel.View()) },
			func() string { return m.viewportModel.View() },
		)
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
		// Run prints the outcome after the alt screen is restored.
		return ""

	default:
		return "no view found..."
	}
}
