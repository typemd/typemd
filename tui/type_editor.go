package tui

import (
	"sort"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
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
