package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// focusLeftMsg signals the parent model to move focus back to the sidebar.
type focusLeftMsg struct{}

// typeDeletedMsg signals that the current type was deleted.
type typeDeletedMsg struct{}

// openTemplateMsg signals the parent model to open a template in panelTemplate mode.
type openTemplateMsg struct {
	TypeName     string
	TemplateName string
}

// teMode represents the current mode of the type editor.
type teMode int

const (
	teModeView        teMode = iota // viewing type schema
	teModeEditMeta                  // editing a meta field inline
	teModeEditProp                  // property detail panel
	teModeMove                      // reordering properties
	teModeAddWizard                 // add property wizard
	teModeDeleteProp                // delete property confirmation
	teModeDeleteType                // delete type confirmation
	teModeAddTemplate               // entering new template name
	teModeAddView                   // entering new view name
)

// metaFieldCount is the number of meta fields at the top of the editor:
// Name (0), Plural (1), Emoji (2), Color (3), Unique (4), Description (5)
const metaFieldCount = 6

// typeEditor is an independent sub-model for editing type schemas.
type typeEditor struct {
	schema   *core.TypeSchema
	typeName string
	isNew    bool
	vault    *core.Vault

	// Internal state
	cursor int    // unified cursor across meta fields and properties
	mode   teMode

	// Inline edit (meta fields)
	editInput textinput.Model
	editField int // which cursor index is being edited

	// Property detail panel
	propDetail    *propDetailPanel
	propDetailIdx int // index in schema.Properties being edited

	// Move mode
	moveFrom int // original index of the property being moved

	// Add property wizard
	wizard *addPropWizard

	// Template management
	templates     []string         // cached template names for this type
	tmplCursor    int              // cursor within templates list
	tmplNameInput textinput.Model  // name input for adding template

	// View management
	views         []string        // cached view names for this type
	viewNameInput textinput.Model // name input for adding view

	// Layout
	width        int
	height       int
	scrollOffset int

	// Status
	saveErr string
}

// wizardStep tracks the current step of the add property wizard.
type wizardStep int

const (
	wizStepName    wizardStep = iota // Step 1: enter property name
	wizStepType                      // Step 2: select property type
	wizStepOptions                   // Step 2b: enter options for select/multi_select
	wizStepRelation                  // Step 3: relation config
)

// addPropWizard holds state for the multi-step add property wizard.
type addPropWizard struct {
	step     wizardStep
	nameInput textinput.Model
	propName  string
	propType  string

	// Step 2: type selection
	typeCursor int
	typeList   []string

	// Step 2b: options
	optionsInput textinput.Model

	// Step 3: relation config
	relTargetCursor int
	relTargets      []string
	relMultiple     bool
	relBidir        bool
	relInverseInput textinput.Model
	relFieldCursor  int // 0=target, 1=multiple, 2=bidir, 3=inverse
}

// propertyTypeList is derived from core to stay in sync.
var propertyTypeList = core.ValidPropertyTypeNames()

// newTypeEditor creates a type editor for the given schema.
func newTypeEditor(schema *core.TypeSchema, typeName string, isNew bool, vault *core.Vault) *typeEditor {
	ti := textinput.New()
	ti.CharLimit = 100

	tmplInput := textinput.New()
	tmplInput.Placeholder = "template name"
	tmplInput.CharLimit = 50

	viewInput := textinput.New()
	viewInput.Placeholder = "view name"
	viewInput.CharLimit = 50

	var templates []string
	if vault != nil {
		templates, _ = vault.ListTemplates(typeName)
	}

	var views []string
	if vault != nil {
		viewConfigs, _ := vault.ListViews(typeName)
		for _, v := range viewConfigs {
			views = append(views, v.Name)
		}
	}

	return &typeEditor{
		schema:        schema,
		typeName:      typeName,
		isNew:         isNew,
		vault:         vault,
		editInput:     ti,
		templates:     templates,
		tmplNameInput: tmplInput,
		views:         views,
		viewNameInput: viewInput,
	}
}

// orderedProperties returns properties split into pinned and unpinned,
// preserving their original indices in the schema.Properties slice.
func (te *typeEditor) orderedProperties() (pinned []int, unpinned []int) {
	for i, p := range te.schema.Properties {
		if p.Pin > 0 {
			pinned = append(pinned, i)
		} else {
			unpinned = append(unpinned, i)
		}
	}
	sort.Slice(pinned, func(i, j int) bool {
		return te.schema.Properties[pinned[i]].Pin < te.schema.Properties[pinned[j]].Pin
	})
	return
}

// maxPinValue returns the highest pin value across all properties.
func maxPinValue(props []core.Property) int {
	max := 0
	for _, p := range props {
		if p.Pin > max {
			max = p.Pin
		}
	}
	return max
}

// Sentinel value for the "+ Add Property" row in displayItems.
const addPropertySentinel = -100

// Sentinel values for template management in displayItems.
const addTemplateSentinel = -200
const templateSentinelBase = -300 // templates use -300, -301, -302, ...

