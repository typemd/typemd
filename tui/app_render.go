package tui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"

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

	m.updatePropsContent()
}

// updatePropsContent refreshes only the properties panel viewport content.
func (m *model) updatePropsContent() {
	var propsContent string
	if m.aiState != aiIdle {
		propsContent = renderPropertiesWithAI(m.selected, m.displayProps, *m)
	} else if m.propEdit != nil && m.focus == focusProps {
		propsContent = m.propEdit.Render(true)
	} else if m.propEdit != nil {
		propsContent = m.propEdit.Render(false)
	} else {
		propsContent = renderProperties(m.selected, m.displayProps)
	}
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
		if vm.detailObject != nil && m.selected != nil {
			emoji := ""
			if vm.schema != nil {
				emoji = vm.schema.Emoji
			}
			titleText = renderTitleContent(m.selected, m.selected.Type, emoji, m.width-bdr-2)
		}

		bodyH := contentH - titlePanelHeight

		var bodyContent string
		if vm.detailObject != nil {
			// Object detail within view mode — full-width body only
			m.bodyViewport.SetWidth(m.width - bdr - 2)
			m.bodyViewport.SetHeight(bodyH - bdr)
			m.bodyViewport.SetContent(renderBody(m.selected, m.width-bdr-2, m.displayProps))

			bodyStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(m.width - bdr).
				Height(bodyH).
				MaxHeight(bodyH)

			bodyContent = bodyStyle.Render(m.bodyViewport.View())
		} else if vm.HasEditor() {
			// Split: table on left, editor on right
			totalContent := m.width - bdr
			editorW := totalContent * 2 / 5
			tableW := totalContent - editorW

			vm.SetSize(tableW-bdr, bodyH-bdr)
			vm.SetEditorSize(editorW-bdr, bodyH-bdr)

			tableStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(tableW).
				Height(bodyH).
				MaxHeight(bodyH)

			editorStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorFocusBorder).
				Width(editorW).
				Height(bodyH).
				MaxHeight(bodyH)

			bodyContent = lipgloss.JoinHorizontal(lipgloss.Top,
				tableStyle.Render(vm.View()),
				editorStyle.Render(vm.EditorView()),
			)
		} else if vm.HasPreview() {
			// Split: table on left, preview on right
			// Match title: Width(m.width-bdr) fills screen, so total content = m.width-bdr
			totalContent := m.width - bdr
			previewW := totalContent / 2
			tableW := totalContent - previewW

			vm.SetSize(tableW-bdr, bodyH-bdr)
			vm.previewWidth = previewW - bdr

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
		screen := m.toast.Overlay(panels+"\n"+helpBar, m.width, m.height)
		v := tea.NewView(screen)
		v.AltScreen = true
		return v
	}

	// Full-width stats mode takes over the entire screen
	if m.rightPanel == panelStats && m.statsMode != nil {
		contentH := m.height - 3
		if contentH < 0 {
			contentH = 0
		}
		bdr := 2
		stm := m.statsMode

		titleStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(m.width - bdr).
			Height(titlePanelHeight)
		titleText := stm.titleContent()

		bodyH := contentH - titlePanelHeight

		var bodyContent string
		if stm.screen == statsDetail && stm.typeLayout == "popup" {
			// Popup layout: overview behind, detail as popup overlay
			stm.SetSize(m.width-bdr-2, bodyH-bdr)
			overviewContent := stm.viewOverview()

			bodyStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(m.width - bdr).
				Height(bodyH).
				MaxHeight(bodyH)

			background := lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render(" Vault Statistics"),
				bodyStyle.Render(overviewContent),
			)
			background += "\n  " + stm.HelpBar()

			// Render popup overlay
			popupW := m.width * 2 / 3
			if popupW < 40 {
				popupW = min(40, m.width-4)
			}
			// Set width to popup inner width so bar charts scale correctly
			stm.SetSize(popupW-bdr-2, bodyH-bdr)
			popupContent := stm.viewDetail()
			screen := renderOverlayPopup(background, popupContent, m.width, m.height, popupW)

			if m.showHelp {
				screen = renderStatsHelp(screen, m.width, m.height, stm.screen)
			}

			screen = m.toast.Overlay(screen, m.width, m.height)
			v := tea.NewView(screen)
			v.AltScreen = true
			return v
		}

		// Fullscreen layout
		bodyStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(m.width - bdr).
			Height(bodyH).
			MaxHeight(bodyH)

		stm.SetSize(m.width-bdr-2, bodyH-bdr)
		bodyContent = bodyStyle.Render(stm.View())

		panels := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(titleText),
			bodyContent,
		)

		helpBar := "  " + stm.HelpBar()
		screen := panels + "\n" + helpBar

		if m.showHelp {
			screen = renderStatsHelp(screen, m.width, m.height, stm.screen)
		}

		screen = m.toast.Overlay(screen, m.width, m.height)
		v := tea.NewView(screen)
		v.AltScreen = true
		return v
	}

	// Schema explore — full-width panel
	if m.rightPanel == panelSchemaExplore && m.schemaExplorer != nil {
		contentH := m.height - 3
		if contentH < 0 {
			contentH = 0
		}
		bdr := 2
		bodyStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(m.width - bdr).
			Height(contentH).
			MaxHeight(contentH)

		m.schemaExplorer.SetSize(m.width-bdr-2, contentH-bdr)
		content := bodyStyle.Render(m.schemaExplorer.View())

		helpBar := "  [AI EXPLORE]  ↑↓: navigate  |  enter: accept  |  s: skip  |  esc: exit"
		screen := content + "\n" + helpBar

		if m.showHelp {
			screen = renderHelp(screen, m.width, m.height, m.readOnly)
		}

		screen = m.toast.Overlay(screen, m.width, m.height)
		v := tea.NewView(screen)
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
				name := row.Object.GetName()
				maxWidth := leftW - 3 - runewidth.StringWidth(row.Object.Type) - 1 // 3 = leading spaces, 1 = "/"
				if maxWidth > 0 {
					name = runewidth.Truncate(name, maxWidth, "…")
				}
				line := fmt.Sprintf("   %s/%s", row.Object.Type, name)
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
		titleText := titlePrefix(te.schema.Emoji, te.typeName)
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
				propsBorderColor := colorFocusBorder
				if m.propEdit != nil && m.propEdit.isEditing() {
					propsBorderColor = colorEditBorder
				}
				propsStyle = propsStyle.BorderForeground(propsBorderColor)
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
			} else if m.rename != nil {
				titleContent = renderRenameTitleContent(m.rename)
				titleStyle = titleStyle.BorderForeground(colorEditBorder)
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
	} else if m.rename != nil {
		helpBar = "  [RENAME]  enter: save  esc: cancel"
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
	} else if m.aiState == aiActionPicker {
		helpBar = "  [AI]  ↑↓: navigate  |  enter: select  |  esc: cancel"
	} else if m.aiState == aiLoadingDescribe || m.aiState == aiLoadingTags {
		helpBar = "  [AI]  Generating..."
	} else if m.aiState == aiPreviewDescribe {
		helpBar = "  [AI]  tab: accept  |  esc: reject"
	} else if m.aiState == aiShowingTags {
		helpBar = "  [AI TAGS]  ↑↓: navigate  |  space: toggle  |  enter: apply  |  esc: cancel"
	} else if m.propEdit != nil && m.propEdit.isPicking() {
		helpBar = "  [PICK]  ↑↓: navigate  |  enter: select  |  esc: cancel"
	} else if m.propEdit != nil && m.propEdit.mode == propModeDateSegment {
		helpBar = "  [DATE]  ←→: segment  |  ↑↓: adjust  |  c: calendar  |  enter: confirm  |  esc: cancel"
	} else if m.propEdit != nil && m.propEdit.mode == propModeDateCalendar {
		helpBar = "  [CAL]  ←→↑↓: navigate  |  H/L: month  |  t: today  |  c: segments  |  enter: confirm  |  esc: cancel"
	} else if m.propEdit != nil && m.propEdit.isEditing() {
		helpBar = "  [EDIT]  enter: confirm  |  esc: cancel"
	} else if m.focus == focusProps && m.propEdit != nil {
		helpBar = "  [PROPS]  ↑↓: navigate  |  enter: edit  |  esc: back  |  tab: switch"
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

	// Help overlay using Layer/Compositor (background remains visible)
	if m.showHelp {
		screen = renderHelp(screen, m.width, m.height, m.readOnly)
	}

	// Overlay popup if type editor has one active
	if m.rightPanel == panelTypeEditor && m.typeEditor != nil {
		if overlay := m.typeEditor.Overlay(m.width, m.height); overlay != "" {
			screen = overlay
		}
	}

	// AI action picker overlay
	if m.aiState == aiActionPicker {
		screen = renderAIActionPicker(screen, m.width, m.height, m.aiActionCursor)
	}

	// View picker overlay
	if m.viewPicker != nil {
		screen = m.viewPicker.View(m.width, m.height)
	}

	// Toast overlay (bottom-right)
	screen = m.toast.Overlay(screen, m.width, m.height)

	v := tea.NewView(screen)
	v.AltScreen = true
	return v
}
