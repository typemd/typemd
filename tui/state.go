package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/typemd/typemd/core"
	"gopkg.in/yaml.v3"
)

// SessionState represents the TUI session state that is persisted across restarts.
type SessionState struct {
	SelectedObjectID string   `yaml:"selected_object_id,omitempty"`
	SelectedTypeName string   `yaml:"selected_type_name,omitempty"` // type header cursor was on
	ExpandedGroups   []string `yaml:"expanded_groups,omitempty"`
	ScrollOffset     int      `yaml:"scroll_offset,omitempty"`
	Focus            string   `yaml:"focus,omitempty"`
	LeftPanelWidth   int      `yaml:"left_panel_width,omitempty"`
	PropsPanelWidth  int      `yaml:"props_panel_width,omitempty"`
	PropsVisible     bool     `yaml:"props_visible"`

	// View mode state (present only when TUI was in view mode on exit)
	ViewTypeName      string   `yaml:"view_type_name,omitempty"`
	ViewName          string   `yaml:"view_name,omitempty"`
	ViewCursor        int      `yaml:"view_cursor,omitempty"`
	ViewScroll        int      `yaml:"view_scroll,omitempty"`
	ViewExpandedGroups []string `yaml:"view_expanded_groups,omitempty"`

	// Stats mode state (present only when TUI was in stats mode on exit).
	// StatsActive is the canonical marker of "TUI was in stats mode on exit";
	// the cursor/scroll fields default to 0 and would otherwise be stripped by
	// `omitempty`, making the presence-of-stats-mode signal ambiguous.
	StatsActive   bool   `yaml:"stats_active,omitempty"`
	StatsCursor   int    `yaml:"stats_cursor,omitempty"`
	StatsScroll   int    `yaml:"stats_scroll,omitempty"`
	StatsTypeName string `yaml:"stats_type_name,omitempty"`
}

const stateFileName = "tui-state.yaml"

// stateFilePath returns the path to the TUI state file for the given vault root.
func stateFilePath(vaultRoot string) string {
	return filepath.Join(vaultRoot, ".typemd", stateFileName)
}

// loadSessionState reads the TUI session state from disk.
// Returns a zero-value SessionState if the file is missing, unreadable, or invalid.
func loadSessionState(vaultRoot string) SessionState {
	data, err := os.ReadFile(stateFilePath(vaultRoot))
	if err != nil {
		return SessionState{}
	}
	var state SessionState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return SessionState{}
	}
	return state
}

// saveSessionState writes the TUI session state to disk.
// Errors are silently ignored — state persistence is best-effort.
func saveSessionState(vaultRoot string, state SessionState) {
	data, err := yaml.Marshal(&state)
	if err != nil {
		return
	}
	_ = os.WriteFile(stateFilePath(vaultRoot), data, 0644)
}

// captureState extracts the current TUI state into a SessionState for persistence.
func (m model) captureState() SessionState {
	state := SessionState{
		ScrollOffset:    m.scrollOffset,
		Focus:           focusPanelToString(m.focus),
		LeftPanelWidth:  m.leftW,
		PropsPanelWidth: m.propsWidth,
		PropsVisible:    m.propsVisible,
	}

	if m.selected != nil {
		state.SelectedObjectID = m.selected.ID
	} else if m.typeEditor != nil {
		state.SelectedTypeName = m.typeEditor.typeName
	}

	for _, g := range m.groups {
		if g.Expanded {
			state.ExpandedGroups = append(state.ExpandedGroups, g.Name)
		}
	}

	// Capture view mode state when active
	if m.rightPanel == panelView && m.viewMode != nil {
		state.ViewTypeName = m.viewMode.typeName
		state.ViewName = m.viewMode.viewName
		state.ViewCursor = m.viewMode.cursor
		state.ViewScroll = m.viewMode.scroll
		state.ViewExpandedGroups = m.viewMode.expandedGroupLabels()
	}

	// Capture stats mode state when active
	if m.rightPanel == panelStats && m.statsMode != nil {
		state.StatsActive = true
		state.StatsCursor = m.statsMode.cursor
		state.StatsScroll = m.statsMode.scroll
		if m.statsMode.screen == statsDetail {
			state.StatsTypeName = m.statsMode.detailType
		}
	}

	return state
}

