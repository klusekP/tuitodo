// Command tuitodo is the TUI entrypoint: wires config, file repo, app service, and Bubble Tea.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/tuitodo/internal/app"
	"github.com/klusekP/tuitodo/internal/config"
	"github.com/klusekP/tuitodo/internal/storage/jsonfile"
	"github.com/klusekP/tuitodo/internal/tui"
)

func main() {
	dataPath := config.DefaultDataPath()
	repo := jsonfile.New(dataPath)
	svc := app.NewService(repo)

	m := tui.New(svc)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Startup error:", err)
		os.Exit(1)
	}
}
