package ui

import (
	"fmt"
	"nptui/netplan"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	app    *tview.Application
	pages  *tview.Pages
	config *netplan.NetworkConfig
}

// NewApp creates a new application
func NewApp() *App {
	app := tview.NewApplication()
	
	// Set better color scheme
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorBlue
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorGreen
	tview.Styles.BorderColor = tcell.ColorGreen
	tview.Styles.TitleColor = tcell.ColorYellow
	tview.Styles.GraphicsColor = tcell.ColorGreen
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorYellow
	tview.Styles.TertiaryTextColor = tcell.ColorGreen
	tview.Styles.InverseTextColor = tcell.ColorBlack
	tview.Styles.ContrastSecondaryTextColor = tcell.ColorWhite
	
	return &App{
		app:   app,
		pages: tview.NewPages(),
	}
}

// Run starts the application
func (a *App) Run() error {
	// Load config
	config, err := netplan.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load netplan config: %v", err)
	}
	a.config = config

	// Show main menu
	a.showMainMenu()

	// Set root and run
	a.app.SetRoot(a.pages, true)
	return a.app.Run()
}

// showMainMenu displays the main menu
func (a *App) showMainMenu() {
	menu := tview.NewList()
	menu.SetBorder(true)
	menu.SetTitle(" Netplan TUI - Main Menu ")
	menu.SetTitleAlign(tview.AlignCenter)
	
	// Set list colors for better visibility
	menu.SetMainTextColor(tcell.ColorWhite)
	menu.SetSecondaryTextColor(tcell.ColorYellow)
	menu.SetSelectedTextColor(tcell.ColorBlack)
	menu.SetSelectedBackgroundColor(tcell.ColorGreen)
	menu.SetShortcutColor(tcell.ColorDarkCyan)

	menu.AddItem("Edit Network Interfaces", "Configure network adapters", '1', func() {
		a.showInterfaceList()
	})
	menu.AddItem("Apply Configuration", "Apply netplan changes", '2', func() {
		a.applyConfig()
	})
	menu.AddItem("Quit", "Exit the program", 'q', func() {
		a.app.Stop()
	})

	menu.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Items are already handled by their callback
	})

	// Add footer with help text
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(menu, 0, 1, true).
		AddItem(a.createFooter("↑↓: Navigate | Enter: Select | q: Quit"), 1, 0, false)

	a.pages.AddPage("main", flex, true, true)
}

