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

	width  int
	height int
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

	// Query objects using vault facade, combining type filter with view filter rules
	filter := append(core.TypeFilter(vm.typeName), vm.view.Filter...)
	objects, err := vm.vault.QueryObjects(filter, vm.view.Sort...)
	if err != nil {
		vm.objects = nil
	} else {
		vm.objects = objects
	}

	// Build groups
	vm.buildGroups()
}

// buildGroups organizes objects by group_by property or as a flat list.
func (vm *viewMode) buildGroups() {
	if vm.view.GroupBy == "" {
		// Flat list — single unnamed group
		vm.groups = []viewGroup{{
			Label:    "",
			Objects:  vm.objects,
			Expanded: true,
		}}
		return
	}

	// Group by property value
	groupMap := make(map[string][]*core.Object)
	var groupOrder []string
	for _, obj := range vm.objects {
		val := ""
		if v, ok := obj.Properties[vm.view.GroupBy]; ok && v != nil {
			val = fmt.Sprintf("%v", v)
		}
		if val == "" {
			val = "(none)"
		}
		if _, exists := groupMap[val]; !exists {
			groupOrder = append(groupOrder, val)
		}
		groupMap[val] = append(groupMap[val], obj)
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

// visibleRows returns the flattened list of displayable rows.
func (vm *viewMode) visibleRows() []viewRow {
	var rows []viewRow
	for gi, g := range vm.groups {
		if vm.view.GroupBy != "" {
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
	return vm.detailObject == nil
}

// Update handles messages for the view mode.
func (vm *viewMode) Update(msg tea.Msg) (*viewMode, tea.Cmd) {
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
			// Toggle preview panel
			if vm.previewObject != nil {
				vm.previewObject = nil
			} else {
				vm.updatePreview()
			}
		}
	}
	return vm, nil
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

// viewColumns returns the property names to display as columns.
// Uses pinned properties first, then unpinned, up to a reasonable limit.
func (vm *viewMode) viewColumns() []string {
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
	// Sort pinned by pin value
	sort.Slice(pinned, func(i, j int) bool {
		return pinned[i].pin < pinned[j].pin
	})
	for _, p := range pinned {
		cols = append(cols, p.name)
	}
	cols = append(cols, unpinned...)

	// Limit columns based on available width (rough: 15 chars per col)
	maxCols := (vm.width - 22) / 14 // 22 for name column + padding
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

// formatPropValue formats a property value for table display.
func formatPropValue(obj *core.Object, propName string) string {
	val, ok := obj.Properties[propName]
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case bool:
		if v {
			return "✓"
		}
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// View renders the full-width view list with property columns.
func (vm *viewMode) View() string {
	rows := vm.visibleRows()
	if len(rows) == 0 {
		return "  (no objects)"
	}

	cols := vm.viewColumns()
	nameW := 20
	colW := 12

	// Adjust name width based on actual content
	for _, row := range rows {
		if row.object != nil {
			name := row.object.GetName()
			if len([]rune(name)) > nameW {
				nameW = len([]rune(name))
			}
		}
	}
	if nameW > 30 {
		nameW = 30
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

	var b strings.Builder

	// Column header
	header := fmt.Sprintf("  %-*s", nameW, "NAME")
	for _, col := range cols {
		header += fmt.Sprintf("  %-*s", colW, strings.ToUpper(truncate(col, colW)))
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

			line := fmt.Sprintf("  %-*s", nameW, truncate(name, nameW))
			for _, col := range cols {
				val := formatPropValue(row.object, col)
				if val == "" {
					line += "  " + dimStyle.Render(fmt.Sprintf("%-*s", colW, "·"))
				} else {
					line += fmt.Sprintf("  %-*s", colW, truncate(val, colW))
				}
			}

			if isCurrent {
				b.WriteString(highlightStyle.Render(line) + "\n")
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

	var b strings.Builder
	b.WriteString(" " + name + "\n")
	b.WriteString(" " + strings.Repeat("─", len([]rune(name))+2) + "\n")

	// Show key properties
	if vm.schema != nil {
		for _, p := range vm.schema.Properties {
			val := formatPropValue(vm.previewObject, p.Name)
			if val != "" {
				label := p.Name
				if p.Emoji != "" {
					label = p.Emoji + " " + label
				}
				b.WriteString(fmt.Sprintf(" %s: %s\n", label, val))
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
	if vm.detailObject != nil {
		return "esc: back to list"
	}
	if vm.previewObject != nil {
		return "↑/↓: navigate  enter: open  p: close preview  esc: exit view"
	}
	return "↑/↓: navigate  enter: open  p: preview  esc: exit view"
}
