package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/core/ai"
	tea "charm.land/bubbletea/v2"
)

// schemaExplorerScreen tracks the schema explorer's internal state.
type schemaExplorerScreen int

const (
	seSelectType  schemaExplorerScreen = iota // type selection
	seLoading                                 // waiting for AI analysis
	seSuggestions                              // showing suggestions
	seSummary                                 // showing completion summary
)

// suggestionStatus tracks each suggestion's resolution.
type suggestionStatus int

const (
	suggestionPending  suggestionStatus = iota
	suggestionAccepted
	suggestionSkipped
)

// schemaExplorer is the sub-model for AI schema exploration.
type schemaExplorer struct {
	vault       *core.Vault
	screen      schemaExplorerScreen
	types       []string       // available type names
	typeCursor  int            // cursor in type list
	typeName    string         // selected type
	suggestions []ai.SchemaSuggestion
	statuses    []suggestionStatus
	cursor      int            // cursor in suggestion list
	err         error
	width       int
	height      int
}

// aiExploreResultMsg carries the AI schema exploration result.
type aiExploreResultMsg struct {
	Exploration *ai.SchemaExploration
	Err         error
}

func newSchemaExplorer(vault *core.Vault) *schemaExplorer {
	types := vault.ListTypes()
	return &schemaExplorer{
		vault:  vault,
		screen: seSelectType,
		types:  types,
	}
}

func (se *schemaExplorer) SetSize(w, h int) {
	se.width = w
	se.height = h
}

func (se *schemaExplorer) Update(msg tea.KeyPressMsg) (*schemaExplorer, tea.Cmd) {
	switch se.screen {
	case seSelectType:
		return se.updateTypeSelect(msg)
	case seSuggestions:
		return se.updateSuggestions(msg)
	case seSummary:
		// Any key exits
		return nil, nil
	}
	return se, nil
}

func (se *schemaExplorer) updateTypeSelect(msg tea.KeyPressMsg) (*schemaExplorer, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if se.typeCursor > 0 {
			se.typeCursor--
		}
	case "down", "j":
		if se.typeCursor < len(se.types)-1 {
			se.typeCursor++
		}
	case "enter":
		if se.typeCursor >= 0 && se.typeCursor < len(se.types) {
			se.typeName = se.types[se.typeCursor]
			se.screen = seLoading
			return se, se.triggerExplore()
		}
	case "esc":
		return nil, nil
	}
	return se, nil
}

func (se *schemaExplorer) updateSuggestions(msg tea.KeyPressMsg) (*schemaExplorer, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		se.moveCursorUp()
	case "down", "j":
		se.moveCursorDown()
	case "enter":
		// Accept current suggestion
		if se.cursor >= 0 && se.cursor < len(se.suggestions) && se.statuses[se.cursor] == suggestionPending {
			if err := se.applySuggestion(se.suggestions[se.cursor]); err != nil {
				se.err = err
			} else {
				se.statuses[se.cursor] = suggestionAccepted
			}
			se.moveCursorDown()
		}
		if se.allResolved() {
			se.screen = seSummary
		}
	case "s":
		// Skip current suggestion
		if se.cursor >= 0 && se.cursor < len(se.suggestions) && se.statuses[se.cursor] == suggestionPending {
			se.statuses[se.cursor] = suggestionSkipped
			se.moveCursorDown()
		}
		if se.allResolved() {
			se.screen = seSummary
		}
	case "esc":
		se.screen = seSummary
	}
	return se, nil
}

func (se *schemaExplorer) moveCursorUp() {
	for i := se.cursor - 1; i >= 0; i-- {
		if se.statuses[i] == suggestionPending {
			se.cursor = i
			return
		}
	}
}

func (se *schemaExplorer) moveCursorDown() {
	for i := se.cursor + 1; i < len(se.suggestions); i++ {
		if se.statuses[i] == suggestionPending {
			se.cursor = i
			return
		}
	}
}

func (se *schemaExplorer) allResolved() bool {
	for _, s := range se.statuses {
		if s == suggestionPending {
			return false
		}
	}
	return true
}

func (se *schemaExplorer) triggerExplore() tea.Cmd {
	return func() tea.Msg {
		svc := se.vault.AIService()
		if svc == nil {
			return aiExploreResultMsg{Err: errAIUnavailable}
		}

		schema, _ := se.vault.LoadType(se.typeName)
		schemaCtx := schemaToAIContext(se.typeName, schema)

		// Sample objects
		cfg := se.vault.Config()
		sampleCount := 10
		bodyTruncate := 500
		if cfg != nil && cfg.AI.Explore.SampleCount > 0 {
			sampleCount = cfg.AI.Explore.SampleCount
		}
		if cfg != nil && cfg.AI.Explore.BodyTruncate > 0 {
			bodyTruncate = cfg.AI.Explore.BodyTruncate
		}

		filter := core.TypeFilter(se.typeName)
		allObjects, _ := se.vault.QueryObjects(filter)

		var objects []ai.ObjectContext
		for i, obj := range allObjects {
			if i >= sampleCount {
				break
			}
			objCtx := objectToAIContext(obj)
			if len(objCtx.Body) > bodyTruncate {
				objCtx.Body = objCtx.Body[:bodyTruncate]
			}
			objects = append(objects, objCtx)
		}

		result, err := svc.ExploreSchema(context.Background(), schemaCtx, objects)
		return aiExploreResultMsg{Exploration: result, Err: err}
	}
}

