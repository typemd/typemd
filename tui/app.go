package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	typeEditor  *typeEditor // non-nil when rightPanel == panelTypeEditor
	newTypeName  textinput.Model
	newTypeMode  bool // true when entering new type name in sidebar
	newObjectName textinput.Model
	newObjectMode bool   // true when entering new object name
	newObjectType string // type for the new object

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

// newBodyTextarea creates a configured textarea for body editing.
func newBodyTextarea() textarea.Model {
	ta := textarea.New()
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.Prompt = " " // single-space indent matching view mode; must be set before SetWidth
	// Remove textarea's own border — the panel border is provided by lipgloss
	noBase := lipgloss.NewStyle()
	s := ta.Styles()
	s.Focused.Base = noBase
	s.Blurred.Base = noBase
	ta.SetStyles(s)
	return ta
}

// titlePanelHeight is the total height of the title panel (1 content line + 2 border lines).
const titlePanelHeight = 3

// resizeBodyTextarea updates the body textarea dimensions to match the current layout.
func (m *model) resizeBodyTextarea() {
	h := m.height - 3 // help bar + borders
	if m.selected != nil {
		h -= titlePanelHeight // title panel takes vertical space
	}
	if h < 0 {
		h = 0
	}
	m.bodyTextarea.SetWidth(m.bodyWidth())
	m.bodyTextarea.SetHeight(h)
}

// selectedTypeEmoji returns the emoji for the currently selected object's type.
func (m model) selectedTypeEmoji() string {
	if m.selected == nil {
		return ""
	}
	for _, g := range m.groups {
		if g.Name == m.selected.Type {
			return g.Emoji
		}
	}
	return ""
}

