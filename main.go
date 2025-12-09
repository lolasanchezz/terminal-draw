package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type model struct {
	mouseX int
	mouseY int
	width  int
	height int
	matrix [][]string
	click  string
}
type initCmd struct{}

func main() {
	m := model{}
	if _, err := tea.NewProgram(m, tea.WithMouseAllMotion(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return initCmd{} }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var _ tea.Cmd
	m.click = "not click"
	switch msg := msg.(type) {
	case initCmd:
		w, h, err := term.GetSize(os.Stdout.Fd())
		if err != nil {
			log.Fatal(err)
		}
		m.width = w
		m.height = h
		if w > 0 && h > 0 {
			m.matrix = make([][]string, h)
			for i := range h {
				m.matrix[i] = make([]string, w)
				for j := range w {
					m.matrix[i][j] = " "
					if j == w-1 {
						m.matrix[i][j] = "\n"
					}
				}
			}
		}
	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y
		switch msg.Action {
		case tea.MouseActionPress:
			m.click = "\npressing"
		case tea.MouseActionMotion:
			m.click = "\nmoving"
			m.matrix[m.mouseY][m.mouseX] = "#"
		case tea.MouseActionRelease:
			m.click = "\nrelease"
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	final := ""
	for i := range m.height {
		for j := range m.width {
			final += m.matrix[i][j]
		}
	}
	return final
	/*
		return "\nmouse x: " +
			strconv.Itoa(m.mouseX) +
			"\nmouse y:" +
			strconv.Itoa(m.mouseY) +
			"\nwidth " +
			strconv.Itoa(m.width) +
			"\nheight " +
			strconv.Itoa(m.height) +
			m.click
	*/

}
