package cmd

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

// disambiguateItem represents a single candidate in the disambiguation picker.
type disambiguateItem struct {
	name string // human-readable name from the object
	id   string // full object ID
}

// disambiguatePicker is a Bubble Tea model for selecting from ambiguous matches.
type disambiguatePicker struct {
	items    []disambiguateItem
	cursor   int
	done     bool
	selected bool // true if user selected; false if cancelled
}

func newDisambiguatePicker(items []disambiguateItem) disambiguatePicker {
	return disambiguatePicker{items: items}
}

func (m disambiguatePicker) Init() tea.Cmd {
	return nil
}

func (m disambiguatePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.done = true
			m.selected = true
			return m, tea.Quit
		case "q", "esc":
			m.done = true
			m.selected = false
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m disambiguatePicker) View() tea.View {
	if m.done {
		return tea.NewView("")
	}

	var b strings.Builder
	b.WriteString("Multiple matches found:\n\n")

	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		if item.name != "" && item.name != item.id {
			fmt.Fprintf(&b, "%s%s\n    %s\n", cursor, item.name, item.id)
		} else {
			fmt.Fprintf(&b, "%s%s\n", cursor, item.id)
		}
	}

	b.WriteString("\n↑↓ move  enter select  esc cancel\n")
	return tea.NewView(b.String())
}

// selectedID returns the ID of the selected item, or "" if cancelled.
func (m disambiguatePicker) selectedID() string {
	if !m.selected || m.cursor >= len(m.items) {
		return ""
	}
	return m.items[m.cursor].id
}

// buildDisambiguateItems creates picker items from match IDs, looking up
// object names from the vault for display.
func buildDisambiguateItems(vault *core.Vault, matches []string) []disambiguateItem {
	items := make([]disambiguateItem, len(matches))
	for i, id := range matches {
		items[i] = disambiguateItem{id: id}
		if obj, err := vault.GetObject(id); err == nil {
			items[i].name = obj.GetName()
		} else {
			// Fall back to slug without ULID for a readable display name.
			if parsed, parseErr := core.ParseObjectID(id); parseErr == nil {
				items[i].name = core.StripULID(parsed.Filename)
			}
		}
	}
	return items
}