func (m model) Init() tea.Cmd {
	if m.vault != nil {
		return watchObjects(m.vault.ObjectsDir())
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case m.newObjectMode:
			return updateNewObject(m, msg)
		case m.newTypeMode:
			return updateNewType(m, msg)
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

// refreshData syncs the index from disk and reloads all objects, preserving cursor position when possible.
func (m *model) refreshData() {
	if m.vault == nil {
		return
	}

	// Sync filesystem to DB first
	m.vault.SyncIndex()

	objects, err := m.vault.QueryObjects("")
	if err != nil {
		return
	}

	// Remember selected object ID to restore selection
	var selectedID string
	if m.selected != nil {
		selectedID = m.selected.ID
	}

	m.groups = buildGroups(objects, m.vault)
	m.searchResults = nil

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

// startNewObject enters the new object name input mode for a specific type.
func (m *model) startNewObject(groupIndex int) {
	if m.readOnly || groupIndex >= len(m.groups) {
		return
	}
	ti := textinput.New()
	ti.Placeholder = "object name"
	ti.CharLimit = 100
	ti.Focus()
	m.newObjectName = ti
	m.newObjectMode = true
	m.newObjectType = m.groups[groupIndex].Name
}

// startNewType enters the new type name input mode.
func (m *model) startNewType() {
	if m.readOnly {
		return
	}
	ti := textinput.New()
	ti.Placeholder = "type name"
	ti.CharLimit = 50
	ti.Focus()
	m.newTypeName = ti
	m.newTypeMode = true
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

// adjustScroll updates scrollOffset so cursor is always visible.
func (m *model) adjustScroll() {
	contentH := m.height - 3
	m.scrollOffset = adjustScrollOffset(m.cursor, m.scrollOffset, contentH)
}

// resizePanel adjusts the focused panel width by delta characters.
func (m *model) resizePanel(delta int) {
	switch m.focus {
	case focusLeft:
		m.leftW += delta
		if m.leftW < 20 {
			m.leftW = 20
		}
		if m.leftW > 50 {
			m.leftW = 50
		}
	case focusBody:
		// Body has no dedicated width field; grow body = shrink props
		if m.propsVisible {
			m.propsWidth -= delta
			if m.propsWidth < 20 {
				m.propsWidth = 20
			}
			if m.propsWidth > 40 {
				m.propsWidth = 40
			}
		} else {
			// Props hidden; grow body = shrink left
			m.leftW -= delta
			if m.leftW < 20 {
				m.leftW = 20
			}
			if m.leftW > 50 {
				m.leftW = 50
			}
		}
	case focusProps:
		m.propsWidth += delta
		if m.propsWidth < 20 {
			m.propsWidth = 20
		}
		if m.propsWidth > 40 {
			m.propsWidth = 40
		}
	}
	// Recalculate dependent widths
	m.bodyViewport.SetWidth(m.bodyWidth())
	m.propsViewport.SetWidth(m.propsWidth)
	m.bodyTextarea.SetWidth(m.bodyWidth())
	m.updateDetail()
}

// softWrapLines wraps each line individually, preserving leading indentation on continuation lines.
func softWrapLines(content string, width int) string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		if lipgloss.Width(line) <= width {
			result = append(result, line)
			continue
		}
		// Detect leading whitespace
		trimmed := strings.TrimLeft(line, " ")
		indent := line[:len(line)-len(trimmed)]
		wrapped := lipgloss.NewStyle().Width(width - lipgloss.Width(indent)).Render(trimmed)
		for i, wl := range strings.Split(wrapped, "\n") {
			if i == 0 {
				result = append(result, indent+wl)
			} else {
				result = append(result, indent+wl)
			}
		}
	}
	return strings.Join(result, "\n")
}

// updateDetail refreshes viewport contents with current selected object.
func (m *model) updateDetail() {
	bodyContent := renderBody(m.selected, m.bodyViewport.Width(), m.displayProps)
	if m.softWrap && m.bodyViewport.Width() > 0 {
		bodyContent = softWrapLines(bodyContent, m.bodyViewport.Width())
	}
	m.bodyViewport.SetContent(bodyContent)

	propsContent := renderProperties(m.selected, m.displayProps)
	if m.softWrap && m.propsViewport.Width() > 0 {
		propsContent = softWrapLines(propsContent, m.propsViewport.Width())
	}
	m.propsViewport.SetContent(propsContent)
}

// defaultLeftWidth calculates the default left panel width.
func (m model) defaultLeftWidth() int {
	w := m.width * 2 / 5
	if w < 20 {
		w = 20
	}
	if w > 50 {
		w = 50
	}
	return w
}

// leftWidth returns the current width for the left panel.
func (m model) leftWidth() int {
	if m.leftW > 0 {
		return m.leftW
	}
	return m.defaultLeftWidth()
}

// defaultPropsWidth calculates the default properties panel width.
func (m model) defaultPropsWidth() int {
	remaining := m.width - m.leftWidth() - 6 // 6 = borders for 3 panels
	w := remaining * 3 / 10                   // 30% of remaining
	if w < 20 {
		w = 20
	}
	if w > 40 {
		w = 40
	}
	return w
}

// hasTitlePanel returns true when the right side should show a title panel.
func (m model) hasTitlePanel() bool {
	return m.selected != nil || (m.rightPanel == panelTypeEditor && m.typeEditor != nil)
}

// bodyWidth calculates the body panel width from remaining space.
func (m model) bodyWidth() int {
	borders := 4 // left panel border (2) + body panel border (2)
	if m.propsVisible {
		borders += 2 // props panel border
	}
	w := m.width - m.leftWidth() - borders
	if m.propsVisible {
		w -= m.propsWidth
	}
	if w < 10 {
		w = 10
	}
	return w
}

// shouldAutoHideProps returns true if terminal is too narrow for three panels.
func (m model) shouldAutoHideProps() bool {
	minTotal := 20 + 10 + 20 + 6 // minLeft + minBody + minProps + borders
	return m.width < minTotal
}

func (m model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	// Help overlay takes over the entire screen
	if m.showHelp {
		v := tea.NewView(renderHelp(m.width, m.height, m.readOnly))
		v.AltScreen = true
		return v
	}

	leftW := m.leftWidth()
	bodyW := m.bodyWidth()
	// In lipgloss v2, Width()/Height() set the TOTAL size including border.
	// Internal widths (leftW, bodyW, contentH) remain content-area sizes;
	// we add the border size (+2) when passing to the panel style.
	contentH := m.height - 3 // content area: terminal minus help-bar minus borders
	if contentH < 0 {
		contentH = 0
	}
	bdr := 2 // left+right or top+bottom border size

	// When an object is selected, the title panel takes vertical space from body/props
	hasTitlePanel := m.hasTitlePanel()
	bodyPropsContentH := contentH
	if hasTitlePanel {
		bodyPropsContentH = contentH - titlePanelHeight
		if bodyPropsContentH < 0 {
			bodyPropsContentH = 0
		}
	}

	leftPanelH := contentH + bdr    // left panel spans full height
	bodyPropsPanelH := bodyPropsContentH + bdr // body/props panels are shorter when title exists

	// Styles — MaxHeight clamps viewport content that overflows after line wrapping.
	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(leftW + bdr).
		Height(leftPanelH).
		MaxHeight(leftPanelH)
	bodyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(bodyW + bdr).
		Height(bodyPropsPanelH).
		MaxHeight(bodyPropsPanelH)

	// Focus highlighting (edit mode uses distinct border color)
	activeBorderColor := colorFocusBorder
	if m.editMode {
		activeBorderColor = colorEditBorder
	}
	switch m.focus {
	case focusLeft:
		leftStyle = leftStyle.BorderForeground(activeBorderColor)
	case focusBody:
		bodyStyle = bodyStyle.BorderForeground(activeBorderColor)
	}

	// Left panel content
	var leftContent string
	if m.searchResults != nil {
		rows := searchResultRows(m.searchResults)
		if len(rows) == 0 {
			leftContent = "  (no results)"
		} else {
			var lines []string
			for i, row := range rows {
				line := fmt.Sprintf("   %s/%s", row.Object.Type, row.Object.GetName())
				if i == m.cursor {
					style := highlightStyle
					line = style.Render(line)
				}
				lines = append(lines, line)
			}
			leftContent = strings.Join(lines, "\n")
		}
	} else {
		leftContent = renderList(m.groups, m.cursor, m.scrollOffset, m.focus == focusLeft, leftW, contentH)
		if m.newObjectMode {
			leftContent += "\n New " + m.newObjectType + ": " + m.newObjectName.View()
		} else if m.newTypeMode {
			leftContent += "\n New type: " + m.newTypeName.View()
		}
	}

	var rightSide string

	if m.rightPanel == panelTypeEditor && m.typeEditor != nil {
		// Type editor uses full right-side width (no props panel)
		editorW := m.width - m.leftWidth() - 4 // left border + body border
		if editorW < 10 {
			editorW = 10
		}
		editorStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(editorW + bdr).
			Height(bodyPropsPanelH).
			MaxHeight(bodyPropsPanelH)
		if m.focus != focusLeft {
			editorStyle = editorStyle.BorderForeground(activeBorderColor)
		}

		// Title panel for type editor
		titleW := m.width - leftW - bdr
		titleStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(titleW).
			Height(titlePanelHeight)
		te := m.typeEditor
		emojiPrefix := ""
		if te.schema.Emoji != "" {
			emojiPrefix = padEmoji(te.schema.Emoji) + " "
		}
		titleText := fmt.Sprintf(" %s%s", emojiPrefix, te.typeName)
		titleContent := titleStyle.Render(titleText)

		// Adjust editor panel height for title
		editorH := bodyPropsPanelH
		editorStyle = editorStyle.Height(editorH).MaxHeight(editorH)

		rightSide = lipgloss.JoinVertical(lipgloss.Left,
			titleContent,
			editorStyle.Render(te.View()),
		)
	} else {
		// Object detail view (existing behavior)
		var bodyPanelContent string
		if m.editMode && m.focus == focusBody {
			bodyPanelContent = m.bodyTextarea.View()
		} else {
			bodyPanelContent = m.bodyViewport.View()
		}

		rightSide = bodyStyle.Render(bodyPanelContent)

		// Properties panel (optional)
		if m.propsVisible {
			propsStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(m.propsWidth + bdr).
				Height(bodyPropsPanelH).
				MaxHeight(bodyPropsPanelH)
			if m.focus == focusProps {
				propsStyle = propsStyle.BorderForeground(activeBorderColor)
			}
			rightSide = lipgloss.JoinHorizontal(lipgloss.Top,
				rightSide,
				propsStyle.Render(m.propsViewport.View()),
			)
		}

		// Title panel above body+props
		if hasTitlePanel {
			titleW := m.width - leftW - bdr
			titleStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(titleW).
				Height(titlePanelHeight)
			titleContent := renderTitleContent(m.selected, m.selected.Type, m.selectedTypeEmoji(), titleW-bdr)
			rightSide = lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render(titleContent),
				rightSide,
			)
		}
	}

	// Compose left + right
	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent),
		rightSide,
	)

	// Help bar
	var helpBar string
	if m.newObjectMode {
		helpBar = "  [NEW OBJECT]  enter: create  esc: cancel"
	} else if m.newTypeMode {
		helpBar = "  [NEW TYPE]  enter: create  esc: cancel"
	} else if m.searchMode {
		helpBar = "  / " + m.searchInput.View()
	} else if m.rightPanel == panelTypeEditor && m.typeEditor != nil && m.focus != focusLeft {
		helpBar = m.typeEditor.HelpBar()
	} else if m.saveConflict {
		helpBar = "  [CONFLICT]  " + m.saveErr
	} else if m.saveErr != "" {
		helpBar = "  [ERROR]  " + m.saveErr
	} else if m.editMode {
		helpBar = "  [EDIT]  esc: exit edit mode"
	} else {
		modeLabel := "VIEW"
		if m.readOnly {
			modeLabel = "READONLY"
		}
		if m.searchResults != nil {
			helpBar = fmt.Sprintf("  [%s]  Search results  |  esc: clear  |  ↑↓: navigate  |  tab: switch  |  q: quit", modeLabel)
		} else {
			helpBar = fmt.Sprintf("  [%s]  ?/h: help  |  /: search  |  q: quit", modeLabel)
		}
	}

	screen := panels + "\n" + helpBar

	// Overlay popup if type editor has one active
	if m.rightPanel == panelTypeEditor && m.typeEditor != nil {
		if overlay := m.typeEditor.Overlay(m.width, m.height); overlay != "" {
			screen = overlay
		}
	}

	v := tea.NewView(screen)
	v.AltScreen = true
	return v
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

	objects, err := v.QueryObjects("")
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
	if selected != nil {
		initialRightPanel = panelObject
	} else if selectedID == "" && savedState.SelectedTypeName != "" {
		// Cursor on a type header — open type editor
		if ts, err := v.LoadType(savedState.SelectedTypeName); err == nil {
			initialTypeEditor = newTypeEditor(ts, savedState.SelectedTypeName, false, v)
			initialRightPanel = panelTypeEditor
		}
	}

	m := model{
		vault:         v,
		focus:         focusLeft, // always start with focus on sidebar
		rightPanel:    initialRightPanel,
		typeEditor:    initialTypeEditor,
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
