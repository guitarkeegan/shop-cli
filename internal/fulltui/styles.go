package fulltui

import (
	"fmt"
	"strings"

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
	// Outer padding around the split screen.
	docStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Left pane (forms).
	leftPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1)

	// Right pane (live preview).
	rightPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1)

	// Help text style for preview controls
	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

// previewHelp sits above the preview, followed by a blank line.
const (
	previewHelp      = "j/k or ↑/↓ to scroll"
	previewHelpLines = 2
)

// paneLayout is the size of each pane's content area, excluding borders.
type paneLayout struct {
	leftWidth, rightWidth, height int
}

// layout splits the window into two panes that, with their borders and the
// outer padding, exactly fill it. lipgloss adds borders outside a style's
// Width and Height, so they are subtracted here.
func (m model) layout() paneLayout {
	frameW, frameH := docStyle.GetFrameSize()
	borderW := leftPaneStyle.GetHorizontalBorderSize() + rightPaneStyle.GetHorizontalBorderSize()

	contentW := max(m.width-frameW-borderW, 0)
	return paneLayout{
		leftWidth:  contentW / 2,
		rightWidth: contentW - contentW/2,
		height:     max(m.height-frameH-leftPaneStyle.GetVerticalBorderSize(), 0),
	}
}

// previewSize is the viewport size that fits inside the right pane, below
// the help text.
func (m model) previewSize() (width, height int) {
	l := m.layout()
	return max(l.rightWidth-rightPaneStyle.GetHorizontalPadding(), 0),
		max(l.height-rightPaneStyle.GetVerticalPadding()-previewHelpLines, 0)
}

// clip cuts s to at most n lines, so a tall form can't push the layout
// past the bottom of the window.
func clip(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}

// splitScreen renders the left view beside the right one, filling the window.
func (m model) splitScreen(views ...func() string) string {
	if len(views) < 2 {
		return "split screen needs two views!"
	}

	if m.width == 0 {
		return "loading..."
	}

	l := m.layout()

	left := leftPaneStyle.
		Width(l.leftWidth).
		Height(l.height).
		Render(clip(views[0](), l.height))

	right := rightPaneStyle.
		Width(l.rightWidth).
		Height(l.height).
		Render(clip(helpTextStyle.Render(previewHelp)+"\n\n"+views[1](), l.height))

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
}