// Sentinel values for view management in displayItems.
const addViewSentinel = -250
const viewSentinelBase = -400 // views use -400, -401, -402, ...

// displayItems returns the flat list of cursor-addressable items.
// Items 0..5 are meta fields, then pinned properties, then unpinned properties,
// then the "+ Add Property" action row.
// Section separators are not included (they're visual only).
func (te *typeEditor) displayItems() []int {
	pinned, unpinned := te.orderedProperties()
	items := make([]int, 0, metaFieldCount+len(pinned)+len(unpinned)+1+len(te.templates)+1)
	// Meta fields use negative sentinel values: -1=Name, -2=Plural, -3=Emoji, -4=Color, -5=Unique, -6=Description
	for i := 0; i < metaFieldCount; i++ {
		items = append(items, -(i + 1))
	}
	items = append(items, pinned...)
	items = append(items, unpinned...)
	items = append(items, addPropertySentinel)
	// Template items
	for i := range te.templates {
		items = append(items, templateSentinelBase-i)
	}
	items = append(items, addTemplateSentinel)
	// View items
	for i := range te.views {
		items = append(items, viewSentinelBase-i)
	}
	items = append(items, addViewSentinel)
	return items
}

// totalItems returns the total number of cursor-addressable items.
func (te *typeEditor) totalItems() int {
	return metaFieldCount + len(te.schema.Properties) + 1 + len(te.templates) + 1 + len(te.views) + 1
	// meta + properties + addProperty + templates + addTemplate + views + addView
}

// save persists the current schema to disk.
func (te *typeEditor) save() {
	if te.vault == nil {
		return
	}
	if err := te.vault.SaveType(te.schema); err != nil {
		te.saveErr = err.Error()
	} else {
		te.saveErr = ""
	}
}

// Update handles messages for the type editor.
func (te *typeEditor) Update(msg tea.Msg) (*typeEditor, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return te, nil
	}

	switch te.mode {
	case teModeView:
		return te.updateView(keyMsg)
	case teModeEditMeta:
		return te.updateEdit(keyMsg)
	case teModeEditProp:
		return te.updatePropDetail(keyMsg)
	case teModeMove:
		return te.updateMove(keyMsg)
	case teModeAddWizard:
		return te.updateAddWizard(keyMsg)
	case teModeDeleteProp:
		return te.updateDeleteProp(keyMsg)
	case teModeDeleteType:
		return te.updateDeleteType(keyMsg)
	case teModeAddTemplate:
		return te.updateAddTemplate(keyMsg)
	case teModeAddView:
		return te.updateAddView(keyMsg)
	}
	return te, nil
}

func (te *typeEditor) updateView(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	total := te.totalItems()

	switch msg.String() {
	case "esc":
		// Signal parent to move focus back to sidebar (not close editor)
		return te, tea.Sequence(func() tea.Msg { return focusLeftMsg{} })

	case "up", "k":
		if te.cursor > 0 {
			te.cursor--
		}

	case "down", "j":
		if te.cursor < total-1 {
			te.cursor++
		}

	case "e":
		te.startEdit()

	case "p":
		te.togglePin()

	case "m":
		te.startMove()

	case "enter":
		items := te.displayItems()
		if te.cursor < len(items) {
			item := items[te.cursor]
			switch {
			case item == addPropertySentinel:
				te.startAddWizard()
			case item == addTemplateSentinel:
				te.startAddTemplate()
			case item == addViewSentinel:
				te.startAddView()
			case item <= viewSentinelBase:
				// View item — signal parent to open view mode
				viewIdx := viewSentinelBase - item
				if viewIdx >= 0 && viewIdx < len(te.views) {
					viewName := te.views[viewIdx]
					return te, func() tea.Msg {
						return openViewMsg{
							TypeName: te.typeName,
							ViewName: viewName,
						}
					}
				}
			case item <= templateSentinelBase:
				// Template item — signal parent to open template
				tmplIdx := templateSentinelBase - item
				if tmplIdx >= 0 && tmplIdx < len(te.templates) {
					return te, tea.Sequence(func() tea.Msg {
						return openTemplateMsg{
							TypeName:     te.typeName,
							TemplateName: te.templates[tmplIdx],
						}
					})
				}
			case item >= 0: // property
				te.openPropDetail()
			}
		}

	case "a":
		te.startAddWizard()

	case "D": // Shift+D = delete type
		te.mode = teModeDeleteType

	case "d":
		te.startDelete()
	}

	return te, nil
}

