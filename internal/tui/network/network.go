package network

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/system"
)

type Model struct {
	stats       []system.NetworkStats
	connections []system.NetworkConnection
	selected    int
	loading     bool
	err         error
}

func New() Model {
	return Model{loading: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchStats, fetchConnections)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selected < len(m.stats)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "r":
			return m, tea.Batch(fetchStats, fetchConnections)
		}

	case statsMsg:
		m.stats = msg.stats
		m.loading = false
		m.err = msg.err
		return m, nil

	case connectionsMsg:
		m.connections = msg.connections
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
		return "Loading network info..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1)

	header := headerStyle.Render("Network Monitor (r: refresh)")

	var sections []string
	sections = append(sections, header, "")

	statsSection := renderStats(m.stats, m.selected)
	sections = append(sections, statsSection)

	if len(m.connections) > 0 {
		sections = append(sections, "", renderConnections(m.connections))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func renderStats(stats []system.NetworkStats, selected int) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1)

	title := lipgloss.NewStyle().Bold(true).Render("Network Interfaces")
	var lines []string
	lines = append(lines, title, "")

	headerRow := fmt.Sprintf("%-15s %12s %12s %8s %8s",
		"Interface", "Bytes Sent", "Bytes Recv", "Errors", "Dropped")
	lines = append(lines, headerRow)
	lines = append(lines, "")

	for i, stat := range stats {
		style := lipgloss.NewStyle()
		if i == selected {
			style = style.Background(lipgloss.Color("63")).Foreground(lipgloss.Color("230"))
		}

		line := fmt.Sprintf("%-15s %12s %12s %8d %8d",
			stat.Interface,
			formatBytes(stat.BytesSent),
			formatBytes(stat.BytesRecv),
			stat.ErrorsIn+stat.ErrorsOut,
			stat.DropIn+stat.DropOut)
		lines = append(lines, style.Render(line))
	}

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func renderConnections(conns []system.NetworkConnection) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1)

	title := lipgloss.NewStyle().Bold(true).Render("Active Connections")
	var lines []string
	lines = append(lines, title, "")

	headerRow := fmt.Sprintf("%-8s %-20s %-20s %-8s",
		"Status", "Local", "Remote", "PID")
	lines = append(lines, headerRow)
	lines = append(lines, "")

	for _, conn := range conns[:min(20, len(conns))] {
		local := fmt.Sprintf("%s:%d", conn.LocalAddr, conn.LocalPort)
		remote := fmt.Sprintf("%s:%d", conn.RemoteAddr, conn.RemotePort)
		line := fmt.Sprintf("%-8s %-20s %-20s %-8d",
			conn.Status, local, remote, conn.PID)
		lines = append(lines, line)
	}

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

type statsMsg struct {
	stats []system.NetworkStats
	err   error
}

type connectionsMsg struct {
	connections []system.NetworkConnection
	err         error
}

func fetchStats() tea.Msg {
	stats, err := system.GetNetworkStats()
	return statsMsg{
		stats: stats,
		err:   err,
	}
}

func fetchConnections() tea.Msg {
	connections, err := system.GetNetworkConnections()
	return connectionsMsg{
		connections: connections,
		err:         err,
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
