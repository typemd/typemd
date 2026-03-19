package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ── Property Detail Panel ────────────────────────────────────────────────────

// propDetailPanel holds state for the property metadata editing panel.
type propDetailPanel struct {
	cursor     int // 0=emoji (more fields in future)
	emojiInput textinput.Model
	editing    bool // currently in text input mode
}

func newPropDetailPanel(prop *core.Property) *propDetailPanel {
	ei := textinput.New()
	ei.CharLimit = 20
	ei.SetValue(prop.Emoji)
	return &propDetailPanel{
		emojiInput: ei,
	}
}

func (te *typeEditor) openPropDetail() {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	item := items[te.cursor]
	if item < 0 || item == addPropertySentinel {
		return
	}
	te.propDetailIdx = item
	te.propDetail = newPropDetailPanel(&te.schema.Properties[item])
	te.mode = teModeEditProp
}

func (te *typeEditor) updatePropDetail(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	pd := te.propDetail
	if pd == nil {
		te.mode = teModeView
		return te, nil
	}

	if pd.editing {
		switch msg.String() {
		case "enter":
			pd.editing = false
			pd.emojiInput.Blur()
			// Apply value
			te.schema.Properties[te.propDetailIdx].Emoji = pd.emojiInput.Value()
			te.save()
			return te, nil
		case "esc":
			pd.editing = false
			pd.emojiInput.Blur()
			// Revert
			pd.emojiInput.SetValue(te.schema.Properties[te.propDetailIdx].Emoji)
			return te, nil
		}
		var cmd tea.Cmd
		pd.emojiInput, cmd = pd.emojiInput.Update(msg)
		return te, cmd
	}

	switch msg.String() {
	case "esc":
		te.propDetail = nil
		te.mode = teModeView
	case "enter", "e":
		pd.editing = true
		pd.emojiInput.Focus()
		return te, pd.emojiInput.Focus()
	case "up", "k":
		// future: navigate between fields
	case "down", "j":
		// future: navigate between fields
	}
	return te, nil
}

// Overlay returns a popup string if a modal is active, or empty string if not.
func (te *typeEditor) Overlay(width, height int) string {
	if te.mode != teModeEditProp || te.propDetail == nil {
		return ""
	}
	return te.renderPropPopup(width, height)
}

func (te *typeEditor) renderPropPopup(termW, termH int) string {
	pd := te.propDetail
	p := te.schema.Properties[te.propDetailIdx]

	var b strings.Builder

	if pd.editing {
		b.WriteString(fmt.Sprintf("Emoji: %s", pd.emojiInput.View()))
	} else {
		val := p.Emoji
		if val == "" {
			val = "(none)"
		} else {
			val = padEmoji(val)
		}
		line := fmt.Sprintf("Emoji: %s", val)
		if pd.cursor == 0 {
			line = highlightStyle.Render(line)
		}
		b.WriteString(line)
	}

	b.WriteString("\n")

	// future: description field here

	if pd.editing {
		b.WriteString("\nenter: save  esc: cancel")
	} else {
		b.WriteString("\nenter: edit  esc: back")
	}

	popupW := 36
	if popupW > termW-10 {
		popupW = termW - 10
	}

	titleStyle := lipgloss.NewStyle().Bold(true)
	title := titleStyle.Render(fmt.Sprintf("%s (%s)", p.Name, p.Type))
	fullContent := fmt.Sprintf("%s\n──────────────────\n%s", title, b.String())

	return renderPopup(fullContent, termW, termH, popupW)
}
