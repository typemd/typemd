package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"charm.land/lipgloss/v2"
)

var (
	aiPreviewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true)
	aiLabelNewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true)
	aiLabelExistingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("2"))
	aiCursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))
)

// renderAIActionPicker renders the AI action picker as a centered overlay popup.
func renderAIActionPicker(background string, termW, termH, cursor int) string {
	var b strings.Builder
	b.WriteString("AI Assist\n")
	b.WriteString("─────────\n\n")

	for i, a := range aiActions {
		prefix := "  "
		if i == cursor {
			prefix = aiCursorStyle.Render("▸ ")
		}
		b.WriteString(fmt.Sprintf("%s%s\n", prefix, a.label))
	}

	return renderOverlayPopup(background, b.String(), termW, termH, 30)
}

// renderPropertiesWithAI renders properties panel with AI state overlaid.
func renderPropertiesWithAI(obj *core.Object, displayProps []core.DisplayProperty, m model) string {
	if obj == nil && len(displayProps) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	unpinned := unpinnedProps(displayProps)

	// Show AI description preview at top if in preview state
	if m.aiState == aiPreviewDescribe && m.aiPreviewDesc != "" {
		b.WriteString(fmt.Sprintf(" description: %s\n", aiPreviewStyle.Render(m.aiPreviewDesc)))
	}

	if len(unpinned) == 0 && m.aiState != aiPreviewDescribe {
		b.WriteString(" (none)\n")
		return b.String()
	}

	for _, p := range unpinned {
		// Skip description in the normal list if we're showing AI preview for it
		if p.Key == "description" && m.aiState == aiPreviewDescribe {
			continue
		}

		b.WriteString(fmt.Sprintf(" %s\n", p.Format()))
	}

	// Show tag popup below properties
	if m.aiState == aiShowingTags && len(m.aiTagItems) > 0 {
		b.WriteString("\n")
		b.WriteString(" AI Suggested Tags\n")
		b.WriteString(" ─────────────────\n")
		for i, item := range m.aiTagItems {
			cursor := "  "
			if i == m.aiTagCursor {
				cursor = aiCursorStyle.Render("▸ ")
			}

			checkbox := "☐"
			if item.Selected {
				checkbox = "☑"
			}

			label := item.Name
			if item.IsNew {
				label = aiLabelNewStyle.Render(item.Name + " ★ new")
			} else {
				label = aiLabelExistingStyle.Render(item.Name)
			}

			reason := ""
			if item.Reason != "" {
				reason = fmt.Sprintf(" — %s", item.Reason)
			}

			b.WriteString(fmt.Sprintf("%s%s %s%s\n", cursor, checkbox, label, reason))
		}
	}

	return b.String()
}
