package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type model struct {
	mouseX   int
	mouseY   int
	width    int
	height   int
	matrix   [][]string
	clicking bool
	brush    string
	toolbar  toolbarModel
}
type initCmd struct{}

var cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(0))

func main() {
	m := model{}
	m.brush = "#"
	if _, err := tea.NewProgram(m, tea.WithMouseAllMotion(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg { return initCmd{} }
}

func makeMatrix(w int, h int) [][]string {
	matrix := make([][]string, h)
	for i := range h {
		matrix[i] = make([]string, w)
		for j := range w {
			matrix[i][j] = " "
			if j == w-1 {
				matrix[i][j] = "\n"
			}
		}
	}
	return matrix
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var _ tea.Cmd
	switch msg := msg.(type) {

	case initCmd:
		w, h, err := term.GetSize(os.Stdout.Fd())
		if err != nil {
			log.Fatal(err)
		}
		m.width = w
		m.height = h

		padding := 2
		height := 6
		m.toolbar.height = height
		//	m.toolbar.width = w - (2 * padding)
		toolbarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(m.width-(2*padding)).
			Padding(height-5, 0).
			Align(lipgloss.Center)

		if w > 0 && h > 0 {
			m.matrix = makeMatrix(m.width, m.height-m.toolbar.height)
		}

	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y
		switch msg.Action {
		case tea.MouseActionPress:
			//checking whether the action was within toolbar or easel
			if msg.X > m.width+1 {
			}
			m.clicking = true
		case tea.MouseActionRelease:
			m.clicking = false
		case tea.MouseActionMotion:
			//checking whether the action was within toolbar or easel
			if msg.X < len(m.matrix[0]) { //just in case
				if msg.Y < m.toolbar.height+1 {

				} else {
					if m.clicking {
						if msg.Y < len(m.matrix) { // you never know ... could be a weird mouse spasm that would cause out of bounds and crash, so this is just in case
							m.matrix[m.mouseY+-+m.toolbar.height][m.mouseX] = cursorStyle.Render(m.brush)
						}
					}
				}
			}

		}

	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		}
		switch msg.String() {
		case "r":
			cursorStyle = cursorStyle.Foreground(lipgloss.Color("#ff0000"))
		case "w":
			cursorStyle = cursorStyle.Foreground(lipgloss.Color("#ffffffff"))
		case "e":
			m.brush = " "
		case "b":
			m.brush = "#"
		case "c":
			m.matrix = makeMatrix(m.width, m.height-m.toolbar.height)
		}
	}
	return m, nil
}

func (m model) View() string {
	final := ""
	for i := range len(m.matrix) {
		for j := range len(m.matrix[0]) {
			final += m.matrix[i][j]
		}
	}
	return lipgloss.JoinVertical(lipgloss.Center, m.toolbarView(), final)
}

type colorOption struct {
	color string
	val   string
}

var toolbar = struct {
	colors       []colorOption
	strokes      []string
	width        []string
	padding      string
	interPadding string
}{
	colors: []colorOption{
		{color: "red", val: "#ff0000"},
		{color: "blue", val: "#0000ff"},
		{color: "green", val: "#00ff00"},
	},
	strokes:      []string{"#", ".", "-"},
	width:        []string{"◼", "◼◼", "◼◼◼"},
	padding:      "    ",
	interPadding: " ",
}

type toolbarModel struct {
	width  int
	height int
}

var toolbarStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

func (m model) toolbarView() string {
	// colors

	colorChar := "⬤"
	colorStr := toolbar.padding

	for _, color := range toolbar.colors {
		colorStr = colorStr + toolbar.interPadding + toolbar.interPadding + lipgloss.NewStyle().Foreground(lipgloss.Color(color.val)).Render(colorChar)
	}
	colorStr += toolbar.padding
	strokeStr := toolbar.padding
	for _, stroke := range toolbar.strokes {
		strokeStr = strokeStr + toolbar.interPadding + stroke
	}
	strokeStr += toolbar.padding

	widthStr := toolbar.padding
	for _, width := range toolbar.width {
		widthStr = widthStr + toolbar.interPadding + width
	}
	widthStr += toolbar.padding

	return toolbarStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, colorStr, strokeStr, widthStr))
}
