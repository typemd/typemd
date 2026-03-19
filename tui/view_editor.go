package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// viewEditorChangedMsg signals the parent that the view config was modified.
type viewEditorChangedMsg struct{}

// viewEditorDeletedMsg signals the parent to delete the view and exit view mode.
type viewEditorDeletedMsg struct {
	TypeName string
	ViewName string
}

// ── Sections and modes ──────────────────────────────────────────────────────

type veSection int

const (
	veSectionFilter veSection = iota
	veSectionSort
	veSectionGroup
	veSectionCount // must be last
)

type veMode int

const (
	veModeView       veMode = iota // browsing rules
	veModePickProp                 // choosing a property (text + list)
	veModePickOp                   // choosing an operator (list)
	veModeInputValue               // entering a filter value
	veModePickDir                  // choosing sort direction
	veModeConfirmDel               // confirm view deletion
)

// veAction tracks whether we're adding or editing a rule.
type veAction int

const (
	veActionAdd  veAction = iota
	veActionEdit
)

// ── View editor struct ──────────────────────────────────────────────────────

type viewEditor struct {
	typeName string
	viewName string
	view     *core.ViewConfig
	schema   *core.TypeSchema
	vault    *core.Vault

	// Section and cursor
	section veSection
	cursor  int // within current section (0 = first rule, last = "+ Add" row)

	// Mode
	mode   veMode
	action veAction
	editIdx int // index of rule being edited (for veActionEdit)

	// Property picker state
	propInput    textinput.Model
	propList     []string // all available property names
	propFiltered []string // filtered by propInput text
	propCursor   int

	// Operator picker state
	opList   []string
	opCursor int

	// Value input state
	valueInput textinput.Model

	// Sort direction picker state
	dirOptions []string
	dirCursor  int

	// In-progress rule being built
	pendingProp string
	pendingOp   string

	// Layout
	width  int
	height int
}

// newViewEditor creates a view editor for the given view.
func newViewEditor(typeName, viewName string, view *core.ViewConfig, schema *core.TypeSchema, vault *core.Vault) *viewEditor {
	pi := textinput.New()
	pi.Prompt = "> "
	pi.CharLimit = 100

	vi := textinput.New()
	vi.Prompt = "Value: "
	vi.CharLimit = 200

	ve := &viewEditor{
		typeName:   typeName,
		viewName:   viewName,
		view:       view,
		schema:     schema,
		vault:      vault,
		mode:       veModeView,
		dirOptions: []string{"asc", "desc"},
	}

	ve.propInput = pi
	ve.valueInput = vi
	ve.buildPropList()

	return ve
}

// buildPropList builds the full property list from schema + system properties.
func (ve *viewEditor) buildPropList() {
	// System properties first
	props := []string{"name", "description", "tags", "created_at", "updated_at"}

	// Schema properties
	if ve.schema != nil {
		for _, p := range ve.schema.Properties {
			props = append(props, p.Name)
		}
	}

	ve.propList = props
	ve.filterProps()
}

// filterProps updates propFiltered based on current propInput text.
func (ve *viewEditor) filterProps() {
	query := strings.ToLower(ve.propInput.Value())
	if query == "" {
		ve.propFiltered = make([]string, len(ve.propList))
		copy(ve.propFiltered, ve.propList)
		return
	}
	ve.propFiltered = nil
	for _, p := range ve.propList {
		if strings.Contains(strings.ToLower(p), query) {
			ve.propFiltered = append(ve.propFiltered, p)
		}
	}
	if ve.propCursor >= len(ve.propFiltered) {
		ve.propCursor = max(0, len(ve.propFiltered)-1)
	}
}

// ── Section item counts ─────────────────────────────────────────────────────

func (ve *viewEditor) sectionItemCount(s veSection) int {
	switch s {
	case veSectionFilter:
		return len(ve.view.Filter) + 1 // +1 for "+ Add Filter"
	case veSectionSort:
		return len(ve.view.Sort) + 1
	case veSectionGroup:
		return len(ve.view.GroupBy) + 1
	}
	return 0
}

func (ve *viewEditor) sectionRuleCount(s veSection) int {
	switch s {
	case veSectionFilter:
		return len(ve.view.Filter)
	case veSectionSort:
		return len(ve.view.Sort)
	case veSectionGroup:
		return len(ve.view.GroupBy)
	}
	return 0
}

// ── Update ──────────────────────────────────────────────────────────────────

