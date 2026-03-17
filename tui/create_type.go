package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// createTypeField represents which field is focused in the type creation form.
type createTypeField int

const (
	createTypeFieldEmoji  createTypeField = iota // emoji text input
	createTypeFieldName                          // name text input (default focus)
	createTypeFieldPlural                        // plural text input
)

// createTypeState holds all state for the type creation flow.
type createTypeState struct {
	field      createTypeField
	emojiInput textinput.Model
	nameInput  textinput.Model
	pluralInput textinput.Model
	errMsg     string
}

// initCreateTypeInput creates a configured text input for type creation.
func initCreateTypeInput(placeholder string, charLimit int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = charLimit
	return ti
}

// startCreateType initializes the type creation flow.
func (m *model) startCreateType() tea.Cmd {
	if m.readOnly {
		return nil
	}

	emojiInput := initCreateTypeInput("emoji", 4)
	nameInput := initCreateTypeInput("type name", 50)
	pluralInput := initCreateTypeInput("plural", 50)

	nameInput.Focus()

	m.createType = &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  emojiInput,
		nameInput:   nameInput,
		pluralInput: pluralInput,
	}
	return textinput.Blink
}

// nextCreateTypeField cycles to the next field.
func (cts *createTypeState) nextCreateTypeField() {
	switch cts.field {
	case createTypeFieldEmoji:
		cts.emojiInput.Blur()
		cts.field = createTypeFieldName
		cts.nameInput.Focus()
	case createTypeFieldName:
		cts.nameInput.Blur()
		cts.field = createTypeFieldPlural
		cts.pluralInput.Focus()
	case createTypeFieldPlural:
		cts.pluralInput.Blur()
		cts.field = createTypeFieldEmoji
		cts.emojiInput.Focus()
	}
}

// updateCreateType handles key events during type creation.
func updateCreateType(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	cts := m.createType
	if cts == nil {
		return m, nil
	}

	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(cts.nameInput.Value())
		if name == "" {
			return m, nil
		}

		// Check if type already exists
		if m.vault != nil {
			existing := m.vault.ListTypes()
			for _, t := range existing {
				if t == name {
					cts.errMsg = "type \"" + name + "\" already exists"
					return m, nil
				}
			}
		}

		// Create type schema with provided fields
		schema := &core.TypeSchema{
			Name:   name,
			Emoji:  strings.TrimSpace(cts.emojiInput.Value()),
			Plural: strings.TrimSpace(cts.pluralInput.Value()),
		}
		if err := m.vault.SaveType(schema); err != nil {
			cts.errMsg = err.Error()
			return m, nil
		}

		m.createType = nil

		// Rebuild groups directly (avoid refreshData which calls selectCurrentRow)
		m.rebuildGroups()

		// Open type editor for new type and move cursor to its header
		if ts, err := m.vault.LoadType(name); err == nil {
			m.typeEditor = newTypeEditor(ts, name, true, m.vault)
			m.rightPanel = panelTypeEditor
			m.selected = nil
			m.focus = focusBody
			m.moveCursorToType(name)
		}
		return m, nil

	case "esc":
		m.createType = nil
		return m, nil

	case "tab":
		cts.nextCreateTypeField()
		return m, textinput.Blink
	}

	// Clear error on any input
	cts.errMsg = ""

	// Route to the focused field's text input
	var cmd tea.Cmd
	switch cts.field {
	case createTypeFieldEmoji:
		cts.emojiInput, cmd = cts.emojiInput.Update(msg)
	case createTypeFieldName:
		cts.nameInput, cmd = cts.nameInput.Update(msg)
	case createTypeFieldPlural:
		cts.pluralInput, cmd = cts.pluralInput.Update(msg)
	}
	return m, cmd
}

// syncTypeGroupMeta updates a type group's display metadata from a schema.
func (m *model) syncTypeGroupMeta(typeName string, schema *core.TypeSchema) {
	if schema == nil {
		return
	}
	for i := range m.groups {
		if m.groups[i].Name == typeName {
			m.groups[i].Emoji = schema.Emoji
			m.groups[i].Plural = schema.PluralName()
			break
		}
	}
}

// moveCursorToType moves the sidebar cursor to the header row for the given type name.
func (m *model) moveCursorToType(typeName string) {
	rows := m.currentRows()
	for i, row := range rows {
		if row.Kind == rowHeader && row.GroupIndex < len(m.groups) && m.groups[row.GroupIndex].Name == typeName {
			m.cursor = i
			m.adjustScroll()
			break
		}
	}
}

// renderCreateTypeTitleContent renders the title panel content during type creation.
func renderCreateTypeTitleContent(cts *createTypeState) string {
	if cts == nil {
		return ""
	}

	var b strings.Builder

	// Emoji field
	emoji := strings.TrimSpace(cts.emojiInput.Value())
	if cts.field == createTypeFieldEmoji {
		b.WriteString(fmt.Sprintf(" [%s]", cts.emojiInput.View()))
	} else if emoji != "" {
		b.WriteString(fmt.Sprintf(" %s", padEmoji(emoji)))
	} else {
		b.WriteString(" ")
	}

	b.WriteString(" new type · ")

	// Name field
	if cts.field == createTypeFieldName {
		b.WriteString(cts.nameInput.View())
	} else {
		name := cts.nameInput.Value()
		if name == "" {
			name = "..."
		}
		b.WriteString(name)
	}

	// Plural field
	b.WriteString("  plural: ")
	if cts.field == createTypeFieldPlural {
		b.WriteString(cts.pluralInput.View())
	} else {
		plural := cts.pluralInput.Value()
		if plural == "" {
			plural = "..."
		}
		b.WriteString(plural)
	}

	// Error message
	if cts.errMsg != "" {
		b.WriteString(fmt.Sprintf("  ✗ %s", cts.errMsg))
	}

	return b.String()
}

// renderCreateTypePreview renders a read-only type schema preview for the right panel.
func renderCreateTypePreview(cts *createTypeState) string {
	if cts == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(" Type Schema Preview\n")
	b.WriteString(" ───────────────────\n")

	name := strings.TrimSpace(cts.nameInput.Value())
	if name == "" {
		name = "(not set)"
	}
	b.WriteString(fmt.Sprintf(" Name:   %s\n", name))

	plural := strings.TrimSpace(cts.pluralInput.Value())
	if plural == "" {
		plural = "(not set)"
	}
	b.WriteString(fmt.Sprintf(" Plural: %s\n", plural))

	emoji := strings.TrimSpace(cts.emojiInput.Value())
	if emoji == "" {
		emoji = "(not set)"
	}
	b.WriteString(fmt.Sprintf(" Emoji:  %s\n", emoji))

	b.WriteString(" Unique: no\n")
	b.WriteString("\n")
	b.WriteString(" Properties: (none)\n")

	return b.String()
}

// renderCreateTypeHelpBar returns the help bar text for the type creation state.
func renderCreateTypeHelpBar() string {
	return "  [NEW TYPE]  tab: next field  enter: create  esc: cancel"
}
