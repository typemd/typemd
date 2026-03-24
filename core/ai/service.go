package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AIService provides domain-specific AI operations.
type AIService struct {
	provider Provider
	config   ServiceConfig
}

// ServiceConfig holds configuration for the AIService.
type ServiceConfig struct {
	Model          string
	Language       string
	DescribePrompt string
	TagPrompt      string
	ExplorePrompt  string
}

// ObjectContext holds the relevant data from an object for prompt assembly.
type ObjectContext struct {
	Name       string
	Properties map[string]any
	Body       string
}

// SchemaContext holds type schema information for prompt assembly.
type SchemaContext struct {
	TypeName          string
	TypeDescription   string
	PropertyNames     []string
	PropertyTypes     map[string]string
	PropertyDescriptions map[string]string
}

// TagInfo holds information about an existing tag.
type TagInfo struct {
	Name        string
	Description string
}

// SuggestedTag represents a single AI-suggested tag.
type SuggestedTag struct {
	Name   string `json:"name"`
	IsNew  bool   `json:"is_new"`
	Reason string `json:"reason"`
}

// TagSuggestion holds the result of an AI tag suggestion.
type TagSuggestion struct {
	Tags []SuggestedTag `json:"tags"`
}

// Suggestion type constants.
const (
	SuggestionAdd    = "add"
	SuggestionModify = "modify"
	SuggestionRemove = "remove"
)

// SchemaSuggestion represents a single schema modification suggestion.
type SchemaSuggestion struct {
	Type         string `json:"type"` // "add", "modify", "remove"
	PropertyName string `json:"property_name"`
	PropertyType string `json:"property_type,omitempty"`
	Reason       string `json:"reason"`
	Description  string `json:"description,omitempty"`
}

// SchemaExploration holds the result of AI schema analysis.
type SchemaExploration struct {
	Suggestions []SchemaSuggestion `json:"suggestions"`
}

// describeResponse is the JSON schema response for Describe.
type describeResponse struct {
	Description string `json:"description"`
}

// NewAIService creates a new AIService with the given provider and config.
func NewAIService(provider Provider, cfg ServiceConfig) *AIService {
	return &AIService{
		provider: provider,
		config:   cfg,
	}
}

// withLanguage appends a language instruction to the system prompt if configured.
func (s *AIService) withLanguage(systemPrompt string) string {
	if s.config.Language == "" {
		return systemPrompt
	}
	return systemPrompt + "\n\nIMPORTANT: Respond in " + s.config.Language + "."
}

// Describe generates a description for an object using AI.
func (s *AIService) Describe(ctx context.Context, obj ObjectContext, schema SchemaContext) (string, error) {
	prompt := s.buildDescribePrompt(obj, schema)
	systemPrompt := s.config.DescribePrompt
	if systemPrompt == "" {
		systemPrompt = DefaultDescribePrompt
	}

	jsonSchema := json.RawMessage(`{"type":"object","properties":{"description":{"type":"string"}},"required":["description"]}`)

	resp, err := s.provider.Complete(ctx, &CompletionRequest{
		SystemPrompt: s.withLanguage(systemPrompt),
		UserPrompt:   prompt,
		JSONSchema:   jsonSchema,
		Model:        s.config.Model,
	})
	if err != nil {
		return "", fmt.Errorf("ai describe: %w", err)
	}

	var result describeResponse
	if err := json.Unmarshal(resp.JSONResult, &result); err != nil {
		return "", fmt.Errorf("ai describe: parse response: %w", err)
	}

	return result.Description, nil
}

// SuggestTags suggests tags for an object using AI.
func (s *AIService) SuggestTags(ctx context.Context, obj ObjectContext, schema SchemaContext, existingTags []TagInfo) (*TagSuggestion, error) {
	prompt := s.buildTagPrompt(obj, schema, existingTags)
	systemPrompt := s.config.TagPrompt
	if systemPrompt == "" {
		systemPrompt = DefaultTagPrompt
	}

	jsonSchema := json.RawMessage(`{"type":"object","properties":{"tags":{"type":"array","items":{"type":"object","properties":{"name":{"type":"string"},"is_new":{"type":"boolean"},"reason":{"type":"string"}},"required":["name","is_new","reason"]}}},"required":["tags"]}`)

	resp, err := s.provider.Complete(ctx, &CompletionRequest{
		SystemPrompt: s.withLanguage(systemPrompt),
		UserPrompt:   prompt,
		JSONSchema:   jsonSchema,
		Model:        s.config.Model,
	})
	if err != nil {
		return nil, fmt.Errorf("ai suggest tags: %w", err)
	}

	var result TagSuggestion
	if err := json.Unmarshal(resp.JSONResult, &result); err != nil {
		return nil, fmt.Errorf("ai suggest tags: parse response: %w", err)
	}

	return &result, nil
}

