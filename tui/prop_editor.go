package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// propEditMode represents the current editing state of the property editor.
type propEditMode int

const (
	propModeNavigate          propEditMode = iota // cursor navigation, no editing
	propModeTextInput                              // textinput active (string, number, date, datetime, url)
	propModeSelectPick                             // select option picker open
	propModeMultiPick                              // multi_select option picker open
	propModeRelationPick                           // single-value relation picker open
	propModeRelationMultiPick                      // multi-value relation picker open
	propModeDateSegment                            // date picker segmented input mode
	propModeDateCalendar                           // date picker calendar mode
)

// propEditor manages cursor navigation and inline editing in the properties panel.
type propEditor struct {
	// Display state
	items  []propItem // editable + display items (filtered from DisplayProperty)
	cursor int        // index into items (only stops on editable items)

	// Edit state
	mode      propEditMode
	textInput textinput.Model
	editIndex int // index of item currently being edited

	// Select/multi-select picker state
	pickerOptions []core.Option // available options for current select/multi_select
	pickerCursor  int           // cursor within picker
	pickerChecked []bool        // for multi_select: which options are checked

	// Relation picker state
	relCandidates     []relationCandidate // all candidates (unfiltered)
	relFiltered       []relationCandidate // candidates matching search
	relSearch         string              // current search text
	relChecked        []bool              // for multi-value: which candidates are checked
	relInitialChecked []bool              // snapshot at picker open for diff calculation

	// Date picker state
	datePicker *datePicker
}

// propItem represents a single property row in the editor.
type propItem struct {
	dp       core.DisplayProperty
	editable bool
	schema   *core.Property // schema property (for options, type info); nil for system props
}

// resolvedType returns the effective property type, preferring schema type over display type.
func (pi *propItem) resolvedType() string {
	if pi.schema != nil {
		return pi.schema.Type
	}
	return pi.dp.Type
}

// newPropEditor creates a property editor from display properties and schema.
func newPropEditor(displayProps []core.DisplayProperty, schema *core.TypeSchema) *propEditor {
	pe := &propEditor{
		mode: propModeNavigate,
	}
	pe.updateItems(displayProps, schema)
	return pe
}

// updateItems rebuilds the item list from display properties.
func (pe *propEditor) updateItems(displayProps []core.DisplayProperty, schema *core.TypeSchema) {
	pe.items = nil
	for _, dp := range displayProps {
		// Skip pinned (displayed in body) and name (displayed in title)
		if dp.Pin > 0 || dp.Key == core.NameProperty {
			continue
		}

		item := propItem{dp: dp}

		// Determine editability
		item.editable = isPropertyEditable(dp)

		// Find schema property for type info and options
		if schema != nil {
			item.schema = schema.FindProperty(dp.Key)
		}

		pe.items = append(pe.items, item)
	}

	// Clamp cursor
	if pe.cursor >= len(pe.items) {
		pe.cursor = len(pe.items) - 1
	}
	if pe.cursor < 0 {
		pe.cursor = 0
	}
	// Snap to nearest editable
	pe.snapToEditable(1)
}

// isPropertyEditable determines if a display property can be edited inline.
func isPropertyEditable(dp core.DisplayProperty) bool {
	// Immutable system properties
	if core.IsImmutableSystemProperty(dp.Key) {
		return false
	}
	// Reverse relations and backlinks are always read-only
	if dp.IsReverse || dp.IsBacklink {
		return false
	}
	// Local properties (not in schema, not system) are read-only
	if dp.IsLocal {
		return false
	}
	// Forward relations and tags are editable via relation picker
	// All other property types are editable via their respective editors
	return true
}

// snapToEditable moves the cursor to the nearest editable item in the given direction.
func (pe *propEditor) snapToEditable(dir int) {
	if len(pe.items) == 0 {
		return
	}
	// If current is editable, stay
	if pe.cursor >= 0 && pe.cursor < len(pe.items) && pe.items[pe.cursor].editable {
		return
	}
	// Search in direction
	for i := 1; i < len(pe.items); i++ {
		idx := pe.cursor + i*dir
		if idx >= 0 && idx < len(pe.items) && pe.items[idx].editable {
			pe.cursor = idx
			return
		}
	}
	// Try opposite direction
	for i := 1; i < len(pe.items); i++ {
		idx := pe.cursor - i*dir
		if idx >= 0 && idx < len(pe.items) && pe.items[idx].editable {
			pe.cursor = idx
			return
		}
	}
}

