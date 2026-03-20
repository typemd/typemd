package tui

import (
	"strings"

	"github.com/typemd/typemd/core"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// renameState holds state for the object rename flow.
type renameState struct {
	nameInput textinput.Model
	typeName  string
	emoji     string
}

// startRename initializes the rename flow for the currently selected object.
func (m *model) startRename() tea.Cmd {
	if m.readOnly || m.selected == nil {
		return nil
	}

	m.rename = &renameState{
		nameInput: initNameInput(m.selected.GetName()),
		typeName:  m.selected.Type,
		emoji:     m.selectedTypeEmoji(),
	}
	return textinput.Blink
}

// renderRenameTitleContent renders the title panel content during rename mode.
func renderRenameTitleContent(rs *renameState) string {
	if rs == nil {
		return ""
	}
	return titlePrefix(rs.emoji, rs.typeName) + " · " + rs.nameInput.View()
}

// updateRename handles key events during rename mode.
func updateRename(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	rs := m.rename
	if rs == nil {
		return m, nil
	}

	switch msg.String() {
	case "enter":
		newName := strings.TrimSpace(rs.nameInput.Value())
		if newName != "" && newName != m.selected.GetName() {
			m.selected.Properties[core.NameProperty] = newName
			m.dirty = true
			m.saveObject()
			m.displayProps, _ = m.vault.BuildDisplayProperties(m.selected)
			m.updateDetail()
			m.rebuildGroups()
			m.moveCursorToObject(m.selected)
		}
		m.rename = nil
		return m, nil

	case "esc":
		m.rename = nil
		return m, nil
	}

	var cmd tea.Cmd
	rs.nameInput, cmd = rs.nameInput.Update(msg)
	return m, cmd
}