func (te *typeEditor) startEdit() {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	item := items[te.cursor]
	if item == addPropertySentinel || item == addTemplateSentinel || item == addViewSentinel || item <= templateSentinelBase {
		return
	}

	switch item {
	case -1: // Name — not editable
		return
	case -2: // Plural
		te.mode = teModeEditMeta
		te.editField = te.cursor
		te.editInput.SetValue(te.schema.Plural)
		te.editInput.Focus()
	case -3: // Emoji
		te.mode = teModeEditMeta
		te.editField = te.cursor
		te.editInput.SetValue(te.schema.Emoji)
		te.editInput.Focus()
	case -4: // Color
		te.mode = teModeEditMeta
		te.editField = te.cursor
		te.editInput.SetValue(te.schema.Color)
		te.editInput.Focus()
	case -5: // Unique — toggle
		te.schema.Unique = !te.schema.Unique
		te.save()
	case -6: // Description
		te.mode = teModeEditMeta
		te.editField = te.cursor
		te.editInput.SetValue(te.schema.Description)
		te.editInput.Focus()
	default: // Property — e does nothing, use enter for detail panel
		return
	}
}

func (te *typeEditor) updateEdit(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "enter":
		te.confirmEdit()
		return te, nil
	case "esc":
		te.mode = teModeView
		te.editInput.Blur()
		return te, nil
	}

	var cmd tea.Cmd
	te.editInput, cmd = te.editInput.Update(msg)
	return te, cmd
}

func (te *typeEditor) confirmEdit() {
	items := te.displayItems()
	if te.editField >= len(items) {
		te.mode = teModeView
		return
	}
	item := items[te.editField]
	val := te.editInput.Value()

	switch item {
	case -2: // Plural
		te.schema.Plural = val
	case -3: // Emoji
		te.schema.Emoji = val
	case -4: // Color
		te.schema.Color = val
	case -6: // Description
		te.schema.Description = val
	}
	te.save()

	te.mode = teModeView
	te.editInput.Blur()
}

func (te *typeEditor) togglePin() {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	item := items[te.cursor]
	if item < 0 {
		return // meta fields, not a property
	}

	prop := &te.schema.Properties[item]
	if prop.Pin > 0 {
		// Unpin
		prop.Pin = 0
	} else {
		// Pin: assign max(existing pins) + 1
		prop.Pin = maxPinValue(te.schema.Properties) + 1
	}
	te.save()
}

func (te *typeEditor) startMove() {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	item := items[te.cursor]
	if item < 0 {
		return // can't move meta fields
	}
	te.mode = teModeMove
	te.moveFrom = te.cursor
}

func (te *typeEditor) updateMove(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		te.moveProperty(-1)
	case "down", "j":
		te.moveProperty(1)
	case "enter", "esc":
		te.mode = teModeView
		te.save()
	}
	return te, nil
}

func (te *typeEditor) moveProperty(dir int) {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	curItem := items[te.cursor]
	if curItem < 0 {
		return
	}

	newCursor := te.cursor + dir
	if newCursor < metaFieldCount || newCursor >= len(items) {
		return // can't move beyond boundaries
	}
	targetItem := items[newCursor]
	if targetItem < 0 {
		return // can't swap with meta fields
	}

	// Swap in the schema.Properties slice
	te.schema.Properties[curItem], te.schema.Properties[targetItem] =
		te.schema.Properties[targetItem], te.schema.Properties[curItem]

	// Handle cross-section pin adjustment
	a := &te.schema.Properties[curItem]
	b := &te.schema.Properties[targetItem]

	// If a was pinned and b was not (or vice versa), adjust pins
	if a.Pin > 0 && b.Pin == 0 {
		// a moved from pinned into unpinned area
		a.Pin = 0
		if b.Pin == 0 {
			maxPin := 0
			for _, p := range te.schema.Properties {
				if p.Pin > maxPin {
					maxPin = p.Pin
				}
			}
			b.Pin = maxPin + 1
		}
	} else if a.Pin == 0 && b.Pin > 0 {
		// a moved from unpinned into pinned area
		maxPin := 0
		for _, p := range te.schema.Properties {
			if p.Pin > maxPin {
				maxPin = p.Pin
			}
		}
		a.Pin = maxPin + 1
		b.Pin = 0
	}

	te.cursor = newCursor
}

func (te *typeEditor) startDelete() {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	item := items[te.cursor]
	if item < 0 {
		return // can't delete meta fields, sentinels, or templates
	}
	te.mode = teModeDeleteProp
}

func (te *typeEditor) updateDeleteProp(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "y":
		items := te.displayItems()
		if te.cursor < len(items) {
			item := items[te.cursor]
			if item >= 0 && item < len(te.schema.Properties) {
				te.schema.Properties = append(te.schema.Properties[:item], te.schema.Properties[item+1:]...)
				te.save()
				// Clamp cursor
				total := te.totalItems()
				if te.cursor >= total {
					te.cursor = total - 1
				}
			}
		}
		te.mode = teModeView
	case "n", "esc":
		te.mode = teModeView
	}
	return te, nil
}

func (te *typeEditor) updateDeleteType(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "y":
		if te.vault != nil {
			if err := te.vault.DeleteType(te.typeName); err != nil {
				te.saveErr = err.Error()
				te.mode = teModeView
				return te, nil
			}
		}
		te.mode = teModeView
		return te, tea.Sequence(func() tea.Msg { return typeDeletedMsg{} })
	case "n", "esc":
		te.mode = teModeView
	}
	return te, nil
}

