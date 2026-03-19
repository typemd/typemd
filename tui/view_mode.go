package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// openViewMsg signals the parent model to enter view mode.
type openViewMsg struct {
	TypeName string
	ViewName string
}

// viewMode is a sub-model for the full-width view display.
type viewMode struct {
	typeName string
	viewName string
	view     *core.ViewConfig
	schema   *core.TypeSchema
	vault    *core.Vault

	objects []*core.Object
	groups  []viewGroup
	cursor  int
	scroll  int

	// When viewing an object from the view list
	detailObject *core.Object

	// Preview: the object under cursor for side panel preview
	previewObject *core.Object

	// View editor: shown as right split panel (mutually exclusive with preview)
	editor *viewEditor

	width        int
	height       int
	previewWidth int // content width of preview panel (0 = not in split mode)
}

// viewGroup represents a group of objects in the view list.
type viewGroup struct {
	Label    string
	Objects  []*core.Object
	Expanded bool
}

// newViewMode creates a new viewMode for the given type and view.
func newViewMode(typeName, viewName string, vault *core.Vault) *viewMode {
	vm := &viewMode{
		typeName: typeName,
		viewName: viewName,
		vault:    vault,
	}
	vm.load()
	return vm
}

// load fetches the view config, schema, and objects.
func (vm *viewMode) load() {
	vc, err := vm.vault.LoadView(vm.typeName, vm.viewName)
	if err != nil {
		vm.view = vm.vault.DefaultView(vm.typeName)
	} else {
		vm.view = vc
	}

	vm.schema, _ = vm.vault.LoadType(vm.typeName)
	vm.queryAndGroup()
}

// queryAndGroup re-queries objects using the current view config and rebuilds groups.
func (vm *viewMode) queryAndGroup() {
	filter := append(core.TypeFilter(vm.typeName), vm.view.Filter...)
	objects, err := vm.vault.QueryObjects(filter, vm.view.Sort...)
	if err != nil {
		vm.objects = nil
	} else {
		vm.objects = objects
	}
	vm.buildGroups()
}

// buildGroups organizes objects by group_by rules or as a flat list.
// Supports multi-level grouping with compound labels joined by " · ".
func (vm *viewMode) buildGroups() {
	if len(vm.view.GroupBy) == 0 {
		// Flat list — single unnamed group
		vm.groups = []viewGroup{{
			Label:    "",
			Objects:  vm.objects,
			Expanded: true,
		}}
		return
	}

	// Build compound group key for each object
	groupMap := make(map[string][]*core.Object)
	var groupOrder []string
	for _, obj := range vm.objects {
		key := vm.compoundGroupKey(obj)
		if _, exists := groupMap[key]; !exists {
			groupOrder = append(groupOrder, key)
		}
		groupMap[key] = append(groupMap[key], obj)
	}

	vm.groups = make([]viewGroup, 0, len(groupOrder))
	for _, label := range groupOrder {
		vm.groups = append(vm.groups, viewGroup{
			Label:    label,
			Objects:  groupMap[label],
			Expanded: true,
		})
	}
}

// compoundGroupKey builds a compound group label from all GroupBy rules.
// Values are joined with " · ". Missing values become "(none)".
func (vm *viewMode) compoundGroupKey(obj *core.Object) string {
	parts := make([]string, len(vm.view.GroupBy))
	for i, rule := range vm.view.GroupBy {
		val := vm.displayPropValue(obj, rule.Property)
		if val == "" {
			val = "(none)"
		}
		parts[i] = val
	}
	return strings.Join(parts, " · ")
}

// visibleRows returns the flattened list of displayable rows.
func (vm *viewMode) visibleRows() []viewRow {
	var rows []viewRow
	for gi, g := range vm.groups {
		if len(vm.view.GroupBy) > 0 {
			rows = append(rows, viewRow{isHeader: true, groupIdx: gi, label: g.Label})
		}
		if g.Expanded {
			for _, obj := range g.Objects {
				rows = append(rows, viewRow{groupIdx: gi, object: obj})
			}
		}
	}
	return rows
}

type viewRow struct {
	isHeader bool
	groupIdx int
	label    string
	object   *core.Object
}

// SetSize updates the viewport dimensions.
func (vm *viewMode) SetSize(width, height int) {
	vm.width = width
	vm.height = height
}

// CanQuit returns true if the view mode can safely exit.
func (vm *viewMode) CanQuit() bool {
	if vm.editor != nil {
		return vm.editor.CanQuit()
	}
	return vm.detailObject == nil
}

