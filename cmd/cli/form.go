package main

import "github.com/charmbracelet/huh"

type model struct {
	confirm bool
}

func createForm(model *model) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Create the baseline project in the current directory?").
				Description("This currently copies the go-web-starter-basline template as-is.").
				Value(&model.confirm),
		),
	)

	return form
}