// ── Property Detail Panel ────────────────────────────────────────────────────

// propDetailPanel holds state for the property metadata editing panel.
type propDetailPanel struct {
	cursor     int // 0=emoji (more fields in future)
	emojiInput textinput.Model
	editing    bool // currently in text input mode
}

func newPropDetailPanel(prop *core.Property) *propDetailPanel {
	ei := textinput.New()
	ei.CharLimit = 20
	ei.SetValue(prop.Emoji)
	return &propDetailPanel{
		emojiInput: ei,
	}
}

func (te *typeEditor) openPropDetail() {
	items := te.displayItems()
	if te.cursor >= len(items) {
		return
	}
	item := items[te.cursor]
	if item < 0 || item == addPropertySentinel {
		return
	}
	te.propDetailIdx = item
	te.propDetail = newPropDetailPanel(&te.schema.Properties[item])
	te.mode = teModeEditProp
}

func (te *typeEditor) updatePropDetail(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	pd := te.propDetail
	if pd == nil {
		te.mode = teModeView
		return te, nil
	}

	if pd.editing {
		switch msg.String() {
		case "enter":
			pd.editing = false
			pd.emojiInput.Blur()
			// Apply value
			te.schema.Properties[te.propDetailIdx].Emoji = pd.emojiInput.Value()
			te.save()
			return te, nil
		case "esc":
			pd.editing = false
			pd.emojiInput.Blur()
			// Revert
			pd.emojiInput.SetValue(te.schema.Properties[te.propDetailIdx].Emoji)
			return te, nil
		}
		var cmd tea.Cmd
		pd.emojiInput, cmd = pd.emojiInput.Update(msg)
		return te, cmd
	}

	switch msg.String() {
	case "esc":
		te.propDetail = nil
		te.mode = teModeView
	case "enter", "e":
		pd.editing = true
		pd.emojiInput.Focus()
		return te, pd.emojiInput.Focus()
	case "up", "k":
		// future: navigate between fields
	case "down", "j":
		// future: navigate between fields
	}
	return te, nil
}

// Overlay returns a popup string if a modal is active, or empty string if not.
func (te *typeEditor) Overlay(width, height int) string {
	if te.mode != teModeEditProp || te.propDetail == nil {
		return ""
	}
	return te.renderPropPopup(width, height)
}

func (te *typeEditor) renderPropPopup(termW, termH int) string {
	pd := te.propDetail
	p := te.schema.Properties[te.propDetailIdx]

	var b strings.Builder

	if pd.editing {
		b.WriteString(fmt.Sprintf("Emoji: %s", pd.emojiInput.View()))
	} else {
		val := p.Emoji
		if val == "" {
			val = "(none)"
		} else {
			val = padEmoji(val)
		}
		line := fmt.Sprintf("Emoji: %s", val)
		if pd.cursor == 0 {
			line = highlightStyle.Render(line)
		}
		b.WriteString(line)
	}

	b.WriteString("\n")

	// future: description field here

	if pd.editing {
		b.WriteString("\nenter: save  esc: cancel")
	} else {
		b.WriteString("\nenter: edit  esc: back")
	}

	popupW := 36
	if popupW > termW-10 {
		popupW = termW - 10
	}

	titleStyle := lipgloss.NewStyle().Bold(true)
	title := titleStyle.Render(fmt.Sprintf("%s (%s)", p.Name, p.Type))
	fullContent := fmt.Sprintf("%s\n──────────────────\n%s", title, b.String())

	return renderPopup(fullContent, termW, termH, popupW)
}

// ── Add Property Wizard ─────────────────────────────────────────────────────

func (te *typeEditor) startAddWizard() {
	ni := textinput.New()
	ni.Placeholder = "property name"
	ni.CharLimit = 50
	ni.Focus()

	oi := textinput.New()
	oi.Placeholder = "option1, option2, ..."
	oi.CharLimit = 200

	ii := textinput.New()
	ii.Placeholder = "inverse property name"
	ii.CharLimit = 50

	// Gather target types
	var targets []string
	if te.vault != nil {
		targets = te.vault.ListTypes()
	}

	te.wizard = &addPropWizard{
		step:            wizStepName,
		nameInput:       ni,
		typeList:        propertyTypeList,
		optionsInput:    oi,
		relTargets:      targets,
		relInverseInput: ii,
	}
	te.mode = teModeAddWizard
}

func (te *typeEditor) updateAddWizard(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	if te.wizard == nil {
		te.mode = teModeView
		return te, nil
	}
	wiz := te.wizard

	switch wiz.step {
	case wizStepName:
		return te.updateWizardName(msg)
	case wizStepType:
		return te.updateWizardType(msg)
	case wizStepOptions:
		return te.updateWizardOptions(msg)
	case wizStepRelation:
		return te.updateWizardRelation(msg)
	}
	return te, nil
}

