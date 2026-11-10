package featureservice

import "github.com/charmbracelet/huh"

func buildForm(m Model) *huh.Form {
	options := make([]huh.Option[string], 0, len(m.featMap))
	for feature := range m.featMap {
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
