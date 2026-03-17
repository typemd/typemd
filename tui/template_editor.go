package tui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// templateDeletedMsg signals that a template was deleted.
type templateDeletedMsg struct {
	TypeName     string
	TemplateName string
}

// tmplMode represents the current mode of the template editor.
type tmplMode int

const (
	tmplModeView     tmplMode = iota // read-only display
	tmplModeEditBody                 // textarea body editing
	tmplModeEditProp                 // single property value editing
	tmplModeDelete                   // delete confirmation
)

// tmplFocus tracks which sub-panel is focused within the template editor.
type tmplFocus int

const (
	tmplFocusBody  tmplFocus = iota // body viewport focused
	tmplFocusProps                  // props viewport focused
)

// tmplProp is a computed property list item for display.
type tmplProp struct {
	Name   string
	Value  string // string representation of current value
	Source string // "template" if from template, "schema" if placeholder
}

// templateEditor is an independent sub-model for viewing and editing templates.
type templateEditor struct {
	typeName     string
	templateName string
	template     *core.Template
	schema       *core.TypeSchema
	vault        *core.Vault

	mode tmplMode

	// Internal focus
	focus tmplFocus

	// View mode
	bodyViewport  viewport.Model
	propsViewport viewport.Model

	// Edit body mode
	bodyTextarea  textarea.Model
	bodyEditStart string // snapshot for cancel

	// Property editing
	propsList     []tmplProp // computed property list
	propsCursor   int
	propEditInput textinput.Model
	propEditIdx   int // index in propsList being edited

	// Layout
	width        int
	height       int
	bodyW        int
	propsW       int
	propsVisible bool

	// State
	dirty   bool
	saveErr string
}

// newTemplateEditor creates a template editor for the given template.
func newTemplateEditor(typeName, templateName string, tmpl *core.Template, schema *core.TypeSchema, vault *core.Vault) *templateEditor {
	bodyVP := viewport.New()
	propsVP := viewport.New()

	bodyTA := newBodyTextarea()

	ti := textinput.New()
	ti.CharLimit = 200

	te := &templateEditor{
		typeName:      typeName,
		templateName:  templateName,
		template:      tmpl,
		schema:        schema,
		vault:         vault,
		mode:          tmplModeView,
		focus:         tmplFocusBody,
		bodyViewport:  bodyVP,
		propsViewport: propsVP,
		bodyTextarea:  bodyTA,
		propEditInput: ti,
	}

	te.buildPropsList()
	te.refreshBodyViewport()
	te.propsViewport.SetContent(te.renderProps())
	return te
}

// buildPropsList computes the property list for display.
func (te *templateEditor) buildPropsList() {
	var props []tmplProp
	seen := make(map[string]bool)

	// Mutable system props from template
	for _, name := range []string{"name", "description", "tags"} {
		if val, ok := te.template.Properties[name]; ok {
			props = append(props, tmplProp{Name: name, Value: fmt.Sprintf("%v", val), Source: "template"})
			seen[name] = true
		}
	}

	// Schema-defined properties
	if te.schema != nil {
		for _, p := range te.schema.Properties {
			if seen[p.Name] {
				continue
			}
			if val, ok := te.template.Properties[p.Name]; ok {
				props = append(props, tmplProp{Name: p.Name, Value: fmt.Sprintf("%v", val), Source: "template"})
			} else {
				props = append(props, tmplProp{Name: p.Name, Value: "", Source: "schema"})
			}
			seen[p.Name] = true
		}
	}

	// Template props not in schema (non-system)
	for name, val := range te.template.Properties {
		if seen[name] || core.IsImmutableSystemProperty(name) || core.IsSystemProperty(name) {
			continue
		}
		props = append(props, tmplProp{Name: name, Value: fmt.Sprintf("%v", val), Source: "template"})
	}

	te.propsList = props
}

// Update handles messages for the template editor.
func (te *templateEditor) Update(msg tea.Msg) (*templateEditor, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return te, nil
	}

	switch te.mode {
	case tmplModeView:
		return te.updateView(keyMsg)
	case tmplModeEditBody:
		return te.updateEditBody(keyMsg)
	case tmplModeEditProp:
		return te.updateEditProp(keyMsg)
	case tmplModeDelete:
		return te.updateDelete(keyMsg)
	}
	return te, nil
}

