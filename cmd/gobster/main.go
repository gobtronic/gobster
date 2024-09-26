package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/gobtronic/gobster/cmd/gobster/load"
)

func main() {
	log.SetReportTimestamp(false)

	p := tea.NewProgram(load.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Error("An error occured while initializing gobster", "err", err)
		os.Exit(1)
	}
}
