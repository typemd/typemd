package tui

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
	"github.com/typemd/typemd/tui/widget"
)

// titlePanelHeight is the total height of the title panel (1 content line + 2 border lines).
const titlePanelHeight = 3

// newBodyTextarea creates a configured textarea for body editing.
func newBodyTextarea() textarea.Model {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.Prompt = " " // single-space indent matching view mode; must be set before SetWidth
	// Remove textarea's own border — the panel border is provided by lipgloss
	noBase := lipgloss.NewStyle()
	s := ta.Styles()
	s.Focused.Base = noBase
	s.Blurred.Base = noBase
	ta.SetStyles(s)
	return ta
}

// resizeBodyTextarea updates the body textarea dimensions to match the current layout.
func (m *model) resizeBodyTextarea() {
	h := m.height - 3 // help bar + borders
	if m.selected != nil {
		h -= titlePanelHeight // title panel takes vertical space
	}
	if h < 0 {
		h = 0
	}
	m.bodyTextarea.SetWidth(m.bodyWidth())
	m.bodyTextarea.SetHeight(h)
}

// adjustScroll updates scrollOffset so cursor is always visible.
func (m *model) adjustScroll() {
	contentH := m.height - 3
	m.scrollOffset = widget.AdjustScroll(m.cursor, m.scrollOffset, contentH)
}

// resizePanel adjusts the focused panel width by delta characters.
func (m *model) resizePanel(delta int) {
	switch m.focus {
	case focusLeft:
		m.leftW += delta
		if m.leftW < 20 {
			m.leftW = 20
		}
		if m.leftW > 50 {
			m.leftW = 50
		}
	case focusBody:
		// Body has no dedicated width field; grow body = shrink props
		if m.propsVisible {
			m.propsWidth -= delta
			if m.propsWidth < 20 {
				m.propsWidth = 20
			}
			if m.propsWidth > 40 {
				m.propsWidth = 40
			}
		} else {
			// Props hidden; grow body = shrink left
			m.leftW -= delta
			if m.leftW < 20 {
				m.leftW = 20
			}
			if m.leftW > 50 {
				m.leftW = 50
			}
		}
	case focusProps:
		m.propsWidth += delta
		if m.propsWidth < 20 {
			m.propsWidth = 20
		}
		if m.propsWidth > 40 {
			m.propsWidth = 40
		}
	}
	// Recalculate dependent widths
	m.bodyViewport.SetWidth(m.bodyWidth())
	m.propsViewport.SetWidth(m.propsWidth)
	m.bodyTextarea.SetWidth(m.bodyWidth())
	m.updateDetail()
}

// defaultLeftWidth calculates the default left panel width.
func (m model) defaultLeftWidth() int {
	w := m.width * 2 / 5
	if w < 20 {
		w = 20
	}
	if w > 50 {
		w = 50
	}
	return w
}

// leftWidth returns the current width for the left panel.
func (m model) leftWidth() int {
	if m.leftW > 0 {
		return m.leftW
	}
	return m.defaultLeftWidth()
}

// defaultPropsWidth calculates the default properties panel width.
func (m model) defaultPropsWidth() int {
	remaining := m.width - m.leftWidth() - 6 // 6 = borders for 3 panels
	w := remaining * 3 / 10                   // 30% of remaining
	if w < 20 {
		w = 20
	}
	if w > 40 {
		w = 40
	}
	return w
}

// hasTitlePanel returns true when the right side should show a title panel.
func (m model) hasTitlePanel() bool {
	return m.selected != nil ||
		(m.rightPanel == panelTypeEditor && m.typeEditor != nil) ||
		(m.rightPanel == panelTemplate && m.tmplEditor != nil) ||
		m.create != nil || m.createType != nil
}

// bodyWidth calculates the body panel width from remaining space.
func (m model) bodyWidth() int {
	borders := 4 // left panel border (2) + body panel border (2)
	if m.propsVisible {
		borders += 2 // props panel border
	}
	w := m.width - m.leftWidth() - borders
	if m.propsVisible {
		w -= m.propsWidth
	}
	if w < 10 {
		w = 10
	}
	return w
}

// shouldAutoHideProps returns true if terminal is too narrow for three panels.
func (m model) shouldAutoHideProps() bool {
	minTotal := 20 + 10 + 20 + 6 // minLeft + minBody + minProps + borders
	return m.width < minTotal
}
