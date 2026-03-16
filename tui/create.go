package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// createStep represents the current step in the object creation flow.
type createStep int

const (
	createStepTemplate createStep = iota // selecting a template
	createStepName                       // entering object name
)

// createMode represents the type of creation flow.
type createMode int

const (
	createModeSingle createMode = iota + 1 // n: Create & Edit (single object)
	createModeBatch                        // N: Quick Create (batch)
)

// flashDismissMsg is sent by tea.Tick to auto-dismiss the flash message.
// seq identifies which flash it belongs to, preventing early dismissal of newer flashes.
type flashDismissMsg struct {
	seq int
}

// createState holds all state for the object creation flow.
type createState struct {
	mode             createMode
	step             createStep
	typeName         string
	templates        []string       // available templates for this type
	cursor           int            // template selection cursor
	template         string         // selected template name (empty = none)
	nameInput        textinput.Model
	flash            string       // success message (batch mode)
	flashSeq         int          // monotonic counter for flash dismiss correlation
	errMsg           string       // validation/creation error
	lastObj          *core.Object // last created object (for batch mode exit)
	nameTemplateSkip bool         // true when name template auto-skip is active
}

// noneOption is the label for the "no template" option in the template selection list.
const noneOption = "(none)"

// templateOptions returns the template list plus a "(none)" option at the end.
func (cs *createState) templateOptions() []string {
	opts := make([]string, len(cs.templates))
	copy(opts, cs.templates)
	return append(opts, noneOption)
}

// selectedTemplateName returns the template name based on cursor position.
// Returns empty string if "(none)" is selected (last item in the list).
func (cs *createState) selectedTemplateName() string {
	if cs.cursor >= 0 && cs.cursor < len(cs.templates) {
		return cs.templates[cs.cursor]
	}
	return ""
}

// initNameInput creates and configures the text input for name entry.
func initNameInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "object name"
	ti.CharLimit = 100
	ti.Focus()
	return ti
}

// tryNameTemplateSkip checks if the type has a name template and auto-creates in single mode.
// Returns the tea.Cmd and true if auto-skip was triggered, or nil and false otherwise.
func (m *model) tryNameTemplateSkip(cs *createState) (tea.Cmd, bool) {
	if cs.mode != createModeSingle || m.vault == nil {
		return nil, false
	}
	schema, err := m.vault.LoadType(cs.typeName)
	if err != nil || schema.NameTemplate == "" {
		return nil, false
	}
	cs.nameTemplateSkip = true
	return m.executeCreate(cs), true
}

// startCreate initializes the creation flow for the given type group and mode.
func (m *model) startCreate(groupIndex int, mode createMode) tea.Cmd {
	if m.readOnly || groupIndex >= len(m.groups) {
		return nil
	}

	typeName := m.groups[groupIndex].Name

	// Fetch available templates
	var templates []string
	if m.vault != nil {
		templates, _ = m.vault.ListTemplates(typeName)
	}

	cs := &createState{
		mode:     mode,
		typeName: typeName,
	}

	switch len(templates) {
	case 0:
		// No templates — skip to name step
		cs.step = createStepName
		cs.nameInput = initNameInput()
	case 1:
		// Single template — auto-select
		cs.template = templates[0]
		cs.step = createStepName
		cs.nameInput = initNameInput()
	default:
		// Multiple templates — show selection
		cs.templates = templates
		cs.step = createStepTemplate
		cs.cursor = 0
	}

	// Check for name template auto-skip in single mode
	if cs.step == createStepName {
		if cmd, ok := m.tryNameTemplateSkip(cs); ok {
			return cmd
		}
	}

	m.create = cs
	return textinput.Blink
}

// tryStartCreate attempts to start object creation from the current cursor position.
// Returns the model and command, or false if the cursor is not on a valid row.
func (m *model) tryStartCreate(mode createMode) (tea.Cmd, bool) {
	if m.readOnly || m.focus != focusLeft {
		return nil, false
	}
	rows := m.currentRows()
	if m.cursor >= 0 && m.cursor < len(rows) {
		row := rows[m.cursor]
		if row.Kind == rowHeader || row.Kind == rowObject {
			return m.startCreate(row.GroupIndex, mode), true
		}
	}
	return nil, false
}

