package cmd

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

// starterItem represents a single starter type in the picker.
type starterItem struct {
	name        string
	emoji       string
	description string
	selected    bool
}

// starterPicker is a Bubble Tea model for selecting starter types during init.
type starterPicker struct {
	items    []starterItem
	cursor   int
	done     bool
	quitting bool
}

func newStarterPicker() starterPicker {
	starters := core.StarterTypes()
	items := make([]starterItem, len(starters))
	for i, st := range starters {
		items[i] = starterItem{
			name:        st.Name,
			emoji:       st.Emoji,
			description: st.Description,
			selected:    true, // all selected by default
		}
	}
	return starterPicker{items: items}
}

func (m starterPicker) Init() tea.Cmd {
	return nil
}

func (m starterPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "space":
			m.items[m.cursor].selected = !m.items[m.cursor].selected
		case "a":
			for i := range m.items {
				m.items[i].selected = true
			}
		case "n":
			for i := range m.items {
				m.items[i].selected = false
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		case "q", "esc":
			for i := range m.items {
				m.items[i].selected = false
			}
			m.done = true
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m starterPicker) View() tea.View {
	if m.done {
		return tea.NewView("")
	}

	var b strings.Builder
	b.WriteString("Add starter types to your vault:\n\n")

	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		check := "[ ]"
		if item.selected {
			check = "[✓]"
		}
		b.WriteString(fmt.Sprintf("%s%s %s %-8s %s\n", cursor, check, item.emoji, item.name, item.description))
	}

	b.WriteString("\n↑↓ move  space toggle  a all  n none  enter confirm  esc skip\n")
	return tea.NewView(b.String())
}

// selectedNames returns the names of all selected starter types.
func (m starterPicker) selectedNames() []string {
	var names []string
	for _, item := range m.items {
		if item.selected {
			names = append(names, item.name)
		}
	}
	return names
}

// selectedItems returns the selected starter items with full metadata.
func (m starterPicker) selectedItems() []starterItem {
	var items []starterItem
	for _, item := range m.items {
		if item.selected {
			items = append(items, item)
		}
	}
	return items
}