// moveUp moves cursor to the previous editable item.
func (pe *propEditor) moveUp() {
	for i := pe.cursor - 1; i >= 0; i-- {
		if pe.items[i].editable {
			pe.cursor = i
			return
		}
	}
}

// moveDown moves cursor to the next editable item.
func (pe *propEditor) moveDown() {
	for i := pe.cursor + 1; i < len(pe.items); i++ {
		if pe.items[i].editable {
			pe.cursor = i
			return
		}
	}
}

// currentItem returns the item under the cursor, or nil if out of bounds.
func (pe *propEditor) currentItem() *propItem {
	if pe.cursor >= 0 && pe.cursor < len(pe.items) {
		return &pe.items[pe.cursor]
	}
	return nil
}

// activateEdit starts editing the current property based on its type.
// vault is needed for relation picker candidate loading.
func (pe *propEditor) activateEdit(vault *core.Vault) tea.Cmd {
	item := pe.currentItem()
	if item == nil || !item.editable {
		return nil
	}

	// Relations (schema-defined or tags system property)
	if item.dp.IsRelation || item.dp.Key == core.TagsProperty {
		return pe.activateRelationPicker(item, vault)
	}

	switch item.resolvedType() {
	case "checkbox":
		return nil // direct toggle handled in Update
	case "select":
		return pe.activateSelectPicker(item)
	case "multi_select":
		return pe.activateMultiPicker(item)
	case "date":
		return pe.activateDatePicker(item)
	default:
		return pe.activateTextInput(item)
	}
}

// activateTextInput opens a textinput for the current property.
func (pe *propEditor) activateTextInput(item *propItem) tea.Cmd {
	pe.mode = propModeTextInput
	pe.editIndex = pe.cursor

	ti := textinput.New()
	ti.CharLimit = 500
	ti.SetWidth(40)

	// Pre-fill with current value
	currentVal := formatEditValue(item.dp)
	ti.SetValue(currentVal)
	ti.Focus()
	ti.CursorEnd()

	pe.textInput = ti
	return textinput.Blink
}

// activateSelectPicker opens the select option picker.
func (pe *propEditor) activateSelectPicker(item *propItem) tea.Cmd {
	if item.schema == nil || len(item.schema.Options) == 0 {
		return nil
	}
	pe.mode = propModeSelectPick
	pe.editIndex = pe.cursor
	pe.pickerOptions = item.schema.Options
	pe.pickerCursor = 0

	// Position cursor on current value
	currentVal := fmt.Sprintf("%v", item.dp.Value)
	for i, opt := range pe.pickerOptions {
		if opt.Value == currentVal {
			pe.pickerCursor = i
			break
		}
	}
	return nil
}

// activateMultiPicker opens the multi-select option picker.
func (pe *propEditor) activateMultiPicker(item *propItem) tea.Cmd {
	if item.schema == nil || len(item.schema.Options) == 0 {
		return nil
	}
	pe.mode = propModeMultiPick
	pe.editIndex = pe.cursor
	pe.pickerOptions = item.schema.Options
	pe.pickerCursor = 0

	// Initialize checked state from current value
	pe.pickerChecked = make([]bool, len(pe.pickerOptions))
	selected := currentMultiSelectValues(item.dp.Value)
	for i, opt := range pe.pickerOptions {
		for _, s := range selected {
			if opt.Value == s {
				pe.pickerChecked[i] = true
				break
			}
		}
	}
	return nil
}

func (pe *propEditor) activateDatePicker(item *propItem) tea.Cmd {
	pe.mode = propModeDateSegment
	pe.editIndex = pe.cursor
	pe.datePicker = newDatePicker(parseDatePickerValue(item.dp.Value))
	return nil
}

// cancelEdit cancels the current edit and returns to navigation mode.
func (pe *propEditor) cancelEdit() {
	pe.mode = propModeNavigate
	pe.textInput = textinput.Model{}
	pe.relCandidates = nil
	pe.relFiltered = nil
	pe.relSearch = ""
	pe.relChecked = nil
	pe.relInitialChecked = nil
	pe.datePicker = nil
}

// isEditing returns true if the editor is in any editing state.
func (pe *propEditor) isEditing() bool {
	return pe.mode != propModeNavigate
}

