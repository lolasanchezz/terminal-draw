package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	mouseX int
	mouseY int
}

func main() {
	m := model{}
	if _, err := tea.NewProgram(m, tea.WithMouseAllMotion(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var _ tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return strconv.Itoa(m.mouseX) + strconv.Itoa(m.mouseY)
}
