package main

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
)

type model struct {
	folder  string
	confirm bool
}

func createForm(model *model) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project folder name").
				Description("The baseline project will be created inside this folder in the current directory.").
				Placeholder("my-project").
				Value(&model.folder).
				Validate(func(value string) error {
					if strings.TrimSpace(value) == "" {
						return errors.New("project folder name is required")
					}

					return nil
				}),
			huh.NewConfirm().
				Title("Create the baseline project in this folder?").
				Description("This currently copies the go-web-starter-basline template as-is.").
				Value(&model.confirm),
		),
	)

	return form
}