func (ve *viewEditor) Update(msg tea.Msg) (*viewEditor, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return ve, nil
	}

	switch ve.mode {
	case veModeView:
		return ve.updateView(keyMsg)
	case veModePickProp:
		return ve.updatePickProp(keyMsg)
	case veModePickOp:
		return ve.updatePickOp(keyMsg)
	case veModeInputValue:
		return ve.updateInputValue(keyMsg)
	case veModePickDir:
		return ve.updatePickDir(keyMsg)
	case veModeConfirmDel:
		return ve.updateConfirmDel(keyMsg)
	}
	return ve, nil
}

func (ve *viewEditor) updateView(msg tea.KeyPressMsg) (*viewEditor, tea.Cmd) {
	switch msg.String() {
	case "tab":
		ve.section = (ve.section + 1) % veSectionCount
		ve.cursor = 0
	case "up", "k":
		if ve.cursor > 0 {
			ve.cursor--
		}
	case "down", "j":
		if ve.cursor < ve.sectionItemCount(ve.section)-1 {
			ve.cursor++
		}
	case "enter":
		ruleCount := ve.sectionRuleCount(ve.section)
		if ve.cursor == ruleCount {
			// "+ Add" row
			return ve.startAdd()
		}
		// Edit existing rule
		return ve.startEdit()
	case "x", "d":
		return ve.deleteRule()
	case "D":
		ve.mode = veModeConfirmDel
	case "K":
		return ve.moveRule(-1)
	case "J":
		return ve.moveRule(1)
	}
	return ve, nil
}

func (ve *viewEditor) updatePickProp(msg tea.KeyPressMsg) (*viewEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		ve.mode = veModeView
		ve.propInput.Blur()
		return ve, nil
	case "enter":
		if len(ve.propFiltered) > 0 && ve.propCursor < len(ve.propFiltered) {
			ve.pendingProp = ve.propFiltered[ve.propCursor]
			ve.propInput.Blur()
			return ve.afterPropPicked()
		}
		return ve, nil
	case "up":
		if ve.propCursor > 0 {
			ve.propCursor--
		}
		return ve, nil
	case "down":
		if ve.propCursor < len(ve.propFiltered)-1 {
			ve.propCursor++
		}
		return ve, nil
	}
	// Forward to text input
	var cmd tea.Cmd
	ve.propInput, cmd = ve.propInput.Update(msg)
	ve.filterProps()
	return ve, cmd
}

func (ve *viewEditor) updatePickOp(msg tea.KeyPressMsg) (*viewEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Go back to property picker
		ve.mode = veModePickProp
		return ve, ve.propInput.Focus()
	case "enter":
		if ve.opCursor < len(ve.opList) {
			ve.pendingOp = ve.opList[ve.opCursor]
			// Skip value input for is_empty/is_not_empty
			if ve.pendingOp == "is_empty" || ve.pendingOp == "is_not_empty" {
				return ve.commitFilterRule("")
			}
			ve.mode = veModeInputValue
			ve.valueInput.SetValue("")
			if ve.action == veActionEdit {
				// Pre-populate value from existing rule
				if ve.editIdx < len(ve.view.Filter) {
					ve.valueInput.SetValue(ve.view.Filter[ve.editIdx].Value)
				}
			}
			return ve, ve.valueInput.Focus()
		}
		return ve, nil
	case "up", "k":
		if ve.opCursor > 0 {
			ve.opCursor--
		}
	case "down", "j":
		if ve.opCursor < len(ve.opList)-1 {
			ve.opCursor++
		}
	}
	return ve, nil
}

func (ve *viewEditor) updateInputValue(msg tea.KeyPressMsg) (*viewEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		ve.valueInput.Blur()
		ve.mode = veModePickOp
		return ve, nil
	case "enter":
		val := ve.valueInput.Value()
		ve.valueInput.Blur()
		return ve.commitFilterRule(val)
	}
	var cmd tea.Cmd
	ve.valueInput, cmd = ve.valueInput.Update(msg)
	return ve, cmd
}

func (ve *viewEditor) updatePickDir(msg tea.KeyPressMsg) (*viewEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Go back to property picker
		ve.mode = veModePickProp
		return ve, ve.propInput.Focus()
	case "enter":
		dir := ve.dirOptions[ve.dirCursor]
		return ve.commitSortRule(dir)
	case "up", "k", "down", "j", " ":
		ve.dirCursor = 1 - ve.dirCursor // toggle between 0 and 1
	}
	return ve, nil
}

