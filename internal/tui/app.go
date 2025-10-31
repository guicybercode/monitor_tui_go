package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guicybercode/systui/internal/tui/dashboard"
	"github.com/guicybercode/systui/internal/tui/logs"
	"github.com/guicybercode/systui/internal/tui/network"
	"github.com/guicybercode/systui/internal/tui/packages"
	"github.com/guicybercode/systui/internal/tui/processes"
	"github.com/guicybercode/systui/internal/tui/services"
)

type View int

const (
	ViewDashboard View = iota
	ViewProcesses
	ViewServices
	ViewNetwork
	ViewEditor
	ViewPackages
	ViewLogs
)

type App struct {
	currentView View
	dashboard   dashboard.Model
	processes   processes.Model
	services    services.Model
	network     network.Model
	packages    packages.Model
	logs        logs.Model
	width       int
	height      int
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

func NewApp() *App {
	return &App{
		currentView: ViewDashboard,
		dashboard:   dashboard.New(),
		processes:   processes.New(),
		services:    services.New(),
		network:     network.New(),
		packages:    packages.New(),
		logs:        logs.New(),
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.dashboard.Init(),
		a.processes.Init(),
		a.services.Init(),
		a.network.Init(),
		a.packages.Init(),
		a.logs.Init(),
	)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			a.currentView = ViewDashboard
		case "2":
			a.currentView = ViewProcesses
		case "3":
			a.currentView = ViewServices
		case "4":
			a.currentView = ViewNetwork
		case "5":
			a.currentView = ViewPackages
		case "6":
			a.currentView = ViewLogs
		case "q", "ctrl+c":
			return a, tea.Quit
		}
	}

	var cmd tea.Cmd

	switch a.currentView {
	case ViewDashboard:
		var m tea.Model
		m, cmd = a.dashboard.Update(msg)
		a.dashboard = m.(dashboard.Model)
	case ViewProcesses:
		var m tea.Model
		m, cmd = a.processes.Update(msg)
		a.processes = m.(processes.Model)
	case ViewServices:
		var m tea.Model
		m, cmd = a.services.Update(msg)
		a.services = m.(services.Model)
	case ViewNetwork:
		var m tea.Model
		m, cmd = a.network.Update(msg)
		a.network = m.(network.Model)
	case ViewPackages:
		var m tea.Model
		m, cmd = a.packages.Update(msg)
		a.packages = m.(packages.Model)
	case ViewLogs:
		var m tea.Model
		m, cmd = a.logs.Update(msg)
		a.logs = m.(logs.Model)
	}

	return a, cmd
}

func (a *App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	var content string

	switch a.currentView {
	case ViewDashboard:
		content = a.dashboard.View()
	case ViewProcesses:
		content = a.processes.View()
	case ViewServices:
		content = a.services.View()
	case ViewNetwork:
		content = a.network.View()
	case ViewPackages:
		content = a.packages.View()
	case ViewLogs:
		content = a.logs.View()
	}

	header := titleStyle.Render("SysTUI - System Monitor")
	menu := a.renderMenu()

	contentHeight := a.height - 5
	renderedContent := lipgloss.NewStyle().
		Width(a.width).
		Height(contentHeight).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		menu,
		renderedContent,
		helpStyle.Render("Press 1-6 to switch views, q to quit"),
	)
}

func (a *App) renderMenu() string {
	menuItems := []string{"[1] Dashboard", "[2] Processes", "[3] Services", "[4] Network", "[5] Packages", "[6] Logs"}
	menu := ""
	for i, item := range menuItems {
		if i > 0 {
			menu += " | "
		}
		style := lipgloss.NewStyle()
		if int(a.currentView) == i {
			style = style.Bold(true).Foreground(lipgloss.Color("63"))
		}
		menu += style.Render(item)
	}
	return menu
}
