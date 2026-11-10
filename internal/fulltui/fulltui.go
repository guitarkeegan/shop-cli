package fulltui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var dbg = func() func(format string, as ...any) {
	if os.Getenv("DEBUG") == "" {
		return func(format string, as ...any) {}
	}
	return func(format string, as ...any) {
		f, err := os.OpenFile("shop-cli-logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return
		}
		defer f.Close()
		fmt.Fprintf(f, format+"\n", as...)
	}
}()

// Run starts the interactive file builder. It returns the path of the
// created file, or "" if the user quit without creating one.
func Run() (string, error) {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return "", err
	}
	m, ok := final.(model)
	if !ok {
		return "", nil
	}
	return m.created.path, m.created.err
}