// ExploreSchema analyzes objects to suggest schema improvements.
func (s *AIService) ExploreSchema(ctx context.Context, schema SchemaContext, objects []ObjectContext) (*SchemaExploration, error) {
	prompt := s.buildExplorePrompt(schema, objects)
	systemPrompt := s.config.ExplorePrompt
	if systemPrompt == "" {
		systemPrompt = DefaultExplorePrompt
	}

	jsonSchema := json.RawMessage(`{"type":"object","properties":{"suggestions":{"type":"array","items":{"type":"object","properties":{"type":{"type":"string","enum":["add","modify","remove"]},"property_name":{"type":"string"},"property_type":{"type":"string"},"reason":{"type":"string"},"description":{"type":"string"}},"required":["type","property_name","reason"]}}},"required":["suggestions"]}`)

	resp, err := s.provider.Complete(ctx, &CompletionRequest{
		SystemPrompt: s.withLanguage(systemPrompt),
		UserPrompt:   prompt,
		JSONSchema:   jsonSchema,
		Model:        s.config.Model,
	})
	if err != nil {
		return nil, fmt.Errorf("ai explore schema: %w", err)
	}

	var result SchemaExploration
	if err := json.Unmarshal(resp.JSONResult, &result); err != nil {
		return nil, fmt.Errorf("ai explore schema: parse response: %w", err)
	}

	return &result, nil
}

func (s *AIService) buildDescribePrompt(obj ObjectContext, schema SchemaContext) string {
	var b strings.Builder

	fmt.Fprintf(&b, "## Object\n\n")
	fmt.Fprintf(&b, "**Type:** %s\n", schema.TypeName)
	fmt.Fprintf(&b, "**Name:** %s\n\n", obj.Name)

	if len(obj.Properties) > 0 {
		fmt.Fprintf(&b, "### Properties\n\n")
		for k, v := range obj.Properties {
			fmt.Fprintf(&b, "- **%s:** %v\n", k, v)
		}
		b.WriteString("\n")
	}

	if obj.Body != "" {
		fmt.Fprintf(&b, "### Body\n\n%s\n\n", obj.Body)
	}

	if schema.TypeDescription != "" || len(schema.PropertyDescriptions) > 0 {
		fmt.Fprintf(&b, "## Schema Context\n\n")
		if schema.TypeDescription != "" {
			fmt.Fprintf(&b, "**Type description:** %s\n\n", schema.TypeDescription)
		}
		if len(schema.PropertyDescriptions) > 0 {
			fmt.Fprintf(&b, "### Property Descriptions\n\n")
			for name, desc := range schema.PropertyDescriptions {
				fmt.Fprintf(&b, "- **%s:** %s\n", name, desc)
			}
		}
	}

	return b.String()
}

func (s *AIService) buildTagPrompt(obj ObjectContext, schema SchemaContext, existingTags []TagInfo) string {
	var b strings.Builder

	fmt.Fprintf(&b, "## Object\n\n")
	fmt.Fprintf(&b, "**Type:** %s\n", schema.TypeName)
	fmt.Fprintf(&b, "**Name:** %s\n\n", obj.Name)

	if len(obj.Properties) > 0 {
		fmt.Fprintf(&b, "### Properties\n\n")
		for k, v := range obj.Properties {
			fmt.Fprintf(&b, "- **%s:** %v\n", k, v)
		}
		b.WriteString("\n")
	}

	if obj.Body != "" {
		fmt.Fprintf(&b, "### Body\n\n%s\n\n", obj.Body)
	}

	if len(existingTags) > 0 {
		fmt.Fprintf(&b, "## Existing Tags\n\n")
		for _, tag := range existingTags {
			if tag.Description != "" {
				fmt.Fprintf(&b, "- **%s:** %s\n", tag.Name, tag.Description)
			} else {
				fmt.Fprintf(&b, "- %s\n", tag.Name)
			}
		}
		b.WriteString("\n")
	}

	if schema.TypeDescription != "" || len(schema.PropertyDescriptions) > 0 {
		fmt.Fprintf(&b, "## Schema Context\n\n")
		if schema.TypeDescription != "" {
			fmt.Fprintf(&b, "**Type description:** %s\n\n", schema.TypeDescription)
		}
	}

	return b.String()
}

func (s *AIService) buildExplorePrompt(schema SchemaContext, objects []ObjectContext) string {
	var b strings.Builder

	fmt.Fprintf(&b, "## Type Schema: %s\n\n", schema.TypeName)
	if schema.TypeDescription != "" {
		fmt.Fprintf(&b, "**Description:** %s\n\n", schema.TypeDescription)
	}

	if len(schema.PropertyNames) > 0 {
		fmt.Fprintf(&b, "### Current Properties\n\n")
		for _, name := range schema.PropertyNames {
			propType := schema.PropertyTypes[name]
			desc := schema.PropertyDescriptions[name]
			if desc != "" {
				fmt.Fprintf(&b, "- **%s** (%s): %s\n", name, propType, desc)
			} else {
				fmt.Fprintf(&b, "- **%s** (%s)\n", name, propType)
			}
		}
		b.WriteString("\n")
	}

	fmt.Fprintf(&b, "## Sample Objects (%d)\n\n", len(objects))
	for i, obj := range objects {
		fmt.Fprintf(&b, "### Object %d: %s\n\n", i+1, obj.Name)
		if len(obj.Properties) > 0 {
			for k, v := range obj.Properties {
				fmt.Fprintf(&b, "- **%s:** %v\n", k, v)
			}
			b.WriteString("\n")
		}
		if obj.Body != "" {
			fmt.Fprintf(&b, "**Body excerpt:**\n%s\n\n", obj.Body)
		}
	}

	return b.String()
}