func (ve *viewEditor) updateConfirmDel(msg tea.KeyPressMsg) (*viewEditor, tea.Cmd) {
	switch msg.String() {
	case "y":
		ve.mode = veModeView
		return ve, tea.Sequence(func() tea.Msg {
			return viewEditorDeletedMsg{TypeName: ve.typeName, ViewName: ve.viewName}
		})
	case "n", "esc":
		ve.mode = veModeView
	}
	return ve, nil
}

// ── Rule operations ─────────────────────────────────────────────────────────

func (ve *viewEditor) startAdd() (*viewEditor, tea.Cmd) {
	ve.action = veActionAdd
	ve.propInput.SetValue("")
	ve.propCursor = 0
	ve.filterProps()
	ve.mode = veModePickProp
	return ve, ve.propInput.Focus()
}

func (ve *viewEditor) startEdit() (*viewEditor, tea.Cmd) {
	ve.action = veActionEdit
	ve.editIdx = ve.cursor

	// Pre-populate property picker with current value
	var currentProp string
	switch ve.section {
	case veSectionFilter:
		if ve.cursor < len(ve.view.Filter) {
			currentProp = ve.view.Filter[ve.cursor].Property
		}
	case veSectionSort:
		if ve.cursor < len(ve.view.Sort) {
			currentProp = ve.view.Sort[ve.cursor].Property
		}
	case veSectionGroup:
		if ve.cursor < len(ve.view.GroupBy) {
			currentProp = ve.view.GroupBy[ve.cursor].Property
		}
	}

	ve.propInput.SetValue(currentProp)
	ve.propCursor = 0
	ve.filterProps()
	// Try to position cursor on current prop
	for i, p := range ve.propFiltered {
		if p == currentProp {
			ve.propCursor = i
			break
		}
	}

	ve.mode = veModePickProp
	return ve, ve.propInput.Focus()
}

func (ve *viewEditor) afterPropPicked() (*viewEditor, tea.Cmd) {
	switch ve.section {
	case veSectionFilter:
		// Need operator picker
		ve.opList = ve.operatorsForProp(ve.pendingProp)
		ve.opCursor = 0
		if ve.action == veActionEdit && ve.editIdx < len(ve.view.Filter) {
			// Pre-select current operator
			for i, op := range ve.opList {
				if op == ve.view.Filter[ve.editIdx].Operator {
					ve.opCursor = i
					break
				}
			}
		}
		ve.mode = veModePickOp
		return ve, nil
	case veSectionSort:
		// Need direction picker
		ve.dirCursor = 0
		if ve.action == veActionEdit && ve.editIdx < len(ve.view.Sort) {
			if ve.view.Sort[ve.editIdx].Direction == "desc" {
				ve.dirCursor = 1
			}
		}
		ve.mode = veModePickDir
		return ve, nil
	case veSectionGroup:
		// Group rule only needs property — commit immediately
		return ve.commitGroupRule()
	}
	return ve, nil
}

func (ve *viewEditor) commitFilterRule(value string) (*viewEditor, tea.Cmd) {
	rule := core.FilterRule{
		Property: ve.pendingProp,
		Operator: ve.pendingOp,
		Value:    value,
	}
	if ve.action == veActionEdit && ve.editIdx < len(ve.view.Filter) {
		ve.view.Filter[ve.editIdx] = rule
	} else {
		ve.view.Filter = append(ve.view.Filter, rule)
	}
	ve.mode = veModeView
	ve.save()
	return ve, tea.Sequence(func() tea.Msg { return viewEditorChangedMsg{} })
}

func (ve *viewEditor) commitSortRule(direction string) (*viewEditor, tea.Cmd) {
	rule := core.SortRule{
		Property:  ve.pendingProp,
		Direction: direction,
	}
	if ve.action == veActionEdit && ve.editIdx < len(ve.view.Sort) {
		ve.view.Sort[ve.editIdx] = rule
	} else {
		ve.view.Sort = append(ve.view.Sort, rule)
	}
	ve.mode = veModeView
	ve.save()
	return ve, tea.Sequence(func() tea.Msg { return viewEditorChangedMsg{} })
}

func (ve *viewEditor) commitGroupRule() (*viewEditor, tea.Cmd) {
	rule := core.GroupRule{Property: ve.pendingProp}
	if ve.action == veActionEdit && ve.editIdx < len(ve.view.GroupBy) {
		ve.view.GroupBy[ve.editIdx] = rule
	} else {
		ve.view.GroupBy = append(ve.view.GroupBy, rule)
	}
	ve.mode = veModeView
	ve.save()
	return ve, tea.Sequence(func() tea.Msg { return viewEditorChangedMsg{} })
}

