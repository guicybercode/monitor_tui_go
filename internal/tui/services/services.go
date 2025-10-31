package services

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/system"
)

type Model struct {
	services []system.ServiceInfo
	selected int
	loading  bool
	err      error
}

func New() Model {
	return Model{loading: true}
}

func (m Model) Init() tea.Cmd {
	return fetchServices
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.selected < len(m.services)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "s":
			if len(m.services) > 0 && m.selected < len(m.services) {
				name := m.services[m.selected].Name
				err := system.StartService(name)
				if err == nil {
					return m, fetchServices
				}
			}
		case "x":
			if len(m.services) > 0 && m.selected < len(m.services) {
				name := m.services[m.selected].Name
				err := system.StopService(name)
				if err == nil {
					return m, fetchServices
				}
			}
		case "t":
			if len(m.services) > 0 && m.selected < len(m.services) {
				name := m.services[m.selected].Name
				err := system.RestartService(name)
				if err == nil {
					return m, fetchServices
				}
			}
		case "r":
			return m, fetchServices
		}

	case servicesMsg:
		m.services = msg.services
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
		return "Loading services..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Padding(0, 1)

	header := headerStyle.Render("Services (j/k: navigate, s: start, x: stop, t: restart, r: refresh)")

	var lines []string
	lines = append(lines, header, "")

	headerRow := fmt.Sprintf("%-35s %-12s %s",
		"Name", "State", "Description")
	lines = append(lines, headerRow)
	lines = append(lines, "")

	for i, svc := range m.services {
		style := lipgloss.NewStyle()
		if i == m.selected {
			style = style.Background(lipgloss.Color("63")).Foreground(lipgloss.Color("230"))
		}

		name := svc.Name
		if len(name) > 35 {
			name = name[:32] + "..."
		}

		desc := svc.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}

		line := fmt.Sprintf("%-35s %-12s %s",
			name, svc.State, desc)
		lines = append(lines, style.Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

type servicesMsg struct {
	services []system.ServiceInfo
	err      error
}

func fetchServices() tea.Msg {
	services, err := system.GetServices()
	return servicesMsg{
		services: services,
		err:      err,
	}
}
