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

func previewCmd(render func() (string, error)) tea.Cmd {
	return func() tea.Msg {
		preview, err := render()
		return previewMsg{Preview: preview, Err: err}
	}
}

func showFilePreviewCmd(ctx context.Context, opts ...renderservice.SectionServiceOption) tea.Cmd {
	return previewCmd(func() (string, error) { return renderservice.CreateSectionPreview(ctx, opts...) })
}

func showBlockPreviewCmd(ctx context.Context, opts ...renderservice.RenderBlockServiceOption) tea.Cmd {
	return previewCmd(func() (string, error) { return renderservice.CreateBlockPreview(ctx, opts...) })
}

func showSnippetPreviewCmd(ctx context.Context, opts ...renderservice.SnippetOption) tea.Cmd {
	return previewCmd(func() (string, error) { return renderservice.CreateSnippetPreview(ctx, opts...) })
}

// fileCreatedMsg reports the result of writing a file. The program quits on
// receiving it, and Run reports the path or error once the screen is restored.
type fileCreatedMsg struct {
	path string
	err  error
}

func createFileCmd(create func() (string, error)) tea.Cmd {
	return func() tea.Msg {
		path, err := create()
		return fileCreatedMsg{path: path, err: err}
	}
}

func createSectionCmd(ctx context.Context, opts ...renderservice.SectionServiceOption) tea.Cmd {
	return createFileCmd(func() (string, error) { return renderservice.CreateSection(ctx, opts...) })
}

func createBlockCmd(ctx context.Context, opts ...renderservice.RenderBlockServiceOption) tea.Cmd {
	return createFileCmd(func() (string, error) { return renderservice.CreateBlock(ctx, opts...) })
}

func createSnippetCmd(ctx context.Context, opts ...renderservice.SnippetOption) tea.Cmd {
	return createFileCmd(func() (string, error) { return renderservice.CreateSnippet(ctx, opts...) })
}
