package fulltui

type state string

const (
	selectingFile              state = "selectingFile"
	includingSettings          state = "includingSettings"
	continueSettings           state = "continueSettings"
	includingBlocks            state = "includingBlocks"
	continueBlockSettings      state = "continueBlockSettings"
	additionalBlocksPrompt     state = "additionalBlocksPrompt"
	includingBlockFileSettings state = "includingBlockFileSettings"
	continueBlockFileSettings  state = "continueBlockFileSettings"
	includingSnippetDoc        state = "includingSnippetDoc"
	continueSnippetParams      state = "continueSnippetParams"
	previewingFile             state = "previewingFile"
	creatingFile               state = "creatingfile"
	quitting                   state = "quitting"
)
