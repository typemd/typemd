package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

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

func (m model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	// Full-width view mode takes over the entire screen
	if m.rightPanel == panelView && m.viewMode != nil {
		contentH := m.height - 3
		if contentH < 0 {
			contentH = 0
		}
		bdr := 2
		vm := m.viewMode

		titleStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(m.width - bdr).
			Height(titlePanelHeight)
		titleText := vm.titleContent()

		bodyH := contentH - titlePanelHeight

		var bodyContent string
		if vm.HasPreview() {
			// Split: table on left, preview on right
			previewW := m.width * 2 / 5 // 40% for preview
			tableW := m.width - previewW - bdr - 2 // -2 for gap between panels

			vm.SetSize(tableW-bdr, bodyH-bdr)

			tableStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(tableW).
				Height(bodyH).
				MaxHeight(bodyH)

			previewStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("8")).
				Width(previewW).
				Height(bodyH).
				MaxHeight(bodyH)

			bodyContent = lipgloss.JoinHorizontal(lipgloss.Top,
				tableStyle.Render(vm.View()),
				previewStyle.Render(vm.PreviewContent()),
			)
		} else {
			// Full-width table
			vm.SetSize(m.width-bdr-2, bodyH-bdr)

			bodyStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(m.width - bdr).
				Height(bodyH).
				MaxHeight(bodyH)

			bodyContent = bodyStyle.Render(vm.View())
		}

		panels := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(titleText),
			bodyContent,
		)

		helpBar := "  " + vm.HelpBar()
		v := tea.NewView(panels + "\n" + helpBar)
		v.AltScreen = true
		return v
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
	}

	var rightSide string

	if m.rightPanel == panelTemplate && m.tmplEditor != nil {
		// Template detail view — uses full right-side width like type editor
		te := m.tmplEditor
		editorW := m.width - m.leftWidth() - 4
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

		// Title panel
		titleW := m.width - leftW - bdr
		titleStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(titleW).
			Height(titlePanelHeight)
		titleText := te.titleContent(titleW - bdr)

		rightSide = lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(titleText),
			editorStyle.Render(te.View()),
		)
	} else if m.rightPanel == panelTypeEditor && m.typeEditor != nil && m.createType == nil {
		// Type editor uses full right-side width (no props panel)
		editorW := m.width - m.leftWidth() - 4 // left border + body border
		if editorW < 10 {
			editorW = 10
		}
		// Set type editor content height for scroll support
		m.typeEditor.height = bodyPropsContentH
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
		if m.createType != nil {
			bodyPanelContent = renderCreateTypePreview(m.createType)
		} else if m.editMode && m.focus == focusBody {
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
			var titleContent string
			if m.create != nil {
				titleContent = renderCreateTitleContent(m.create, titleW-bdr)
				titleStyle = titleStyle.BorderForeground(colorFocusBorder)
			} else if m.createType != nil {
				titleContent = renderCreateTypeTitleContent(m.createType)
				titleStyle = titleStyle.BorderForeground(colorFocusBorder)
			} else if m.selected != nil {
				titleContent = renderTitleContent(m.selected, m.selected.Type, m.selectedTypeEmoji(), titleW-bdr)
			}
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
	if m.create != nil {
		helpBar = renderCreateHelpBar(m.create)
	} else if m.createType != nil {
		helpBar = renderCreateTypeHelpBar()
	} else if m.searchMode {
		helpBar = "  / " + m.searchInput.View()
	} else if m.rightPanel == panelTemplate && m.tmplEditor != nil && m.focus != focusLeft {
		helpBar = m.tmplEditor.HelpBar()
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

	// View picker overlay
	if m.viewPicker != nil {
		screen = m.viewPicker.View(m.width, m.height)
	}

	v := tea.NewView(screen)
	v.AltScreen = true
	return v
}