func (te *templateEditor) updateView(msg tea.KeyPressMsg) (*templateEditor, tea.Cmd) {
	// When props panel is focused, route most keys there
	if te.focus == tmplFocusProps {
		switch msg.String() {
		case "tab":
			te.focus = tmplFocusBody
			return te, nil
		case "d":
			te.mode = tmplModeDelete
			return te, nil
		default:
			return te.updatePropsNav(msg)
		}
	}

	// Body panel focused
	switch msg.String() {
	case "e":
		te.mode = tmplModeEditBody
		te.bodyEditStart = te.template.Body
		te.bodyTextarea.SetValue(te.template.Body)
		te.bodyTextarea.Focus()
		return te, te.bodyTextarea.Focus()

	case "d":
		te.mode = tmplModeDelete
		return te, nil

	case "tab":
		if te.propsVisible {
			te.focus = tmplFocusProps
		}
		return te, nil

	default:
		return te.updateBodyNav(msg)
	}
}

func (te *templateEditor) updateBodyNav(msg tea.KeyPressMsg) (*templateEditor, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		te.bodyViewport.ScrollDown(1)
	case "k", "up":
		te.bodyViewport.ScrollUp(1)
	}
	return te, nil
}

func (te *templateEditor) updatePropsNav(msg tea.KeyPressMsg) (*templateEditor, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if te.propsCursor < len(te.propsList)-1 {
			te.propsCursor++
			te.propsViewport.SetContent(te.renderProps())
		}
	case "k", "up":
		if te.propsCursor > 0 {
			te.propsCursor--
			te.propsViewport.SetContent(te.renderProps())
		}
	case "enter":
		if len(te.propsList) > 0 && te.propsCursor < len(te.propsList) {
			te.propEditIdx = te.propsCursor
			te.propEditInput.SetValue(te.propsList[te.propsCursor].Value)
			te.propEditInput.Focus()
			te.mode = tmplModeEditProp
			te.propsViewport.SetContent(te.renderProps())
			return te, te.propEditInput.Focus()
		}
	}
	return te, nil
}

func (te *templateEditor) updateEditBody(msg tea.KeyPressMsg) (*templateEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Save
		te.template.Body = te.bodyTextarea.Value()
		te.bodyTextarea.Blur()
		te.dirty = true
		te.saveTemplate()
		te.mode = tmplModeView
		te.refreshBodyViewport()
		return te, nil

	case "ctrl+c":
		// Cancel
		te.bodyTextarea.SetValue(te.bodyEditStart)
		te.bodyTextarea.Blur()
		te.mode = tmplModeView
		return te, nil
	}

	var cmd tea.Cmd
	te.bodyTextarea, cmd = te.bodyTextarea.Update(msg)
	return te, cmd
}

func (te *templateEditor) updateEditProp(msg tea.KeyPressMsg) (*templateEditor, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Confirm: update template property
		newVal := te.propEditInput.Value()
		propName := te.propsList[te.propEditIdx].Name
		if newVal == "" {
			// Delete property from template
			delete(te.template.Properties, propName)
		} else {
			if te.template.Properties == nil {
				te.template.Properties = make(map[string]any)
			}
			te.template.Properties[propName] = newVal
		}
		te.propEditInput.Blur()
		te.dirty = true
		te.saveTemplate()
		te.buildPropsList()
		te.propsViewport.SetContent(te.renderProps())
		te.mode = tmplModeView
		return te, nil

	case "esc":
		// Cancel
		te.propEditInput.Blur()
		te.mode = tmplModeView
		te.propsViewport.SetContent(te.renderProps())
		return te, nil
	}

	var cmd tea.Cmd
	te.propEditInput, cmd = te.propEditInput.Update(msg)
	return te, cmd
}

func (te *templateEditor) updateDelete(msg tea.KeyPressMsg) (*templateEditor, tea.Cmd) {
	switch msg.String() {
	case "y":
		if te.vault != nil {
			if err := te.vault.DeleteTemplate(te.typeName, te.templateName); err != nil {
				te.saveErr = err.Error()
				te.mode = tmplModeView
				return te, nil
			}
		}
		te.mode = tmplModeView
		return te, tea.Sequence(func() tea.Msg {
			return templateDeletedMsg{
				TypeName:     te.typeName,
				TemplateName: te.templateName,
			}
		})
	case "n", "esc":
		te.mode = tmplModeView
	}
	return te, nil
}

// saveTemplate persists the current template to disk.
func (te *templateEditor) saveTemplate() {
	if te.vault == nil {
		return
	}
	if err := te.vault.SaveTemplate(te.typeName, te.templateName, te.template); err != nil {
		te.saveErr = err.Error()
	} else {
		te.saveErr = ""
		te.dirty = false
	}
}