func (te *typeEditor) updateWizardName(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(wiz.nameInput.Value())
		if name == "" {
			return te, nil
		}
		// Check for duplicate
		for _, p := range te.schema.Properties {
			if p.Name == name {
				te.saveErr = fmt.Sprintf("property %q already exists", name)
				return te, nil
			}
		}
		// Check reserved system property names
		if core.IsSystemProperty(name) {
			te.saveErr = fmt.Sprintf("%q is a reserved system property", name)
			return te, nil
		}
		te.saveErr = ""
		wiz.propName = name
		wiz.step = wizStepType
		return te, nil
	case "esc":
		te.cancelWizard()
		return te, nil
	}
	var cmd tea.Cmd
	wiz.nameInput, cmd = wiz.nameInput.Update(msg)
	return te, cmd
}

func (te *typeEditor) updateWizardType(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "up", "k":
		if wiz.typeCursor > 0 {
			wiz.typeCursor--
		}
	case "down", "j":
		if wiz.typeCursor < len(wiz.typeList)-1 {
			wiz.typeCursor++
		}
	case "enter":
		wiz.propType = wiz.typeList[wiz.typeCursor]
		switch wiz.propType {
		case "select", "multi_select":
			wiz.step = wizStepOptions
			wiz.optionsInput.Focus()
		case "relation":
			wiz.step = wizStepRelation
		default:
			te.finishWizard()
		}
	case "esc":
		wiz.step = wizStepName
	}
	return te, nil
}

func (te *typeEditor) updateWizardOptions(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "enter":
		te.finishWizard()
	case "esc":
		wiz.step = wizStepType
	default:
		var cmd tea.Cmd
		wiz.optionsInput, cmd = wiz.optionsInput.Update(msg)
		return te, cmd
	}
	return te, nil
}

func (te *typeEditor) updateWizardRelation(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "up", "k":
		switch wiz.relFieldCursor {
		case 0: // target list
			if wiz.relTargetCursor > 0 {
				wiz.relTargetCursor--
			}
		default:
			if wiz.relFieldCursor > 0 {
				wiz.relFieldCursor--
			}
		}
	case "down", "j":
		switch wiz.relFieldCursor {
		case 0: // target list
			if wiz.relTargetCursor < len(wiz.relTargets)-1 {
				wiz.relTargetCursor++
			}
		default:
			if wiz.relFieldCursor < 3 {
				wiz.relFieldCursor++
			}
		}
	case "tab":
		wiz.relFieldCursor = (wiz.relFieldCursor + 1) % 4
		if wiz.relFieldCursor == 3 && !wiz.relBidir {
			wiz.relFieldCursor = 0 // skip inverse if not bidirectional
		}
		if wiz.relFieldCursor == 3 {
			wiz.relInverseInput.Focus()
		} else {
			wiz.relInverseInput.Blur()
		}
	case "enter":
		if wiz.relFieldCursor == 1 { // toggle multiple
			wiz.relMultiple = !wiz.relMultiple
		} else if wiz.relFieldCursor == 2 { // toggle bidirectional
			wiz.relBidir = !wiz.relBidir
		} else if wiz.relFieldCursor == 0 || wiz.relFieldCursor == 3 {
			te.finishWizard()
		}
	case "esc":
		wiz.step = wizStepType
		wiz.relInverseInput.Blur()
	case " ":
		if wiz.relFieldCursor == 1 {
			wiz.relMultiple = !wiz.relMultiple
		} else if wiz.relFieldCursor == 2 {
			wiz.relBidir = !wiz.relBidir
		}
	default:
		if wiz.relFieldCursor == 3 {
			var cmd tea.Cmd
			wiz.relInverseInput, cmd = wiz.relInverseInput.Update(msg)
			return te, cmd
		}
	}
	return te, nil
}

func (te *typeEditor) finishWizard() {
	wiz := te.wizard
	prop := core.Property{
		Name: wiz.propName,
		Type: wiz.propType,
	}

	switch wiz.propType {
	case "select", "multi_select":
		raw := strings.TrimSpace(wiz.optionsInput.Value())
		if raw != "" {
			for _, v := range strings.Split(raw, ",") {
				v = strings.TrimSpace(v)
				if v != "" {
					prop.Options = append(prop.Options, core.Option{Value: v})
				}
			}
		}
	case "relation":
		if len(wiz.relTargets) > 0 && wiz.relTargetCursor < len(wiz.relTargets) {
			prop.Target = wiz.relTargets[wiz.relTargetCursor]
		}
		prop.Multiple = wiz.relMultiple
		prop.Bidirectional = wiz.relBidir
		if wiz.relBidir {
			prop.Inverse = strings.TrimSpace(wiz.relInverseInput.Value())
		}
	}

	te.schema.Properties = append(te.schema.Properties, prop)
	te.save()
	te.cancelWizard()
}

func (te *typeEditor) cancelWizard() {
	te.wizard = nil
	te.mode = teModeView
	te.saveErr = ""
}

