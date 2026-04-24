// Package tui is the Bubble Tea user interface. It depends on [app.Service], not on storage.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/tuitodo/internal/app"
)

type mode int

const (
	modeList mode = iota
	modeAdd
	modeEdit
)

// Model is the Bubble Tea program state: UI chrome + a reference to the task [app.Service].
type Model struct {
	svc    *app.Service
	cursor int
	offset int
	mode   mode
	input  textinput.Model
	width  int
	height int
	status string
}

// New builds a TUI [Model] for the given service. Initial load errors become the first status line.
func New(svc *app.Service) Model {
	m := Model{svc: svc}
	if err := svc.InitError(); err != nil {
		m.status = "Load error: " + err.Error()
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

// NewTextInput is exported for tests; normal path uses it inside Model when entering add/edit.
func newTextInput(placeholder, value string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(value)
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = width
	return ti
}

func (m *Model) setSaveErr(err error) {
	if err != nil {
		m.status = "Save error: " + err.Error()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.mode == modeAdd || m.mode == modeEdit {
			innerW, _ := m.innerSize()
			m.input.Width = max(10, innerW-4)
		}
		m.clampOffset()
		return m, nil
	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			return m.updateList(msg)
		case modeAdd, modeEdit:
			return m.updateInput(msg)
		}
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	_, listH := m.listArea()
	items := m.svc.Items()
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(items)-1 {
			m.cursor++
		}
	case "pgup":
		m.cursor -= listH
		if m.cursor < 0 {
			m.cursor = 0
		}
	case "pgdown":
		m.cursor += listH
		if m.cursor > len(items)-1 {
			m.cursor = len(items) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
	case "g", "home":
		m.cursor = 0
	case "G", "end":
		if len(items) > 0 {
			m.cursor = len(items) - 1
		}
	case " ", "enter", "x":
		if len(items) > 0 {
			done, err := m.svc.Toggle(m.cursor)
			m.setSaveErr(err)
			if err == nil {
				if done {
					m.status = "Task completed"
				} else {
					m.status = "Task marked as pending"
				}
			}
		}
	case "a", "n", "+":
		m.mode = modeAdd
		innerW, _ := m.innerSize()
		m.input = newTextInput("New task...", "", max(10, innerW-4))
		m.status = ""
		return m, textinput.Blink
	case "e":
		if len(items) > 0 {
			m.mode = modeEdit
			innerW, _ := m.innerSize()
			m.input = newTextInput("Edit task...", items[m.cursor].Title, max(10, innerW-4))
			m.status = ""
			return m, textinput.Blink
		}
	case "d", "delete", "backspace":
		if len(items) > 0 {
			title, err := m.svc.Remove(m.cursor)
			m.setSaveErr(err)
			if err == nil {
				if m.cursor >= m.svc.Len() && m.cursor > 0 {
					m.cursor--
				}
				m.status = fmt.Sprintf("Deleted: %s", truncate(title, 30))
			}
		}
	case "c":
		cleared, err := m.svc.ClearDone()
		m.setSaveErr(err)
		if err == nil {
			if m.cursor >= m.svc.Len() {
				m.cursor = max(0, m.svc.Len()-1)
			}
			m.status = fmt.Sprintf("Cleared %d completed", cleared)
		}
	}
	m.clampOffset()
	return m, nil
}

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = modeList
		m.input.Blur()
		m.clampOffset()
		return m, nil
	case tea.KeyEnter:
		title := strings.TrimSpace(m.input.Value())
		if title == "" {
			m.mode = modeList
			m.clampOffset()
			return m, nil
		}
		switch m.mode {
		case modeAdd:
			idx, err := m.svc.Add(title)
			m.setSaveErr(err)
			if err == nil {
				m.cursor = idx
				m.status = "Task added"
			}
		case modeEdit:
			err := m.svc.SetTitle(m.cursor, title)
			m.setSaveErr(err)
			if err == nil {
				m.status = "Task updated"
			}
		}
		m.mode = modeList
		m.input.Blur()
		m.clampOffset()
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
