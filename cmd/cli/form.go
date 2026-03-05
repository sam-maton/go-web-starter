package main

import "github.com/charmbracelet/huh"

type model struct {
	folder string
	module string
	extras []string
}

func createForm(model *model) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter the project folder name").Value(&model.folder),
			huh.NewInput().Title("Enter the Go module name").Value(&model.module),
			huh.NewMultiSelect[string]().Options(
				huh.NewOption("Authentication", "auth").Selected(false),
				huh.NewOption("Air", "air").Selected(false),
			).Value(&model.extras),
		),
	)

	return form
}
