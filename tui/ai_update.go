package tui

import (
	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	tea "charm.land/bubbletea/v2"
)

// handleAIGenerate processes ctrl+g — shows the AI action picker popup.
func handleAIGenerate(m model) (model, tea.Cmd) {
	if m.vault == nil || m.vault.AIService() == nil {
		return m, nil
	}
	if m.selected == nil || m.rightPanel != panelObject {
		return m, nil
	}
	if m.readOnly {
		return m, nil
	}

	m.aiState = aiActionPicker
	m.aiActionCursor = 0
	return m, nil
}

// updateAIActionPicker handles key events in the AI action picker popup.
func updateAIActionPicker(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.aiActionCursor > 0 {
			m.aiActionCursor--
		}
	case "down", "j":
		if m.aiActionCursor < len(aiActions)-1 {
			m.aiActionCursor++
		}
	case "enter":
		action := aiActions[m.aiActionCursor].action
		switch action {
		case aiActionDescribe:
			m.aiState = aiLoadingDescribe
			return m, triggerAIDescribe(m.vault, m.selected)
		case aiActionTag:
			m.aiState = aiLoadingTags
			return m, triggerAITags(m.vault, m.selected)
		}
	case "esc":
		m.aiState = aiIdle
	}
	return m, nil
}

// handleSchemaExplore processes ctrl+e for schema exploration.
func handleSchemaExplore(m model) (model, tea.Cmd) {
	if m.vault == nil || m.vault.AIService() == nil {
		return m, nil
	}
	if m.focus != focusLeft {
		return m, nil
	}
	se := newSchemaExplorer(m.vault)
	se.SetSize(m.width-2, m.height-3)
	m.schemaExplorer = se
	m.rightPanel = panelSchemaExplore
	m.focus = focusBody
	return m, nil
}

// updateAIDescribeResult handles the AI describe response.
func updateAIDescribeResult(m model, msg aiDescribeResultMsg) (model, tea.Cmd) {
	if msg.Err != nil {
		cmd := m.toast.Show(widget.ToastError, []widget.ToastItem{{Message: msg.Err.Error()}})
		m.aiState = aiIdle
		return m, cmd
	}
	m.aiState = aiPreviewDescribe
	m.aiPreviewDesc = msg.Description
	m.updateDetail() // refresh properties panel to show preview
	return m, nil
}

// updateAITagResult handles the AI tag suggestion response.
func updateAITagResult(m model, msg aiTagResultMsg) (model, tea.Cmd) {
	if msg.Err != nil {
		cmd := m.toast.Show(widget.ToastError, []widget.ToastItem{{Message: msg.Err.Error()}})
		m.aiState = aiIdle
		return m, cmd
	}
	if msg.Suggestion == nil || len(msg.Suggestion.Tags) == 0 {
		m.aiState = aiIdle
		return m, nil
	}

	// Convert suggestions to popup items, filtering out already-assigned tags
	existingTags := getObjectTags(m.selected)
	var items []tagPopupItem
	for _, tag := range msg.Suggestion.Tags {
		if _, assigned := existingTags[tag.Name]; assigned {
			continue
		}
		items = append(items, tagPopupItem{
			Name:   tag.Name,
			IsNew:  tag.IsNew,
			Reason: tag.Reason,
		})
	}

	if len(items) == 0 {
		m.aiState = aiIdle
		return m, nil
	}

	m.aiState = aiShowingTags
	m.aiTagItems = items
	m.aiTagCursor = 0
	m.updateDetail()
	return m, nil
}

// getObjectTags returns a set of tag names currently assigned to the object.
func getObjectTags(obj *core.Object) map[string]struct{} {
	tags := make(map[string]struct{})
	if obj == nil {
		return tags
	}
	if tagList, ok := obj.Properties["tags"]; ok {
		switch v := tagList.(type) {
		case []any:
			for _, t := range v {
				if s, ok := t.(string); ok {
					tags[s] = struct{}{}
				}
			}
		case []string:
			for _, t := range v {
				tags[t] = struct{}{}
			}
		}
	}
	return tags
}

// updateAIPreview handles key events during AI description preview.
func updateAIPreview(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		// Accept the AI description
		if m.selected != nil && m.aiPreviewDesc != "" {
			m.selected.Properties["description"] = m.aiPreviewDesc
			m.dirty = true
		}
		m.aiState = aiIdle
		m.aiPreviewDesc = ""
		m.updateDetail() // refresh after clearing AI state so normal style renders
		if m.dirty {
			m.saveObject()
		}
	case "esc":
		// Reject
		m.aiState = aiIdle
		m.aiPreviewDesc = ""
		m.updateDetail()
	}
	return m, nil
}

// updateAITagPopup handles key events during the tag suggestion popup.
func updateAITagPopup(m model, msg tea.KeyPressMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.aiTagCursor > 0 {
			m.aiTagCursor--
		}
	case "down", "j":
		if m.aiTagCursor < len(m.aiTagItems)-1 {
			m.aiTagCursor++
		}
	case " ", "space":
		// Toggle selection
		if m.aiTagCursor >= 0 && m.aiTagCursor < len(m.aiTagItems) {
			m.aiTagItems[m.aiTagCursor].Selected = !m.aiTagItems[m.aiTagCursor].Selected
		}
	case "enter":
		// Apply selected tags
		return applySelectedTags(m)
	case "esc":
		// Cancel
		m.aiState = aiIdle
		m.aiTagItems = nil
	}
	m.updateDetail()
	return m, nil
}

// aiTagError resets AI state and shows an error toast. Used by applySelectedTags error paths.
func aiTagError(m model, err error) (model, tea.Cmd) {
	cmd := m.toast.Show(widget.ToastError, []widget.ToastItem{{Message: err.Error()}})
	m.aiState = aiIdle
	m.aiTagItems = nil
	return m, cmd
}

// applySelectedTags applies the selected tags from the popup to the object.
func applySelectedTags(m model) (model, tea.Cmd) {
	if m.selected == nil || m.vault == nil {
		m.aiState = aiIdle
		m.aiTagItems = nil
		return m, nil
	}

	for _, item := range m.aiTagItems {
		if !item.Selected {
			continue
		}
		var tagID string
		if item.IsNew {
			tagObj, err := m.vault.Objects.Create("tag", item.Name, "")
			if err != nil {
				return aiTagError(m, err)
			}
			tagID = tagObj.ID
		} else {
			resolved, err := m.vault.ResolveID("tag/" + item.Name)
			if err != nil {
				return aiTagError(m, err)
			}
			tagID = resolved
		}
		if err := m.vault.LinkObjects(m.selected.ID, "tags", tagID); err != nil {
			return aiTagError(m, err)
		}
		m.dirty = true
	}

	if m.dirty {
		// Reload the object to get updated tags
		if obj, err := m.vault.GetObject(m.selected.ID); err == nil {
			m.selected = obj
		}
		m.updateDetail()
		m.saveObject()
	}

	m.aiState = aiIdle
	m.aiTagItems = nil
	return m, nil
}
