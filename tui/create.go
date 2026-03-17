package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// createField represents which field is focused in the title panel creation form.
type createField int

const (
	createFieldName     createField = iota // name text input (default)
	createFieldTemplate                    // template cycling selector
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
// Both name input and template selector are active simultaneously in the title panel.
type createState struct {
	mode      createMode
	field     createField    // which field is focused
	typeName  string
	emoji     string         // type emoji for title panel rendering
	schema    *core.TypeSchema // cached schema (loaded once in startCreate)
	templates []string       // available templates for this type (empty = no templates)
	cursor    int            // current template index (into templates + "(none)")
	nameInput textinput.Model
	flash     string         // success message (batch mode)
	flashSeq  int            // monotonic counter for flash dismiss correlation
	errMsg    string         // validation/creation error
	lastObj   *core.Object   // last created object (for batch mode exit)

	// Preview state — updated when template changes
	previewBody  string
	previewProps []core.DisplayProperty

	// Template cache to avoid repeated disk reads
	templateCache map[string]*core.Template
}

// noneOption is the label for the "no template" option in the template selector.
const noneOption = "(none)"

// templateOptionCount returns the total number of template options including "(none)".
func (cs *createState) templateOptionCount() int {
	return len(cs.templates) + 1
}

// hasTemplateSelector returns true if the template selector should be interactive.
func (cs *createState) hasTemplateSelector() bool {
	return len(cs.templates) >= 2
}

// currentTemplateName returns the display name of the currently selected template.
func (cs *createState) currentTemplateName() string {
	if cs.cursor >= 0 && cs.cursor < len(cs.templates) {
		return cs.templates[cs.cursor]
	}
	return noneOption
}

// selectedTemplateName returns the template name for object creation.
// Returns empty string if "(none)" is selected.
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

// startCreate initializes the creation flow for the given type group and mode.
func (m *model) startCreate(groupIndex int, mode createMode) tea.Cmd {
	if m.readOnly || groupIndex >= len(m.groups) {
		return nil
	}

	g := m.groups[groupIndex]
	typeName := g.Name

	// Fetch available templates
	var templates []string
	if m.vault != nil {
		templates, _ = m.vault.ListTemplates(typeName)
	}

	// Load schema once (cached for the entire creation flow)
	var schema *core.TypeSchema
	if m.vault != nil {
		schema, _ = m.vault.LoadType(typeName)
	}

	cs := &createState{
		mode:          mode,
		field:         createFieldName,
		typeName:      typeName,
		emoji:         g.Emoji,
		schema:        schema,
		nameInput:     initNameInput(),
		templateCache: make(map[string]*core.Template),
	}

	switch len(templates) {
	case 0:
		// No templates
	case 1:
		// Single template — auto-select, show as static label
		cs.templates = templates
		cs.cursor = 0
	default:
		// Multiple templates — interactive selector
		cs.templates = templates
		cs.cursor = 0
	}

	// Pre-fill name from name template if applicable
	if schema != nil && schema.NameTemplate != "" {
		prefill := core.EvaluateNameTemplate(schema.NameTemplate, time.Now())
		cs.nameInput.SetValue(prefill)
		cs.nameInput.CursorEnd()
	}

	m.create = cs
	m.updateCreatePreview()
	return textinput.Blink
}

// tryStartCreate attempts to start object creation from the current cursor position.
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

// updateCreatePreview loads the current template and updates preview state.
func (m *model) updateCreatePreview() {
	cs := m.create
	if cs == nil || m.vault == nil {
		return
	}

	tmplName := cs.selectedTemplateName()
	if tmplName != "" {
		tmpl := cs.loadTemplate(m.vault, tmplName)
		if tmpl != nil {
			cs.previewBody = tmpl.Body
			cs.previewProps = buildPreviewProps(cs.schema, tmpl.Properties)
		} else {
			cs.previewBody = ""
			cs.previewProps = buildPreviewProps(cs.schema, nil)
		}
	} else {
		cs.previewBody = ""
		cs.previewProps = buildPreviewProps(cs.schema, nil)
	}

	// Update viewports with preview content
	m.bodyViewport.SetContent(renderPreviewBody(cs.previewBody))
	m.propsViewport.SetContent(renderProperties(nil, cs.previewProps))
}

// loadTemplate loads a template with caching.
func (cs *createState) loadTemplate(vault *core.Vault, name string) *core.Template {
	if tmpl, ok := cs.templateCache[name]; ok {
		return tmpl
	}
	tmpl, err := vault.LoadTemplate(cs.typeName, name)
	if err != nil {
		return nil
	}
	cs.templateCache[name] = tmpl
	return tmpl
}

// buildPreviewProps builds display properties from a schema with optional template overrides.
func buildPreviewProps(schema *core.TypeSchema, overrides map[string]any) []core.DisplayProperty {
	if schema == nil {
		return nil
	}
	var props []core.DisplayProperty
	for _, p := range schema.Properties {
		dp := core.DisplayProperty{
			Key:   p.Name,
			Emoji: p.Emoji,
			Pin:   p.Pin,
		}
		if overrides != nil {
			if val, ok := overrides[p.Name]; ok {
				dp.Value = fmt.Sprintf("%v", val)
				props = append(props, dp)
				continue
			}
		}
		if p.Default != nil {
			dp.Value = fmt.Sprintf("%v", p.Default)
		}
		props = append(props, dp)
	}
	return props
}

// renderPreviewBody renders a template body for preview in the body viewport.
func renderPreviewBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return " (empty)\n"
	}
	var b strings.Builder
	for _, line := range strings.Split(body, "\n") {
		b.WriteString(fmt.Sprintf(" %s\n", line))
	}
	return b.String()
}