// ── Template Management ─────────────────────────────────────────────────────

func (te *typeEditor) startAddTemplate() {
	te.tmplNameInput.SetValue("")
	te.tmplNameInput.Focus()
	te.mode = teModeAddTemplate
	te.saveErr = ""
}

func (te *typeEditor) updateAddTemplate(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(te.tmplNameInput.Value())
		if name == "" {
			return te, nil
		}
		// Check for duplicate
		for _, t := range te.templates {
			if t == name {
				te.saveErr = fmt.Sprintf("template %q already exists", name)
				return te, nil
			}
		}
		// Create empty template
		if te.vault != nil {
			tmpl := &core.Template{
				Name:       name,
				Properties: make(map[string]any),
			}
			if err := te.vault.SaveTemplate(te.typeName, name, tmpl); err != nil {
				te.saveErr = err.Error()
				return te, nil
			}
		}
		te.saveErr = ""
		te.refreshTemplates()
		te.tmplNameInput.Blur()
		te.mode = teModeView
		return te, nil
	case "esc":
		te.tmplNameInput.Blur()
		te.mode = teModeView
		te.saveErr = ""
		return te, nil
	}
	var cmd tea.Cmd
	te.tmplNameInput, cmd = te.tmplNameInput.Update(msg)
	return te, cmd
}

func (te *typeEditor) refreshTemplates() {
	if te.vault != nil {
		te.templates, _ = te.vault.ListTemplates(te.typeName)
	}
}

func (te *typeEditor) refreshViews() {
	if te.vault != nil {
		views, _ := te.vault.ListViews(te.typeName)
		te.views = make([]string, len(views))
		for i, v := range views {
			te.views[i] = v.Name
		}
	}
}

func (te *typeEditor) startAddView() {
	te.viewNameInput.SetValue("")
	te.viewNameInput.Focus()
	te.mode = teModeAddView
	te.saveErr = ""
}

func (te *typeEditor) updateAddView(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(te.viewNameInput.Value())
		if name == "" {
			return te, nil
		}
		// Check for duplicate
		for _, v := range te.views {
			if v == name {
				te.saveErr = fmt.Sprintf("view %q already exists", name)
				return te, nil
			}
		}
		// Create view with default config
		if te.vault != nil {
			view := &core.ViewConfig{
				Name:   name,
				Layout: core.ViewLayoutList,
				Sort:   []core.SortRule{{Property: "name", Direction: "asc"}},
			}
			if err := te.vault.SaveView(te.typeName, view); err != nil {
				te.saveErr = err.Error()
				return te, nil
			}
		}
		te.saveErr = ""
		te.refreshViews()
		te.viewNameInput.Blur()
		te.mode = teModeView
		return te, nil
	case "esc":
		te.viewNameInput.Blur()
		te.mode = teModeView
		te.saveErr = ""
		return te, nil
	}
	var cmd tea.Cmd
	te.viewNameInput, cmd = te.viewNameInput.Update(msg)
	return te, cmd
}