// applySessionState applies a saved session state to groups (mutating their
// Expanded fields) and returns the resolved cursor position and selected object ID.
// It handles fallback logic when saved state references objects or types that
// no longer exist.
func applySessionState(state SessionState, groups []typeGroup) (cursor int, selectedID string) {
	// Apply expanded groups
	expandedSet := make(map[string]bool, len(state.ExpandedGroups))
	for _, name := range state.ExpandedGroups {
		expandedSet[name] = true
	}

	hasExpanded := false
	for i := range groups {
		if expandedSet[groups[i].Name] {
			groups[i].Expanded = true
			hasExpanded = true
		}
	}

	// If no saved groups matched, fall back to expanding first group
	if !hasExpanded && len(groups) > 0 {
		groups[0].Expanded = true
	}

	// Find the saved type header by name
	if state.SelectedTypeName != "" && state.SelectedObjectID == "" {
		rows := visibleRows(groups)
		for i, row := range rows {
			if row.Kind == rowHeader && row.GroupIndex < len(groups) && groups[row.GroupIndex].Name == state.SelectedTypeName {
				return i, ""
			}
		}
	}

	// Find the saved object by ID
	if state.SelectedObjectID != "" {
		rows := visibleRows(groups)
		for i, row := range rows {
			if row.Kind == rowObject && row.Object != nil && row.Object.ID == state.SelectedObjectID {
				return i, row.Object.ID
			}
		}

		// Object not found — try fallback to same type group
		if objType, _, ok := strings.Cut(state.SelectedObjectID, "/"); ok {
			// Ensure the type group is expanded so we can find objects in it
			for i := range groups {
				if groups[i].Name == objType && !groups[i].Expanded {
					groups[i].Expanded = true
				}
			}
			rows = visibleRows(groups)
			for i, row := range rows {
				if row.Kind == rowObject && row.Object != nil && row.Object.Type == objType {
					return i, row.Object.ID
				}
			}
		}
	}

	// Final fallback: first object in first expanded group
	rows := visibleRows(groups)
	for i, row := range rows {
		if row.Kind == rowObject && row.Object != nil {
			return i, row.Object.ID
		}
	}

	return 0, ""
}

// restoreViewMode attempts to restore view mode from saved session state.
// Returns a non-nil *viewMode if restoration succeeds, nil otherwise (fallback to sidebar).
//
// Fallback chain (tui-session-state spec R10):
//   - Saved viewName resolves → use it with saved cursor/scroll/expanded groups.
//   - Saved viewName missing but "default" exists → fall back to default view, reset cursor/scroll.
//   - Neither saved nor default resolve → return nil (sidebar mode).
func restoreViewMode(state SessionState, v *core.Vault) *viewMode {
	if state.ViewTypeName == "" || state.ViewName == "" {
		return nil
	}

	// Check type exists before creating view mode
	if _, err := v.LoadType(state.ViewTypeName); err != nil {
		return nil
	}

	// When the saved view is gone, fall back to "default" so the user lands
	// on a valid view rather than seeing the stale name in the title
	// (tui-session-state spec R10 S1). The implicit default view is always
	// available, so no further fallback is needed. Also skip the cursor /
	// scroll / expanded-groups restore below — they reference the old view
	// and would be misleading when applied against the default.
	_, saveErr := v.LoadView(state.ViewTypeName, state.ViewName)
	if saveErr != nil && state.ViewName != "default" {
		return newViewMode(state.ViewTypeName, "default", v)
	}

	vm := newViewMode(state.ViewTypeName, state.ViewName, v)

	// Apply expanded groups from state
	if len(state.ViewExpandedGroups) > 0 {
		expandedSet := make(map[string]bool, len(state.ViewExpandedGroups))
		for _, label := range state.ViewExpandedGroups {
			expandedSet[label] = true
		}
		for i := range vm.groups {
			if vm.groups[i].Label != "" {
				vm.groups[i].Expanded = expandedSet[vm.groups[i].Label]
			}
		}
	}

	// Clamp cursor and scroll to valid range
	totalRows := len(vm.visibleRows())
	if totalRows > 0 {
		vm.cursor = min(state.ViewCursor, totalRows-1)
		vm.scroll = min(state.ViewScroll, totalRows-1)
	}

	return vm
}

// restoreStatsMode attempts to restore stats mode from saved session state.
// Returns a non-nil *statsMode if restoration succeeds, nil otherwise.
// Stats state takes precedence over view state when both are present.
func restoreStatsMode(state SessionState, v *core.Vault) *statsMode {
	// Presence of stats_active is the canonical marker. For backwards
	// compatibility with state files written before stats_active existed,
	// still restore when any of the other stats_* fields are non-zero.
	hasStatsState := state.StatsActive ||
		state.StatsTypeName != "" ||
		state.StatsCursor > 0 ||
		state.StatsScroll > 0

	if !hasStatsState {
		return nil
	}

	layout := ""
	if cfg := v.Config(); cfg != nil {
		layout = cfg.TUI.StatsTypeLayout
	}

	sm := newStatsMode(v, layout)

	// Clamp cursor to valid range
	if sm.typeCount() > 0 {
		sm.cursor = min(state.StatsCursor, sm.typeCount()-1)
		sm.scroll = min(state.StatsScroll, sm.typeCount()-1)
	}

	// Restore type detail if saved
	if state.StatsTypeName != "" {
		// Verify the type still exists
		if _, err := v.LoadType(state.StatsTypeName); err == nil {
			sm.loadTypeStats(state.StatsTypeName)
			sm.screen = statsDetail
		}
		// If type no longer exists, stay on vault overview (fallback)
	}

	return sm
}

// focusPanelToString converts a focusPanel value to its string representation.
func focusPanelToString(f focusPanel) string {
	switch f {
	case focusBody:
		return "body"
	case focusProps:
		return "props"
	default:
		return "left"
	}
}

// stringToFocusPanel converts a string to a focusPanel value.
func stringToFocusPanel(s string) focusPanel {
	switch s {
	case "body":
		return focusBody
	case "props":
		return focusProps
	default:
		return focusLeft
	}
}
