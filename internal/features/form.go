package features

import (
	"maps"
	"slices"

	"github.com/charmbracelet/huh"
)

func buildForm(m Model) *huh.Form {
	features := slices.Sorted(maps.Keys(m.featMap))
	options := make([]huh.Option[string], 0, len(features))
	for _, feature := range features {
		options = append(options, huh.NewOption(feature, feature))
	}

	var selected string
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("feature").
				Title("Select a feature").
				Options(options...).
				Value(&selected),
		),
	)
}
