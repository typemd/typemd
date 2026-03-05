package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/MilesChou/typemd/core"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusPanel int

const (
	focusLeft focusPanel = iota
	focusRight
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

	// Right panel
	viewport  viewport.Model
	relations []core.Relation
	schema    *core.TypeSchema

	// Search
	searchMode    bool
	searchInput   textinput.Model
	searchResults []*core.Object

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
		m.refreshData()
		// Restart watcher for next change
		return m, watchObjects(m.vault.ObjectsDir())

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rightWidth := m.width - m.leftWidth() - 4 // borders
		contentHeight := m.height - 3              // help bar + borders
		if rightWidth < 0 {
			rightWidth = 0
		}
		if contentHeight < 0 {
			contentHeight = 0
		}
		m.viewport.Width = rightWidth
		m.viewport.Height = contentHeight
		m.updateDetail()
		return m, nil

	case tea.KeyMsg:
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

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "/":
			m.searchMode = true
			m.searchInput.Focus()
			return m, textinput.Blink

		case "tab":
			if m.focus == focusLeft {
				m.focus = focusRight
			} else {
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
			} else {
				m.viewport.LineUp(1)
			}
			return m, nil

		case "down", "j":
			if m.focus == focusLeft {
				rows := m.currentRows()
				m.cursor = clampCursor(m.cursor+1, len(rows))
				m.adjustScroll()
				m.selectCurrentRow()
			} else {
				m.viewport.LineDown(1)
			}
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

// selectCurrentRow updates the selected object based on current cursor position.
// Re-reads the object from disk to get the latest body and properties.
func (m *model) selectCurrentRow() {
	rows := m.currentRows()
	if m.cursor >= 0 && m.cursor < len(rows) {
		row := rows[m.cursor]
		if !row.IsHeader && row.Object != nil {
			if m.vault != nil {
				if obj, err := m.vault.GetObject(row.Object.ID); err == nil {
					// Fill in missing schema-defined properties
					if schema, err := m.vault.LoadType(obj.Type); err == nil {
						for _, p := range schema.Properties {
							if _, ok := obj.Properties[p.Name]; !ok {
								obj.Properties[p.Name] = nil
							}
						}
						m.schema = schema
					} else {
						m.schema = nil
					}
					m.selected = obj
				} else {
					m.selected = row.Object
					m.schema = nil
				}
				m.relations, _ = m.vault.ListRelations(m.selected.ID)
			} else {
				m.selected = row.Object
				m.schema = nil
			}
			m.updateDetail()
			return
		}
	}
}

// adjustScroll updates scrollOffset so cursor is always visible.
func (m *model) adjustScroll() {
	contentH := m.height - 3
	m.scrollOffset = adjustScrollOffset(m.cursor, m.scrollOffset, contentH)
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

// updateDetail refreshes the viewport content with current selected object.
func (m *model) updateDetail() {
	content := renderDetail(m.selected, m.relations, m.schema)
	if m.softWrap && m.viewport.Width > 0 {
		content = softWrapLines(content, m.viewport.Width)
	}
	m.viewport.SetContent(content)
}

// leftWidth returns the width allocated for the left panel.
func (m model) leftWidth() int {
	w := m.width * 2 / 5
	if w < 20 {
		w = 20
	}
	if w > 50 {
		w = 50
	}
	return w
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	leftW := m.leftWidth()
	rightW := m.width - leftW - 4 // account for borders
	contentH := m.height - 3      // help bar + borders
	if contentH < 0 {
		contentH = 0
	}

	// Styles
	leftBorder := lipgloss.RoundedBorder()
	rightBorder := lipgloss.RoundedBorder()

	leftStyle := lipgloss.NewStyle().
		Border(leftBorder).
		Width(leftW).
		Height(contentH)
	rightStyle := lipgloss.NewStyle().
		Border(rightBorder).
		Width(rightW).
		Height(contentH)

	if m.focus == focusLeft {
		leftStyle = leftStyle.BorderForeground(lipgloss.Color("63"))
	} else {
		rightStyle = rightStyle.BorderForeground(lipgloss.Color("63"))
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
				line := fmt.Sprintf("   %s/%s", row.Object.Type, row.Object.Filename)
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

	// Right panel content
	rightContent := m.viewport.View()

	// Compose panels
	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(leftContent),
		rightStyle.Render(rightContent),
	)

	// Help bar
	var helpBar string
	if m.searchMode {
		helpBar = "  / " + m.searchInput.View()
	} else if m.searchResults != nil {
		helpBar = "  Search results  |  esc: clear  |  ↑↓: navigate  |  tab: switch  |  q: quit"
	} else {
		wrapLabel := "off"
		if m.softWrap {
			wrapLabel = "on"
		}
		helpBar = fmt.Sprintf("  ↑↓/jk: navigate  |  enter: toggle  |  tab: switch  |  /: search  |  w: wrap(%s)  |  q: quit", wrapLabel)
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
	if err := v.Open(); err != nil {
		return fmt.Errorf("open vault: %w", err)
	}
	defer v.Close()

	objects, err := v.QueryObjects("")
	if err != nil {
		return fmt.Errorf("query objects: %w", err)
	}

	groups := buildGroups(objects)

	// Auto-select first object
	var selected *core.Object
	var relations []core.Relation
	var schema *core.TypeSchema
	rows := visibleRows(groups)
	for _, row := range rows {
		if !row.IsHeader && row.Object != nil {
			selected = row.Object
			relations, _ = v.ListRelations(selected.ID)
			schema, _ = v.LoadType(selected.Type)
			break
		}
	}

	vp := viewport.New(0, 0)
	vp.SetContent(renderDetail(selected, relations, schema))

	m := model{
		vault:       v,
		focus:       focusLeft,
		groups:      groups,
		cursor:      0,
		selected:    selected,
		viewport:    vp,
		relations:   relations,
		schema:      schema,
		searchInput: initSearchInput(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
