package main

import (
	"fmt"
	"os"

	"nptui/ui"
)

func main() {
	// Check if running as root
	if os.Geteuid() != 0 {
		fmt.Println("This program must be run as root to modify network configuration.")
		fmt.Println("Please run with: sudo nptui")
		os.Exit(1)
	}

	app := ui.NewApp()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