// showInterfaceList shows the list of network interfaces
func (a *App) showInterfaceList() {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Network Interfaces ")
	list.SetTitleAlign(tview.AlignCenter)
	
	// Set list colors for better visibility
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSecondaryTextColor(tcell.ColorYellow)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorGreen)
	list.SetShortcutColor(tcell.ColorDarkCyan)

	// Get available interfaces
	interfaces, err := netplan.GetInterfaces()
	if err != nil {
		a.showError(fmt.Sprintf("Failed to get interfaces: %v", err))
		return
	}

	// Add interfaces to list
	for _, iface := range interfaces {
		config := a.config.GetInterfaceConfig(iface)
		secondary := config.FormatConfig()
		
		// Capture variable in closure
		ifaceName := iface
		list.AddItem(ifaceName, secondary, 0, func() {
			a.showInterfaceEdit(ifaceName)
		})
	}

	list.AddItem("Back", "Return to main menu", 'b', func() {
		a.pages.SwitchToPage("main")
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(a.createFooter("↑↓: Navigate | Enter: Edit | b: Back"), 1, 0, false)

	a.pages.AddPage("interfaces", flex, true, true)
}

// showInterfaceEdit shows the interface edit form
func (a *App) showInterfaceEdit(iface string) {
	config := a.config.GetInterfaceConfig(iface)
	
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(fmt.Sprintf(" Configure Interface: %s ", iface))
	form.SetTitleAlign(tview.AlignCenter)

	// Configuration method
	configMethod := "dhcp"
	if !config.DHCP4 && len(config.Addresses) > 0 {
		configMethod = "static"
	}

	// Form fields
	ipAddress := ""
	gateway := ""
	dns := ""
	
	if len(config.Addresses) > 0 {
		ipAddress = config.Addresses[0]
	}
	if config.Gateway4 != "" {
		gateway = config.Gateway4
	}
	if config.Nameservers != nil && len(config.Nameservers.Addresses) > 0 {
		dns = config.Nameservers.Addresses[0]
	}

	// Create input fields first
	ipField := tview.NewInputField().
		SetLabel("IP Address/CIDR").
		SetText(ipAddress).
		SetFieldWidth(30)
	if configMethod == "dhcp" {
		ipField.SetDisabled(true)
	}
	
	gwField := tview.NewInputField().
		SetLabel("Gateway").
		SetText(gateway).
		SetFieldWidth(30)
	if configMethod == "dhcp" {
		gwField.SetDisabled(true)
	}
	
	dnsField := tview.NewInputField().
		SetLabel("DNS Server").
		SetText(dns).
		SetFieldWidth(30)
	if configMethod == "dhcp" {
		dnsField.SetDisabled(true)
	}

	// Create a toggle field for DHCP/Static selection (better than dropdown)
	configField := tview.NewInputField().
		SetLabel("Configuration").
		SetText(map[string]string{"dhcp": "DHCP", "static": "Static"}[configMethod]).
		SetFieldWidth(20).
		SetFieldBackgroundColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorYellow)
	
	// Make it read-only but toggleable
	configField.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		return false // Don't accept any input
	})
	
	// Toggle on Space or Enter
	configField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter || event.Rune() == ' ' {
			// Toggle between DHCP and Static
			if configMethod == "dhcp" {
				configMethod = "static"
				configField.SetText("Static")
				// Enable static fields
				ipField.SetDisabled(false)
				gwField.SetDisabled(false)
				dnsField.SetDisabled(false)
			} else {
				configMethod = "dhcp"
				configField.SetText("DHCP")
				// Disable static fields
				ipField.SetDisabled(true)
				gwField.SetDisabled(true)
				dnsField.SetDisabled(true)
			}
			return nil
		}
		return event
	})
	
	form.AddFormItem(configField)
	
	// Set field colors
	ipField.SetFieldBackgroundColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorYellow)
	gwField.SetFieldBackgroundColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorYellow)
	dnsField.SetFieldBackgroundColor(tcell.ColorBlue).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorYellow)
	
	form.AddFormItem(ipField)
	form.AddFormItem(gwField)
	form.AddFormItem(dnsField)

	// Set form button style
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetButtonBackgroundColor(tcell.ColorBlue)
	form.SetButtonTextColor(tcell.ColorWhite)
	form.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(tcell.ColorGreen).
		Foreground(tcell.ColorBlack))
	
	form.AddButton("Save", func() {
		newConfig := netplan.EthernetConfig{}
		
		if configMethod == "dhcp" {
			newConfig.DHCP4 = true
		} else {
			newConfig.DHCP4 = false
			
			ip := ipField.GetText()
			if ip != "" {
				newConfig.Addresses = []string{ip}
			}
			
			gw := gwField.GetText()
			if gw != "" {
				newConfig.Gateway4 = gw
			}
			
			dnsText := dnsField.GetText()
			if dnsText != "" {
				newConfig.Nameservers = &netplan.DNS{
					Addresses: []string{dnsText},
				}
			}
		}
		
		a.config.SetInterfaceConfig(iface, newConfig)
		
		if err := netplan.SaveConfig(a.config); err != nil {
			a.showError(fmt.Sprintf("Failed to save config: %v", err))
			return
		}
		
		a.showInfo("Configuration saved! Use 'Apply Configuration' to activate.")
		a.showInterfaceList()
	})

	form.AddButton("Cancel", func() {
		a.showInterfaceList()
	})

	// Add Esc handler at Flex level to avoid interfering with form navigation
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).
		AddItem(a.createFooter("Tab: Navigate | Space/Enter: Toggle | Esc: Back"), 1, 0, false)
	
	// Set input capture on flex instead of form for better compatibility
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.showInterfaceList()
			return nil
		}
		return event
	})

	a.pages.AddPage("edit", flex, true, true)
}

// applyConfig applies the netplan configuration
func (a *App) applyConfig() {
	modal := tview.NewModal().
		SetText("Apply netplan configuration?\nThis will activate the network changes.").
		AddButtons([]string{"Apply", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Apply" {
				if err := netplan.ApplyConfig(); err != nil {
					a.showError(fmt.Sprintf("Failed to apply configuration: %v", err))
				} else {
					a.showInfo("Configuration applied successfully!\nNetwork changes are now active.")
				}
			} else {
				a.pages.SwitchToPage("main")
			}
		})

	a.pages.AddPage("apply", modal, true, true)
}

// showError shows an error dialog
func (a *App) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.SwitchToPage("main")
		})

	modal.SetBackgroundColor(tcell.ColorRed)
	modal.SetTextColor(tcell.ColorWhite)
	
	a.pages.AddPage("error", modal, true, true)
}

// showInfo shows an info dialog
func (a *App) showInfo(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("info")
		})

	a.pages.AddPage("info", modal, true, true)
}

// createFooter creates a footer with help text
func (a *App) createFooter(text string) *tview.TextView {
	footer := tview.NewTextView().
		SetText(text).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	footer.SetBackgroundColor(tcell.ColorDarkBlue)
	footer.SetTextColor(tcell.ColorWhite)
	return footer
}

