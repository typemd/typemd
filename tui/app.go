package tui

import (
	"fmt"
	"os"
	"time"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type focusPanel int

const (
	focusLeft focusPanel = iota
	focusBody
	focusProps
)

type rightPanelMode int

const (
	panelEmpty      rightPanelMode = iota // no content selected
	panelObject                           // object detail view (existing behavior)
	panelTypeEditor                       // type editor view
	panelTemplate                         // template detail view
	panelView                             // full-width view mode
)

type typeGroup struct {
	Name     string
	Plural   string
	Emoji    string
	Objects  []*core.Object
	Expanded bool
}

type model struct {
	vault *core.Vault
	focus focusPanel

	// Right panel mode
	rightPanel  rightPanelMode
	typeEditor  *typeEditor          // non-nil when rightPanel == panelTypeEditor
	tmplEditor  *templateEditor      // non-nil when rightPanel == panelTemplate
	viewMode    *viewMode            // non-nil when rightPanel == panelView
	viewPicker  *viewPicker          // non-nil when view selection popup is active
	createType   *createTypeState    // non-nil when type creation flow is active
	create       *createState        // non-nil when object creation flow is active

	// Left panel
	groups       []typeGroup
	cursor       int
	scrollOffset int
	selected     *core.Object
	leftW        int // adjustable width for left panel (0 = use default)

	// Body panel (center)
	bodyViewport  viewport.Model
	bodyTextarea  textarea.Model

	// Properties panel (right)
	propsViewport viewport.Model
	propsWidth    int  // adjustable width for properties panel
	propsVisible  bool // toggle visibility

	// Shared detail state
	displayProps []core.DisplayProperty

	// Edit mode
	editMode      bool
	bodyEditStart string // textarea.Value() snapshot taken at edit entry (sanitized)

	// Save state
	dirty          bool      // unsaved in-memory changes
	saveErr        string    // last save error (shown in status bar)
	skipNextReload bool      // suppress next fileChangedMsg (triggered by our own save)
	loadedModTime  time.Time // file mtime when object was last loaded
	saveConflict   bool      // concurrent external edit detected; awaiting user decision

	// Search
	searchMode    bool
	searchInput   textinput.Model
	searchResults []*core.Object

	// Help
	showHelp bool

	// Settings
	readOnly bool
	softWrap bool

	// Layout
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	if m.vault != nil {
		return watchObjects(m.vault.ObjectsDir())
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// View picker needs to receive ALL message types (including huh internal msgs)
	if m.viewPicker != nil {
		vp, cmd := m.viewPicker.Update(msg)
		m.viewPicker = vp
		return m, cmd
	}

	switch msg := msg.(type) {
	case focusLeftMsg:
		m.focus = focusLeft
		return m, nil

	case typeDeletedMsg:
		m.typeEditor = nil
		m.rightPanel = panelEmpty
		m.focus = focusLeft
		m.refreshData()
		return m, nil

	case openTemplateMsg:
		// Transition from type editor to template detail view
		if m.vault != nil {
			tmpl, err := m.vault.LoadTemplate(msg.TypeName, msg.TemplateName)
			if err != nil {
				return m, nil
			}
			schema, _ := m.vault.LoadType(msg.TypeName)
			te := newTemplateEditor(msg.TypeName, msg.TemplateName, tmpl, schema, m.vault)
			// Size viewports so content is visible
			contentH := m.height - 3 - titlePanelHeight
			if contentH < 0 {
				contentH = 0
			}
			editorW := m.width - m.leftWidth() - 4
			if editorW < 10 {
				editorW = 10
			}
			te.SetSize(editorW, contentH, m.defaultPropsWidth(), true)
			m.tmplEditor = te
			m.rightPanel = panelTemplate
			m.focus = focusBody
		}
		return m, nil

	case templateDeletedMsg:
		// Template already deleted by templateEditor; clean up and return to type editor
		m.tmplEditor = nil
		m.rightPanel = panelTypeEditor
		m.focus = focusBody
		if m.typeEditor != nil {
			m.typeEditor.refreshTemplates()
		}
		return m, nil

	case viewEditorDeletedMsg:
		// View was deleted by the editor; exit view mode
		if m.vault != nil {
			_ = m.vault.DeleteView(msg.TypeName, msg.ViewName)
		}
		m.viewMode = nil
		m.rightPanel = panelTypeEditor
		m.focus = focusBody
		if m.typeEditor != nil {
			m.typeEditor.refreshViews()
		}
		return m, nil

	case openViewMsg:
		// Transition to full-width view mode
		if m.vault != nil {
			vm := newViewMode(msg.TypeName, msg.ViewName, m.vault)
			vm.SetSize(m.width-2, m.height-3)
			m.viewMode = vm
			m.rightPanel = panelView
			m.focus = focusBody
		}
		return m, nil

	case flashDismissMsg:
		if m.create != nil && msg.seq == m.create.flashSeq {
			m.create.flash = ""
		}
		return m, nil

	case fileChangedMsg:
		if m.skipNextReload {
			m.skipNextReload = false
			return m, watchObjects(m.vault.ObjectsDir())
		}
		m.refreshData()
		return m, watchObjects(m.vault.ObjectsDir())

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentHeight := m.height - 3 // help bar + borders
		if contentHeight < 0 {
			contentHeight = 0
		}

		// Initialize panel widths if not set
		if m.leftW == 0 {
			m.leftW = m.defaultLeftWidth()
		}
		if m.propsWidth == 0 {
			m.propsWidth = m.defaultPropsWidth()
		}

		// Auto-hide on narrow terminals
		if m.shouldAutoHideProps() {
			m.propsVisible = false
		}

		// Body/props panels are shorter when title panel is shown
		hasTitlePanel := m.hasTitlePanel()
		bodyPropsH := contentHeight
		if hasTitlePanel {
			bodyPropsH -= titlePanelHeight
			if bodyPropsH < 0 {
				bodyPropsH = 0
			}
		}

		// Update viewport sizes
		m.bodyViewport.SetWidth(m.bodyWidth())
		m.bodyViewport.SetHeight(bodyPropsH)
		m.propsViewport.SetWidth(m.propsWidth)
		m.propsViewport.SetHeight(bodyPropsH)
		m.resizeBodyTextarea()

		// Resize template editor if active (always show props for templates)
		if m.tmplEditor != nil {
			editorW := m.width - m.leftWidth() - 4
			if editorW < 10 {
				editorW = 10
			}
			m.tmplEditor.SetSize(editorW, bodyPropsH, m.defaultPropsWidth(), true)
		}

		m.updateDetail()
		return m, nil

	case tea.KeyPressMsg:
		// Mode priority: help > search > conflict > typeEditor > edit > normal
		switch {
		case m.showHelp:
			return updateHelp(m, msg)
		case m.searchMode:
			var cmd tea.Cmd
			m, cmd = updateSearch(m, msg)
			if !m.searchMode && m.searchResults != nil {
				m.selectCurrentRow()
			}
			return m, cmd
		case m.saveConflict:
			return updateConflict(m, msg)
		case m.create != nil:
			return updateCreate(m, msg)
		case m.createType != nil:
			return updateCreateType(m, msg)
		case m.rightPanel == panelView && m.viewMode != nil:
			// q/ctrl+c quits globally
			if (msg.String() == "q" || msg.String() == "ctrl+c") && m.viewMode.CanQuit() {
				if m.vault != nil {
					saveSessionState(m.vault.Root, m.captureState())
				}
				return m, tea.Quit
			}
			// Esc navigates back: editor → detail → list → sidebar
			if msg.String() == "esc" {
				if m.viewMode.HasEditor() {
					// Let viewMode.Update handle closing the editor
					vm, cmd := m.viewMode.Update(msg)
					m.viewMode = vm
					return m, cmd
				}
				if m.viewMode.detailObject != nil {
					// Return from object detail to view list
					m.viewMode.detailObject = nil
					m.selected = nil
					m.rightPanel = panelView
					return m, nil
				}
				m.viewMode = nil
				m.rightPanel = panelEmpty
				m.focus = focusLeft
				m.selectCurrentRow()
				return m, nil
			}
			vm, cmd := m.viewMode.Update(msg)
			m.viewMode = vm
			// If an object was selected, load it into the normal detail view
			if vm != nil && vm.detailObject != nil && m.selected != vm.detailObject {
				obj, err := m.vault.GetObject(vm.detailObject.ID)
				if err == nil {
					m.applyLoadedObject(obj)
					m.rightPanel = panelObject
				}
			}
			return m, cmd
		case m.rightPanel == panelTemplate && m.tmplEditor != nil && m.focus != focusLeft:
			// q/ctrl+c quits globally unless in an interactive mode
			if (msg.String() == "q" || msg.String() == "ctrl+c") && m.tmplEditor.CanQuit() {
				if m.vault != nil {
					saveSessionState(m.vault.Root, m.captureState())
				}
				return m, tea.Quit
			}
			// Esc in view mode returns to type editor
			if msg.String() == "esc" && m.tmplEditor.CanQuit() {
				m.tmplEditor = nil
				m.rightPanel = panelTypeEditor
				m.focus = focusBody
				return m, nil
			}
			te, cmd := m.tmplEditor.Update(msg)
			m.tmplEditor = te
			return m, cmd
		case m.rightPanel == panelTypeEditor && m.typeEditor != nil && m.focus != focusLeft:
			// q/ctrl+c quits globally unless in an interactive mode
			if (msg.String() == "q" || msg.String() == "ctrl+c") && m.typeEditor.CanQuit() {
				if m.vault != nil {
					saveSessionState(m.vault.Root, m.captureState())
				}
				return m, tea.Quit
			}
			te, cmd := m.typeEditor.Update(msg)
			m.typeEditor = te
			m.syncTypeGroupMeta(te.typeName, te.schema)
			return m, cmd
		case m.editMode:
			return updateEdit(m, msg)
		default:
			return updateNormal(m, msg)
		}
	}
	// Route remaining messages (e.g. cursor blink) to textarea when in body edit mode
	if m.editMode && m.focus == focusBody {
		var cmd tea.Cmd
		m.bodyTextarea, cmd = m.bodyTextarea.Update(msg)
		return m, cmd
	}
	return m, nil
}

// rebuildGroups reloads objects and rebuilds type groups from the vault.
func (m *model) rebuildGroups() {
	if m.vault == nil {
		return
	}
	objects, err := m.vault.QueryObjects(nil)
	if err != nil {
		return
	}
	m.groups = buildGroups(objects, m.vault)
	m.searchResults = nil
}

// refreshData syncs the index from disk and reloads all objects, preserving cursor position when possible.
func (m *model) refreshData() {
	if m.vault == nil {
		return
	}

	// Sync filesystem to DB first
	m.vault.SyncIndex()

	m.rebuildGroups()

	// Remember selected object ID to restore selection
	var selectedID string
	if m.selected != nil {
		selectedID = m.selected.ID
	}

	// Try to restore cursor to previously selected object
	rows := visibleRows(m.groups)
	m.cursor = 0
	for i, row := range rows {
		if row.Kind == rowObject && row.Object != nil && row.Object.ID == selectedID {
			m.cursor = i
			break
		}
	}

	m.selectCurrentRow()
}

// currentRows returns the appropriate rows based on whether search results are active.
// currentTypeName returns the type name at the cursor position in the sidebar.
func (m *model) currentTypeName() string {
	rows := m.currentRows()
	if m.cursor < 0 || m.cursor >= len(rows) {
		return ""
	}
	row := rows[m.cursor]
	if row.Kind == rowNewType {
		return ""
	}
	if row.GroupIndex >= 0 && row.GroupIndex < len(m.groups) {
		return m.groups[row.GroupIndex].Name
	}
	return ""
}

func (m *model) currentRows() []listRow {
	if m.searchResults != nil {
		return searchResultRows(m.searchResults)
	}
	return visibleRows(m.groups)
}

// refreshLoadedModTime updates loadedModTime from the file's current mtime.
func (m *model) refreshLoadedModTime(obj *core.Object) {
	objPath := m.vault.ObjectPath(obj.Type, obj.Filename)
	if info, err := os.Stat(objPath); err == nil {
		m.loadedModTime = info.ModTime()
	}
}

// applyLoadedObject sets the selected object and updates displayProps and loadedModTime.
// Called after a successful GetObject to avoid duplicating this pattern.
func (m *model) applyLoadedObject(obj *core.Object) {
	m.selected = obj
	m.displayProps, _ = m.vault.BuildDisplayProperties(obj)
	m.refreshLoadedModTime(obj)
}

// selectCurrentRow updates the selected object based on current cursor position.
// Re-reads the object from disk to get the latest body and properties.
func (m *model) selectCurrentRow() {
	rows := m.currentRows()
	if m.cursor >= 0 && m.cursor < len(rows) {
		row := rows[m.cursor]
		switch row.Kind {
		case rowObject:
			if row.Object != nil {
				if m.vault != nil {
					if obj, err := m.vault.GetObject(row.Object.ID); err == nil {
						m.applyLoadedObject(obj)
					} else {
						m.selected = row.Object
						m.displayProps = nil
					}
				} else {
					m.selected = row.Object
					m.displayProps = nil
				}
				m.rightPanel = panelObject
				m.typeEditor = nil
				m.dirty = false
				m.saveErr = ""
				m.saveConflict = false
				m.updateDetail()
			}
		case rowHeader:
			if m.vault != nil {
				g := m.groups[row.GroupIndex]
				if ts, err := m.vault.LoadType(g.Name); err == nil {
					m.typeEditor = newTypeEditor(ts, g.Name, false, m.vault)
					m.rightPanel = panelTypeEditor
					m.selected = nil
				}
			}
		}
	}
}

// doSave executes the actual vault write and resets save state on success.
// Shared by saveObject and forceSave.
func (m *model) doSave() {
	if m.readOnly {
		return
	}
	if err := m.vault.SaveObject(m.selected); err != nil {
		m.saveErr = fmt.Sprintf("Save failed: %v", err)
		return
	}
	// Update loadedModTime so subsequent saves don't trigger a false conflict.
	m.refreshLoadedModTime(m.selected)
	m.dirty = false
	m.saveErr = ""
	m.saveConflict = false
	m.skipNextReload = true
}

// saveObject attempts to save the selected object to disk.
// Sets saveConflict if a concurrent external edit is detected.
// Sets saveErr on failure. On success, clears dirty and sets skipNextReload.
func (m *model) saveObject() {
	if m.selected == nil || m.vault == nil {
		return
	}
	objPath := m.vault.ObjectPath(m.selected.Type, m.selected.Filename)
	if info, err := os.Stat(objPath); err == nil {
		if info.ModTime().After(m.loadedModTime) {
			m.saveConflict = true
			m.saveErr = "File changed externally. 'y' to overwrite · 'n' to reload · esc to cancel"
			return
		}
	}
	m.doSave()
}

// forceSave saves the selected object ignoring concurrent edit detection.
func (m *model) forceSave() {
	if m.selected == nil || m.vault == nil {
		return
	}
	m.doSave()
}

// reloadFromDisk reloads the selected object from disk, discarding local changes.
func (m *model) reloadFromDisk() {
	if m.selected == nil || m.vault == nil {
		return
	}
	if obj, err := m.vault.GetObject(m.selected.ID); err == nil {
		m.applyLoadedObject(obj)
		m.updateDetail()
	}
	m.dirty = false
	m.saveErr = ""
	m.saveConflict = false
}

func Start(vaultPath string, readOnly bool, reindex bool) error {
	if vaultPath == "" {
		var err error
		vaultPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
	}

	v := core.NewVault(vaultPath)
	loadTheme(vaultPath)
	if err := v.Open(); err != nil {
		return fmt.Errorf("open vault: %w", err)
	}
	defer v.Close()

	if reindex {
		if _, err := v.SyncIndex(); err != nil {
			return fmt.Errorf("reindex: %w", err)
		}
	}

	objects, err := v.QueryObjects(nil)
	if err != nil {
		return fmt.Errorf("query objects: %w", err)
	}

	groups := buildGroups(objects, v)

	// Load saved session state (or zero value if missing/corrupt)
	savedState := loadSessionState(vaultPath)

	// Apply saved state to groups and resolve cursor position
	initialCursor, selectedID := applySessionState(savedState, groups)

	// Select the resolved object and capture its mtime for conflict detection
	var selected *core.Object
	var displayProps []core.DisplayProperty
	var initialModTime time.Time
	if selectedID != "" {
		if obj, err := v.GetObject(selectedID); err == nil {
			selected = obj
			displayProps, _ = v.BuildDisplayProperties(selected)
			objPath := v.ObjectPath(selected.Type, selected.Filename)
			if info, err := os.Stat(objPath); err == nil {
				initialModTime = info.ModTime()
			}
		}
	}

	bodyVP := viewport.New()
	bodyVP.SetContent(renderBody(selected, 0, displayProps))
	propsVP := viewport.New()
	propsVP.SetContent(renderProperties(selected, displayProps))

	bodyTA := newBodyTextarea()

	// Note: focus is always reset to focusLeft on startup for consistent UX
	_ = savedState.Focus

	// Restore panel widths (0 = use default, applied later on WindowSizeMsg)
	leftW := savedState.LeftPanelWidth
	propsWidth := savedState.PropsPanelWidth
	propsVisible := savedState.PropsVisible

	// Determine initial right panel mode and type editor
	var initialRightPanel rightPanelMode
	var initialTypeEditor *typeEditor
	var initialViewMode *viewMode

	// Try to restore view mode first (takes precedence over sidebar selection)
	if vm := restoreViewMode(savedState, v); vm != nil {
		initialRightPanel = panelView
		initialViewMode = vm
	} else if selected != nil {
		initialRightPanel = panelObject
	} else if selectedID == "" && savedState.SelectedTypeName != "" {
		// Cursor on a type header — open type editor
		if ts, err := v.LoadType(savedState.SelectedTypeName); err == nil {
			initialTypeEditor = newTypeEditor(ts, savedState.SelectedTypeName, false, v)
			initialRightPanel = panelTypeEditor
		}
	}

	initialFocus := focusLeft
	if initialViewMode != nil {
		initialFocus = focusBody
	}

	m := model{
		vault:         v,
		focus:         initialFocus,
		rightPanel:    initialRightPanel,
		typeEditor:    initialTypeEditor,
		viewMode:      initialViewMode,
		groups:        groups,
		cursor:        initialCursor,
		scrollOffset:  savedState.ScrollOffset,
		selected:      selected,
		bodyViewport:  bodyVP,
		bodyTextarea:  bodyTA,
		propsViewport: propsVP,
		leftW:         leftW,
		propsWidth:    propsWidth,
		propsVisible:  propsVisible,
		readOnly:      readOnly,
		softWrap:      true,
		displayProps:  displayProps,
		loadedModTime: initialModTime,
		searchInput:   initSearchInput(),
	}

	p := tea.NewProgram(m)
	_, err = p.Run()
	return err
}
