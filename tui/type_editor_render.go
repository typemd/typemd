package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	tea "charm.land/bubbletea/v2"
)

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

// ── View Rendering ──────────────────────────────────────────────────────────

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

	if te.mode == teModeDeleteTemplate {
		writeBlank()
		item := items[te.cursor]
		tmplIdx := templateSentinelBase - item
		if tmplIdx >= 0 && tmplIdx < len(te.templates) {
			writeLine(fmt.Sprintf(" Delete template '%s'? [y/n]", te.templates[tmplIdx]))
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
	te.scrollOffset = widget.AdjustScroll(cursorLine, te.scrollOffset, visibleH)

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

