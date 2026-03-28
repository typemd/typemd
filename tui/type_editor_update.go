package tui

import (
	tea "charm.land/bubbletea/v2"
)

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
	switch {
	case item <= templateSentinelBase:
		tmplIdx := templateSentinelBase - item
		if tmplIdx >= 0 && tmplIdx < len(te.templates) {
			te.mode = teModeDeleteTemplate
		}
	case item >= 0:
		te.mode = teModeDeleteProp
	}
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

func (te *typeEditor) updateDeleteTemplate(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	switch msg.String() {
	case "y":
		items := te.displayItems()
		if te.cursor < len(items) {
			item := items[te.cursor]
			tmplIdx := templateSentinelBase - item
			if tmplIdx >= 0 && tmplIdx < len(te.templates) {
				tmplName := te.templates[tmplIdx]
				if te.vault != nil {
					if err := te.vault.DeleteTemplate(te.typeName, tmplName); err != nil {
						te.saveErr = err.Error()
						te.mode = teModeView
						return te, nil
					}
				}
				te.refreshTemplates()
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
