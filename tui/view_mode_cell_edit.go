package tui

import (
	"fmt"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// cellEditMode represents the current editing widget type.
type cellEditMode int

const (
	cellModeTextInput    cellEditMode = iota // textinput for string, number, datetime, url
	cellModeSelectPick                       // single-select option picker
	cellModeMultiPick                        // multi-select option picker
	cellModeDateSegment                      // date picker segmented input mode
	cellModeDateCalendar                     // date picker calendar mode
)

// cellEdit tracks the state of an inline cell edit in the table view.
// A nil *cellEdit means no edit is active.
type cellEdit struct {
	rowIdx    int          // row index in visibleRows
	colIdx    int          // column index (0 = NAME, 1+ = property)
	propName  string       // property name being edited (or "name" for NAME column)
	propType  string       // property type (string, number, select, etc.)
	obj       *core.Object // the object being edited
	mode     cellEditMode

	textInput     textinput.Model
	pickerOptions []core.Option
	pickerCursor  int
	pickerChecked []bool // for multi_select: which options are checked

	datePicker *datePicker
}

type viewCellToastMsg struct {
	Level   widget.ToastLevel
	Message string
}

// viewCellSavedMsg signals that a cell edit was saved successfully.
// The parent uses this to trigger a reload of the view data.
type viewCellSavedMsg struct{}

// isCellReadOnly returns true if the column at colIdx is not editable.
func (vm *viewMode) isCellReadOnly(colIdx int, cols []string) bool {
	if colIdx == 0 {
		// NAME column — editable (it's the "name" system property, which is user-authored)
		return false
	}
	propIdx := colIdx - 1
	if propIdx >= len(cols) {
		return true
	}
	propName := cols[propIdx]

	// Immutable system properties (created_at, updated_at)
	if core.IsImmutableSystemProperty(propName) {
		return true
	}

	// Tags are not inline-editable
	if propName == core.TagsProperty {
		return true
	}

	// Relations are not inline-editable
	if vm.schema != nil {
		if p := vm.schema.FindProperty(propName); p != nil {
			if p.Type == "relation" {
				return true
			}
		}
	}

	return false
}

// activateCellEdit initializes inline editing for the given cell.
// Returns a tea.Cmd if editing was activated (e.g. textinput blink),
// or nil if the cell is read-only or a checkbox was toggled directly.
func (vm *viewMode) activateCellEdit(row viewRow, colIdx int, cols []string) tea.Cmd {
	if row.object == nil {
		return nil
	}
	if vm.isCellReadOnly(colIdx, cols) {
		return nil
	}

	obj := row.object

	// Determine property name and type
	var propName, propType string
	var options []core.Option

	if colIdx == 0 {
		// NAME column
		propName = core.NameProperty
		propType = "string"
	} else {
		propIdx := colIdx - 1
		if propIdx >= len(cols) {
			return nil
		}
		propName = cols[propIdx]

		// Look up property type from schema
		if vm.schema != nil {
			if p := vm.schema.FindProperty(propName); p != nil {
				propType = p.Type
				options = p.Options
			}
		}
		// System property "description" is a string
		if propName == core.DescriptionProperty {
			propType = "string"
		}
		if propType == "" {
			propType = "string"
		}
	}

	// Checkbox: toggle directly without entering edit mode
	if propType == "checkbox" {
		currentVal := false
		if b, ok := obj.Properties[propName].(bool); ok {
			currentVal = b
		}
		if obj.Properties == nil {
			obj.Properties = make(map[string]any)
		}
		obj.Properties[propName] = !currentVal
		return vm.saveCellObject(obj)
	}

	ce := &cellEdit{
		rowIdx:   vm.cursor,
		colIdx:   colIdx,
		propName: propName,
		propType: propType,
		obj:      obj,
	}

	switch propType {
	case "date":
		ce.mode = cellModeDateSegment
		ce.datePicker = newDatePicker(parseDatePickerValue(obj.Properties[propName]))
		vm.cellEdit = ce
		return nil

	case "select":
		if len(options) == 0 {
			return nil
		}
		ce.mode = cellModeSelectPick
		ce.pickerOptions = options
		ce.pickerCursor = 0
		// Position cursor on current value
		currentVal := fmt.Sprintf("%v", obj.Properties[propName])
		for i, opt := range options {
			if opt.Value == currentVal {
				ce.pickerCursor = i
				break
			}
		}
		vm.cellEdit = ce
		return nil

	case "multi_select":
		if len(options) == 0 {
			return nil
		}
		ce.mode = cellModeMultiPick
		ce.pickerOptions = options
		ce.pickerCursor = 0
		ce.pickerChecked = make([]bool, len(options))
		selected := currentMultiSelectValues(obj.Properties[propName])
		for i, opt := range options {
			for _, s := range selected {
				if opt.Value == s {
					ce.pickerChecked[i] = true
					break
				}
			}
		}
		vm.cellEdit = ce
		return nil

	default:
		// string, number, date, datetime, url — use textinput
		ce.mode = cellModeTextInput
		ti := textinput.New()
		ti.Prompt = ""
		ti.CharLimit = 500

		// Pre-fill with current value
		var currentVal string
		if colIdx == 0 {
			currentVal = obj.GetName()
		} else if v := obj.Properties[propName]; v != nil {
			currentVal = fmt.Sprintf("%v", v)
		}
		ti.SetValue(currentVal)
		ti.Focus()
		ti.CursorEnd()

		ce.textInput = ti
		vm.cellEdit = ce
		return textinput.Blink
	}
}

// cancelCellEdit exits cell editing mode without saving.
func (vm *viewMode) cancelCellEdit() {
	vm.cellEdit = nil
}

// applyCellValue validates and applies the edited value, then saves.
func (vm *viewMode) applyCellValue() tea.Cmd {
	ce := vm.cellEdit
	if ce == nil {
		return nil
	}

	obj := ce.obj
	if obj.Properties == nil {
		obj.Properties = make(map[string]any)
	}

	switch ce.mode {
	case cellModeTextInput:
		input := ce.textInput.Value()

		// Validate before applying
		var options []core.Option
		if vm.schema != nil {
			if p := vm.schema.FindProperty(ce.propName); p != nil {
				options = p.Options
			}
		}
		if err := core.ValidatePropertyValue(ce.propType, options, input); err != nil {
			vm.cellEdit = nil
			return func() tea.Msg {
				return viewCellToastMsg{
					Level:   widget.ToastError,
					Message: fmt.Sprintf("Validation: %v", err),
				}
			}
		}

		obj.Properties[ce.propName] = parseEditedValue(ce.propType, input)

	case cellModeDateSegment, cellModeDateCalendar:
		if ce.datePicker != nil {
			obj.Properties[ce.propName] = ce.datePicker.Value()
		}

	case cellModeSelectPick:
		if ce.pickerCursor >= 0 && ce.pickerCursor < len(ce.pickerOptions) {
			obj.Properties[ce.propName] = ce.pickerOptions[ce.pickerCursor].Value
		}

	case cellModeMultiPick:
		var result []any
		for i, checked := range ce.pickerChecked {
			if checked {
				result = append(result, ce.pickerOptions[i].Value)
			}
		}
		obj.Properties[ce.propName] = result
	}

	vm.cellEdit = nil
	return vm.saveCellObject(obj)
}

// saveCellObject persists the object and returns a viewCellSavedMsg or toast error.
func (vm *viewMode) saveCellObject(obj *core.Object) tea.Cmd {
	if err := vm.vault.SaveObject(obj); err != nil {
		return func() tea.Msg {
			return viewCellToastMsg{
				Level:   widget.ToastError,
				Message: fmt.Sprintf("Save failed: %v", err),
			}
		}
	}
	return func() tea.Msg {
		return viewCellSavedMsg{}
	}
}

// updateCellEdit handles key events while a cell edit is active.
func (vm *viewMode) updateCellEdit(msg tea.KeyPressMsg) (*viewMode, tea.Cmd) {
	ce := vm.cellEdit
	if ce == nil {
		return vm, nil
	}

	switch ce.mode {
	case cellModeTextInput:
		switch msg.String() {
		case "enter":
			cmd := vm.applyCellValue()
			return vm, cmd
		case "esc":
			vm.cancelCellEdit()
			return vm, nil
		default:
			var cmd tea.Cmd
			ce.textInput, cmd = ce.textInput.Update(msg)
			return vm, cmd
		}

	case cellModeDateSegment, cellModeDateCalendar:
		if ce.datePicker != nil {
			consumed, done, confirmed := ce.datePicker.Update(msg)
			if consumed {
				if done {
					if confirmed {
						cmd := vm.applyCellValue()
						return vm, cmd
					}
					vm.cancelCellEdit()
					return vm, nil
				}
				if ce.datePicker.Mode() == datePickerCalendar {
					ce.mode = cellModeDateCalendar
				} else {
					ce.mode = cellModeDateSegment
				}
			}
		}
		return vm, nil

	case cellModeSelectPick:
		switch msg.String() {
		case "up", "k":
			if ce.pickerCursor > 0 {
				ce.pickerCursor--
			}
			return vm, nil
		case "down", "j":
			if ce.pickerCursor < len(ce.pickerOptions)-1 {
				ce.pickerCursor++
			}
			return vm, nil
		case "enter":
			cmd := vm.applyCellValue()
			return vm, cmd
		case "esc":
			vm.cancelCellEdit()
			return vm, nil
		}

	case cellModeMultiPick:
		switch msg.String() {
		case "up", "k":
			if ce.pickerCursor > 0 {
				ce.pickerCursor--
			}
			return vm, nil
		case "down", "j":
			if ce.pickerCursor < len(ce.pickerOptions)-1 {
				ce.pickerCursor++
			}
			return vm, nil
		case " ", "space":
			if ce.pickerCursor >= 0 && ce.pickerCursor < len(ce.pickerChecked) {
				ce.pickerChecked[ce.pickerCursor] = !ce.pickerChecked[ce.pickerCursor]
			}
			return vm, nil
		case "enter":
			cmd := vm.applyCellValue()
			return vm, cmd
		case "esc":
			vm.cancelCellEdit()
			return vm, nil
		}
	}

	return vm, nil
}
