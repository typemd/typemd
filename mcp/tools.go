package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/typemd/typemd/core"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer, vault *core.Vault) {
	searchTool := mcplib.NewTool("search",
		mcplib.WithDescription("Search objects in the vault using full-text search"),
		mcplib.WithString("keyword",
			mcplib.Required(),
			mcplib.Description("Search keyword"),
		),
	)
	s.AddTool(searchTool, searchHandler(vault))

	getObjectTool := mcplib.NewTool("get_object",
		mcplib.WithDescription("Get an object by its full or abbreviated ID. Supports prefix matching (e.g. book/clean-code resolves to the full ULID-suffixed ID)."),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Object ID or prefix (e.g. book/clean-code or book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz)"),
		),
	)
	s.AddTool(getObjectTool, getObjectHandler(vault))

	listTemplatesTool := mcplib.NewTool("list_templates",
		mcplib.WithDescription("List available object templates. Optionally filter by type name."),
		mcplib.WithString("type",
			mcplib.Description("Type name to filter templates (e.g. book). If omitted, lists templates across all types."),
		),
	)
	s.AddTool(listTemplatesTool, listTemplatesHandler(vault))

	getTemplateTool := mcplib.NewTool("get_template",
		mcplib.WithDescription("Get a specific template's content including properties and body."),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Type name (e.g. book)"),
		),
		mcplib.WithString("name",
			mcplib.Required(),
			mcplib.Description("Template name (e.g. reading)"),
		),
	)
	s.AddTool(getTemplateTool, getTemplateHandler(vault))
}

type objectSummary struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Filename string `json:"filename"`
}

func searchHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		keyword, err := request.RequireString("keyword")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		results, err := vault.SearchObjects(keyword)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
		}

		summaries := make([]objectSummary, len(results))
		for i, obj := range results {
			summaries[i] = objectSummary{
				ID:       obj.ID,
				Type:     obj.Type,
				Filename: obj.Filename,
			}
		}

		data, err := json.Marshal(summaries)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("marshal result: %v", err)), nil
		}

		return mcplib.NewToolResultText(string(data)), nil
	}
}

type objectDetail struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Filename   string         `json:"filename"`
	Properties map[string]any `json:"properties"`
	Body       string         `json:"body"`
}

func getObjectHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := request.RequireString("id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		obj, err := vault.ResolveObject(id)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("get object failed: %v", err)), nil
		}

		detail := objectDetail{
			ID:         obj.ID,
			Type:       obj.Type,
			Filename:   obj.Filename,
			Properties: obj.Properties,
			Body:       obj.Body,
		}

		data, err := json.Marshal(detail)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("marshal result: %v", err)), nil
		}

		return mcplib.NewToolResultText(string(data)), nil
	}
}

type templateSummary struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func listTemplatesHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		typeName, _ := request.RequireString("type")

		var typeNames []string
		if typeName != "" {
			typeNames = []string{typeName}
		} else {
			typeNames = vault.ListTypes()
			sort.Strings(typeNames)
		}

		summaries := []templateSummary{}
		for _, tn := range typeNames {
			names, err := vault.ListTemplates(tn)
			if err != nil {
				continue
			}
			for _, name := range names {
				summaries = append(summaries, templateSummary{Type: tn, Name: name})
			}
		}

		data, err := json.Marshal(summaries)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("marshal result: %v", err)), nil
		}

		return mcplib.NewToolResultText(string(data)), nil
	}
}

type templateDetail struct {
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
	Body       string         `json:"body"`
}

func getTemplateHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		typeName, err := request.RequireString("type")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		name, err := request.RequireString("name")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		tmpl, err := vault.LoadTemplate(typeName, name)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("get template failed: %v", err)), nil
		}

		detail := templateDetail{
			Type:       typeName,
			Name:       tmpl.Name,
			Properties: tmpl.Properties,
			Body:       tmpl.Body,
		}

		data, err := json.Marshal(detail)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("marshal result: %v", err)), nil
		}

		return mcplib.NewToolResultText(string(data)), nil
	}
}
