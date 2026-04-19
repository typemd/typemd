package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

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

	listTypesTool := mcplib.NewTool("list_types",
		mcplib.WithDescription("List all available type schemas in the vault, including built-in types (tag, page) and custom types with their schema metadata (plural, emoji, properties)."),
	)
	s.AddTool(listTypesTool, listTypesHandler(vault))

	createObjectTool := mcplib.NewTool("create_object",
		mcplib.WithDescription("Create a new object in the vault."),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Type name (e.g. book, person). Use list_types to see available types."),
		),
		mcplib.WithString("name",
			mcplib.Required(),
			mcplib.Description("Object name used as the filename slug (e.g. clean-code)"),
		),
		mcplib.WithString("template",
			mcplib.Description("Optional template name to apply during creation"),
		),
		mcplib.WithObject("properties",
			mcplib.Description("Optional properties to set on the created object (JSON object)"),
		),
		mcplib.WithString("body",
			mcplib.Description("Optional markdown body content"),
		),
	)
	s.AddTool(createObjectTool, createObjectHandler(vault))

	updateObjectTool := mcplib.NewTool("update_object",
		mcplib.WithDescription("Update an existing object's properties and/or body. Properties are merged (only provided keys are updated). Body replaces entirely if provided."),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Object ID or prefix (e.g. book/clean-code)"),
		),
		mcplib.WithObject("properties",
			mcplib.Description("Properties to merge into the object (only provided keys are updated)"),
		),
		mcplib.WithString("body",
			mcplib.Description("New markdown body content (replaces existing body)"),
		),
	)
	s.AddTool(updateObjectTool, updateObjectHandler(vault))

	linkObjectsTool := mcplib.NewTool("link_objects",
		mcplib.WithDescription("Create a relation between two objects."),
		mcplib.WithString("from_id",
			mcplib.Required(),
			mcplib.Description("Source object ID or prefix"),
		),
		mcplib.WithString("relation",
			mcplib.Required(),
			mcplib.Description("Relation property name defined in the type schema"),
		),
		mcplib.WithString("to_id",
			mcplib.Required(),
			mcplib.Description("Target object ID or prefix"),
		),
	)
	s.AddTool(linkObjectsTool, linkObjectsHandler(vault))

	unlinkObjectsTool := mcplib.NewTool("unlink_objects",
		mcplib.WithDescription("Remove a relation between two objects."),
		mcplib.WithString("from_id",
			mcplib.Required(),
			mcplib.Description("Source object ID or prefix"),
		),
		mcplib.WithString("relation",
			mcplib.Required(),
			mcplib.Description("Relation property name defined in the type schema"),
		),
		mcplib.WithString("to_id",
			mcplib.Required(),
			mcplib.Description("Target object ID or prefix"),
		),
		mcplib.WithBoolean("both",
			mcplib.Description("Remove the relation in both directions (default: false)"),
		),
	)
	s.AddTool(unlinkObjectsTool, unlinkObjectsHandler(vault))

	vaultOverviewTool := mcplib.NewTool("vault_overview",
		mcplib.WithDescription("Summarise the vault's structure in a single call. Returns every registered type with its plural, emoji, description, object count, and a short list of recent objects. Use this to orient yourself to an unfamiliar vault."),
	)
	s.AddTool(vaultOverviewTool, vaultOverviewHandler(vault))

	listObjectsTool := mcplib.NewTool("list_objects",
		mcplib.WithDescription("List object summaries with optional type filter and pagination. Returns id, type, name, and updated_at per object plus a total count. Archived objects are excluded."),
		mcplib.WithString("type",
			mcplib.Description("Optional type filter (e.g. book). Omit to list across all types."),
		),
		mcplib.WithNumber("limit",
			mcplib.Description("Maximum number of summaries to return. Default 50, clamped to 500."),
		),
		mcplib.WithNumber("offset",
			mcplib.Description("Number of summaries to skip. Default 0."),
		),
	)
	s.AddTool(listObjectsTool, listObjectsHandler(vault))

	queryObjectsTool := mcplib.NewTool("query_objects",
		mcplib.WithDescription("Run a structured query using FilterRule semantics. Each filter is {property, operator, value}; operators follow the query index vocabulary (is, is_not, contains, before, after, etc.). Supports sort and pagination."),
		mcplib.WithArray("filters",
			mcplib.Required(),
			mcplib.Description("Filter rules: array of {property, operator, value} objects."),
		),
		mcplib.WithArray("sort",
			mcplib.Description("Optional sort rules: array of {property, direction} with direction 'asc' or 'desc'."),
		),
		mcplib.WithNumber("limit",
			mcplib.Description("Maximum number of summaries to return. Default 50, clamped to 500."),
		),
		mcplib.WithNumber("offset",
			mcplib.Description("Number of summaries to skip. Default 0."),
		),
	)
	s.AddTool(queryObjectsTool, queryObjectsHandler(vault))

	listBacklinksTool := mcplib.NewTool("list_backlinks",
		mcplib.WithDescription("List all references that point to an object. Response splits wiki-link backlinks from typed relation backlinks."),
		mcplib.WithString("id",
			mcplib.Required(),
			mcplib.Description("Target object ID or prefix (e.g. book/clean-code)."),
		),
	)
	s.AddTool(listBacklinksTool, listBacklinksHandler(vault))

	vaultStatsTool := mcplib.NewTool("vault_stats",
		mcplib.WithDescription("Return per-property distribution statistics for a single type: fill counts, fill rates, and the existing per-type analytics (number/select/date distributions)."),
		mcplib.WithString("type",
			mcplib.Required(),
			mcplib.Description("Type name (e.g. book)."),
		),
	)
	s.AddTool(vaultStatsTool, vaultStatsHandler(vault))
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

		return marshalResult(summaries)
	}
}

