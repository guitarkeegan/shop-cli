package featureservice

import (
	"errors"
	"io/fs"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var (
	ErrReadingTemplates = errors.New("error reading templates dir")
)

type Feature = string
type Filepath = string

type Model struct {
	templates fs.FS
	featMap   map[Feature][]Filepath
	form      *huh.Form
}

func makeFeatMap(tmpls fs.FS) (map[Feature][]Filepath, error) {
	featMap := make(map[Feature][]Filepath)

	entries, err := fs.ReadDir(tmpls, ".")
	if err != nil {
		return nil, ErrReadingTemplates
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		feature := entry.Name()
		fs.WalkDir(tmpls, feature, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			featMap[feature] = append(featMap[feature], path)
			return nil
		})
	}

	return featMap, nil

}

func InitialModel(featS fs.FS) Model {
	fserv, err := makeFeatMap(featS)
	if err != nil {
		log.Fatalf("initial model failed: %v\n", err)
	}
	m := Model{
		featMap: fserv,
	}
	m.form = buildForm(m)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.form.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		}
	}

	var formModel tea.Model
	formModel, cmd = m.form.Update(msg)
	m.form = formModel.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}

	return m, cmd
}
