package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	tea "charm.land/bubbletea/v2"
)

// relationCandidate represents a candidate object in the relation picker.
type relationCandidate struct {
	id          string // full object ID (e.g. "person/alan-01abc...")
	displayName string // human-readable name (e.g. "alan" or "person/alan")
	lowerName   string // pre-computed lowercase displayName for search filtering
	fullIndex   int    // index in the full (unfiltered) relCandidates slice
}

// relationLinkedMsg is sent after a successful LinkObjects call.
type relationLinkedMsg struct {
	fromID  string
	relName string
	toID    string
}

// relationUnlinkedMsg is sent after a successful UnlinkObjects call.
type relationUnlinkedMsg struct {
	fromID  string
	relName string
	toID    string
}

// relationLinkErrMsg is sent when a Link/Unlink operation fails.
type relationLinkErrMsg struct {
	err error
}

// activateRelationPicker loads candidates and opens the relation picker.
func (pe *propEditor) activateRelationPicker(item *propItem, vault *core.Vault) tea.Cmd {
	// Determine target type from schema or system property
	var targetType string
	var multiple bool

	if item.schema != nil && item.schema.Type == "relation" {
		targetType = item.schema.Target
		multiple = item.schema.Multiple
	} else if item.dp.Key == core.TagsProperty {
		targetType = core.TagTypeName
		multiple = true
	} else {
		// Try to find the relation property via schema
		return nil
	}

	// Query candidates
	var filters []core.FilterRule
	if targetType != "" {
		filters = append(filters, core.FilterRule{
			Property: "type",
			Operator: "is",
			Value:    targetType,
		})
	}

	results, err := vault.QueryObjects(filters, core.SortRule{Property: "name", Direction: "asc"})
	if err != nil {
		return nil
	}

	// Build candidate list
	candidates := make([]relationCandidate, len(results))
	for i, obj := range results {
		name := obj.DisplayName()
		display := name
		if targetType == "" {
			display = obj.Type + "/" + name
		}
		candidates[i] = relationCandidate{
			id:          obj.ID,
			displayName: display,
			lowerName:   strings.ToLower(display),
			fullIndex:   i,
		}
	}

	pe.editIndex = pe.cursor
	pe.relCandidates = candidates
	pe.relFiltered = candidates
	pe.relSearch = ""
	pe.pickerCursor = 0

	if multiple {
		pe.mode = propModeRelationMultiPick
		// Initialize checked state from current relation values
		currentIDs := currentRelationIDs(item.dp)
		pe.relChecked = make([]bool, len(candidates))
		for i, c := range candidates {
			for _, id := range currentIDs {
				if c.id == id {
					pe.relChecked[i] = true
					break
				}
			}
		}
		// Snapshot initial checked state for diff on confirm
		pe.relInitialChecked = make([]bool, len(pe.relChecked))
		copy(pe.relInitialChecked, pe.relChecked)
	} else {
		pe.mode = propModeRelationPick
		// Position cursor on current value
		currentID := currentSingleRelationID(item.dp)
		// Index 0 is "(none)", candidates start at 1 in the display
		pe.pickerCursor = 0 // default to "(none)"
		for i, c := range candidates {
			if c.id == currentID {
				pe.pickerCursor = i + 1 // +1 for "(none)" at top
				break
			}
		}
	}

	return nil
}

// currentRelationIDs extracts the current relation target IDs from a display property.
func currentRelationIDs(dp core.DisplayProperty) []string {
	switch v := dp.Value.(type) {
	case []any:
		ids := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				ids = append(ids, s)
			}
		}
		return ids
	case []string:
		return v
	case string:
		if v != "" {
			return []string{v}
		}
	}
	return nil
}

// currentSingleRelationID extracts the current single relation target ID.
func currentSingleRelationID(dp core.DisplayProperty) string {
	if dp.Value == nil {
		return ""
	}
	return fmt.Sprintf("%v", dp.Value)
}

// filterRelationCandidates filters candidates by search text (case-insensitive substring).
func (pe *propEditor) filterRelationCandidates() {
	if pe.relSearch == "" {
		pe.relFiltered = pe.relCandidates
		return
	}
	query := strings.ToLower(pe.relSearch)
	pe.relFiltered = nil
	for _, c := range pe.relCandidates {
		if strings.Contains(c.lowerName, query) {
			pe.relFiltered = append(pe.relFiltered, c)
		}
	}
	// Reset cursor if out of bounds
	maxCursor := len(pe.relFiltered)
	if pe.mode == propModeRelationPick {
		maxCursor++ // +1 for "(none)" option
	}
	if pe.pickerCursor >= maxCursor {
		pe.pickerCursor = 0
	}
}

