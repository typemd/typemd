package tui

import (
	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// initSearchInput creates and configures the search text input.
func initSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search objects..."
	ti.CharLimit = 100
	return ti
}

// searchResultRows returns flat list rows from search results (no grouping).
func searchResultRows(results []*core.Object) []listRow {
	var rows []listRow
	for _, obj := range results {
		rows = append(rows, listRow{Kind: rowObject, GroupIndex: -1, Object: obj})
	}
	return rows
}

// updateSearch handles key events in search mode.
// Returns updated model and command.
func updateSearch(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searchMode = false
		m.searchInput.Reset()
		m.searchResults = nil
		m.cursor = 0
		return m, nil
	case "enter":
		keyword := m.searchInput.Value()
		if keyword == "" {
			m.searchMode = false
			m.searchInput.Reset()
			m.searchResults = nil
			return m, nil
		}
		results, err := m.vault.SearchObjects(keyword)
		if err != nil {
			m.searchResults = nil
		} else {
			m.searchResults = results
		}
		m.searchMode = false
		m.cursor = 0
		m.focus = focusLeft
		return m, nil
	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}
}
