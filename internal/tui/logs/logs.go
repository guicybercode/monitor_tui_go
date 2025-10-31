package logs

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/logparser"
)

type Model struct {
	entries  []logparser.LogEntry
	selected int
	logPath  string
	pattern  string
	loading  bool
	err      error
}

func New() Model {
	return Model{
		logPath: "/var/log/syslog",
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchLogs
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selected < len(m.entries)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "r":
			return m, fetchLogs
		}

	case logsMsg:
		m.entries = msg.entries
		m.loading = false
		m.err = msg.err
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if m.loading {
		return "Loading logs..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1)

	header := headerStyle.Render(fmt.Sprintf("Log Viewer: %s (j/k: navigate, r: refresh)", m.logPath))

	var lines []string
	lines = append(lines, header, "")

	headerRow := fmt.Sprintf("%-20s %-8s %s",
		"Timestamp", "Severity", "Message")
	lines = append(lines, headerRow)
	lines = append(lines, "")

	for i, entry := range m.entries {
		style := lipgloss.NewStyle()
		if i == m.selected {
			style = style.Background(lipgloss.Color("63")).Foreground(lipgloss.Color("230"))
		}

		severityColor := lipgloss.Color("241")
		switch entry.Severity {
		case "ERROR":
			severityColor = lipgloss.Color("196")
		case "WARN":
			severityColor = lipgloss.Color("220")
		case "INFO":
			severityColor = lipgloss.Color("39")
		}

		severityStyle := style.Copy().Foreground(severityColor)
		msg := entry.Message
		if len(msg) > 50 {
			msg = msg[:47] + "..."
		}

		line := fmt.Sprintf("%-20s %-8s %s",
			entry.Timestamp, entry.Severity, msg)
		lines = append(lines, severityStyle.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

type logsMsg struct {
	entries []logparser.LogEntry
	err     error
}

func fetchLogs() tea.Msg {
	entries, err := logparser.ParseLogs("/var/log/syslog", "", "", "", "")
	if err != nil {
		return logsMsg{
			entries: []logparser.LogEntry{},
			err:     err,
		}
	}

	return logsMsg{
		entries: entries,
		err:     nil,
	}
}