// renderCreateTitleContent renders the title panel content during creation mode.
func renderCreateTitleContent(cs *createState, width int) string {
	if cs == nil {
		return ""
	}

	// Flash takes priority in batch mode
	if cs.flash != "" {
		prefix := ""
		if cs.emoji != "" {
			prefix = padEmoji(cs.emoji) + " "
		}
		return fmt.Sprintf(" %s%s · %s", prefix, cs.typeName, cs.flash)
	}

	var b strings.Builder

	// Type prefix: "📚 book · "
	if cs.emoji != "" {
		b.WriteString(fmt.Sprintf(" %s %s · ", padEmoji(cs.emoji), cs.typeName))
	} else {
		b.WriteString(fmt.Sprintf(" %s · ", cs.typeName))
	}

	// Name input
	b.WriteString(cs.nameInput.View())

	// Template selector (if templates exist)
	if len(cs.templates) > 0 {
		tmplLabel := cs.currentTemplateName()
		if cs.field == createFieldTemplate {
			b.WriteString(fmt.Sprintf("  [📝 %s]", tmplLabel))
		} else {
			b.WriteString(fmt.Sprintf("  📝 %s", tmplLabel))
		}
	}

	// Error message
	if cs.errMsg != "" {
		b.WriteString(fmt.Sprintf("  ✗ %s", cs.errMsg))
	}

	return b.String()
}

// renderCreateHelpBar returns the help bar text for the current creation state.
func renderCreateHelpBar(cs *createState) string {
	mode := "NEW OBJECT"
	if cs.mode == createModeBatch {
		mode = "QUICK CREATE"
	}

	var hints []string
	if cs.field == createFieldName {
		if cs.hasTemplateSelector() {
			hints = append(hints, "tab: template")
		}
	} else {
		hints = append(hints, "◀▶: switch", "tab: name")
	}
	if cs.mode == createModeBatch {
		hints = append(hints, "enter: create", "esc: done")
	} else {
		hints = append(hints, "enter: create & edit", "esc: cancel")
	}

	return fmt.Sprintf("  [%s]  %s", mode, strings.Join(hints, "  "))
}

// executeCreate creates the object using the current createState settings.
func (m *model) executeCreate(cs *createState) tea.Cmd {
	name := strings.TrimSpace(cs.nameInput.Value())
	if name == "" {
		return nil
	}

	tmpl := cs.selectedTemplateName()
	obj, err := m.vault.Objects.Create(cs.typeName, name, tmpl)
	if err != nil {
		cs.errMsg = err.Error()
		m.create = cs
		return nil
	}

	cs.errMsg = ""
	m.saveErr = ""

	// Rebuild groups from index
	m.rebuildGroups()

	// Select the new object and move cursor
	m.selected = obj
	m.rightPanel = panelObject
	m.typeEditor = nil
	m.displayProps, _ = m.vault.BuildDisplayProperties(obj)
	m.updateDetail()
	m.moveCursorToObject(obj)

	if cs.mode == createModeSingle {
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

// updateCreate dispatches key events to the appropriate field handler.
func updateCreate(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	cs := m.create
	if cs == nil {
		return m, nil
	}

	// Global keys (both fields)
	switch msg.String() {
	case "enter":
		cmd := m.executeCreate(cs)
		return m, cmd

	case "esc":
		if cs.mode == createModeBatch && cs.lastObj != nil {
			m.moveCursorToObject(cs.lastObj)
		}
		m.create = nil
		// Restore normal view
		if m.selected != nil {
			m.displayProps, _ = m.vault.BuildDisplayProperties(m.selected)
			m.updateDetail()
		} else {
			m.bodyViewport.SetContent(renderBody(nil, m.bodyViewport.Width(), nil))
			m.propsViewport.SetContent("")
		}
		return m, nil

	case "tab":
		if cs.hasTemplateSelector() {
			if cs.field == createFieldName {
				cs.field = createFieldTemplate
				cs.nameInput.Blur()
			} else {
				cs.field = createFieldName
				cs.nameInput.Focus()
			}
		}
		return m, nil
	}

	// Field-specific keys
	switch cs.field {
	case createFieldName:
		return updateCreateNameField(m, msg)
	case createFieldTemplate:
		return updateCreateTemplateField(m, msg)
	}

	return m, nil
}

// updateCreateNameField handles key events when the name field is focused.
func updateCreateNameField(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	cs := m.create

	// Clear error on any input
	if cs.errMsg != "" {
		cs.errMsg = ""
	}

	var cmd tea.Cmd
	cs.nameInput, cmd = cs.nameInput.Update(msg)
	return m, cmd
}

// updateCreateTemplateField handles key events when the template field is focused.
func updateCreateTemplateField(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	cs := m.create
	total := cs.templateOptionCount()

	switch msg.String() {
	case "left", "up", "k":
		cs.cursor--
		if cs.cursor < 0 {
			cs.cursor = total - 1
		}
		m.updateCreatePreview()
		return m, nil

	case "right", "down", "j":
		cs.cursor++
		if cs.cursor >= total {
			cs.cursor = 0
		}
		m.updateCreatePreview()
		return m, nil
	}

	return m, nil
}
