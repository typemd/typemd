package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
)

// padEmoji strips the variation selector (U+FE0F) and pads the emoji
// to a consistent 2-cell display width for terminal alignment.
func padEmoji(emoji string) string {
	display := strings.ReplaceAll(emoji, "\uFE0F", "")
	ew := runewidth.StringWidth(display)
	if ew < 2 {
		return display + strings.Repeat(" ", 2-ew)
	}
	return display
}

type rowKind int

const (
	rowHeader    rowKind = iota // type group header
	rowObject                   // object item
	rowNewObject                // "+ New Object" action row (per group)
	rowNewType                  // "+ New Type" action row
)

// listRow represents one visible row in the left panel.
type listRow struct {
	Kind       rowKind
	GroupIndex int
	Object     *core.Object // nil for headers and newType
}

// buildGroups creates type groups from a flat list of objects, sorted by type name.
// If vault is provided, each group's Emoji is populated from the type schema.
func buildGroups(objects []*core.Object, vault *core.Vault) []typeGroup {
	groupMap := make(map[string][]*core.Object)
	for _, obj := range objects {
		groupMap[obj.Type] = append(groupMap[obj.Type], obj)
	}

	// Include all types from the vault, even those with no objects
	if vault != nil {
		for _, name := range vault.ListTypes() {
			if _, ok := groupMap[name]; !ok {
				groupMap[name] = nil
			}
		}
	}

	var groups []typeGroup
	for name, objs := range groupMap {
		var emoji string
		plural := name
		if vault != nil {
			if ts, err := vault.LoadType(name); err == nil {
				emoji = ts.Emoji
				plural = ts.PluralName()
			}
		}
		groups = append(groups, typeGroup{
			Name:     name,
			Plural:   plural,
			Emoji:    emoji,
			Objects:  objs,
			Expanded: false,
		})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}

// visibleRows returns the list of rows currently visible based on expand/collapse state.
// Includes a "+ New Type" row at the bottom.
func visibleRows(groups []typeGroup) []listRow {
	var rows []listRow
	for i, g := range groups {
		rows = append(rows, listRow{Kind: rowHeader, GroupIndex: i})
		if g.Expanded {
			for _, obj := range g.Objects {
				rows = append(rows, listRow{Kind: rowObject, GroupIndex: i, Object: obj})
			}
		}
	}
	rows = append(rows, listRow{Kind: rowNewType})
	return rows
}

// clampCursor ensures cursor stays within valid range.
func clampCursor(cursor, totalRows int) int {
	if cursor < 0 {
		return 0
	}
	if totalRows == 0 {
		return 0
	}
	if cursor >= totalRows {
		return totalRows - 1
	}
	return cursor
}

// adjustScrollOffset returns a new offset so that cursor is visible within viewHeight.
func adjustScrollOffset(cursor, offset, viewHeight int) int {
	if viewHeight <= 0 {
		return 0
	}
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+viewHeight {
		return cursor - viewHeight + 1
	}
	return offset
}

// renderList renders the left panel list with scroll offset support.
func renderList(groups []typeGroup, cursor, scrollOffset int, focused bool, width, height int) string {
	rows := visibleRows(groups)
	if len(rows) == 0 {
		return "  (no objects)"
	}

	end := scrollOffset + height
	if end > len(rows) {
		end = len(rows)
	}
	start := scrollOffset
	if start > len(rows) {
		start = len(rows)
	}

	var lines []string
	for i := start; i < end; i++ {
		row := rows[i]
		var line string
		switch row.Kind {
		case rowHeader:
			g := groups[row.GroupIndex]
			arrow := "▶"
			if g.Expanded {
				arrow = "▼"
			}
			if g.Emoji != "" {
				line = fmt.Sprintf(" %s %s %s (%d)", arrow, padEmoji(g.Emoji), g.Plural, len(g.Objects))
			} else {
				line = fmt.Sprintf(" %s %s (%d)", arrow, g.Plural, len(g.Objects))
			}
		case rowObject:
			line = fmt.Sprintf("   %s", row.Object.GetName())
		case rowNewType:
			line = " + New Type"
		}

		if i == cursor {
			style := highlightStyle.Width(width)
			line = style.Render(line)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
