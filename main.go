package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type stage int

const (
	stageNameInput stage = iota
	stageOptions
	stageDone
)

type model struct {
	stage       stage
	projectName string
	choices     []string
	cursor      int
	selected    map[int]struct{}
}

func initialModel() model {
	return model{
		stage:    stageNameInput,
		choices:  []string{"Option 1", "Option 2", "Option 3"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch m.stage {
		case stageNameInput:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				if strings.TrimSpace(m.projectName) != "" {
					m.stage = stageOptions
				}
			case "backspace":
				if len(m.projectName) > 0 {
					m.projectName = m.projectName[:len(m.projectName)-1]
				}
			default:
				if msg.Text != "" {
					m.projectName += msg.Text
				}
			}

		case stageOptions:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "space":
				if _, ok := m.selected[m.cursor]; ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			case "enter":
				m.stage = stageDone
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var s string

	switch m.stage {
	case stageNameInput:
		s = "What is your project name?\n\n"
		s += "> " + m.projectName + "_\n\n"
		s += "Press Enter to continue.\n"

	case stageOptions:
		s = fmt.Sprintf("Project: %s\n\n", m.projectName)
		s += "Select options (Space to toggle, Enter to confirm):\n\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
		s += "\nPress Enter to confirm.\n"
	}

	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	m, ok := finalModel.(model)
	if !ok || m.stage != stageDone {
		return
	}

	fmt.Printf("\nProject name: %s\n", m.projectName)
	fmt.Println("Selected options:")
	for i, choice := range m.choices {
		if _, ok := m.selected[i]; ok {
			fmt.Printf("  - %s\n", choice)
		}
	}
	if len(m.selected) == 0 {
		fmt.Println("  (none)")
	}
}
