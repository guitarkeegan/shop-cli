// Package fulltui is the full-screen interactive builder for sections,
// blocks and snippets, with a live preview of the generated file.
package fulltui

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
