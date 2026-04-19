package mcp

import (
	"context"
	"fmt"
	"sort"
	"strings"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/typemd/typemd/core"
)

const (
	defaultListLimit = 50
	maxListLimit     = 500
	recentPerType    = 5
	directionAsc     = "asc"
	directionDesc    = "desc"
)

type objectListSummary struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type overviewEntry struct {
	Name        string              `json:"name"`
	Plural      string              `json:"plural,omitempty"`
	Emoji       string              `json:"emoji,omitempty"`
	Description string              `json:"description,omitempty"`
	Count       int                 `json:"count"`
	Recent      []objectListSummary `json:"recent"`
}

type listObjectsResponse struct {
	Total   int                 `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	Objects []objectListSummary `json:"objects"`
}

type backlinkWiki struct {
	FromID      string `json:"from_id"`
	DisplayText string `json:"display_text,omitempty"`
	Target      string `json:"target,omitempty"`
}

type backlinkRelation struct {
	FromID   string `json:"from_id"`
	Relation string `json:"relation"`
}

type listBacklinksResponse struct {
	WikiBacklinks     []backlinkWiki     `json:"wiki_backlinks"`
	RelationBacklinks []backlinkRelation `json:"relation_backlinks"`
}

type vaultStatsProperty struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Filled   int     `json:"filled"`
	Total    int     `json:"total"`
	FillRate float64 `json:"fill_rate"`
	Stats    any     `json:"stats,omitempty"`
}

type vaultStatsResponse struct {
	Type       string               `json:"type"`
	Emoji      string               `json:"emoji,omitempty"`
	Plural     string               `json:"plural,omitempty"`
	Count      int                  `json:"count"`
	Properties []vaultStatsProperty `json:"properties"`
}

func summariseObject(obj *core.Object) objectListSummary {
	return objectListSummary{
		ID:        obj.ID,
		Type:      obj.Type,
		Name:      propertyString(obj, "name"),
		UpdatedAt: propertyString(obj, core.UpdatedAtProperty),
	}
}

func propertyString(obj *core.Object, key string) string {
	if obj == nil || obj.Properties == nil {
		return ""
	}
	v, ok := obj.Properties[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func clampPagination(rawLimit, rawOffset any) (limit, offset int) {
	limit = defaultListLimit
	if n, ok := toInt(rawLimit); ok && n > 0 {
		limit = n
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if n, ok := toInt(rawOffset); ok && n > 0 {
		offset = n
	}
	return limit, offset
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	}
	return 0, false
}

func paginateSummaries(objects []*core.Object, limit, offset int) ([]objectListSummary, int) {
	total := len(objects)
	if offset > total {
		return []objectListSummary{}, total
	}
	end := min(offset+limit, total)
	slice := objects[offset:end]
	summaries := make([]objectListSummary, len(slice))
	for i, obj := range slice {
		summaries[i] = summariseObject(obj)
	}
	return summaries, total
}

func vaultOverviewHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		typeNames := vault.ListTypes()
		sort.Strings(typeNames)

		all, err := vault.QueryObjects(nil, core.QuerySort(core.SortRule{Property: core.UpdatedAtProperty, Direction: directionDesc}))
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("query objects: %v", err)), nil
		}

		byType := make(map[string][]*core.Object, len(typeNames))
		for _, obj := range all {
			byType[obj.Type] = append(byType[obj.Type], obj)
		}

		entries := make([]overviewEntry, 0, len(typeNames))
		for _, name := range typeNames {
			entry := overviewEntry{Name: name, Recent: []objectListSummary{}}

			if schema, err := vault.LoadType(name); err == nil && schema != nil {
				entry.Plural = schema.Plural
				entry.Emoji = schema.Emoji
				entry.Description = schema.Description
			}

			objects := byType[name]
			entry.Count = len(objects)
			cap := min(recentPerType, len(objects))
			for i := range cap {
				entry.Recent = append(entry.Recent, summariseObject(objects[i]))
			}
			entries = append(entries, entry)
		}

		return marshalResult(entries)
	}
}

func listObjectsHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		args := request.GetArguments()
		typeName, _ := args["type"].(string)
		limit, offset := clampPagination(args["limit"], args["offset"])

		var filters []core.FilterRule
		if typeName != "" {
			filters = core.TypeFilter(typeName)
		}

		objects, err := vault.QueryObjects(filters, core.QuerySort(core.SortRule{Property: core.UpdatedAtProperty, Direction: directionDesc}))
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("list objects failed: %v", err)), nil
		}

		summaries, total := paginateSummaries(objects, limit, offset)
		return marshalResult(listObjectsResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			Objects: summaries,
		})
	}
}

func queryObjectsHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		args := request.GetArguments()

		filters, err := parseFilters(args["filters"])
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}
		sorts, err := parseSort(args["sort"])
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		limit, offset := clampPagination(args["limit"], args["offset"])

		opts := []core.QueryOption{}
		if len(sorts) > 0 {
			opts = append(opts, core.QuerySort(sorts...))
		}

		objects, err := vault.QueryObjects(filters, opts...)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("query failed: %v", err)), nil
		}

		summaries, total := paginateSummaries(objects, limit, offset)
		return marshalResult(listObjectsResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			Objects: summaries,
		})
	}
}

func parseFilters(raw any) ([]core.FilterRule, error) {
	if raw == nil {
		return nil, fmt.Errorf("filters is required")
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("filters must be an array")
	}
	rules := make([]core.FilterRule, 0, len(arr))
	for i, entry := range arr {
		m, ok := entry.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("filter %d must be an object", i)
		}
		property, _ := m["property"].(string)
		operator, _ := m["operator"].(string)
		if strings.TrimSpace(property) == "" {
			return nil, fmt.Errorf("filter %d missing 'property'", i)
		}
		if strings.TrimSpace(operator) == "" {
			return nil, fmt.Errorf("filter %d missing 'operator'", i)
		}
		value := ""
		switch v := m["value"].(type) {
		case string:
			value = v
		case nil:
			value = ""
		default:
			value = fmt.Sprintf("%v", v)
		}
		rules = append(rules, core.FilterRule{
			Property: property,
			Operator: operator,
			Value:    value,
		})
	}
	return rules, nil
}

func parseSort(raw any) ([]core.SortRule, error) {
	if raw == nil {
		return nil, nil
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("sort must be an array")
	}
	rules := make([]core.SortRule, 0, len(arr))
	for i, entry := range arr {
		m, ok := entry.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("sort %d must be an object", i)
		}
		property, _ := m["property"].(string)
		direction, _ := m["direction"].(string)
		if strings.TrimSpace(property) == "" {
			return nil, fmt.Errorf("sort %d missing 'property'", i)
		}
		if direction == "" {
			direction = directionAsc
		}
		if direction != directionAsc && direction != directionDesc {
			return nil, fmt.Errorf("sort %d invalid direction %q (expected asc or desc)", i, direction)
		}
		rules = append(rules, core.SortRule{Property: property, Direction: direction})
	}
	return rules, nil
}

func listBacklinksHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := request.RequireString("id")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		obj, err := vault.ResolveObject(id)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("resolve %q: %v", id, err)), nil
		}
		resolvedID := obj.ID

		wiki, err := vault.ListBacklinks(resolvedID)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("list wiki backlinks: %v", err)), nil
		}

		relations, err := vault.ListRelations(resolvedID)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("list relations: %v", err)), nil
		}

		response := listBacklinksResponse{
			WikiBacklinks:     make([]backlinkWiki, 0, len(wiki)),
			RelationBacklinks: []backlinkRelation{},
		}

		for _, wl := range wiki {
			response.WikiBacklinks = append(response.WikiBacklinks, backlinkWiki{
				FromID:      wl.FromID,
				DisplayText: wl.DisplayText,
				Target:      wl.Target,
			})
		}

		for _, r := range relations {
			if r.ToID != resolvedID {
				continue
			}
			response.RelationBacklinks = append(response.RelationBacklinks, backlinkRelation{
				FromID:   r.FromID,
				Relation: r.Name,
			})
		}

		return marshalResult(response)
	}
}

func vaultStatsHandler(vault *core.Vault) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		typeName, err := request.RequireString("type")
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		stats, err := vault.TypeStats(typeName)
		if err != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("stats failed: %v", err)), nil
		}

		response := vaultStatsResponse{
			Type:       stats.TypeName,
			Emoji:      stats.Emoji,
			Plural:     stats.Plural,
			Count:      stats.Count,
			Properties: make([]vaultStatsProperty, 0, len(stats.Properties)),
		}

		for _, p := range stats.Properties {
			rate := 0.0
			if p.Total > 0 {
				rate = float64(p.Filled) / float64(p.Total)
			}
			response.Properties = append(response.Properties, vaultStatsProperty{
				Name:     p.Name,
				Type:     p.Type,
				Filled:   p.Filled,
				Total:    p.Total,
				FillRate: rate,
				Stats:    p.Stats,
			})
		}

		return marshalResult(response)
	}
}

