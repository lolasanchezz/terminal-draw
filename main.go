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
	m := model{
		brush: "#",
		toolbar: toolbarModel{
			hitboxes:        make(map[string][]int),
			visibleElements: []string{"colors", "strokes", "width"},
		},
	}
	f, err := tea.LogToFile("debug.log", "debug")
	f.Truncate(0)
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	if _, err := tea.NewProgram(m, tea.WithMouseAllMotion(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func (m model) Init() tea.Cmd {
	m.toolbar.visibleElements = []string{"colors", "strokes", "width"}
	m.toolbar.hitboxes = make(map[string][]int)
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

		//padding := 2
		height := 5
		m.toolbar.height = height
		//	m.toolbar.width = w - (2 * padding)
		toolbarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 0). // 1 line top, 1 line bottom
			Height(3).
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

			m, _ = m.toolbarUpdate(msg)

			m.clicking = true
		case tea.MouseActionRelease:
			m.clicking = false
		case tea.MouseActionMotion:
			//checking whether the action was within toolbar or easel
			if msg.X < len(m.matrix[0]) { //just in case
				if msg.Y < m.toolbar.height {
					m, _ = m.toolbarUpdate(msg)
				} else {
					if m.clicking {
						if msg.Y < len(m.matrix) { // you never know ... could be a weird mouse spasm that would cause out of bounds and crash, so this is just in case
							m.matrix[m.mouseY-m.toolbar.height][m.mouseX-1] = cursorStyle.Render(m.brush)
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
			m.brush = "█"
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

type toolbarEntry struct {
	name   string
	values []string
}

var toolbar = struct {
	elements     []toolbarEntry
	padding      string
	interPadding string
}{
	elements: []toolbarEntry{
		{name: "colors", values: []string{"#ff0000", "#0000ff", "#00ff00"}},
		{name: "strokes", values: []string{"#", ".", "-", "█"}},
		{name: "width", values: []string{"◼", "◼◼", "◼◼◼"}},
	},
	padding:      "    ",
	interPadding: "   ",
}

type toolbarModel struct {
	width           int
	height          int
	hitboxes        map[string][]int
	visibleElements []string
}

var toolbarStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Background(lipgloss.Color("#485356"))

func (m model) toolbarInit() model {
	m.toolbar.visibleElements = []string{"colors", "strokes", "width"}
	m.toolbar.hitboxes = make(map[string][]int, len(m.toolbar.visibleElements))
	return m
}

func (m model) toolbarView() string {
	// colors

	colorChar := "⬤"
	finalArr := []string{toolbar.padding}
	//special case for colors
	for _, color := range toolbar.elements[0].values {
		finalArr = append(finalArr, (lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(colorChar)), toolbar.interPadding)
	}

	for _, elem := range toolbar.elements[1:] {
		finalArr = append(finalArr, (toolbar.padding + toolbar.padding))
		for _, val := range elem.values {
			finalArr = append(finalArr, (toolbar.interPadding + val))
		}
	}
	finalArr = append(finalArr, toolbar.padding)

	//calculating hitboxes
	finalStr := lipgloss.JoinHorizontal(lipgloss.Center, finalArr...)

	lenOfStrs := lipgloss.Width(finalStr)

	offset := m.width/2 - (lenOfStrs / 2) + lipgloss.Width(toolbar.padding)

	for _, element := range toolbar.elements {
		m.toolbar.hitboxes[element.name] = []int{}
		for _, char := range element.values {
			rendered := ""
			var charLen int
			if element.name == "colors" {
				rendered = toolbar.interPadding +
					lipgloss.NewStyle().Foreground(lipgloss.Color(char)).Render(colorChar)
			} else {
				rendered = toolbar.interPadding + char
			}

			charLen = lipgloss.Width(rendered)
			m.toolbar.hitboxes[element.name] = append(m.toolbar.hitboxes[element.name], offset)
			offset += charLen
		}
		offset += lipgloss.Width(toolbar.padding)*2 + lipgloss.Width(toolbar.interPadding)
	}

	//for TESTING
	if len(m.matrix) != 0 {
		for _, xArr := range m.toolbar.hitboxes {
			for _, x := range xArr {
				m.matrix[0][x] = "*"
			}
		}
	}
	return toolbarStyle.Render(finalStr)
}

func (m model) toolbarUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionPress:
			for name, hitbox := range m.toolbar.hitboxes {
				for i, xCoord := range hitbox {
					if ((msg.X == xCoord) || (msg.X+1 == xCoord) || (msg.X-1 == xCoord)) && (msg.Y == m.toolbar.height/2-1) {
						m.readHitboxes(name, xCoord, i)
					}
				}
			}
		}

	}
	return m, nil
}

func (m *model) readHitboxes(key string, x int, i int) {
	//manual implementations for each thing that has to change within the hitbox
	switch key {
	case "colors":
		cursorStyle = cursorStyle.Foreground(lipgloss.Color(toolbar.elements[0].values[i]))
	case "strokes":
		m.brush = toolbar.elements[1].values[i]
	}
}
