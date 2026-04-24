package tui

import "github.com/charmbracelet/lipgloss"

// Theme groups lipgloss styles (single place to tweak presentation).
var theme = struct {
	Title, Stats                     lipgloss.Style
	SelectedRow, SelectedAccent, Done lipgloss.Style
	Pending, Help, Status, Empty     lipgloss.Style
	InputLabel, ScrollHint, Separator lipgloss.Style
	App                              lipgloss.Style
}{
	Title: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 2),
	Stats: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A1A1AA")),
	SelectedRow: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#3F3F46")).
		Bold(true),
	SelectedAccent: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EE6FF8")).
		Bold(true),
	Done: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Strikethrough(true),
	Pending: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E4E4E7")),
	Help: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")),
	Status: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Italic(true),
	Empty: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true).
		Align(lipgloss.Center),
	InputLabel: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")),
	ScrollHint: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#71717A")).
		Italic(true),
	Separator: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3F3F46")),
	App: lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")),
}
