package fulltui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func ShopTheme() *huh.Theme {
	t := huh.ThemeBase()

	// White/black base with a warm grey ramp + gold accent
	var (
		// backgrounds / surfaces
		base     = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"} // "canvas"
		surface  = lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#111111"} // cards / subtle fill
		overlay1 = lipgloss.AdaptiveColor{Light: "#D0D0D0", Dark: "#2A2A2A"} // borders / separators
		overlay0 = lipgloss.AdaptiveColor{Light: "#B0B0B0", Dark: "#3A3A3A"} // placeholder / faint

		// text
		text     = lipgloss.AdaptiveColor{Light: "#111111", Dark: "#F2F2F2"}
		subtext1 = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#D6D6D6"}
		subtext0 = lipgloss.AdaptiveColor{Light: "#555555", Dark: "#B8B8B8"}

		// accents / status
		accent = lipgloss.AdaptiveColor{Light: "#F9BD09", Dark: "#F9BD09"} // gold
		green  = lipgloss.AdaptiveColor{Light: "#1F7A3A", Dark: "#2BD06A"}
		red    = lipgloss.AdaptiveColor{Light: "#B00020", Dark: "#FF4D4D"}

		// cursor (keep high contrast)
		cursor = lipgloss.AdaptiveColor{Light: "#111111", Dark: "#F9BD09"}
	)

	// Focused
	t.Focused.Base = t.Focused.Base.BorderForeground(overlay1)
	t.Focused.Card = t.Focused.Base

	t.Focused.Title = t.Focused.Title.Foreground(accent)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(accent)
	t.Focused.Directory = t.Focused.Directory.Foreground(accent)

	t.Focused.Description = t.Focused.Description.Foreground(subtext0)

	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)

	// selectors / nav
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(accent)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(accent)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(accent)

	// options
	t.Focused.Option = t.Focused.Option.Foreground(text)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(accent)

	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)

	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(text)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(text)

	// buttons
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(base).Background(accent)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(text).Background(surface)

	// text input
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(cursor)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(overlay0)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(accent)

	// Blurred
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base

	// Help (de-emphasized)
	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(subtext0)
	t.Help.ShortKey = t.Help.ShortKey.Foreground(subtext0)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(subtext1)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(subtext0)
	t.Help.FullKey = t.Help.FullKey.Foreground(subtext0)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(subtext1)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(subtext0)

	// Group
	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}

var pageTitle = lipgloss.NewStyle().
	Bold(true).
	Border(lipgloss.DoubleBorder(), true).
	Foreground(lipgloss.Color("#f9bd09")).
	Background(lipgloss.Color("#000")).
	AlignHorizontal(lipgloss.Left)

func withPageTitle(t string, view string) string {
	return fmt.Sprintf("%s\n%s", pageTitle.Render(t), view)
}

var (
	// Overall doc style (optional, just for padding/margins)
	docStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Left pane (main app)
	leftPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderRight(true).
			Padding(0, 1)

	// Right pane (preview)
	rightPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1)

	// Help text style for preview controls
	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

// Takes in two Views, the first will be the one on the left.
// The second will be on the right.
func (m model) splitScreen(views ...func() string) string {
	if len(views) < 2 {
		return "split screen needs two views!"
	}

	if m.width == 0 {
		return "loading..."
	}

	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	left := leftPaneStyle.
		Width(leftWidth).
		Height(m.height - 2). // small fudge for padding/borders; tweak as needed
		Render(views[0]())

	// Add help text at the top of the preview pane
	helpText := helpTextStyle.Render("j/k or ↑/↓ to scroll")
	previewContent := fmt.Sprintf("%s\n\n%s", helpText, views[1]())

	right := rightPaneStyle.
		Width(rightWidth).
		Height(m.height - 2).
		Render(previewContent)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		right,
	)

	// Optional outer padding
	return docStyle.Render(content)
}
