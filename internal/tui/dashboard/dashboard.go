package dashboard

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/system"
)

type tickMsg time.Time

type Model struct {
	cpu     *system.CPUMetrics
	mem     *system.MemoryMetrics
	disk    *system.DiskMetrics
	network []system.NetworkStats
	loading bool
	err     error
}

func New() Model {
	return Model{loading: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), fetchMetrics())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m, tea.Batch(tick(), fetchMetrics())

	case metricsMsg:
		m.cpu = msg.cpu
		m.mem = msg.mem
		m.disk = msg.disk
		m.network = msg.network
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
		return "Loading metrics..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	var sections []string

	if m.cpu != nil {
		sections = append(sections, renderCPU(m.cpu))
	}

	if m.mem != nil {
		sections = append(sections, renderMemory(m.mem))
	}

	if m.disk != nil {
		sections = append(sections, renderDisk(m.disk))
	}

	if len(m.network) > 0 {
		sections = append(sections, renderNetwork(m.network))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func renderCPU(cpu *system.CPUMetrics) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1).
		Width(40)

	title := lipgloss.NewStyle().Bold(true).Render("CPU")
	usage := fmt.Sprintf("Usage: %.2f%%", cpu.Usage)
	cores := fmt.Sprintf("Cores: %d", cpu.Cores)
	model := fmt.Sprintf("Model: %s", truncate(cpu.Model, 30))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		usage,
		cores,
		model,
	)

	return boxStyle.Render(content)
}

func renderMemory(mem *system.MemoryMetrics) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1).
		Width(40)

	title := lipgloss.NewStyle().Bold(true).Render("Memory")
	total := fmt.Sprintf("Total: %s", formatBytes(mem.Total))
	used := fmt.Sprintf("Used: %s (%.2f%%)", formatBytes(mem.Used), mem.UsedPercent)
	available := fmt.Sprintf("Available: %s", formatBytes(mem.Available))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		total,
		used,
		available,
	)

	return boxStyle.Render(content)
}

func renderDisk(disk *system.DiskMetrics) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1).
		Width(40)

	title := lipgloss.NewStyle().Bold(true).Render("Disk")
	total := fmt.Sprintf("Total: %s", formatBytes(disk.Total))
	used := fmt.Sprintf("Used: %s (%.2f%%)", formatBytes(disk.Used), disk.UsedPercent)
	free := fmt.Sprintf("Free: %s", formatBytes(disk.Free))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		total,
		used,
		free,
	)

	return boxStyle.Render(content)
}

func renderNetwork(stats []system.NetworkStats) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(1).
		Width(50)

	title := lipgloss.NewStyle().Bold(true).Render("Network")
	var lines []string
	lines = append(lines, title, "")

	for _, stat := range stats[:min(5, len(stats))] {
		line := fmt.Sprintf("%s: ↑ %s ↓ %s",
			stat.Interface,
			formatBytes(stat.BytesSent),
			formatBytes(stat.BytesRecv))
		lines = append(lines, line)
	}

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

type metricsMsg struct {
	cpu     *system.CPUMetrics
	mem     *system.MemoryMetrics
	disk    *system.DiskMetrics
	network []system.NetworkStats
	err     error
}

func fetchMetrics() tea.Cmd {
	return func() tea.Msg {
		cpu, err1 := system.GetCPUMetrics()
		mem, err2 := system.GetMemoryMetrics()
		disk, err3 := system.GetDiskMetrics()
		network, err4 := system.GetNetworkStats()

		var err error
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			err = fmt.Errorf("failed to fetch metrics")
		}

		return metricsMsg{
			cpu:     cpu,
			mem:     mem,
			disk:    disk,
			network: network,
			err:     err,
		}
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
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

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
