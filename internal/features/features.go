package features

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type FeatureService struct {
	Templates fs.FS
}

func New(templates fs.FS) *FeatureService {
	return &FeatureService{
		Templates: templates,
	}
}

func (fs *FeatureService) Run() error {
	m, err := tea.NewProgram(InitialModel(fs.Templates)).Run()
	if err != nil {
		return err
	}
	model := m.(Model)
	selected, ok := model.form.Get("feature").(string)
	if !ok {
		return fmt.Errorf("form did not return a string")
	}
	files, ok := model.featMap[selected]
	if !ok {
		return fmt.Errorf("feature %q not found", selected)
	}
	if err = CreateFeature(fs.Templates, files); err != nil {
		return err
	}

	return nil
}

// renders the files, prompts the user of the files already exist
func CreateFeature(templates fs.FS, filepaths []Filepath) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}
	for _, required := range []string{"sections", "src"} {
		info, err := os.Stat(filepath.Join(cwd, required))
		if err != nil || !info.IsDir() {
			return fmt.Errorf("required directory %q not found in current directory", required)
		}
	}

	for _, fp := range filepaths {
		// fp is like "feature-name/sections/file.liquid" — strip the feature prefix for the destination
		destRel := fp[strings.Index(fp, "/")+1:]
		destPath := filepath.Join(cwd, destRel)

		_, err := os.Stat(destPath)
		if err == nil {
			action, err := promptConflict(destRel)
			if err != nil {
				return err
			}
			switch action {
			case "skip":
				continue
			case "override":
				if err := writeFile(templates, fp, destPath); err != nil {
					return err
				}
			case "keep both":
				newPath := resolveKeepBoth(destPath)
				if err := writeFile(templates, fp, newPath); err != nil {
					return err
				}
			}
			continue
		}

		// file doesn't exist, create it
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("could not create directories for %s: %w", destPath, err)
		}
		if err := writeFile(templates, fp, destPath); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(templates fs.FS, srcPath, destPath string) error {
	data, err := fs.ReadFile(templates, srcPath)
	if err != nil {
		return fmt.Errorf("could not read template %s: %w", srcPath, err)
	}
	return os.WriteFile(destPath, data, 0644)
}

func resolveKeepBoth(destPath string) string {
	ext := filepath.Ext(destPath)
	base := strings.TrimSuffix(destPath, ext)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s(%d)%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

func promptConflict(fp string) (string, error) {
	var result string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("%q already exists", fp)).
				Options(
					huh.NewOption("Skip", "skip"),
					huh.NewOption("Override", "override"),
					huh.NewOption("Keep both", "keep both"),
				).
				Value(&result),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return result, nil
}
