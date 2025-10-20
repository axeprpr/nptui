package ui

import (
	"fmt"
	"nptui/netplan"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11")). // Yellow
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("10")). // Green
			Foreground(lipgloss.Color("0")).  // Black
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")) // White

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Yellow
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")). // White
			Background(lipgloss.Color("4"))   // Blue

	disabledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Gray

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")). // White
			Background(lipgloss.Color("4")).  // Blue
			Padding(0, 1)
)

type screen int

const (
	mainMenuScreen screen = iota
	interfaceListScreen
	interfaceEditScreen
)

type model struct {
	screen      screen
	cursor      int
	config      *netplan.NetworkConfig
	interfaces  []string
	selectedIf  string
	configMode  string // "dhcp" or "static"
	ipAddress   string
	gateway     string
	dns         string
	editField   int // 0=config, 1=ip, 2=gateway, 3=dns
	message     string
	err         error
}

func initialModel() model {
	config, err := netplan.LoadConfig()
	if err != nil {
		return model{err: err}
	}

	return model{
		screen: mainMenuScreen,
		cursor: 0,
		config: config,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.screen {
		case mainMenuScreen:
			return m.updateMainMenu(msg)
		case interfaceListScreen:
			return m.updateInterfaceList(msg)
		case interfaceEditScreen:
			return m.updateInterfaceEdit(msg)
		}

	case tea.WindowSizeMsg:
		return m, nil
	}

	return m, nil
}

func (m model) updateMainMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 2 {
			m.cursor++
		}
	case "enter", " ":
		switch m.cursor {
		case 0: // Edit Network Interfaces
			interfaces, err := netplan.GetInterfaces()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.interfaces = interfaces
			m.screen = interfaceListScreen
			m.cursor = 0
		case 1: // Apply Configuration
			if err := netplan.ApplyConfig(); err != nil {
				m.message = fmt.Sprintf("Error: %v", err)
			} else {
				m.message = "Configuration applied successfully!"
			}
		case 2: // Quit
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateInterfaceList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc", "q", "b":
		m.screen = mainMenuScreen
		m.cursor = 0
		m.message = ""
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.interfaces) {
			m.cursor++
		}
	case "enter", " ":
		if m.cursor < len(m.interfaces) {
			m.selectedIf = m.interfaces[m.cursor]
			cfg := m.config.GetInterfaceConfig(m.selectedIf)

			// Initialize edit fields
			if cfg.DHCP4 {
				m.configMode = "dhcp"
			} else {
				m.configMode = "static"
			}

			m.ipAddress = ""
			if len(cfg.Addresses) > 0 {
				m.ipAddress = cfg.Addresses[0]
			}

			m.gateway = cfg.GetGateway()

			m.dns = ""
			if cfg.Nameservers != nil && len(cfg.Nameservers.Addresses) > 0 {
				m.dns = cfg.Nameservers.Addresses[0]
			}

			m.screen = interfaceEditScreen
			m.editField = 0
			m.message = ""
		} else {
			// Back option
			m.screen = mainMenuScreen
			m.cursor = 0
		}
	}
	return m, nil
}

func (m model) updateInterfaceEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.screen = interfaceListScreen
		m.cursor = 0
		m.message = ""
		return m, nil
	case "up", "k":
		if m.editField > 0 {
			m.editField--
		}
	case "down", "j":
		maxField := 3
		if m.configMode == "dhcp" {
			maxField = 0
		}
		if m.editField < maxField {
			m.editField++
		}
	case "tab":
		maxField := 3
		if m.configMode == "dhcp" {
			maxField = 0
		}
		m.editField = (m.editField + 1) % (maxField + 1)
	case "enter", " ":
		if m.editField == 0 {
			// Toggle DHCP/Static
			if m.configMode == "dhcp" {
				m.configMode = "static"
			} else {
				m.configMode = "dhcp"
			}
		}
	case "ctrl+s":
		// Save configuration
		return m.saveConfig(), nil
	case "backspace":
		if m.editField > 0 && m.configMode == "static" {
			switch m.editField {
			case 1:
				if len(m.ipAddress) > 0 {
					m.ipAddress = m.ipAddress[:len(m.ipAddress)-1]
				}
			case 2:
				if len(m.gateway) > 0 {
					m.gateway = m.gateway[:len(m.gateway)-1]
				}
			case 3:
				if len(m.dns) > 0 {
					m.dns = m.dns[:len(m.dns)-1]
				}
			}
		}
	default:
		// Handle text input
		if m.editField > 0 && m.configMode == "static" {
			if len(msg.String()) == 1 {
				char := msg.String()
				switch m.editField {
				case 1:
					m.ipAddress += char
				case 2:
					m.gateway += char
				case 3:
					m.dns += char
				}
			}
		}
	}
	return m, nil
}

