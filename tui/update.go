package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// updateHelp handles key events when the help overlay is shown.
func updateHelp(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "esc", "?", "h":
		m.showHelp = false
	}
	return m, nil
}

// updateConflict handles key events during save conflict resolution.
func updateConflict(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.forceSave()
	case "n":
		m.reloadFromDisk()
	case "esc":
		m.saveConflict = false
		m.saveErr = ""
	}
	return m, nil
}

// updateEdit handles key events in edit mode.
func updateEdit(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	if msg.String() == "esc" {
		if m.focus == focusBody && m.selected != nil {
			newBody := m.bodyTextarea.Value()
			if newBody != m.bodyEditStart {
				m.selected.Body = newBody
				m.dirty = true
				m.updateDetail()
			}
			m.bodyTextarea.Blur()
		}
		m.editMode = false
		if m.dirty {
			m.saveObject()
		}
		return m, nil
	}
	if m.focus == focusBody {
		var cmd tea.Cmd
		m.bodyTextarea, cmd = m.bodyTextarea.Update(msg)
		return m, cmd
	}
	return m, nil
}

// updateNormal handles key events in normal (non-modal) mode.
func updateNormal(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		if m.vault != nil {
			saveSessionState(m.vault.Root, m.captureState())
		}
		return m, tea.Quit

	case "/":
		m.searchMode = true
		m.searchInput.Focus()
		return m, textinput.Blink

	case "e":
		if m.readOnly {
			return m, nil
		}
		if m.focus == focusBody && m.selected != nil {
			return m, m.enterBodyEditMode()
		}
		if m.focus == focusProps {
			m.editMode = true
		}
		return m, nil

	case "tab":
		switch m.focus {
		case focusLeft:
			if m.rightPanel == panelTypeEditor {
				m.focus = focusBody // focusBody doubles as "right panel focus" for type editor
			} else {
				m.focus = focusBody
			}
		case focusBody:
			if m.rightPanel == panelTypeEditor {
				m.focus = focusLeft
			} else if m.propsVisible {
				m.focus = focusProps
			} else {
				m.focus = focusLeft
			}
		case focusProps:
			m.focus = focusLeft
		}
		return m, nil

	case "w":
		m.softWrap = !m.softWrap
		m.updateDetail()
		return m, nil

	case "n":
		if cmd, ok := m.tryStartCreate(createModeSingle); ok {
			return m, cmd
		}
		return m, nil

	case "N":
		if cmd, ok := m.tryStartCreate(createModeBatch); ok {
			return m, cmd
		}
		return m, nil

	case "esc":
		// Clear search results and return to normal list
		if m.searchResults != nil {
			m.searchResults = nil
			m.cursor = 0
			m.selectCurrentRow()
			return m, nil
		}

	case "up", "k":
		if m.focus == focusLeft {
			rows := m.currentRows()
			m.cursor = clampCursor(m.cursor-1, len(rows))
			m.adjustScroll()
			m.selectCurrentRow()
		} else if m.focus == focusBody {
			m.bodyViewport.ScrollUp(1)
		} else if m.focus == focusProps {
			m.propsViewport.ScrollUp(1)
		}
		return m, nil

	case "down", "j":
		if m.focus == focusLeft {
			rows := m.currentRows()
			m.cursor = clampCursor(m.cursor+1, len(rows))
			m.adjustScroll()
			m.selectCurrentRow()
		} else if m.focus == focusBody {
			m.bodyViewport.ScrollDown(1)
		} else if m.focus == focusProps {
			m.propsViewport.ScrollDown(1)
		}
		return m, nil

	case "]":
		m.resizePanel(+2)
		return m, nil

	case "[":
		m.resizePanel(-2)
		return m, nil

	case "p":
		m.propsVisible = !m.propsVisible
		if !m.propsVisible && m.focus == focusProps {
			m.focus = focusBody
		}
		// Recalculate widths for both panels
		contentHeight := m.height - 3
		if contentHeight < 0 {
			contentHeight = 0
		}
		if m.selected != nil {
			contentHeight -= titlePanelHeight
			if contentHeight < 0 {
				contentHeight = 0
			}
		}
		m.bodyViewport.SetWidth(m.bodyWidth())
		m.propsViewport.SetWidth(m.propsWidth)
		m.propsViewport.SetHeight(contentHeight)
		m.updateDetail()
		return m, nil

	case "?", "h":
		m.showHelp = true
		return m, nil

	case "enter":
		if m.focus == focusLeft {
			rows := m.currentRows()
			if m.cursor >= 0 && m.cursor < len(rows) {
				row := rows[m.cursor]
				switch row.Kind {
				case rowHeader:
					// Enter on header: focus type editor (already opened by cursor movement)
					if m.rightPanel == panelTypeEditor && m.typeEditor != nil {
						m.focus = focusBody
					}
				case rowObject:
					m.selectCurrentRow()
				case rowNewType:
					return m, m.startCreateType()
				}
			}
		}
		return m, nil

	case " ", "space":
		if m.focus == focusLeft {
			rows := m.currentRows()
			if m.cursor >= 0 && m.cursor < len(rows) {
				row := rows[m.cursor]
				switch row.Kind {
				case rowHeader:
					// Space on header: toggle expand/collapse
					m.groups[row.GroupIndex].Expanded = !m.groups[row.GroupIndex].Expanded
					newRows := m.currentRows()
					m.cursor = clampCursor(m.cursor, len(newRows))
					m.adjustScroll()
				case rowObject:
					m.selectCurrentRow()
				case rowNewType:
					return m, m.startCreateType()
				}
			}
		}
		return m, nil
	}
	return m, nil
}
