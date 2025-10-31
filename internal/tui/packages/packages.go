package packages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/system"
)

type Model struct {
	packages       []system.PackageInfo
	packageManager system.PackageManager
	selected       int
	loading        bool
	err            error
}

func New() Model {
	return Model{loading: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(detectPackageManager, fetchPackages)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selected < len(m.packages)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "r":
			return m, fetchPackages
		}

	case packageManagerMsg:
		m.packageManager = msg.pm
		return m, nil

	case packagesMsg:
		m.packages = msg.packages
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
		return "Loading packages..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1)

	pmName := string(m.packageManager)
	header := headerStyle.Render(fmt.Sprintf("Packages (%s) - j/k: navigate, r: refresh", pmName))

	var lines []string
	lines = append(lines, header, "")

	headerRow := fmt.Sprintf("%-30s %-15s %s",
		"Name", "Version", "Description")
	lines = append(lines, headerRow)
	lines = append(lines, "")

	for i, pkg := range m.packages {
		style := lipgloss.NewStyle()
		if i == m.selected {
			style = style.Background(lipgloss.Color("63")).Foreground(lipgloss.Color("230"))
		}

		name := pkg.Name
		if len(name) > 30 {
			name = name[:27] + "..."
		}

		desc := pkg.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}

		line := fmt.Sprintf("%-30s %-15s %s",
			name, pkg.Version, desc)
		lines = append(lines, style.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

type packageManagerMsg struct {
	pm system.PackageManager
}

type packagesMsg struct {
	packages []system.PackageInfo
	err      error
}

func detectPackageManager() tea.Msg {
	pm := system.DetectPackageManager()
	return packageManagerMsg{pm: pm}
}

func fetchPackages() tea.Msg {
	pm := system.DetectPackageManager()
	packages, err := system.ListPackages(pm)
	return packagesMsg{
		packages: packages,
		err:      err,
	}
}