// formatEditValue returns the current value formatted for textinput pre-fill.
func formatEditValue(dp core.DisplayProperty) string {
	if dp.Value == nil {
		return ""
	}
	return dp.FormatValue()
}

// currentMultiSelectValues extracts current multi_select values as strings.
func currentMultiSelectValues(val any) []string {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case []any:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case []string:
		return v
	case string:
		return []string{v}
	default:
		return nil
	}
}

// parseEditedValue converts the textinput string back to the appropriate Go value
// for storing in obj.Properties.
func parseEditedValue(propType string, input string) any {
	switch propType {
	case "number":
		// Try int first, then float
		if i, err := strconv.Atoi(input); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(input, 64); err == nil {
			return f
		}
		return input
	case "checkbox":
		return true // toggle logic handles this
	default:
		return input
	}
}

// multiPickerResult returns the selected values from the multi-select picker.
func (pe *propEditor) multiPickerResult() []any {
	var result []any
	for i, checked := range pe.pickerChecked {
		if checked {
			result = append(result, pe.pickerOptions[i].Value)
		}
	}
	return result
}

// Render renders the properties panel with cursor and editing widgets.
func (pe *propEditor) Render(focused bool) string {
	if len(pe.items) == 0 {
		return " Properties\n ──────────\n (none)\n"
	}

	var b strings.Builder
	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	localSepRendered := false
	for i, item := range pe.items {
		// Render separator before the first local property
		if item.dp.IsLocal && !localSepRendered {
			b.WriteString(dimStyle.Render(" ── Local Properties ──") + "\n")
			localSepRendered = true
		}

		isCursor := focused && i == pe.cursor && item.editable

		// If this item is being edited with a textinput
		if pe.mode == propModeTextInput && i == pe.editIndex {
			b.WriteString(fmt.Sprintf(" %s: %s\n", item.dp.Key, pe.textInput.View()))
			continue
		}

		// If this item has a date picker open
		if (pe.mode == propModeDateSegment || pe.mode == propModeDateCalendar) && i == pe.editIndex && pe.datePicker != nil {
			b.WriteString(fmt.Sprintf(" %s: %s\n", item.dp.Key, pe.datePicker.View()))
			continue
		}

		// If this item has a select/multi-select picker open
		if (pe.mode == propModeSelectPick || pe.mode == propModeMultiPick) && i == pe.editIndex {
			b.WriteString(pe.renderPickerItem(item, isCursor))
			continue
		}

		// If this item has a relation picker open
		if (pe.mode == propModeRelationPick || pe.mode == propModeRelationMultiPick) && i == pe.editIndex {
			b.WriteString(pe.renderRelationPicker(item, isCursor))
			continue
		}

		// Normal display
		prefix := " "
		if isCursor {
			prefix = "▸"
		}

		formatted := item.dp.Format()
		if isCursor {
			formatted = highlightStyle.Render(prefix + " " + formatted)
		} else if !item.editable {
			formatted = dimStyle.Render(prefix + " " + formatted)
		} else {
			formatted = prefix + " " + formatted
		}

		b.WriteString(formatted + "\n")
	}

	return b.String()
}

// renderPickerItem renders a property with its option picker.
func (pe *propEditor) renderPickerItem(item propItem, isCursor bool) string {
	var b strings.Builder

	// Property label
	prefix := " "
	if isCursor {
		prefix = "▸"
	}
	b.WriteString(fmt.Sprintf("%s %s:\n", prefix, item.dp.Key))

	// Options list
	for i, opt := range pe.pickerOptions {
		label := opt.Value
		if opt.Label != "" {
			label = opt.Label
		}

		if pe.mode == propModeMultiPick {
			check := "☐"
			if pe.pickerChecked[i] {
				check = "☑"
			}
			if i == pe.pickerCursor {
				b.WriteString(highlightStyle.Render(fmt.Sprintf("   %s %s", check, label)) + "\n")
			} else {
				b.WriteString(fmt.Sprintf("   %s %s\n", check, label))
			}
		} else {
			if i == pe.pickerCursor {
				b.WriteString(highlightStyle.Render(fmt.Sprintf("   ▸ %s", label)) + "\n")
			} else {
				b.WriteString(fmt.Sprintf("     %s\n", label))
			}
		}
	}

	return b.String()
}
