package fulltui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	renderservice "github.com/guitarkeegan/shop-cli/internal/render-service"
)

type previewMsg struct {
	Preview string
	Err     error
}

func showFilePreviewCmd(ctx context.Context, opts ...renderservice.SectionServiceOption) tea.Cmd {
	return func() tea.Msg {
		preview, err := renderservice.CreateSectionPreview(ctx, opts...)
		return previewMsg{
			Preview: preview,
			Err:     err,
		}
	}
}

func showSnippetPreviewCmd(ctx context.Context, opts ...renderservice.SnippetOption) tea.Cmd {
	return func() tea.Msg {
		preview, err := renderservice.CreateSnippetPreview(ctx, opts...)
		return previewMsg{
			Preview: preview,
			Err:     err,
		}
	}
}

type isFirstUpdateMsg bool

func isFirstUpdateCmd(isFirstUpdate bool) func() tea.Msg {
	return func() tea.Msg {
		return isFirstUpdateMsg(isFirstUpdate)
	}
}

type isFirstBlockMsg bool

func isFirstBlockCmd(isFirstBlock bool) func() tea.Msg {
	return func() tea.Msg {
		return isFirstBlockMsg(isFirstBlock)
	}
}

type isAdditionalUpdateMsg bool

func isAdditionalUpdateCmd(isAdditionalUpdate bool) func() tea.Msg {
	return func() tea.Msg {
		return isAdditionalUpdateMsg(isAdditionalUpdate)
	}
}

type createSectionMsg struct{}

func createSectionCmd(ctx context.Context, opts ...renderservice.SectionServiceOption) tea.Cmd {
	return func() tea.Msg {
		err := renderservice.CreateSection(ctx, opts...)(ctx)
		if err != nil {
			dbg("panic: %s", err)
			panic("create section failed")
		}
		return createSectionMsg{}
	}
}

type createSnippetMsg struct{}

func createSnippetCmd(ctx context.Context, opts ...renderservice.SnippetOption) tea.Cmd {
	return func() tea.Msg {
		err := renderservice.CreateSnippet(ctx, opts...)(ctx)
		if err != nil {
			dbg("panic: %s", err)
			panic("create snippet failed")
		}
		return createSectionMsg{}
	}
}
