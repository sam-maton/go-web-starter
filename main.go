package main

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	questions   []string
	projectName string
	auth        bool
	input       textinput.Model
}

func NewModel(questions []string) model {
	return model{
		questions: questions,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() tea.View {
	v := tea.NewView("Welcome to the project initializer!\n")
	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func main() {
	m := NewModel([]string{
		"What is the name of your project?",
		"Do you want to set up authentication? (yes/no)",
	})
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