// expandedGroupLabels returns the labels of currently expanded groups.
func (vm *viewMode) expandedGroupLabels() []string {
	var labels []string
	for _, g := range vm.groups {
		if g.Expanded && g.Label != "" {
			labels = append(labels, g.Label)
		}
	}
	return labels
}

// HasEditor returns true if the view editor is open.
func (vm *viewMode) HasEditor() bool {
	return vm.editor != nil
}

// EditorView renders the editor panel content.
func (vm *viewMode) EditorView() string {
	if vm.editor == nil {
		return ""
	}
	return vm.editor.View()
}

// SetEditorSize sets the editor panel dimensions.
func (vm *viewMode) SetEditorSize(width, height int) {
	if vm.editor != nil {
		vm.editor.SetSize(width, height)
	}
}

// EditorHelpBar returns the editor's help bar text.
func (vm *viewMode) EditorHelpBar() string {
	if vm.editor == nil {
		return ""
	}
	return vm.editor.HelpBar()
}

// Update handles messages for the view mode.
func (vm *viewMode) Update(msg tea.Msg) (*viewMode, tea.Cmd) {
	// Handle editor messages
	switch msg.(type) {
	case viewEditorChangedMsg:
		// Editor changed the view config — reload objects
		vm.reloadObjects()
		return vm, nil
	case viewEditorDeletedMsg:
		// Handled by parent (app.go)
		return vm, nil
	}

	// Delegate to editor if open
	if vm.editor != nil {
		keyMsg, ok := msg.(tea.KeyPressMsg)
		if ok && keyMsg.String() == "esc" && vm.editor.CanQuit() {
			// Close editor
			vm.editor = nil
			return vm, nil
		}
		var cmd tea.Cmd
		vm.editor, cmd = vm.editor.Update(msg)
		return vm, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if vm.cursor > 0 {
				vm.cursor--
				if vm.previewObject != nil {
					vm.updatePreview()
				}
			}
		case "down", "j":
			rows := vm.visibleRows()
			if vm.cursor < len(rows)-1 {
				vm.cursor++
				if vm.previewObject != nil {
					vm.updatePreview()
				}
			}
		case "enter", " ":
			rows := vm.visibleRows()
			if vm.cursor < len(rows) {
				row := rows[vm.cursor]
				if row.isHeader {
					vm.groups[row.groupIdx].Expanded = !vm.groups[row.groupIdx].Expanded
				} else if row.object != nil {
					vm.detailObject = row.object
				}
			}
		case "p":
			// Toggle preview panel (close editor if open)
			if vm.previewObject != nil {
				vm.previewObject = nil
			} else {
				vm.editor = nil
				vm.updatePreview()
			}
		case "e":
			// Open view editor (close preview)
			vm.previewObject = nil
			vm.editor = newViewEditor(vm.typeName, vm.viewName, vm.view, vm.schema, vm.vault)
		}
	}
	return vm, nil
}

// reloadObjects re-queries objects using the current view config.
func (vm *viewMode) reloadObjects() {
	vm.queryAndGroup()
	// Reset cursor if it's out of bounds
	rows := vm.visibleRows()
	if vm.cursor >= len(rows) {
		vm.cursor = max(0, len(rows)-1)
	}
}

// updatePreview sets the preview to the object under the cursor.
// Uses the in-memory object from the query result (no disk I/O).
func (vm *viewMode) updatePreview() {
	rows := vm.visibleRows()
	if vm.cursor >= 0 && vm.cursor < len(rows) && rows[vm.cursor].object != nil {
		vm.previewObject = rows[vm.cursor].object
	} else {
		vm.previewObject = nil
	}
}