// refreshBodyViewport updates the body viewport content from the template.
func (te *templateEditor) refreshBodyViewport() {
	body := strings.TrimSpace(te.template.Body)
	if body == "" {
		body = "(empty)"
	}
	var b strings.Builder
	for _, line := range strings.Split(body, "\n") {
		b.WriteString(fmt.Sprintf(" %s\n", line))
	}
	te.bodyViewport.SetContent(b.String())
}

// View renders the template editor panel.
func (te *templateEditor) View() string {
	switch te.mode {
	case tmplModeEditBody:
		return te.viewEditBody()
	case tmplModeDelete:
		return te.viewDelete()
	default:
		return te.viewNormal()
	}
}

func (te *templateEditor) viewNormal() string {

	bodyContent := te.bodyViewport.View()
	bodyStyle := lipgloss.NewStyle().
		Width(te.bodyW).
		Height(te.height)
	if te.focus == tmplFocusBody {
		bodyStyle = bodyStyle.Bold(true)
	}

	if !te.propsVisible {
		var b strings.Builder
		b.WriteString(bodyStyle.Render(bodyContent))
		if te.saveErr != "" {
			b.WriteString(fmt.Sprintf("\n [ERROR] %s", te.saveErr))
		}
		return b.String()
	}

	propsContent := te.propsViewport.View()
	propsStyle := lipgloss.NewStyle().
		Width(te.propsW).
		Height(te.height)
	if te.focus == tmplFocusProps {
		propsStyle = propsStyle.Bold(true)
	}

	result := lipgloss.JoinHorizontal(lipgloss.Top,
		bodyStyle.Render(bodyContent),
		propsStyle.Render(propsContent),
	)

	if te.saveErr != "" {
		result += fmt.Sprintf("\n [ERROR] %s", te.saveErr)
	}

	return result
}

func (te *templateEditor) viewEditBody() string {
	return te.bodyTextarea.View()
}

func (te *templateEditor) viewDelete() string {
	return fmt.Sprintf("\n Delete template '%s'? [y/n]\n", te.templateName)
}

// renderProps renders the property list for the properties viewport.
func (te *templateEditor) renderProps() string {
	var b strings.Builder
	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	if len(te.propsList) == 0 {
		b.WriteString(" (none)\n")
		return b.String()
	}

	for i, p := range te.propsList {
		val := p.Value
		if val == "" {
			val = "(empty)"
		}
		line := fmt.Sprintf("%s: %s", p.Name, val)
		if te.mode == tmplModeEditProp && te.propEditIdx == i {
			line = fmt.Sprintf("%s: %s", p.Name, te.propEditInput.View())
		}
		if i == te.propsCursor {
			b.WriteString(" " + highlightStyle.Render(" "+line+" ") + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}

	return b.String()
}

// HelpBar returns the context-sensitive help text for the template editor.
func (te *templateEditor) HelpBar() string {
	switch te.mode {
	case tmplModeView:
		return "  [TEMPLATE]  e: edit  d: delete  tab: switch  esc: back"
	case tmplModeEditBody:
		return "  [EDIT]  esc: save  ctrl+c: cancel"
	case tmplModeEditProp:
		return "  [EDIT PROP]  enter: confirm  esc: cancel"
	case tmplModeDelete:
		return "  [DELETE]  y: confirm  n/esc: cancel"
	}
	return ""
}

// CanQuit returns true when the editor is in a non-interactive state and the app can safely quit.
func (te *templateEditor) CanQuit() bool {
	return te.mode == tmplModeView
}

// SetSize updates layout dimensions and viewport sizes.
func (te *templateEditor) SetSize(width, height, propsW int, propsVisible bool) {
	te.width = width
	te.height = height
	te.propsW = propsW
	te.propsVisible = propsVisible

	if propsVisible {
		te.bodyW = width - propsW
	} else {
		te.bodyW = width
	}
	if te.bodyW < 10 {
		te.bodyW = 10
	}

	te.bodyViewport.SetWidth(te.bodyW)
	te.bodyViewport.SetHeight(height)
	te.propsViewport.SetWidth(propsW)
	te.propsViewport.SetHeight(height)

	// Update textarea dimensions for edit mode
	te.bodyTextarea.SetWidth(te.bodyW)
	te.bodyTextarea.SetHeight(height)
}

// titleContent returns the title string for the parent to render in the title panel.
func (te *templateEditor) titleContent(width int) string {
	title := fmt.Sprintf(" 📝 %s · %s", te.typeName, te.templateName)
	return runewidth.Truncate(title, width, "…")
}