// filteredToFullIndex returns the index in relCandidates for a given index in relFiltered.
func (pe *propEditor) filteredToFullIndex(filteredIdx int) int {
	if filteredIdx < 0 || filteredIdx >= len(pe.relFiltered) {
		return -1
	}
	return pe.relFiltered[filteredIdx].fullIndex
}

// updateRelationPick handles key events for single-value relation picker.
func updateRelationPick(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	maxIdx := len(pe.relFiltered) // 0 = "(none)", 1..N = candidates

	switch msg.String() {
	case "up", "k":
		if pe.pickerCursor > 0 {
			pe.pickerCursor--
		}
		m.updatePropsContent()
		return m, nil, true

	case "down", "j":
		if pe.pickerCursor < maxIdx {
			pe.pickerCursor++
		}
		m.updatePropsContent()
		return m, nil, true

	case "enter":
		item := &pe.items[pe.editIndex]
		oldID := currentSingleRelationID(item.dp)
		if pe.pickerCursor == 0 {
			// "(none)" selected — clear relation
			cmd := clearRelationCmd(m, item.dp)
			pe.cancelEdit()
			return m, cmd, true
		}
		// Select candidate (index - 1 because of "(none)")
		candidate := pe.relFiltered[pe.pickerCursor-1]
		cmd := linkRelationCmd(m, item.dp.Key, candidate.id, oldID)
		pe.cancelEdit()
		return m, cmd, true

	case "esc":
		pe.cancelEdit()
		m.updatePropsContent()
		return m, nil, true

	case "backspace", "ctrl+h":
		if len(pe.relSearch) > 0 {
			pe.relSearch = pe.relSearch[:len(pe.relSearch)-1]
			pe.filterRelationCandidates()
			m.updatePropsContent()
		}
		return m, nil, true

	default:
		// Append printable characters to search
		key := msg.String()
		if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			pe.relSearch += key
			pe.filterRelationCandidates()
			m.updatePropsContent()
			return m, nil, true
		}
	}

	return m, nil, true
}

// updateRelationMultiPick handles key events for multi-value relation picker.
func updateRelationMultiPick(m model, msg tea.KeyPressMsg) (model, tea.Cmd, bool) {
	pe := m.propEdit
	maxIdx := len(pe.relFiltered) - 1

	switch msg.String() {
	case "up", "k":
		if pe.pickerCursor > 0 {
			pe.pickerCursor--
		}
		m.updatePropsContent()
		return m, nil, true

	case "down", "j":
		if pe.pickerCursor < maxIdx {
			pe.pickerCursor++
		}
		m.updatePropsContent()
		return m, nil, true

	case " ", "space":
		if pe.pickerCursor >= 0 && pe.pickerCursor < len(pe.relFiltered) {
			fullIdx := pe.filteredToFullIndex(pe.pickerCursor)
			if fullIdx >= 0 {
				pe.relChecked[fullIdx] = !pe.relChecked[fullIdx]
			}
			m.updatePropsContent()
		}
		return m, nil, true

	case "enter":
		// Compute diff: compare initial vs current checked state
		item := &pe.items[pe.editIndex]
		cmd := applyRelationMultiDiff(m, item.dp.Key, pe)
		pe.cancelEdit()
		return m, cmd, true

	case "esc":
		pe.cancelEdit()
		m.updatePropsContent()
		return m, nil, true

	case "backspace", "ctrl+h":
		if len(pe.relSearch) > 0 {
			pe.relSearch = pe.relSearch[:len(pe.relSearch)-1]
			pe.filterRelationCandidates()
			m.updatePropsContent()
		}
		return m, nil, true

	default:
		key := msg.String()
		if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			pe.relSearch += key
			pe.filterRelationCandidates()
			m.updatePropsContent()
			return m, nil, true
		}
	}

	return m, nil, true
}

// linkRelationCmd creates a tea.Cmd that links a relation.
// For single-value relations, it first unlinks the old value if present.
func linkRelationCmd(m model, relName, toID, oldID string) tea.Cmd {
	vault := m.vault
	fromID := m.selected.ID

	return func() tea.Msg {
		if oldID != "" && oldID != toID {
			_ = vault.UnlinkObjects(fromID, relName, oldID, false)
		}

		if err := vault.LinkObjects(fromID, relName, toID); err != nil {
			return relationLinkErrMsg{err: err}
		}
		return relationLinkedMsg{fromID: fromID, relName: relName, toID: toID}
	}
}

