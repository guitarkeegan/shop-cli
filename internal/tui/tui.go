package tui

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guitarkeegan/shop-cli/internal/filepicker"
)

// will remain open for the lifetime of the program
var dbg = func() func(format string, as ...any) {
	if os.Getenv("DEBUG") == "" {
		return func(format string, as ...any) {}
	}
	file, err := os.OpenFile("logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	return func(format string, as ...any) {
		fmt.Fprintf(file, format+"\n", as...)
	}
}()

type Model struct {
	pwd          string
	filepicker   filepicker.Model
	selectedFile string
	fileContents string
	err          error
	width        int
	height       int
}

func InitialModel(currentDirectory string) Model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".liquid"}
	fp.CurrentDirectory = currentDirectory

	return Model{
		pwd:        currentDirectory,
		filepicker: fp,
	}
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.filepicker.SetHeight(msg.Height - 4)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	contents, err := m.filepicker.HoveredContent()
	if err != nil {
		m.err = err
		m.fileContents = fmt.Sprintf("Error: %v", err)
	} else {
		m.fileContents = contents
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		return m, tea.Quit
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid")
		m.selectedFile = ""
		return m, tea.Quit
	}

	return m, cmd
}

func (m Model) View() string {
	if m.selectedFile != "" {
		return "got the file!"
	}

	// Calculate split dimensions
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth - 2

	// Style for left pane (filepicker)
	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.height - 2).
		MaxHeight(m.height - 2)

	// Style for right pane (preview)
	previewStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height - 6).
		MaxHeight(m.height - 6).
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	// Get filepicker view
	pickerView := leftStyle.Render(m.filepicker.View())

	// Prepare preview content
	previewLines := strings.Split(m.fileContents, "\n")
	maxLines := max(m.height-6, 1)
	if len(previewLines) > maxLines {
		previewLines = previewLines[:maxLines]
		previewLines = append(previewLines, "...")
	}
	truncatedPreview := strings.Join(previewLines, "\n")

	// Render preview pane
	previewView := previewStyle.Render(truncatedPreview)

	// Combine horizontally
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, pickerView, previewView)

	// Add title and help
	title := lipgloss.NewStyle().
		Bold(true).
		Render("Select a file")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\nq: quit | ↑/↓: navigate | enter: select")

	return title + "\n" + mainView + help
}

func Run() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	m := InitialModel(currentDir)
	p := tea.NewProgram(&m, tea.WithAltScreen())

	tui, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("Alas, there's been an error: %v", err)
	}

	mm := tui.(Model)
	return mm.selectedFile, nil
}
