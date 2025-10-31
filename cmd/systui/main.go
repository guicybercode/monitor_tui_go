package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guicybercode/systui/internal/api"
	"github.com/guicybercode/systui/internal/tui"
)

func main() {
	var headless = flag.Bool("headless", false, "Run in headless API mode")
	var port = flag.Int("port", 8080, "API server port")
	flag.Parse()

	if *headless {
		server := api.NewServer(*port)
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if _, err := tea.NewProgram(tui.NewApp(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