// executeCreate creates the object using the current createState settings.
// Returns a tea.Cmd (possibly a Tick for flash dismiss).
func (m *model) executeCreate(cs *createState) tea.Cmd {
	name := ""
	if cs.step == createStepName && !cs.nameTemplateSkip {
		name = strings.TrimSpace(cs.nameInput.Value())
		if name == "" {
			return nil
		}
	}

	obj, err := m.vault.Objects.Create(cs.typeName, name, cs.template)
	if err != nil {
		cs.errMsg = err.Error()
		if cs.step != createStepName {
			// Error during auto-create (name template) — need to show it somewhere
			m.saveErr = err.Error()
			m.create = nil
			return nil
		}
		m.create = cs
		return nil
	}

	cs.errMsg = ""
	m.saveErr = ""

	// Rebuild groups from index (ObjectService.Create already indexed the object,
	// so we skip the expensive SyncIndex filesystem walk).
	objects, qErr := m.vault.QueryObjects("")
	if qErr == nil {
		m.groups = buildGroups(objects, m.vault)
		m.searchResults = nil
	}

	// Select the new object and move cursor
	m.selected = obj
	m.rightPanel = panelObject
	m.typeEditor = nil
	m.displayProps, _ = m.vault.BuildDisplayProperties(obj)
	m.updateDetail()
	m.moveCursorToObject(obj)

	if cs.mode == createModeSingle {
		// Create & Edit: enter body edit mode
		m.create = nil
		return m.enterBodyEditMode()
	}

	// Quick Create: clear input, show flash, stay in create mode
	cs.lastObj = obj
	cs.flashSeq++
	cs.flash = fmt.Sprintf("✓ Created: %s", obj.GetName())
	seq := cs.flashSeq
	cs.nameInput.Reset()
	cs.nameInput.Focus()
	m.create = cs
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return flashDismissMsg{seq: seq}
	})
}

// enterBodyEditMode sets up body editing and returns the focus command.
func (m *model) enterBodyEditMode() tea.Cmd {
	m.focus = focusBody
	m.editMode = true
	m.bodyTextarea.SetValue(m.selected.Body)
	m.bodyEditStart = m.bodyTextarea.Value()
	m.resizeBodyTextarea()
	m.bodyTextarea.CursorEnd()
	return m.bodyTextarea.Focus()
}

// renderCreateUI renders the creation flow UI for the sidebar.
func renderCreateUI(cs *createState) string {
	var lines []string

	// Flash message (batch mode)
	if cs.flash != "" {
		lines = append(lines, " "+cs.flash)
	}

	switch cs.step {
	case createStepTemplate:
		lines = append(lines, " Select template:")
		opts := cs.templateOptions()
		for i, opt := range opts {
			prefix := "   "
			if i == cs.cursor {
				prefix = " > "
			}
			lines = append(lines, prefix+opt)
		}

	case createStepName:
		lines = append(lines, " New "+cs.typeName+": "+cs.nameInput.View())
		if cs.errMsg != "" {
			lines = append(lines, " ✗ "+cs.errMsg)
		}
	}

	return strings.Join(lines, "\n")
}

// renderCreateHelpBar returns the help bar text for the current creation state.
func renderCreateHelpBar(cs *createState) string {
	switch cs.step {
	case createStepTemplate:
		return "  [NEW OBJECT]  ↑↓: select  enter: confirm  esc: cancel"
	case createStepName:
		if cs.mode == createModeBatch {
			return "  [QUICK CREATE]  enter: create  esc: done"
		}
		return "  [NEW OBJECT]  enter: create & edit  esc: cancel"
	}
	return ""
}

// moveCursorToObject moves the sidebar cursor to the given object.
func (m *model) moveCursorToObject(obj *core.Object) {
	rows := m.currentRows()
	for i, row := range rows {
		if row.Kind == rowObject && row.Object != nil && row.Object.ID == obj.ID {
			m.cursor = i
			m.adjustScroll()
			break
		}
	}
}

// updateCreate dispatches key events to the appropriate creation step handler.
func updateCreate(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	if m.create == nil {
		return m, nil
	}

	switch m.create.step {
	case createStepTemplate:
		return updateCreateTemplate(m, msg)
	case createStepName:
		return updateCreateName(m, msg)
	}
	return m, nil
}

// updateCreateTemplate handles key events during template selection.
func updateCreateTemplate(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	cs := m.create
	opts := cs.templateOptions()

	switch msg.String() {
	case "up", "k":
		cs.cursor--
		if cs.cursor < 0 {
			cs.cursor = len(opts) - 1
		}
		return m, nil

	case "down", "j":
		cs.cursor++
		if cs.cursor >= len(opts) {
			cs.cursor = 0
		}
		return m, nil

	case "enter":
		cs.template = cs.selectedTemplateName()
		cs.step = createStepName
		cs.nameInput = initNameInput()

		// Check for name template auto-skip in single mode
		if cmd, ok := m.tryNameTemplateSkip(cs); ok {
			return m, cmd
		}

		return m, textinput.Blink

	case "esc":
		m.create = nil
		return m, nil
	}

	return m, nil
}

// updateCreateName handles key events during name input.
func updateCreateName(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	cs := m.create

	switch msg.String() {
	case "enter":
		cmd := m.executeCreate(cs)
		return m, cmd

	case "esc":
		if cs.mode == createModeBatch && cs.lastObj != nil {
			// Exit batch mode — select last created object
			m.moveCursorToObject(cs.lastObj)
		}
		m.create = nil
		return m, nil
	}

	// Clear error on any other key input
	if cs.errMsg != "" {
		cs.errMsg = ""
	}

	var cmd tea.Cmd
	cs.nameInput, cmd = cs.nameInput.Update(msg)
	return m, cmd
}