func (ve *viewEditor) deleteRule() (*viewEditor, tea.Cmd) {
	ruleCount := ve.sectionRuleCount(ve.section)
	if ve.cursor >= ruleCount {
		return ve, nil // on "+ Add" row, nothing to delete
	}
	switch ve.section {
	case veSectionFilter:
		ve.view.Filter = append(ve.view.Filter[:ve.cursor], ve.view.Filter[ve.cursor+1:]...)
	case veSectionSort:
		ve.view.Sort = append(ve.view.Sort[:ve.cursor], ve.view.Sort[ve.cursor+1:]...)
	case veSectionGroup:
		ve.view.GroupBy = append(ve.view.GroupBy[:ve.cursor], ve.view.GroupBy[ve.cursor+1:]...)
	}
	if ve.cursor >= ve.sectionItemCount(ve.section) {
		ve.cursor = max(0, ve.sectionItemCount(ve.section)-1)
	}
	ve.save()
	return ve, tea.Sequence(func() tea.Msg { return viewEditorChangedMsg{} })
}

func (ve *viewEditor) moveRule(delta int) (*viewEditor, tea.Cmd) {
	ruleCount := ve.sectionRuleCount(ve.section)
	if ve.cursor >= ruleCount {
		return ve, nil // on "+ Add" row
	}
	newIdx := ve.cursor + delta
	if newIdx < 0 || newIdx >= ruleCount {
		return ve, nil // boundary
	}
	switch ve.section {
	case veSectionFilter:
		ve.view.Filter[ve.cursor], ve.view.Filter[newIdx] = ve.view.Filter[newIdx], ve.view.Filter[ve.cursor]
	case veSectionSort:
		ve.view.Sort[ve.cursor], ve.view.Sort[newIdx] = ve.view.Sort[newIdx], ve.view.Sort[ve.cursor]
	case veSectionGroup:
		ve.view.GroupBy[ve.cursor], ve.view.GroupBy[newIdx] = ve.view.GroupBy[newIdx], ve.view.GroupBy[ve.cursor]
	}
	ve.cursor = newIdx
	ve.save()
	return ve, tea.Sequence(func() tea.Msg { return viewEditorChangedMsg{} })
}

func (ve *viewEditor) save() {
	if ve.vault == nil {
		return
	}
	_ = ve.vault.SaveView(ve.typeName, ve.view)
}

// operatorsForProp returns the valid operators for the given property name.
func (ve *viewEditor) operatorsForProp(propName string) []string {
	propType := ve.resolvePropType(propName)
	return core.OperatorsForType(propType)
}

// resolvePropType looks up the property type from the schema.
func (ve *viewEditor) resolvePropType(propName string) string {
	// System property types
	switch propName {
	case "name", "description":
		return "string"
	case "created_at", "updated_at":
		return "datetime"
	case "tags":
		return "relation"
	}
	if ve.schema != nil {
		if p := ve.schema.FindProperty(propName); p != nil {
			return p.Type
		}
	}
	return "string" // fallback
}

// ── View ────────────────────────────────────────────────────────────────────

func (ve *viewEditor) View() string {
	if ve.mode == veModeConfirmDel {
		return fmt.Sprintf("\n  Delete view '%s'? [y/n]\n", ve.viewName)
	}

	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6"))
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	b.WriteString(titleStyle.Render(fmt.Sprintf(" Edit: %s", ve.viewName)) + "\n")
	b.WriteString(dimStyle.Render(" " + strings.Repeat("─", ve.width-4)) + "\n")

	type sectionDef struct {
		label     string
		section   veSection
		count     int
		format    func(int) string
		addLabel  string
	}
	sections := []sectionDef{
		{"📋 Filter", veSectionFilter, len(ve.view.Filter), ve.formatFilterRule, "Add Filter"},
		{"📊 Sort", veSectionSort, len(ve.view.Sort), ve.formatSortRule, "Add Sort"},
		{"📁 Group By", veSectionGroup, len(ve.view.GroupBy), ve.formatGroupRule, "Add Group"},
	}

	for _, sec := range sections {
		isActive := ve.section == sec.section
		label := sec.label
		if isActive {
			label = sectionStyle.Render(label)
		} else {
			label = dimStyle.Render(label)
		}
		b.WriteString("\n " + label + "\n")
		ve.renderRuleSection(&b, isActive, sec.count, sec.format, sec.addLabel, activeStyle, addStyle)
	}

	// Show picker/input overlay if in a sub-mode
	if ve.mode != veModeView {
		b.WriteString("\n")
		ve.renderSubMode(&b, activeStyle)
	}

	return b.String()
}

