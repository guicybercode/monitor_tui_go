package editor

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	content  []string
	cursorX  int
	cursorY  int
	filePath string
	loaded   bool
	err      error
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursorY > 0 {
				m.cursorY--
			}
		case "down", "j":
			if m.cursorY < len(m.content)-1 {
				m.cursorY++
			}
		case "left", "h":
			if m.cursorX > 0 {
				m.cursorX--
			}
		case "right", "l":
			if m.cursorX < len(m.getCurrentLine())-1 {
				m.cursorX++
			}
		}

	case loadFileMsg:
		m.content = msg.lines
		m.filePath = msg.path
		m.loaded = true
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if !m.loaded {
		return "Press 'l' to load a file"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1)

	header := headerStyle.Render(fmt.Sprintf("Editor: %s (q: quit, arrow keys: navigate)", m.filePath))

	var lines []string
	lines = append(lines, header, "")

	for i, line := range m.content {
		if i == m.cursorY {
			beforeCursor := line[:min(m.cursorX, len(line))]
			atCursor := ""
			if m.cursorX < len(line) {
				atCursor = string(line[m.cursorX])
			}
			afterCursor := ""
			if m.cursorX+1 < len(line) {
				afterCursor = line[m.cursorX+1:]
			}
			cursorStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("230")).
				Foreground(lipgloss.Color("63"))
			line = beforeCursor + cursorStyle.Render(atCursor) + afterCursor
		}
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) getCurrentLine() string {
	if m.cursorY >= 0 && m.cursorY < len(m.content) {
		return m.content[m.cursorY]
	}
	return ""
}

type loadFileMsg struct {
	lines []string
	path  string
	err   error
}

func LoadFile(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := os.ReadFile(path)
		if err != nil {
			return loadFileMsg{
				lines: []string{},
				path:  path,
				err:   err,
			}
		}

		lines := []string{}
		currentLine := ""
		for _, b := range data {
			if b == '\n' {
				lines = append(lines, currentLine)
				currentLine = ""
			} else {
				currentLine += string(b)
			}
		}
		if currentLine != "" {
			lines = append(lines, currentLine)
		}

		return loadFileMsg{
			lines: lines,
			path:  path,
			err:   nil,
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