// View renders the type editor panel.
func (te *typeEditor) View() string {
	if te.mode == teModeAddWizard && te.wizard != nil {
		return te.viewWizard()
	}

	var b strings.Builder
	items := te.displayItems()
	pinned, unpinned := te.orderedProperties()
	lineNum := 0       // tracks current line number
	cursorLine := 0    // line where cursor item is rendered
	writeLine := func(s string) {
		b.WriteString(s + "\n")
		lineNum++
	}
	writeBlank := func() {
		b.WriteString("\n")
		lineNum++
	}

	// Meta fields
	metaLabels := []string{"Name", "Plural", "Emoji", "Color", "Unique", "Description"}
	emojiDisplay := te.schema.Emoji
	if emojiDisplay != "" {
		emojiDisplay = padEmoji(emojiDisplay)
	}
	metaValues := []string{
		te.schema.Name,
		te.schema.Plural,
		emojiDisplay,
		te.schema.Color,
		formatBool(te.schema.Unique),
		te.schema.Description,
	}



	// Helper to find cursor position for a sentinel/item value
	findCursorPos := func(target int) int {
		for i, item := range items {
			if item == target {
				return i
			}
		}
		return -1
	}

	for i := 0; i < metaFieldCount; i++ {
		if te.cursor == i {
			cursorLine = lineNum
		}
		if (te.mode == teModeEditMeta) && te.editField == i {
			writeLine(fmt.Sprintf("  %s: %s", metaLabels[i], te.editInput.View()))
		} else {
			val := metaValues[i]
			if val == "" {
				val = "(empty)"
			}
			lineContent := fmt.Sprintf("%s: %s", metaLabels[i], val)
			if te.cursor == i {
				writeLine(" " + highlightStyle.Render(" "+lineContent+" "))
			} else {
				writeLine("  " + lineContent)
			}
		}
	}

	writeBlank()

	// Pinned section — only shown when there are pinned properties
	if len(pinned) > 0 {
		writeLine(" ── Pinned (Header) ──")
		for _, idx := range pinned {
			pos := findCursorPos(idx)
			if te.cursor == pos {
				cursorLine = lineNum
			}
			te.renderPropertyRow(&b, items, idx)
			lineNum++
		}
		writeBlank()
	}

	// Properties section
	writeLine(" ── Properties ──")
	if len(unpinned) == 0 && len(pinned) > 0 {
		writeLine("  (none)")
	}
	for _, idx := range unpinned {
		pos := findCursorPos(idx)
		if te.cursor == pos {
			cursorLine = lineNum
		}
		te.renderPropertyRow(&b, items, idx)
		lineNum++
	}

	// "+ Add Property" row
	addPos := findCursorPos(addPropertySentinel)
	if te.cursor == addPos {
		cursorLine = lineNum
		writeLine(" " + highlightStyle.Render(" + Add Property "))
	} else {
		writeLine("  + Add Property")
	}

	// Templates section
	writeBlank()
	writeLine(" ── Templates ──")
	if len(te.templates) == 0 {
		writeLine("  (none)")
	} else {
		for tmplI, tmplName := range te.templates {
			pos := findCursorPos(templateSentinelBase - tmplI)
			if te.cursor == pos {
				cursorLine = lineNum
			}
			lineContent := fmt.Sprintf("📝 %s", tmplName)
			if te.cursor == pos {
				writeLine(" " + highlightStyle.Render(" "+lineContent+" "))
			} else {
				writeLine("  " + lineContent)
			}
		}
	}

	// "+ Add Template" row
	addTmplPos := findCursorPos(addTemplateSentinel)
	if te.mode == teModeAddTemplate {
		cursorLine = lineNum
		writeLine(fmt.Sprintf("  + %s", te.tmplNameInput.View()))
	} else if te.cursor == addTmplPos {
		cursorLine = lineNum
		writeLine(" " + highlightStyle.Render(" + Add Template "))
	} else {
		writeLine("  + Add Template")
	}

	// Views section
	writeBlank()
	writeLine(" ── Views ──")
	if len(te.views) == 0 {
		writeLine("  (default only)")
	} else {
		for viewI, viewName := range te.views {
			pos := findCursorPos(viewSentinelBase - viewI)
			if te.cursor == pos {
				cursorLine = lineNum
			}
			lineContent := fmt.Sprintf("🔍 %s", viewName)
			if te.cursor == pos {
				writeLine(" " + highlightStyle.Render(" "+lineContent+" "))
			} else {
				writeLine("  " + lineContent)
			}
		}
	}

	// "+ Add View" row
	addViewPos := findCursorPos(addViewSentinel)
	if te.mode == teModeAddView {
		cursorLine = lineNum
		writeLine(fmt.Sprintf("  + %s", te.viewNameInput.View()))
	} else if te.cursor == addViewPos {
		cursorLine = lineNum
		writeLine(" " + highlightStyle.Render(" + Add View "))
	} else {
		writeLine("  + Add View")
	}

	// Delete confirmation
	if te.mode == teModeDeleteProp {
		writeBlank()
		propIdx := items[te.cursor]
		if propIdx >= 0 && propIdx < len(te.schema.Properties) {
			writeLine(fmt.Sprintf(" Delete property '%s'? [y/n]", te.schema.Properties[propIdx].Name))
		}
	}

	if te.mode == teModeDeleteType {
		writeBlank()
		writeLine(fmt.Sprintf(" Delete type '%s'? [y/n]", te.typeName))
	}

	// Error
	if te.saveErr != "" {
		writeBlank()
		writeLine(fmt.Sprintf(" [ERROR] %s", te.saveErr))
	}

	return te.applyScroll(b.String(), cursorLine)
}

// applyScroll trims the rendered content to fit within the available height,
// keeping the cursor-highlighted line visible. cursorLine is the 0-based line
// index where the cursor item is rendered.
func (te *typeEditor) applyScroll(content string, cursorLine int) string {
	if te.height <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	// Remove trailing empty line from final \n
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	visibleH := te.height
	if len(lines) <= visibleH {
		return content
	}

	// Adjust scroll offset to keep cursor visible
	if cursorLine < te.scrollOffset {
		te.scrollOffset = cursorLine
	}
	if cursorLine >= te.scrollOffset+visibleH {
		te.scrollOffset = cursorLine - visibleH + 1
	}
	if te.scrollOffset < 0 {
		te.scrollOffset = 0
	}

	end := te.scrollOffset + visibleH
	if end > len(lines) {
		end = len(lines)
	}

	return strings.Join(lines[te.scrollOffset:end], "\n")
}

