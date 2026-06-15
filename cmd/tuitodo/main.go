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

// runProgram is the program runner, extracted for testability.
var runProgram = func(m tui.Model) error {
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Startup error:", err)
		os.Exit(1)
	}
}

// run is the program logic, extracted for testability.
func run() error {
	dataPath := config.DefaultDataPath()
	repo := jsonfile.New(dataPath)
	svc := app.NewService(repo)
	m := tui.New(svc)
	return runProgram(m)
}