func (ve *viewEditor) formatFilterRule(i int) string {
	rule := ve.view.Filter[i]
	text := fmt.Sprintf("  %s %s", rule.Property, rule.Operator)
	if rule.Value != "" {
		text += " " + rule.Value
	}
	return text
}

func (ve *viewEditor) formatSortRule(i int) string {
	rule := ve.view.Sort[i]
	arrow := "↑"
	if rule.Direction == "desc" {
		arrow = "↓"
	}
	return fmt.Sprintf("  %s %s", rule.Property, arrow)
}

func (ve *viewEditor) formatGroupRule(i int) string {
	return fmt.Sprintf("  %s", ve.view.GroupBy[i].Property)
}

// renderRuleSection renders a list of rules with a highlighted cursor and an "+ Add" row.
func (ve *viewEditor) renderRuleSection(b *strings.Builder, active bool, ruleCount int, formatRule func(int) string, addLabel string, activeStyle, addStyle lipgloss.Style) {
	for i := 0; i < ruleCount; i++ {
		text := formatRule(i)
		if active && ve.cursor == i && ve.mode == veModeView {
			b.WriteString(" " + activeStyle.Render(text) + "\n")
		} else {
			b.WriteString(" " + text + "\n")
		}
	}
	addText := "  + " + addLabel
	if active && ve.cursor == ruleCount && ve.mode == veModeView {
		b.WriteString(" " + activeStyle.Render(addText) + "\n")
	} else {
		b.WriteString(" " + addStyle.Render(addText) + "\n")
	}
}

func (ve *viewEditor) renderSubMode(b *strings.Builder, activeStyle lipgloss.Style) {
	switch ve.mode {
	case veModePickProp:
		b.WriteString(" " + ve.propInput.View() + "\n")
		visibleH := ve.height - 20
		if visibleH < 3 {
			visibleH = 3
		}
		scrolled := widget.AdjustScroll(ve.propCursor, 0, visibleH)
		for i := scrolled; i < len(ve.propFiltered) && i < scrolled+visibleH; i++ {
			p := ve.propFiltered[i]
			if i == ve.propCursor {
				b.WriteString("   " + activeStyle.Render(p) + "\n")
			} else {
				b.WriteString("   " + p + "\n")
			}
		}
	case veModePickOp:
		b.WriteString(" Select operator:\n")
		for i, op := range ve.opList {
			if i == ve.opCursor {
				b.WriteString("   " + activeStyle.Render(op) + "\n")
			} else {
				b.WriteString("   " + op + "\n")
			}
		}
	case veModeInputValue:
		b.WriteString(" " + ve.valueInput.View() + "\n")
	case veModePickDir:
		b.WriteString(" Select direction:\n")
		for i, d := range ve.dirOptions {
			label := d
			if d == "asc" {
				label = "↑ asc"
			} else {
				label = "↓ desc"
			}
			if i == ve.dirCursor {
				b.WriteString("   " + activeStyle.Render(label) + "\n")
			} else {
				b.WriteString("   " + label + "\n")
			}
		}
	}
}

// ── Help bar ────────────────────────────────────────────────────────────────

func (ve *viewEditor) HelpBar() string {
	switch ve.mode {
	case veModeView:
		return "  [VIEW EDITOR]  ↑↓: navigate  J/K: move  tab: section  enter: edit  x: delete  D: delete view  esc: close"
	case veModePickProp:
		return "  [SELECT PROPERTY]  type to filter  ↑↓: navigate  enter: select  esc: cancel"
	case veModePickOp:
		return "  [SELECT OPERATOR]  ↑↓: navigate  enter: select  esc: back"
	case veModeInputValue:
		return "  [ENTER VALUE]  enter: confirm  esc: back"
	case veModePickDir:
		return "  [SELECT DIRECTION]  ↑↓/space: toggle  enter: confirm  esc: back"
	case veModeConfirmDel:
		return "  [DELETE VIEW]  y: confirm  n/esc: cancel"
	}
	return ""
}

// ── Layout ──────────────────────────────────────────────────────────────────

func (ve *viewEditor) SetSize(width, height int) {
	ve.width = width
	ve.height = height
}

// CanQuit returns true when the editor is in a non-interactive state.
func (ve *viewEditor) CanQuit() bool {
	return ve.mode == veModeView
}
