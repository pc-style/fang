package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pcstyle/kawaii-shell/internal/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("🌸 Kawaii Shell v0.1.0 - Making terminals adorable! ✨")
		return
	}

	// Create the main Bubble Tea application
	app := ui.NewApp()

	// Initialize Bubble Tea program
	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start the program
	if _, err := p.Run(); err != nil {
		log.Fatalf("🥺 Oops! Something went wrong: %v", err)
	}
}