// clearRelationCmd creates a tea.Cmd that clears a single-value relation.
func clearRelationCmd(m model, dp core.DisplayProperty) tea.Cmd {
	vault := m.vault
	fromID := m.selected.ID
	relName := dp.Key

	oldID := currentSingleRelationID(dp)
	if oldID == "" {
		return nil // nothing to clear
	}

	return func() tea.Msg {
		if err := vault.UnlinkObjects(fromID, relName, oldID, false); err != nil {
			return relationLinkErrMsg{err: err}
		}
		return relationUnlinkedMsg{fromID: fromID, relName: relName, toID: oldID}
	}
}

// applyRelationMultiDiff computes link/unlink diff and returns a batched command.
func applyRelationMultiDiff(m model, relName string, pe *propEditor) tea.Cmd {
	vault := m.vault
	fromID := m.selected.ID

	var toLink, toUnlink []string
	for i, c := range pe.relCandidates {
		wasChecked := pe.relInitialChecked[i]
		isChecked := pe.relChecked[i]
		if !wasChecked && isChecked {
			toLink = append(toLink, c.id)
		} else if wasChecked && !isChecked {
			toUnlink = append(toUnlink, c.id)
		}
	}

	if len(toLink) == 0 && len(toUnlink) == 0 {
		return nil // no changes
	}

	return func() tea.Msg {
		var errCount int
		var lastErr error
		for _, id := range toUnlink {
			if err := vault.UnlinkObjects(fromID, relName, id, false); err != nil {
				lastErr = err
				errCount++
			}
		}
		for _, id := range toLink {
			if err := vault.LinkObjects(fromID, relName, id); err != nil {
				lastErr = err
				errCount++
			}
		}
		if lastErr != nil {
			if errCount > 1 {
				return relationLinkErrMsg{err: fmt.Errorf("%d operations failed, last: %w", errCount, lastErr)}
			}
			return relationLinkErrMsg{err: lastErr}
		}
		return relationLinkedMsg{fromID: fromID, relName: relName}
	}
}

// handleRelationResult handles the result of a link/unlink operation.
func handleRelationResult(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case relationLinkedMsg, relationUnlinkedMsg:
		// Reload object from disk to get updated properties
		if m.selected != nil && m.vault != nil {
			if obj, err := m.vault.GetObject(m.selected.ID); err == nil {
				m.applyLoadedObject(obj)
				m.updateDetail()
			}
		}
		m.skipNextReload = true
		return m, nil

	case relationLinkErrMsg:
		cmd := m.toast.Show(widget.ToastError, []widget.ToastItem{
			{Message: fmt.Sprintf("Relation error: %v", msg.err)},
		})
		return m, cmd
	}
	return m, nil
}

// renderRelationPicker renders the relation picker inline within the properties panel.
func (pe *propEditor) renderRelationPicker(item propItem, isCursor bool) string {
	var b strings.Builder

	// Property label
	prefix := " "
	if isCursor {
		prefix = "▸"
	}
	b.WriteString(fmt.Sprintf("%s %s:\n", prefix, item.dp.Key))

	// Search input
	searchDisplay := pe.relSearch
	if searchDisplay == "" {
		searchDisplay = dimStyle.Render("type to filter...")
	}
	b.WriteString(fmt.Sprintf("   🔍 %s\n", searchDisplay))

	if pe.mode == propModeRelationPick {
		// Single-select: "(none)" option at top
		if pe.pickerCursor == 0 {
			b.WriteString(highlightStyle.Render("   ▸ (none)") + "\n")
		} else {
			b.WriteString("     (none)\n")
		}

		for i, c := range pe.relFiltered {
			displayIdx := i + 1 // +1 for "(none)"
			if displayIdx == pe.pickerCursor {
				b.WriteString(highlightStyle.Render(fmt.Sprintf("   ▸ %s", c.displayName)) + "\n")
			} else {
				b.WriteString(fmt.Sprintf("     %s\n", c.displayName))
			}
		}
	} else {
		// Multi-select with checkmarks
		for i, c := range pe.relFiltered {
			fullIdx := pe.filteredToFullIndex(i)
			check := "☐"
			if fullIdx >= 0 && pe.relChecked[fullIdx] {
				check = "☑"
			}
			if i == pe.pickerCursor {
				b.WriteString(highlightStyle.Render(fmt.Sprintf("   %s %s", check, c.displayName)) + "\n")
			} else {
				b.WriteString(fmt.Sprintf("   %s %s\n", check, c.displayName))
			}
		}
	}

	return b.String()
}

// isPicking returns true if the editor is in a picker mode (select, multi, relation).
func (pe *propEditor) isPicking() bool {
	switch pe.mode {
	case propModeSelectPick, propModeMultiPick, propModeRelationPick, propModeRelationMultiPick:
		return true
	}
	return false
}