type objectDetail struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Filename   string         `json:"filename"`
	Properties map[string]any `json:"properties"`
	Body       string         `json:"body"`
}

type typeSummary struct {
	Name       string             `json:"name"`
	Plural     string             `json:"plural,omitempty"`
	Emoji      string             `json:"emoji,omitempty"`
	Properties []propertySummary  `json:"properties,omitempty"`
}

type propertySummary struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Target   string `json:"target,omitempty"`
	Multiple bool   `json:"multiple,omitempty"`
}

func marshalResult(v any) (*mcplib.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("marshal result: %v", err)), nil
	}
	return mcplib.NewToolResultText(string(data)), nil
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

		return marshalResult(detail)
	}
}

func listTypesHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		names := vault.ListTypes()
		summaries := make([]typeSummary, 0, len(names))

		for _, name := range names {
			schema, err := vault.LoadType(name)
			if err != nil {
				continue
			}
			ts := typeSummary{
				Name:   schema.Name,
				Plural: schema.Plural,
				Emoji:  schema.Emoji,
			}
			for _, p := range schema.Properties {
				ts.Properties = append(ts.Properties, propertySummary{
					Name:     p.Name,
					Type:     p.Type,
					Target:   p.Target,
					Multiple: p.Multiple,
				})
			}
			summaries = append(summaries, ts)
		}

		return marshalResult(summaries)
	}
}

func createObjectHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		typeName, err := request.RequireString("type")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		name, err := request.RequireString("name")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		args := request.GetArguments()
		template, _ := args["template"].(string)

		obj, err := vault.NewObject(typeName, name, template)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("create object failed: %v", err)), nil
		}

		props, _ := args["properties"].(map[string]any)
		body, hasBody := args["body"].(string)

		if len(props) > 0 || hasBody {
			maps.Copy(obj.Properties, props)
			if hasBody {
				obj.Body = body
			}
			if err := vault.SaveObject(obj); err != nil {
				return mcplib.NewToolResultError(fmt.Sprintf("save object failed: %v", err)), nil
			}
		}

		return marshalResult(objectSummary{
			ID:       obj.ID,
			Type:     obj.Type,
			Filename: obj.Filename,
		})
	}
}

func updateObjectHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := request.RequireString("id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		obj, err := vault.ResolveObject(id)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("object not found: %v", err)), nil
		}

		args := request.GetArguments()
		props, _ := args["properties"].(map[string]any)
		body, hasBody := args["body"].(string)

		if len(props) == 0 && !hasBody {
			return marshalResult(map[string]string{"id": obj.ID})
		}

		maps.Copy(obj.Properties, props)
		if hasBody {
			obj.Body = body
		}

		if err := vault.SaveObject(obj); err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("update object failed: %v", err)), nil
		}

		return marshalResult(map[string]string{"id": obj.ID})
	}
}

func linkObjectsHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		fromID, err := request.RequireString("from_id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		relation, err := request.RequireString("relation")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		toID, err := request.RequireString("to_id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		if err := vault.LinkObjects(fromID, relation, toID); err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("link failed: %v", err)), nil
		}

		return mcplib.NewToolResultText(`{"status":"linked"}`), nil
	}
}

func unlinkObjectsHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		fromID, err := request.RequireString("from_id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		relation, err := request.RequireString("relation")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		toID, err := request.RequireString("to_id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		args := request.GetArguments()
		both, _ := args["both"].(bool)

		if err := vault.UnlinkObjects(fromID, relation, toID, both); err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("unlink failed: %v", err)), nil
		}

		return mcplib.NewToolResultText(`{"status":"unlinked"}`), nil
	}
}
