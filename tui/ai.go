package tui

import (
	"context"
	"fmt"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/core/ai"
	tea "charm.land/bubbletea/v2"
)

// aiState tracks the current AI operation in progress.
type aiState int

const (
	aiIdle            aiState = iota
	aiActionPicker             // showing action picker popup
	aiLoadingDescribe          // waiting for AI describe result
	aiPreviewDescribe          // showing inline description preview
	aiLoadingTags              // waiting for AI tag suggestions
	aiShowingTags              // showing tag suggestion popup
	aiError                    // showing error message
)

// aiAction represents a selectable AI action.
type aiAction int

const (
	aiActionDescribe aiAction = iota
	aiActionTag
)

var aiActions = []struct {
	action aiAction
	label  string
}{
	{aiActionDescribe, "Generate description"},
	{aiActionTag, "Suggest tags"},
}

// aiDescribeResultMsg carries the AI-generated description.
type aiDescribeResultMsg struct {
	Description string
	Err         error
}

// aiTagResultMsg carries the AI-suggested tags.
type aiTagResultMsg struct {
	Suggestion *ai.TagSuggestion
	Err        error
}

// tagPopupItem represents a single item in the tag suggestion popup.
type tagPopupItem struct {
	Name     string
	IsNew    bool
	Reason   string
	Selected bool
}

// unpinnedProps returns the displayProps that are shown in the properties panel
// (not pinned, not name). This matches the logic in renderProperties.
func unpinnedProps(displayProps []core.DisplayProperty) []core.DisplayProperty {
	var result []core.DisplayProperty
	for _, p := range displayProps {
		if p.Pin == 0 && p.Key != core.NameProperty {
			result = append(result, p)
		}
	}
	return result
}

// triggerAIDescribe starts an AI description generation as a tea.Cmd.
func triggerAIDescribe(vault *core.Vault, obj *core.Object) tea.Cmd {
	return func() tea.Msg {
		svc := vault.AIService()
		if svc == nil {
			return aiDescribeResultMsg{Err: errAIUnavailable}
		}

		objCtx := objectToAIContext(obj)
		schema, _ := vault.LoadType(obj.Type)
		schemaCtx := schemaToAIContext(obj.Type, schema)

		desc, err := svc.Describe(context.Background(), objCtx, schemaCtx)
		return aiDescribeResultMsg{Description: desc, Err: err}
	}
}

// triggerAITags starts an AI tag suggestion as a tea.Cmd.
func triggerAITags(vault *core.Vault, obj *core.Object) tea.Cmd {
	return func() tea.Msg {
		svc := vault.AIService()
		if svc == nil {
			return aiTagResultMsg{Err: errAIUnavailable}
		}

		objCtx := objectToAIContext(obj)
		schema, _ := vault.LoadType(obj.Type)
		schemaCtx := schemaToAIContext(obj.Type, schema)
		existingTags := gatherExistingTags(vault)

		suggestion, err := svc.SuggestTags(context.Background(), objCtx, schemaCtx, existingTags)
		return aiTagResultMsg{Suggestion: suggestion, Err: err}
	}
}

// objectToAIContext converts a core.Object to an ai.ObjectContext.
func objectToAIContext(obj *core.Object) ai.ObjectContext {
	ctx := ai.ObjectContext{
		Name:       obj.GetName(),
		Properties: make(map[string]any),
		Body:       obj.Body,
	}
	for k, v := range obj.Properties {
		if k != core.NameProperty {
			ctx.Properties[k] = v
		}
	}
	return ctx
}

// schemaToAIContext converts a core.TypeSchema to an ai.SchemaContext.
func schemaToAIContext(typeName string, schema *core.TypeSchema) ai.SchemaContext {
	ctx := ai.SchemaContext{TypeName: typeName}
	if schema != nil {
		ctx.TypeDescription = schema.Description
		ctx.PropertyTypes = make(map[string]string)
		ctx.PropertyDescriptions = make(map[string]string)
		for _, p := range schema.Properties {
			ctx.PropertyNames = append(ctx.PropertyNames, p.Name)
			ctx.PropertyTypes[p.Name] = p.Type
			if p.Description != "" {
				ctx.PropertyDescriptions[p.Name] = p.Description
			}
		}
	}
	return ctx
}

// gatherExistingTags collects all tag objects from the vault.
func gatherExistingTags(vault *core.Vault) []ai.TagInfo {
	filter := core.TypeFilter(core.TagTypeName)
	results, _ := vault.QueryObjects(filter)
	var tags []ai.TagInfo
	for _, obj := range results {
		tag := ai.TagInfo{Name: obj.GetName()}
		if desc, ok := obj.Properties["description"].(string); ok && desc != "" {
			tag.Description = desc
		}
		tags = append(tags, tag)
	}
	return tags
}

var errAIUnavailable = fmt.Errorf("AI service not available")