// viewColumns returns the property names to display as columns/inline values.
// If ViewConfig.Columns is set, uses those exactly. Otherwise:
// - Table layout: all schema properties (pinned first, then unpinned)
// - List layout: empty (name only)
func (vm *viewMode) viewColumns() []string {
	// Explicit columns take priority
	if vm.view != nil && len(vm.view.Columns) > 0 {
		return vm.view.Columns
	}

	// List layout defaults to no columns (name only)
	if vm.view == nil || vm.view.Layout != core.ViewLayoutTable {
		return nil
	}

	// Table layout defaults to all schema properties
	if vm.schema == nil {
		return nil
	}
	var cols []string
	// Pinned properties first (sorted by pin value)
	type pinEntry struct {
		name string
		pin  int
	}
	var pinned []pinEntry
	var unpinned []string
	for _, p := range vm.schema.Properties {
		if p.Pin > 0 {
			pinned = append(pinned, pinEntry{p.Name, p.Pin})
		} else {
			unpinned = append(unpinned, p.Name)
		}
	}
	sort.Slice(pinned, func(i, j int) bool {
		return pinned[i].pin < pinned[j].pin
	})
	for _, p := range pinned {
		cols = append(cols, p.name)
	}
	cols = append(cols, unpinned...)

	// Limit columns based on available width (rough: 15 chars per col)
	maxCols := (vm.width - 22) / 14
	if maxCols < 1 {
		maxCols = 1
	}
	if len(cols) > maxCols {
		cols = cols[:maxCols]
	}
	return cols
}

// truncate shortens a string to fit within maxLen cells, adding ellipsis if needed.
func truncate(s string, maxLen int) string {
	return runewidth.Truncate(s, maxLen, "…")
}

// padRight pads a string with spaces to fill exactly width display cells.
// Delegates to runewidth.FillRight for CJK/emoji-aware padding.
func padRight(s string, width int) string {
	return runewidth.FillRight(s, width)
}

// displayPropValue formats a property value for table display using the
// unified DisplayProperty.FormatValue() pipeline.
func (vm *viewMode) displayPropValue(obj *core.Object, propName string) string {
	val, ok := obj.Properties[propName]
	if !ok || val == nil {
		return ""
	}
	dp := core.DisplayProperty{
		Key:   propName,
		Value: val,
	}
	if vm.schema != nil {
		if p := vm.schema.FindProperty(propName); p != nil {
			dp.Type = p.Type
			dp.IsRelation = p.Type == "relation"
		}
	}
	return dp.FormatValue()
}

// View renders the view mode content based on the configured layout.
func (vm *viewMode) View() string {
	rows := vm.visibleRows()
	if len(rows) == 0 {
		return "  (no objects)"
	}

	switch vm.view.Layout {
	case core.ViewLayoutTable:
		return vm.viewTable(rows)
	default:
		return vm.viewList(rows)
	}
}

