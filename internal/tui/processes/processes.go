package processes

import (
	"fmt"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/system"
)

type Model struct {
	processes []system.ProcessInfo
	selected  int
	loading   bool
	err       error
}

func New() Model {
	return Model{loading: true}
}

func (m Model) Init() tea.Cmd {
	return fetchProcesses
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selected < len(m.processes)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "d":
			if len(m.processes) > 0 && m.selected < len(m.processes) {
				pid := m.processes[m.selected].PID
				err := system.KillProcess(pid, syscall.SIGTERM)
				if err == nil {
					return m, fetchProcesses
				}
			}
		case "r":
			return m, fetchProcesses
		}

	case processesMsg:
		m.processes = msg.processes
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
		return "Loading processes..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1)

	header := headerStyle.Render("Processes (j/k: navigate, d: kill, r: refresh)")

	var lines []string
	lines = append(lines, header, "")

	headerRow := fmt.Sprintf("%-8s %-20s %8s %8s %-12s",
		"PID", "Name", "CPU%", "Mem%", "User")
	lines = append(lines, headerRow)
	lines = append(lines, "")

	for i, proc := range m.processes {
		style := lipgloss.NewStyle()
		if i == m.selected {
			style = style.Background(lipgloss.Color("63")).Foreground(lipgloss.Color("230"))
		}

		name := proc.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		line := fmt.Sprintf("%-8d %-20s %8.2f %8.2f %-12s",
			proc.PID, name, proc.CPUPercent, proc.MemPercent, proc.User)
		lines = append(lines, style.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

type processesMsg struct {
	processes []system.ProcessInfo
	err       error
}

func fetchProcesses() tea.Msg {
	processes, err := system.GetProcesses()
	return processesMsg{
		processes: processes,
		err:       err,
	}
}
