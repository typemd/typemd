package tui

import (
	"fmt"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	tea "charm.land/bubbletea/v2"
)

// updatePropEditor handles key events when the properties panel is focused.
// Returns the updated model and command. The second bool indicates if the event was consumed.
func updatePropEditor(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	if pe == nil {
		return m, nil, false
	}

	switch pe.mode {
	case propModeTextInput:
		return updatePropTextInput(m, msg)
	case propModeSelectPick:
		return updatePropSelectPick(m, msg)
	case propModeMultiPick:
		return updatePropMultiPick(m, msg)
	case propModeRelationPick:
		return updateRelationPick(m, msg)
	case propModeRelationMultiPick:
		return updateRelationMultiPick(m, msg)
	case propModeDateSegment, propModeDateCalendar:
		return updatePropDateEdit(m, msg)
	default:
		return updatePropNavigate(m, msg)
	}
}

// updatePropNavigate handles key events in cursor navigation mode.
func updatePropNavigate(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	switch msg.String() {
	case "up", "k":
		pe.moveUp()
		m.updatePropsContent()
		return m, nil, true

	case "down", "j":
		pe.moveDown()
		m.updatePropsContent()
		return m, nil, true

	case "enter":
		item := pe.currentItem()
		if item == nil || !item.editable {
			return m, nil, true
		}
		if item.resolvedType() == "checkbox" {
			return toggleCheckbox(m)
		}
		cmd := pe.activateEdit(m.vault)
		m.updatePropsContent()
		return m, cmd, true

	case " ", "space":
		item := pe.currentItem()
		if item == nil || !item.editable {
			return m, nil, true
		}
		if item.resolvedType() == "checkbox" {
			return toggleCheckbox(m)
		}
		return m, nil, true

	case "esc":
		m.focus = focusLeft
		m.updatePropsContent()
		return m, nil, true
	}

	return m, nil, false
}

// updatePropTextInput handles key events during textinput editing.
func updatePropTextInput(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	switch msg.String() {
	case "enter":
		item := &pe.items[pe.editIndex]
		propType := item.resolvedType()
		input := pe.textInput.Value()

		var options []core.Option
		if item.schema != nil {
			options = item.schema.Options
		}
		if err := core.ValidatePropertyValue(propType, options, input); err != nil {
			toastCmd := m.toast.Show(widget.ToastError, []widget.ToastItem{{Message: fmt.Sprintf("Validation: %v", err)}})
			m.updatePropsContent()
			return m, toastCmd, true
		}

		applyPropertyValue(&m, item.dp.Key, parseEditedValue(propType, input))
		pe.cancelEdit()
		m.updateDetail()
		return m, nil, true

	case "esc":
		pe.cancelEdit()
		m.updatePropsContent()
		return m, nil, true

	default:
		var cmd tea.Cmd
		pe.textInput, cmd = pe.textInput.Update(msg)
		m.updatePropsContent()
		return m, cmd, true
	}
}

// updatePropSelectPick handles key events in select option picker.
func updatePropSelectPick(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	switch msg.String() {
	case "up", "k":
		if pe.pickerCursor > 0 {
			pe.pickerCursor--
		}
		m.updatePropsContent()
		return m, nil, true

	case "down", "j":
		if pe.pickerCursor < len(pe.pickerOptions)-1 {
			pe.pickerCursor++
		}
		m.updatePropsContent()
		return m, nil, true

	case "enter":
		item := &pe.items[pe.editIndex]
		applyPropertyValue(&m, item.dp.Key, pe.pickerOptions[pe.pickerCursor].Value)
		pe.cancelEdit()
		m.updateDetail()
		return m, nil, true

	case "esc":
		pe.cancelEdit()
		m.updatePropsContent()
		return m, nil, true
	}
	return m, nil, true
}

// updatePropMultiPick handles key events in multi-select option picker.
func updatePropMultiPick(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	switch msg.String() {
	case "up", "k":
		if pe.pickerCursor > 0 {
			pe.pickerCursor--
		}
		m.updatePropsContent()
		return m, nil, true

	case "down", "j":
		if pe.pickerCursor < len(pe.pickerOptions)-1 {
			pe.pickerCursor++
		}
		m.updatePropsContent()
		return m, nil, true

	case " ", "space":
		pe.pickerChecked[pe.pickerCursor] = !pe.pickerChecked[pe.pickerCursor]
		m.updatePropsContent()
		return m, nil, true

	case "enter":
		item := &pe.items[pe.editIndex]
		applyPropertyValue(&m, item.dp.Key, pe.multiPickerResult())
		pe.cancelEdit()
		m.updateDetail()
		return m, nil, true

	case "esc":
		pe.cancelEdit()
		m.updatePropsContent()
		return m, nil, true
	}
	return m, nil, true
}

// updatePropDateEdit handles key events during date picker editing.
func updatePropDateEdit(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	if pe.dateEdit == nil {
		pe.cancelEdit()
		m.updatePropsContent()
		return m, nil, true
	}

	consumed, done, confirmed := pe.dateEdit.Update(msg)
	if !consumed {
		return m, nil, false
	}

	if done {
		if confirmed {
			item := &pe.items[pe.editIndex]
			applyPropertyValue(&m, item.dp.Key, pe.dateEdit.Value())
		}
		pe.cancelEdit()
		m.updateDetail()
		return m, nil, true
	}

	// Sync propEditor mode with dateEdit mode
	if pe.dateEdit.Mode() == dateCalendarMode {
		pe.mode = propModeDateCalendar
	} else {
		pe.mode = propModeDateSegment
	}

	m.updatePropsContent()
	return m, nil, true
}

// toggleCheckbox toggles a checkbox property and saves immediately.
func toggleCheckbox(m model) (model, tea.Cmd, bool) {
	pe := m.propEdit
	item := pe.currentItem()
	if item == nil {
		return m, nil, true
	}

	currentVal := false
	if b, ok := item.dp.Value.(bool); ok {
		currentVal = b
	}
	applyPropertyValue(&m, item.dp.Key, !currentVal)
	m.updateDetail()
	return m, nil, true
}

// applyPropertyValue applies a property value and saves the object.
func applyPropertyValue(m *model, key string, value any) {
	if m.selected == nil {
		return
	}
	if m.selected.Properties == nil {
		m.selected.Properties = make(map[string]any)
	}
	m.selected.Properties[key] = value
	m.dirty = true
	m.saveObject()
	m.displayProps, _ = m.vault.BuildDisplayProperties(m.selected)
	m.refreshPropEditor()
}