func (m model) saveConfig() tea.Model {
	newConfig := netplan.EthernetConfig{}

	if m.configMode == "dhcp" {
		newConfig.DHCP4 = true
	} else {
		newConfig.DHCP4 = false
		if m.ipAddress != "" {
			newConfig.Addresses = []string{m.ipAddress}
		}
		if m.gateway != "" {
			newConfig.SetGateway(m.gateway)
		}
		if m.dns != "" {
			newConfig.Nameservers = &netplan.DNS{
				Addresses: []string{m.dns},
			}
		}
	}

	m.config.SetInterfaceConfig(m.selectedIf, newConfig)

	if err := netplan.SaveConfig(m.config); err != nil {
		m.message = fmt.Sprintf("Error: %v", err)
	} else {
		m.message = "Configuration saved! Press Ctrl+S to apply or Esc to go back."
	}

	return m
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to exit.\n", m.err)
	}

	switch m.screen {
	case mainMenuScreen:
		return m.viewMainMenu()
	case interfaceListScreen:
		return m.viewInterfaceList()
	case interfaceEditScreen:
		return m.viewInterfaceEdit()
	}

	return ""
}

func (m model) viewMainMenu() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Netplan TUI - Main Menu"))
	s.WriteString("\n\n")

	choices := []string{
		"Edit Network Interfaces",
		"Apply Configuration",
		"Quit",
	}

	for i, choice := range choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s.WriteString(selectedStyle.Render(cursor + " " + choice))
		} else {
			s.WriteString(normalStyle.Render(cursor + " " + choice))
		}
		s.WriteString("\n")
	}

	if m.message != "" {
		s.WriteString("\n" + labelStyle.Render(m.message) + "\n")
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render(" ↑↓: Navigate | Enter: Select | q: Quit "))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewInterfaceList() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Network Interfaces"))
	s.WriteString("\n\n")

	for i, iface := range m.interfaces {
		cfg := m.config.GetInterfaceConfig(iface)
		status := cfg.FormatConfig()

		cursor := " "
		if m.cursor == i {
			cursor = ">"
			s.WriteString(selectedStyle.Render(fmt.Sprintf("%s %s", cursor, iface)))
			s.WriteString("\n  " + normalStyle.Render(status) + "\n")
		} else {
			s.WriteString(normalStyle.Render(fmt.Sprintf("%s %s", cursor, iface)))
			s.WriteString("\n  " + disabledStyle.Render(status) + "\n")
		}
		s.WriteString("\n")
	}

	cursor := " "
	if m.cursor == len(m.interfaces) {
		cursor = ">"
		s.WriteString(selectedStyle.Render(cursor + " Back"))
	} else {
		s.WriteString(normalStyle.Render(cursor + " Back"))
	}
	s.WriteString("\n\n")

	s.WriteString(helpStyle.Render(" ↑↓: Navigate | Enter: Select | Esc: Back "))
	s.WriteString("\n")

	return s.String()
}

func (m model) viewInterfaceEdit() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render(fmt.Sprintf("Configure Interface: %s", m.selectedIf)))
	s.WriteString("\n\n")

	// Configuration mode
	cursor := " "
	if m.editField == 0 {
		cursor = ">"
	}
	configValue := "DHCP"
	if m.configMode == "static" {
		configValue = "Static"
	}

	if m.editField == 0 {
		s.WriteString(labelStyle.Render(cursor+" Configuration: ") + inputStyle.Render(" "+configValue+" "))
	} else {
		s.WriteString(labelStyle.Render(cursor+" Configuration: ") + normalStyle.Render(configValue))
	}
	s.WriteString("\n\n")

	if m.configMode == "static" {
		// IP Address
		cursor = " "
		if m.editField == 1 {
			cursor = ">"
		}
		if m.editField == 1 {
			s.WriteString(labelStyle.Render(cursor+" IP Address/CIDR: ") + inputStyle.Render(" "+m.ipAddress+" "))
		} else {
			s.WriteString(labelStyle.Render(cursor+" IP Address/CIDR: ") + normalStyle.Render(m.ipAddress))
		}
		s.WriteString("\n\n")

		// Gateway
		cursor = " "
		if m.editField == 2 {
			cursor = ">"
		}
		if m.editField == 2 {
			s.WriteString(labelStyle.Render(cursor+" Gateway: ") + inputStyle.Render(" "+m.gateway+" "))
		} else {
			s.WriteString(labelStyle.Render(cursor+" Gateway: ") + normalStyle.Render(m.gateway))
		}
		s.WriteString("\n\n")

		// DNS
		cursor = " "
		if m.editField == 3 {
			cursor = ">"
		}
		if m.editField == 3 {
			s.WriteString(labelStyle.Render(cursor+" DNS Server: ") + inputStyle.Render(" "+m.dns+" "))
		} else {
			s.WriteString(labelStyle.Render(cursor+" DNS Server: ") + normalStyle.Render(m.dns))
		}
		s.WriteString("\n\n")
	}

	if m.message != "" {
		s.WriteString(labelStyle.Render(m.message) + "\n\n")
	}

	s.WriteString(helpStyle.Render(" ↑↓: Navigate | Space: Toggle | Type to edit | Ctrl+S: Save | Esc: Back "))
	s.WriteString("\n")

	return s.String()
}

// App represents the main application
type App struct {
}

// NewApp creates a new application
func NewApp() *App {
	return &App{}
}

// Run starts the application
func (a *App) Run() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
