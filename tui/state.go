package tui

import (
	"os"
	"path/filepath"
	"strings"

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
		ScrollOffset:  m.scrollOffset,
		Focus:         focusPanelToString(m.focus),
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
