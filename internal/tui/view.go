package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the full-screen layout.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}
	innerW, innerH := m.innerSize()
	_, listH := m.listArea()

	header := m.renderHeader(innerW)
	footer := m.renderFooter(innerW)

	var middle string
	if m.mode == modeAdd || m.mode == modeEdit {
		list := m.renderList(innerW, listH)
		form := m.renderInputForm(innerW)
		middle = lipgloss.JoinVertical(lipgloss.Left, list, "", form)
	} else {
		middle = m.renderList(innerW, listH)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		middle,
		footer,
	)
	sized := lipgloss.NewStyle().
		Width(innerW).
		Height(innerH).
		Render(content)
	return theme.App.Render(sized)
}

func (m Model) renderHeader(width int) string {
	title := theme.Title.Render("✓ TUI Todo")
	done, total := m.svc.Counts()
	stats := theme.Stats.Render(fmt.Sprintf("%d/%d done", done, total))

	tw := lipgloss.Width(title)
	sw := lipgloss.Width(stats)
	gap := width - tw - sw
	if gap < 1 {
		gap = 1
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, title, strings.Repeat(" ", gap), stats)
	sep := theme.Separator.Render(strings.Repeat("─", width))
	return lipgloss.JoinVertical(lipgloss.Left, row, sep)
}

func (m Model) renderList(width, listH int) string {
	items := m.svc.Items()
	if len(items) == 0 {
		empty := theme.Empty.Width(width).Render("No tasks yet. Press 'a' to add your first one.")
		return lipgloss.Place(width, listH, lipgloss.Center, lipgloss.Center, empty)
	}
	start := m.offset
	end := start + listH
	if end > len(items) {
		end = len(items)
	}
	var rows []string
	for i := start; i < end; i++ {
		rows = append(rows, m.renderRow(i, width))
	}
	for len(rows) < listH {
		rows = append(rows, strings.Repeat(" ", width))
	}
	body := strings.Join(rows, "\n")
	if len(items) > listH {
		hint := theme.ScrollHint.Render(fmt.Sprintf(" %d-%d/%d ", start+1, end, len(items)))
		hw := lipgloss.Width(hint)
		if hw < width {
			lines := strings.Split(body, "\n")
			if n := len(lines); n > 0 {
				last := lines[n-1]
				lw := lipgloss.Width(last)
				if lw+hw <= width {
					lines[n-1] = last + strings.Repeat(" ", width-lw-hw) + hint
					body = strings.Join(lines, "\n")
				}
			}
		}
	}
	return body
}

func (m Model) renderRow(i, width int) string {
	t := m.svc.Items()[i]
	isSelected := i == m.cursor && m.mode == modeList

	var cursor string
	if isSelected {
		cursor = theme.SelectedAccent.Render("▸ ")
	} else {
		cursor = "  "
	}
	var box string
	if t.Done {
		box = "[x] "
	} else {
		box = "[ ] "
	}
	prefix := cursor + box
	pw := lipgloss.Width(prefix)
	tmax := width - pw
	if tmax < 1 {
		tmax = 1
	}
	title := truncate(t.Title, tmax)

	var titleStyled string
	switch {
	case t.Done:
		titleStyled = theme.Done.Render(title)
	case isSelected:
		titleStyled = theme.SelectedAccent.Render(title)
	default:
		titleStyled = theme.Pending.Render(title)
	}
	line := prefix + titleStyled
	lw := lipgloss.Width(line)
	if lw < width {
		line += strings.Repeat(" ", width-lw)
	}
	if isSelected {
		line = theme.SelectedRow.Width(width).Render(prefix + title)
	}
	return line
}

func (m Model) renderInputForm(width int) string {
	if m.mode != modeAdd && m.mode != modeEdit {
		return ""
	}
	label := "New task:"
	if m.mode == modeEdit {
		label = "Edit task:"
	}
	help := theme.Help.Render("enter: save • esc: cancel")
	block := lipgloss.JoinVertical(
		lipgloss.Left,
		theme.InputLabel.Render(label),
		m.input.View(),
		help,
	)
	return lipgloss.NewStyle().Width(width).Render(block)
}

func (m Model) renderFooter(width int) string {
	help := "a: add  •  e: edit  •  space/enter: toggle  •  d: delete  •  c: clear done  •  j/k ↑/↓: navigate  •  q: quit"
	if lipgloss.Width(help) > width {
		help = truncate(help, width)
	}
	helpLine := theme.Help.Render(help)
	if m.status == "" {
		return helpLine
	}
	return lipgloss.JoinVertical(lipgloss.Left, theme.Status.Render(m.status), helpLine)
}
