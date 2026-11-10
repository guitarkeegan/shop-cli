package fulltui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// The split screen must fit the terminal exactly. If it is taller, the
// terminal scrolls and the top border is cut off; if wider, the right
// border is clipped.
func TestSplitScreenFitsWindow(t *testing.T) {
	for _, size := range [][2]int{{80, 24}, {120, 40}, {200, 60}, {81, 25}} {
		w, h := size[0], size[1]
		t.Run(fmt.Sprintf("%dx%d", w, h), func(t *testing.T) {
			m := model{width: w, height: h}
			m.viewportModel = m.newViewport()
			m.viewportModel.SetContent(strings.Repeat("{% schema %}\n", 200))

			view := m.splitScreen(
				func() string { return "Section -> Settings\nInclude settings?" },
				m.viewportModel.View,
			)

			if got := lipgloss.Height(view); got != h {
				t.Errorf("height = %d, want %d", got, h)
			}
			if got := lipgloss.Width(view); got > w {
				t.Errorf("width = %d, want at most %d", got, w)
			}
			// Line 0 is the outer padding; the panes' top border follows.
			if top := strings.Split(view, "\n")[1]; strings.Count(top, "┌") != 2 {
				t.Errorf("second line should hold both panes' top borders: %q", top)
			}
		})
	}
}