func (te *typeEditor) renderPropertyRow(b *strings.Builder, items []int, propIdx int) {

	p := te.schema.Properties[propIdx]

	// Find cursor position for this property
	cursorPos := -1
	for i, item := range items {
		if item == propIdx {
			cursorPos = i
			break
		}
	}

	isCurrent := te.cursor == cursorPos

	emoji := ""
	if p.Emoji != "" {
		emoji = " " + padEmoji(p.Emoji)
	}

	lineContent := fmt.Sprintf("%s  %s%s", p.Name, p.Type, emoji)
	if isCurrent {
		if te.mode == teModeMove {
			b.WriteString(" " + highlightStyle.Render("↕"+lineContent+" ") + "\n")
		} else {
			b.WriteString(" " + highlightStyle.Render(" "+lineContent+" ") + "\n")
		}
	} else {
		b.WriteString("  " + lineContent + "\n")
	}
}

func (te *typeEditor) viewWizard() string {
	wiz := te.wizard
	var b strings.Builder

	b.WriteString(" Add Property\n")
	b.WriteString(" ──────────────────────\n\n")

	switch wiz.step {
	case wizStepName:
		b.WriteString(fmt.Sprintf(" Step 1 of 3 — Property name\n\n"))
		b.WriteString(fmt.Sprintf(" Name: %s\n", wiz.nameInput.View()))
		if te.saveErr != "" {
			b.WriteString(fmt.Sprintf("\n [ERROR] %s\n", te.saveErr))
		}
		b.WriteString("\n enter: next  esc: cancel\n")

	case wizStepType:
		b.WriteString(fmt.Sprintf(" Step 2 of 3 — Property type for '%s'\n\n", wiz.propName))
		for i, t := range wiz.typeList {
			prefix := "  "
			if i == wiz.typeCursor {
				prefix = " ▸"
			}
			b.WriteString(fmt.Sprintf("%s %s\n", prefix, t))
		}
		b.WriteString("\n ↑↓: select  enter: next  esc: back\n")

	case wizStepOptions:
		b.WriteString(fmt.Sprintf(" Step 2b — Options for '%s' (%s)\n\n", wiz.propName, wiz.propType))
		b.WriteString(fmt.Sprintf(" Options (comma-separated): %s\n", wiz.optionsInput.View()))
		b.WriteString("\n enter: create  esc: back\n")

	case wizStepRelation:
		b.WriteString(fmt.Sprintf(" Step 3 of 3 — Relation config for '%s'\n\n", wiz.propName))

		// Target type
		b.WriteString(" Target type:\n")
		for i, t := range wiz.relTargets {
			prefix := "  "
			if i == wiz.relTargetCursor && wiz.relFieldCursor == 0 {
				prefix = " ▸"
			}
			b.WriteString(fmt.Sprintf("%s %s\n", prefix, t))
		}

		b.WriteString("\n")

		// Multiple toggle
		multiLabel := "no"
		if wiz.relMultiple {
			multiLabel = "yes"
		}
		prefix := "  "
		if wiz.relFieldCursor == 1 {
			prefix = " ▸"
		}
		b.WriteString(fmt.Sprintf("%s Multiple: %s\n", prefix, multiLabel))

		// Bidirectional toggle
		bidirLabel := "no"
		if wiz.relBidir {
			bidirLabel = "yes"
		}
		prefix = "  "
		if wiz.relFieldCursor == 2 {
			prefix = " ▸"
		}
		b.WriteString(fmt.Sprintf("%s Bidirectional: %s\n", prefix, bidirLabel))

		// Inverse name (only if bidirectional)
		if wiz.relBidir {
			prefix = "  "
			if wiz.relFieldCursor == 3 {
				prefix = " ▸"
			}
			b.WriteString(fmt.Sprintf("%s Inverse: %s\n", prefix, wiz.relInverseInput.View()))
		}

		b.WriteString("\n tab: next field  enter: confirm/create  esc: back\n")
	}

	return b.String()
}

// HelpBar returns the context-sensitive help text for the type editor.
func (te *typeEditor) HelpBar() string {
	switch te.mode {
	case teModeView:
		return "  [TYPE]  enter: open  e: edit  a: add  d: delete  m: move  p: pin  esc: back"
	case teModeEditMeta:
		return "  [EDIT]  enter: save  esc: cancel"
	case teModeEditProp:
		if te.propDetail != nil && te.propDetail.editing {
			return "  [EDIT]  enter: save  esc: cancel"
		}
		return "  [PROPERTY]  enter/e: edit  esc: back"
	case teModeMove:
		return "  [MOVE]  ↑↓: reorder  enter/esc: done"
	case teModeAddWizard:
		return "  [ADD PROPERTY]  follow prompts  esc: cancel/back"
	case teModeDeleteProp:
		return "  [DELETE]  y: confirm  n/esc: cancel"
	case teModeDeleteType:
		return "  [DELETE TYPE]  y: confirm  n/esc: cancel"
	case teModeAddTemplate:
		return "  [NEW TEMPLATE]  enter: create  esc: cancel"
	case teModeAddView:
		return "  [NEW VIEW]  enter: create  esc: cancel"
	}
	return ""
}

// CanQuit returns true when the editor is in a non-interactive state and the app can safely quit.
func (te *typeEditor) CanQuit() bool {
	return te.mode == teModeView
}

func formatBool(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