func (se *schemaExplorer) applySuggestion(s ai.SchemaSuggestion) error {
	schema, err := se.vault.LoadType(se.typeName)
	if err != nil {
		return err
	}

	switch s.Type {
	case ai.SuggestionAdd:
		prop := core.Property{
			Name:        s.PropertyName,
			Type:        s.PropertyType,
			Description: s.Description,
		}
		schema.Properties = append(schema.Properties, prop)
	case ai.SuggestionModify:
		for i, p := range schema.Properties {
			if p.Name == s.PropertyName {
				if s.PropertyType != "" {
					schema.Properties[i].Type = s.PropertyType
				}
				if s.Description != "" {
					schema.Properties[i].Description = s.Description
				}
				break
			}
		}
	case ai.SuggestionRemove:
		var filtered []core.Property
		for _, p := range schema.Properties {
			if p.Name != s.PropertyName {
				filtered = append(filtered, p)
			}
		}
		schema.Properties = filtered
	default:
		return fmt.Errorf("unknown suggestion type: %s", s.Type)
	}

	return se.vault.SaveType(schema)
}

func (se *schemaExplorer) View() string {
	switch se.screen {
	case seSelectType:
		return se.viewTypeSelect()
	case seLoading:
		return se.viewLoading()
	case seSuggestions:
		return se.viewSuggestions()
	case seSummary:
		return se.viewSummary()
	}
	return ""
}

func (se *schemaExplorer) viewTypeSelect() string {
	var b strings.Builder
	b.WriteString(" Schema Explore\n")
	b.WriteString(" ══════════════\n\n")
	b.WriteString(" Select a type to analyze:\n\n")

	for i, t := range se.types {
		cursor := "  "
		if i == se.typeCursor {
			cursor = aiCursorStyle.Render("▸ ")
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, t))
	}

	b.WriteString("\n enter: select  |  esc: cancel")
	return b.String()
}

func (se *schemaExplorer) viewLoading() string {
	return fmt.Sprintf(" Schema Explore: %s\n ══════════════\n\n Analyzing objects...", se.typeName)
}

func (se *schemaExplorer) viewSuggestions() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(" Schema Explore: %s\n", se.typeName))
	b.WriteString(" ══════════════\n\n")

	if se.err != nil {
		b.WriteString(fmt.Sprintf(" Error: %s\n", se.err.Error()))
	}

	for i, s := range se.suggestions {
		status := " "
		switch se.statuses[i] {
		case suggestionAccepted:
			status = "✓"
		case suggestionSkipped:
			status = "–"
		}

		cursor := "  "
		if i == se.cursor && se.statuses[i] == suggestionPending {
			cursor = aiCursorStyle.Render("▸ ")
		}

		typeLabel := strings.ToUpper(s.Type)
		line := fmt.Sprintf("%s[%s] %s %s", cursor, status, typeLabel, s.PropertyName)
		if s.PropertyType != "" {
			line += fmt.Sprintf(" (%s)", s.PropertyType)
		}
		b.WriteString(line + "\n")
		b.WriteString(fmt.Sprintf("       %s\n", s.Reason))
		if s.Description != "" {
			b.WriteString(fmt.Sprintf("       desc: %s\n", s.Description))
		}
		b.WriteString("\n")
	}

	b.WriteString(" enter: accept  |  s: skip  |  esc: done")
	return b.String()
}

func (se *schemaExplorer) viewSummary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(" Schema Explore: %s — Summary\n", se.typeName))
	b.WriteString(" ══════════════\n\n")

	accepted := 0
	skipped := 0
	for _, s := range se.statuses {
		switch s {
		case suggestionAccepted:
			accepted++
		case suggestionSkipped:
			skipped++
		}
	}

	b.WriteString(fmt.Sprintf(" %d accepted, %d skipped\n\n", accepted, skipped))

	for i, s := range se.suggestions {
		switch se.statuses[i] {
		case suggestionAccepted:
			b.WriteString(fmt.Sprintf(" ✓ %s %s\n", strings.ToUpper(s.Type), s.PropertyName))
		case suggestionSkipped:
			b.WriteString(fmt.Sprintf(" – %s %s\n", strings.ToUpper(s.Type), s.PropertyName))
		}
	}

	b.WriteString("\n press any key to return")
	return b.String()
}