// viewList renders objects as a simple list: emoji + name + inline column values.
func (vm *viewMode) viewList(rows []viewRow) string {
	cols := vm.viewColumns()

	// Ensure scroll keeps cursor visible
	visibleH := vm.height - 2 // padding
	if visibleH < 1 {
		visibleH = 1
	}
	vm.scroll = widget.AdjustScroll(vm.cursor, vm.scroll, visibleH)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("6"))
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	emoji := ""
	if vm.schema != nil && vm.schema.Emoji != "" {
		emoji = vm.schema.Emoji + " "
	}

	var b strings.Builder
	for i := vm.scroll; i < len(rows) && i < vm.scroll+visibleH; i++ {
		row := rows[i]
		isCurrent := i == vm.cursor

		if row.isHeader {
			headerText := fmt.Sprintf(" ── %s ──", row.label)
			b.WriteString(dimStyle.Render(headerText) + "\n")
		} else if row.object != nil {
			name := row.object.GetName()
			if name == "" {
				name = row.object.Filename
			}

			line := "  " + emoji + name

			// Append inline column values
			for _, col := range cols {
				val := vm.displayPropValue(row.object, col)
				if val != "" {
					line += " · " + val
				}
			}

			if isCurrent {
				b.WriteString(highlightStyle.Render(padRight(line, vm.width)) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
	}

	return b.String()
}

// viewTable renders objects in a columnar table with headers and separators.
func (vm *viewMode) viewTable(rows []viewRow) string {
	cols := vm.viewColumns()
	nameW := 20
	colW := 12

	// Adjust name width based on actual display width
	for _, row := range rows {
		if row.object != nil {
			name := row.object.GetName()
			w := runewidth.StringWidth(name)
			if w > nameW {
				nameW = w
			}
		}
	}
	if nameW > 30 {
		nameW = 30
	}

	// Trim columns to fit within available width
	// Each row: 2 (indent) + nameW + len(cols) * (2 + colW)
	maxCols := (vm.width - 2 - nameW) / (2 + colW)
	if maxCols < 0 {
		maxCols = 0
	}
	if len(cols) > maxCols {
		cols = cols[:maxCols]
	}

	// Ensure scroll keeps cursor visible
	visibleH := vm.height - 4 // header + separator + padding
	if visibleH < 1 {
		visibleH = 1
	}
	vm.scroll = widget.AdjustScroll(vm.cursor, vm.scroll, visibleH)

	highlightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("6"))
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("4")).
		Bold(true)

	// Build sort indicator map
	sortIndicators := make(map[string]string)
	if vm.view != nil {
		for _, s := range vm.view.Sort {
			if s.Direction == "desc" {
				sortIndicators[s.Property] = " ↓"
			} else {
				sortIndicators[s.Property] = " ↑"
			}
		}
	}

	var b strings.Builder

	// Column header
	header := "  " + padRight("NAME", nameW)
	for _, col := range cols {
		label := strings.ToUpper(col) + sortIndicators[col]
		header += "  " + padRight(truncate(label, colW), colW)
	}
	b.WriteString(headerStyle.Render(header) + "\n")
	// Separator
	sep := "  " + strings.Repeat("─", nameW)
	for range cols {
		sep += "──" + strings.Repeat("─", colW)
	}
	b.WriteString(dimStyle.Render(sep) + "\n")

	// Rows
	for i := vm.scroll; i < len(rows) && i < vm.scroll+visibleH; i++ {
		row := rows[i]
		isCurrent := i == vm.cursor

		if row.isHeader {
			headerText := fmt.Sprintf(" ── %s ──", row.label)
			b.WriteString(dimStyle.Render(headerText) + "\n")
		} else if row.object != nil {
			name := row.object.GetName()
			if name == "" {
				name = row.object.Filename
			}

			line := "  " + padRight(truncate(name, nameW), nameW)
			for _, col := range cols {
				val := vm.displayPropValue(row.object, col)
				if val == "" {
					cell := padRight("·", colW)
					if !isCurrent {
						cell = dimStyle.Render(cell)
					}
					line += "  " + cell
				} else {
					line += "  " + padRight(truncate(val, colW), colW)
				}
			}
			if isCurrent {
				b.WriteString(highlightStyle.Render(padRight(line, vm.width)) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
	}

	return b.String()
}

// titleContent returns the title bar text.
func (vm *viewMode) titleContent() string {
	emoji := ""
	typeName := vm.typeName
	if vm.schema != nil && vm.schema.Emoji != "" {
		emoji = vm.schema.Emoji + " "
	}
	return fmt.Sprintf("%s%s · %s (%d)", emoji, typeName, vm.viewName, len(vm.objects))
}

// PreviewContent returns the body text of the preview object, or empty if no preview.
func (vm *viewMode) PreviewContent() string {
	if vm.previewObject == nil {
		return ""
	}
	name := vm.previewObject.GetName()
	if name == "" {
		name = vm.previewObject.Filename
	}

	// Let lipgloss handle name wrapping; only constrain separator width
	maxW := vm.previewWidth
	if maxW < 10 {
		maxW = 40
	}

	var b strings.Builder
	b.WriteString(" " + name + "\n")
	sepW := maxW - 2 // fill most of panel width
	if sepW < 4 {
		sepW = 4
	}
	b.WriteString(" " + strings.Repeat("─", sepW) + "\n")

	// Show key properties
	if vm.schema != nil {
		for _, p := range vm.schema.Properties {
			formatted := vm.displayPropValue(vm.previewObject, p.Name)
			if formatted != "" {
				label := p.Name
				if p.Emoji != "" {
					label = p.Emoji + " " + label
				}
				b.WriteString(fmt.Sprintf(" %s: %s\n", label, formatted))
			}
		}
	}

	body := strings.TrimSpace(vm.previewObject.Body)
	if body != "" {
		b.WriteString("\n")
		// Limit body preview lines
		lines := strings.Split(body, "\n")
		for i, line := range lines {
			if i >= 20 {
				b.WriteString(" …\n")
				break
			}
			b.WriteString(" " + line + "\n")
		}
	}

	return b.String()
}

// HasPreview returns true if the preview panel should be shown.
func (vm *viewMode) HasPreview() bool {
	return vm.previewObject != nil
}

// HelpBar returns the context-sensitive help text.
func (vm *viewMode) HelpBar() string {
	if vm.editor != nil {
		return vm.editor.HelpBar()
	}
	if vm.detailObject != nil {
		return "esc: back to list"
	}
	if vm.previewObject != nil {
		return "↑/↓: navigate  enter: open  p: close preview  e: edit view  esc: exit view"
	}
	return "↑/↓: navigate  enter: open  p: preview  e: edit view  esc: exit view"
}
