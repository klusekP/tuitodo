package tui

// outer frame: border (1) + padding (2) on each side.
const (
	frameWidth  = 6
	frameHeight = 4

	minInnerWidth  = 20
	minInnerHeight = 6
)

// innerSize returns the content size inside the border.
func (m Model) innerSize() (w, h int) {
	w = m.width - frameWidth
	h = m.height - frameHeight
	if w < minInnerWidth {
		w = minInnerWidth
	}
	if h < minInnerHeight {
		h = minInnerHeight
	}
	return w, h
}

// listArea returns (width, height) for the scrollable list in rows.
func (m Model) listArea() (int, int) {
	innerW, innerH := m.innerSize()
	const headerH = 3
	footerH := 2
	if m.status != "" {
		footerH++
	}
	inputH := 0
	if m.mode == modeAdd || m.mode == modeEdit {
		inputH = 4
	}
	listH := innerH - headerH - footerH - inputH
	if listH < 1 {
		listH = 1
	}
	return innerW, listH
}

// clampOffset keeps the cursor visible within the window.
func (m *Model) clampOffset() {
	_, listH := m.listArea()
	n := m.svc.Len()
	if n == 0 {
		m.offset = 0
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+listH {
		m.offset = m.cursor - listH + 1
	}
	maxOff := n - listH
	if maxOff < 0 {
		maxOff = 0
	}
	if m.offset > maxOff {
		m.offset = maxOff
	}
	if m.offset < 0 {
		m.offset = 0
	}
}
