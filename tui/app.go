package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/typemd/typemd/core"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusPanel int

const (
	focusLeft focusPanel = iota
	focusBody
	focusProps
)

type typeGroup struct {
	Name     string
	Objects  []*core.Object
	Expanded bool
}

type model struct {
	vault *core.Vault
	focus focusPanel

	// Left panel
	groups       []typeGroup
	cursor       int
	scrollOffset int
	selected     *core.Object
	leftW        int // adjustable width for left panel (0 = use default)

	// Body panel (center)
	bodyViewport viewport.Model

	// Properties panel (right)
	propsViewport viewport.Model
	propsWidth    int  // adjustable width for properties panel
	propsVisible  bool // toggle visibility

	// Shared detail state
	displayProps []core.DisplayProperty

	// Edit mode
	editMode bool

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
	switch msg := msg.(type) {
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

		// Update viewport sizes
		m.bodyViewport.Width = m.bodyWidth()
		m.bodyViewport.Height = contentHeight
		m.propsViewport.Width = m.propsWidth
		m.propsViewport.Height = contentHeight

		m.updateDetail()
		return m, nil

	case tea.KeyMsg:
		// Help mode gets top priority
		if m.showHelp {
			switch msg.String() {
			case "esc", "?", "h":
				m.showHelp = false
			}
			return m, nil
		}

		// Search mode gets priority
		if m.searchMode {
			var cmd tea.Cmd
			m, cmd = updateSearch(m, msg)
			if !m.searchMode && m.searchResults != nil {
				// Search completed, select first result if available
				m.selectCurrentRow()
			}
			return m, cmd
		}

		// Conflict resolution intercepts y/n/esc
		if m.saveConflict {
			switch msg.String() {
			case "y":
				m.forceSave()
			case "n":
				m.reloadFromDisk()
			case "esc":
				m.saveConflict = false
				m.saveErr = ""
			}
			return m, nil
		}

		// Edit mode intercepts all keys except Esc
		if m.editMode {
			if msg.String() == "esc" {
				m.editMode = false
				if m.dirty {
					m.saveObject()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "/":
			m.searchMode = true
			m.searchInput.Focus()
			return m, textinput.Blink

		case "e":
			if m.focus == focusBody || m.focus == focusProps {
				m.editMode = true
			}
			return m, nil

		case "tab":
			switch m.focus {
			case focusLeft:
				m.focus = focusBody
			case focusBody:
				if m.propsVisible {
					m.focus = focusProps
				} else {
					m.focus = focusLeft
				}
			case focusProps:
				m.focus = focusLeft
			}
			return m, nil

		case "w":
			m.softWrap = !m.softWrap
			m.updateDetail()
			return m, nil

		case "esc":
			// Clear search results and return to normal list
			if m.searchResults != nil {
				m.searchResults = nil
				m.cursor = 0
				m.selectCurrentRow()
				return m, nil
			}

		case "up", "k":
			if m.focus == focusLeft {
				rows := m.currentRows()
				m.cursor = clampCursor(m.cursor-1, len(rows))
				m.adjustScroll()
				m.selectCurrentRow()
			} else if m.focus == focusBody {
				m.bodyViewport.LineUp(1)
			} else if m.focus == focusProps {
				m.propsViewport.LineUp(1)
			}
			return m, nil

		case "down", "j":
			if m.focus == focusLeft {
				rows := m.currentRows()
				m.cursor = clampCursor(m.cursor+1, len(rows))
				m.adjustScroll()
				m.selectCurrentRow()
			} else if m.focus == focusBody {
				m.bodyViewport.LineDown(1)
			} else if m.focus == focusProps {
				m.propsViewport.LineDown(1)
			}
			return m, nil

		case "]":
			m.resizePanel(+2)
			return m, nil

		case "[":
			m.resizePanel(-2)
			return m, nil

		case "p":
			m.propsVisible = !m.propsVisible
			if !m.propsVisible && m.focus == focusProps {
				m.focus = focusBody
			}
			// Recalculate widths for both panels
			contentHeight := m.height - 3
			if contentHeight < 0 {
				contentHeight = 0
			}
			m.bodyViewport.Width = m.bodyWidth()
			m.propsViewport.Width = m.propsWidth
			m.propsViewport.Height = contentHeight
			m.updateDetail()
			return m, nil

		case "?", "h":
			m.showHelp = true
			return m, nil

		case "enter", " ":
			if m.focus == focusLeft {
				rows := m.currentRows()
				if m.cursor >= 0 && m.cursor < len(rows) {
					row := rows[m.cursor]
					if row.IsHeader {
						m.groups[row.GroupIndex].Expanded = !m.groups[row.GroupIndex].Expanded
						// Re-clamp cursor after collapse
						newRows := m.currentRows()
						m.cursor = clampCursor(m.cursor, len(newRows))
						m.adjustScroll()
					}
					m.selectCurrentRow()
				}
			}
			return m, nil
		}
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

	m.groups = buildGroups(objects)
	m.searchResults = nil

	// Try to restore cursor to previously selected object
	rows := visibleRows(m.groups)
	m.cursor = 0
	for i, row := range rows {
		if !row.IsHeader && row.Object != nil && row.Object.ID == selectedID {
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

// applyLoadedObject sets the selected object and updates displayProps and loadedModTime.
// Called after a successful GetObject to avoid duplicating this pattern.
func (m *model) applyLoadedObject(obj *core.Object) {
	m.selected = obj
	m.displayProps, _ = m.vault.BuildDisplayProperties(obj)
	objPath := m.vault.ObjectPath(obj.Type, obj.Filename)
	if info, statErr := os.Stat(objPath); statErr == nil {
		m.loadedModTime = info.ModTime()
	}
}

// selectCurrentRow updates the selected object based on current cursor position.
// Re-reads the object from disk to get the latest body and properties.
func (m *model) selectCurrentRow() {
	rows := m.currentRows()
	if m.cursor >= 0 && m.cursor < len(rows) {
		row := rows[m.cursor]
		if !row.IsHeader && row.Object != nil {
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
			m.dirty = false
			m.saveErr = ""
			m.saveConflict = false
			m.updateDetail()
			return
		}
	}
}

// doSave executes the actual vault write and resets save state on success.
// Shared by saveObject and forceSave.
func (m *model) doSave() {
	if err := m.vault.SaveObject(m.selected); err != nil {
		m.saveErr = fmt.Sprintf("Save failed: %v", err)
		return
	}
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
	m.bodyViewport.Width = m.bodyWidth()
	m.propsViewport.Width = m.propsWidth
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
	bodyContent := renderBody(m.selected, m.bodyViewport.Width)
	if m.softWrap && m.bodyViewport.Width > 0 {
		bodyContent = softWrapLines(bodyContent, m.bodyViewport.Width)
	}
	m.bodyViewport.SetContent(bodyContent)

	propsContent := renderProperties(m.selected, m.displayProps)
	if m.softWrap && m.propsViewport.Width > 0 {
		propsContent = softWrapLines(propsContent, m.propsViewport.Width)
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

// bodyWidth calculates the body panel width from remaining space.
func (m model) bodyWidth() int {
	w := m.width - m.leftWidth() - 6 // borders for 3 panels
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

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Help overlay takes over the entire screen
	if m.showHelp {
		return renderHelp(m.width, m.height)
	}

	leftW := m.leftWidth()
	bodyW := m.bodyWidth()
	contentH := m.height - 3 // help bar + borders
	if contentH < 0 {
		contentH = 0
	}

	// Styles
	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(leftW).
		Height(contentH)
	bodyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(bodyW).
		Height(contentH)

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
				line := fmt.Sprintf("   %s", row.Object.DisplayID())
				if i == m.cursor {
					style := lipgloss.NewStyle().Bold(true).Reverse(true)
					line = style.Render(line)
				}
				lines = append(lines, line)
			}
			leftContent = strings.Join(lines, "\n")
		}
	} else {
		leftContent = renderList(m.groups, m.cursor, m.scrollOffset, m.focus == focusLeft, leftW, contentH)
	}

	// Compose panels
	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent),
		bodyStyle.Render(m.bodyViewport.View()),
	)

	// Properties panel (optional)
	if m.propsVisible {
		propsStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(m.propsWidth).
			Height(contentH)
		if m.focus == focusProps {
			propsStyle = propsStyle.BorderForeground(activeBorderColor)
		}
		panels = lipgloss.JoinHorizontal(lipgloss.Top,
			panels,
			propsStyle.Render(m.propsViewport.View()),
		)
	}

	// Help bar
	var helpBar string
	if m.searchMode {
		helpBar = "  / " + m.searchInput.View()
	} else if m.saveConflict {
		helpBar = "  [CONFLICT]  " + m.saveErr
	} else if m.saveErr != "" {
		helpBar = "  [ERROR]  " + m.saveErr
	} else if m.editMode {
		helpBar = "  [EDIT]  esc: exit edit mode"
	} else if m.searchResults != nil {
		helpBar = "  [VIEW]  Search results  |  esc: clear  |  ↑↓: navigate  |  tab: switch  |  q: quit"
	} else {
		helpBar = "  [VIEW]  ?/h: help  |  /: search  |  q: quit"
	}

	return panels + "\n" + helpBar
}

func Start(vaultPath string) error {
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

	objects, err := v.QueryObjects("")
	if err != nil {
		return fmt.Errorf("query objects: %w", err)
	}

	groups := buildGroups(objects)

	// Expand first group so first object is visible and selectable
	if len(groups) > 0 {
		groups[0].Expanded = true
	}

	// Auto-select first object
	var selected *core.Object
	var displayProps []core.DisplayProperty
	rows := visibleRows(groups)
	for _, row := range rows {
		if !row.IsHeader && row.Object != nil {
			selected = row.Object
			displayProps, _ = v.BuildDisplayProperties(selected)
			break
		}
	}

	bodyVP := viewport.New(0, 0)
	bodyVP.SetContent(renderBody(selected, 0))
	propsVP := viewport.New(0, 0)
	propsVP.SetContent(renderProperties(selected, displayProps))

	// Set cursor to first object (skip header row)
	initialCursor := 0
	for i, row := range rows {
		if !row.IsHeader && row.Object != nil {
			initialCursor = i
			break
		}
	}

	m := model{
		vault:         v,
		focus:         focusLeft,
		groups:        groups,
		cursor:        initialCursor,
		selected:      selected,
		bodyViewport:  bodyVP,
		propsViewport: propsVP,
		propsVisible:  false,
		softWrap:      true,
		displayProps: displayProps,
		searchInput:   initSearchInput(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
