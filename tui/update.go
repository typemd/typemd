package tui

import (
	"github.com/typemd/typemd/tui/widget"
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
	// Property editor handles its own editing when props focused and propEdit is active
	if m.focus == focusProps && m.propEdit != nil && m.propEdit.isEditing() {
		newM, cmd, consumed := updatePropEditor(m, msg)
		if consumed {
			return newM, cmd
		}
	}

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
	// Property editor in editing state intercepts all keys
	if m.focus == focusProps && m.propEdit != nil && m.propEdit.isEditing() {
		newM, cmd, consumed := updatePropEditor(m, msg)
		if consumed {
			return newM, cmd
		}
	}

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

	case "g":
		return handleAIGenerate(m)

	case "ctrl+e":
		return handleSchemaExplore(m)

	case "ctrl+s":
		if m.vault != nil {
			layout := ""
			if cfg := m.vault.Config(); cfg != nil {
				layout = cfg.TUI.StatsTypeLayout
			}
			sm := newStatsMode(m.vault, layout)
			sm.SetSize(m.width-2, m.height-3)
			m.statsMode = sm
			m.rightPanel = panelStats
			m.focus = focusBody
		}
		return m, nil

	case ",":
		if m.vault != nil {
			ce := newConfigEditor(m.vault)
			ce.SetSize(m.width-2, m.height-3)
			m.configEditor = ce
			m.rightPanel = panelConfig
			m.focus = focusBody
		}
		return m, nil

	case "e":
		if m.readOnly {
			return m, nil
		}
		if m.focus == focusBody && m.selected != nil {
			if m.selected.IsLocked() {
				return m, m.toast.Show(widget.ToastWarning, []widget.ToastItem{{Message: "Object is locked. Unlock to edit."}})
			}
			return m, m.enterBodyEditMode()
		}
		return m, nil

	case "r":
		if m.readOnly {
			return m, nil
		}
		if m.selected != nil && m.rightPanel == panelObject {
			if m.selected.IsLocked() {
				return m, m.toast.Show(widget.ToastWarning, []widget.ToastItem{{Message: "Object is locked. Unlock to edit."}})
			}
			return m, m.startRename()
		}
		return m, nil

	case "tab":
		if m.focusMode {
			return m, nil
		}
		prevFocus := m.focus
		switch m.focus {
		case focusLeft:
			if m.rightPanel == panelTypeEditor || m.rightPanel == panelTemplate {
				m.focus = focusBody
			} else {
				m.focus = focusBody
			}
		case focusBody:
			if m.rightPanel == panelTypeEditor {
				m.focus = focusLeft
			} else if m.rightPanel == panelTemplate {
				m.focus = focusBody // template editor handles tab internally for body/props
			} else if m.propsVisible {
				// Prevent entering property editor for locked objects
				if m.selected != nil && m.selected.IsLocked() {
					toastCmd := m.toast.Show(widget.ToastWarning, []widget.ToastItem{{Message: "Object is locked. Unlock to edit."}})
					return m, toastCmd
				}
				m.focus = focusProps
			} else {
				m.focus = focusLeft
			}
		case focusProps:
			m.focus = focusLeft
		}
		// Refresh props content when entering/leaving focusProps to show/hide cursor
		if m.focus == focusProps || prevFocus == focusProps {
			m.updatePropsContent()
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

	case "L":
		if m.readOnly {
			return m, nil
		}
		if m.selected != nil && m.vault != nil && m.rightPanel == panelObject {
			locked := m.selected.IsLocked()
			if err := m.vault.SetLocked(m.selected.ID, !locked); err != nil {
				return m, m.toast.Show(widget.ToastError, []widget.ToastItem{{Message: err.Error()}})
			}
			obj, err := m.vault.GetObject(m.selected.ID)
			if err == nil {
				m.selected = obj
			}
			m.updateDetail()
			msg := "Locked"
			if locked {
				msg = "Unlocked"
			}
			return m, m.toast.Show(widget.ToastInfo, []widget.ToastItem{{Message: msg}})
		}
		return m, nil

	case "v":
		// Enter view mode for the current type
		if m.focus == focusLeft && m.vault != nil {
			typeName := m.currentTypeName()
			if typeName != "" {
				views, _ := m.vault.ListViews(typeName)
				if len(views) <= 1 {
					// No saved views or only default — enter directly
					return m, func() tea.Msg {
						return openViewMsg{TypeName: typeName, ViewName: "default"}
					}
				}
				// Multiple views — show picker
				var names []string
				for _, v := range views {
					names = append(names, v.Name)
				}
				m.viewPicker = newViewPicker(typeName, names)
				return m, m.viewPicker.Init()
			}
		}
		return m, nil

	case "esc":
		// Property editor: return to sidebar
		if m.focus == focusProps && m.propEdit != nil {
			m.focus = focusLeft
			m.updateDetail()
			return m, nil
		}
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
		} else if m.focus == focusProps && m.propEdit != nil {
			newM, cmd, consumed := updatePropNavigate(m, msg)
			if consumed {
				return newM, cmd
			}
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
		} else if m.focus == focusProps && m.propEdit != nil {
			newM, cmd, consumed := updatePropNavigate(m, msg)
			if consumed {
				return newM, cmd
			}
		} else if m.focus == focusProps {
			m.propsViewport.ScrollDown(1)
		}
		return m, nil

	case "=":
		m.resizePanel(+2)
		return m, nil

	case "-":
		m.resizePanel(-2)
		return m, nil

	case ".":
		m.focusMode = !m.focusMode
		if m.focusMode {
			m.focus = focusBody
		}
		m.recalcLayout()
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
		if m.focus == focusProps && m.propEdit != nil {
			newM, cmd, consumed := updatePropNavigate(m, msg)
			if consumed {
				return newM, cmd
			}
		}
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
		if m.focus == focusProps && m.propEdit != nil {
			newM, cmd, consumed := updatePropNavigate(m, msg)
			if consumed {
				return newM, cmd
			}
		}
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
