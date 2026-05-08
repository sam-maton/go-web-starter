package main

import (
"errors"
"strings"

"github.com/charmbracelet/huh"
)

type model struct {
folder     string
moduleName string
enableAuth bool
confirm    bool
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
huh.NewInput().
Title("Go module name").
Description("This module path will replace the scaffold module in go.mod.").
Placeholder("github.com/you/my-project").
Value(&model.moduleName).
Validate(func(value string) error {
if strings.TrimSpace(value) == "" {
return errors.New("go module name is required")
}

return nil
}),
huh.NewConfirm().
Title("Enable authentication?").
Description("Adds user signup, login/logout, and route protection. Can be removed later.").
Value(&model.enableAuth),
huh.NewConfirm().
Title("Create the baseline project in this folder?").
Description("This copies the baseline template and sets the go.mod module path.").
Value(&model.confirm),
),
)

return form
}
