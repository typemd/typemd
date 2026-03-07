package mcp

import (
	"context"
	"encoding/json"
	"fmt"

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
		mcplib.WithDescription("Get an object by its ID (e.g. book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz)"),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Object ID in type/slug-ulid format"),
		),
	)
	s.AddTool(getObjectTool, getObjectHandler(vault))
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

		obj, err := vault.GetObject(id)
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
