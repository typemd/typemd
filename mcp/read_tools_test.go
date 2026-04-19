package mcp

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	mcplib "github.com/mark3labs/mcp-go/mcp"

	"github.com/typemd/typemd/core"
)

func callHandler(t *testing.T, handler func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error), args map[string]any, out any) *mcplib.CallToolResult {
	t.Helper()
	req := mcplib.CallToolRequest{}
	req.Params.Arguments = args
	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		return result
	}
	if out != nil {
		text := result.Content[0].(mcplib.TextContent).Text
		if err := json.Unmarshal([]byte(text), out); err != nil {
			t.Fatalf("unmarshal: %v (text=%s)", err, text)
		}
	}
	return result
}

func createBook(t *testing.T, v *core.Vault, name, updatedAt string) *core.Object {
	t.Helper()
	obj, err := v.NewObject("book", name, "")
	if err != nil {
		t.Fatalf("NewObject(%s): %v", name, err)
	}
	if updatedAt != "" {
		obj.Properties[core.UpdatedAtProperty] = updatedAt
	}
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("SaveObject(%s): %v", name, err)
	}
	return obj
}

// --- vault_overview ---

func TestVaultOverviewHandler_EmptyVault(t *testing.T) {
	v := newTestVault(t, "")

	var entries []overviewEntry
	callHandler(t, vaultOverviewHandler(v), map[string]any{}, &entries)

	if len(entries) < 2 {
		t.Fatalf("expected at least tag + page built-in types, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Count != 0 {
			t.Errorf("%s: expected count=0, got %d", e.Name, e.Count)
		}
		if len(e.Recent) != 0 {
			t.Errorf("%s: expected empty recent, got %d entries", e.Name, len(e.Recent))
		}
	}
}

func TestVaultOverviewHandler_ReportsCountsAndRecent(t *testing.T) {
	v, _ := setupTestVault(t)
	now := time.Now().UTC()
	for i := range 7 {
		ts := now.Add(time.Duration(i) * time.Minute).Format(time.RFC3339)
		createBook(t, v, "book-"+string(rune('a'+i)), ts)
	}

	var entries []overviewEntry
	callHandler(t, vaultOverviewHandler(v), map[string]any{}, &entries)

	var book *overviewEntry
	for i := range entries {
		if entries[i].Name == "book" {
			book = &entries[i]
			break
		}
	}
	if book == nil {
		t.Fatal("expected 'book' in overview")
	}
	// setupTestVault already created clean-code, so we now have 8 book objects.
	if book.Count != 8 {
		t.Errorf("expected book count=8, got %d", book.Count)
	}
	if len(book.Recent) != recentPerType {
		t.Fatalf("expected recent capped at %d, got %d", recentPerType, len(book.Recent))
	}
	// Verify recent list is sorted by updated_at desc.
	for i := 1; i < len(book.Recent); i++ {
		if book.Recent[i-1].UpdatedAt < book.Recent[i].UpdatedAt {
			t.Errorf("recent not sorted desc: [%d]=%s < [%d]=%s", i-1, book.Recent[i-1].UpdatedAt, i, book.Recent[i].UpdatedAt)
		}
	}
}

// --- list_objects ---

func TestListObjectsHandler_TypeFilterAndPagination(t *testing.T) {
	v, _ := setupTestVault(t)
	for i := range 19 {
		createBook(t, v, "book-"+string(rune('a'+i)), time.Now().Add(time.Duration(i)*time.Second).Format(time.RFC3339))
	}

	var resp listObjectsResponse
	callHandler(t, listObjectsHandler(v), map[string]any{
		"type":   "book",
		"limit":  float64(10),
		"offset": float64(0),
	}, &resp)

	if resp.Total != 20 {
		t.Errorf("expected total=20, got %d", resp.Total)
	}
	if len(resp.Objects) != 10 {
		t.Errorf("expected 10 summaries, got %d", len(resp.Objects))
	}
	for _, s := range resp.Objects {
		if s.Type != "book" {
			t.Errorf("expected type=book, got %s", s.Type)
		}
	}
}

func TestListObjectsHandler_UnknownType(t *testing.T) {
	v, _ := setupTestVault(t)

	var resp listObjectsResponse
	callHandler(t, listObjectsHandler(v), map[string]any{"type": "nonexistent"}, &resp)

	if resp.Total != 0 {
		t.Errorf("expected total=0, got %d", resp.Total)
	}
	if len(resp.Objects) != 0 {
		t.Errorf("expected empty list, got %d", len(resp.Objects))
	}
}

func TestListObjectsHandler_LimitClamped(t *testing.T) {
	v, _ := setupTestVault(t)

	var resp listObjectsResponse
	callHandler(t, listObjectsHandler(v), map[string]any{"limit": float64(10000)}, &resp)

	if resp.Limit != maxListLimit {
		t.Errorf("expected clamped limit=%d, got %d", maxListLimit, resp.Limit)
	}
}

// --- query_objects ---

func TestQueryObjectsHandler_PropertyEquality(t *testing.T) {
	v := newTestVault(t, "name: book\nproperties:\n  - name: status\n    type: string\n")
	obj, err := v.NewObject("book", "reading-book", "")
	if err != nil {
		t.Fatalf("NewObject: %v", err)
	}
	obj.Properties["status"] = "reading"
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("SaveObject: %v", err)
	}
	other, err := v.NewObject("book", "done-book", "")
	if err != nil {
		t.Fatalf("NewObject: %v", err)
	}
	other.Properties["status"] = "completed"
	if err := v.SaveObject(other); err != nil {
		t.Fatalf("SaveObject: %v", err)
	}

	var resp listObjectsResponse
	callHandler(t, queryObjectsHandler(v), map[string]any{
		"filters": []any{
			map[string]any{"property": "status", "operator": "is", "value": "reading"},
		},
	}, &resp)

	if resp.Total != 1 {
		t.Fatalf("expected total=1, got %d", resp.Total)
	}
	if !strings.HasPrefix(resp.Objects[0].ID, "book/reading-book-") {
		t.Errorf("expected reading-book, got %s", resp.Objects[0].ID)
	}
}

func TestQueryObjectsHandler_SortAndLimit(t *testing.T) {
	v, _ := setupTestVault(t)
	for i := range 8 {
		createBook(t, v, "sort-book-"+string(rune('a'+i)), time.Now().Add(time.Duration(i)*time.Minute).Format(time.RFC3339))
	}

	var resp listObjectsResponse
	callHandler(t, queryObjectsHandler(v), map[string]any{
		"filters": []any{
			map[string]any{"property": "type", "operator": "is", "value": "book"},
		},
		"sort": []any{
			map[string]any{"property": core.UpdatedAtProperty, "direction": "desc"},
		},
		"limit": float64(5),
	}, &resp)

	if len(resp.Objects) != 5 {
		t.Fatalf("expected 5 results, got %d", len(resp.Objects))
	}
	for i := 1; i < len(resp.Objects); i++ {
		if resp.Objects[i-1].UpdatedAt < resp.Objects[i].UpdatedAt {
			t.Errorf("results not sorted desc: %s < %s", resp.Objects[i-1].UpdatedAt, resp.Objects[i].UpdatedAt)
		}
	}
}

func TestQueryObjectsHandler_InvalidFilter(t *testing.T) {
	v, _ := setupTestVault(t)

	result := callHandler(t, queryObjectsHandler(v), map[string]any{
		"filters": []any{
			map[string]any{"operator": "is", "value": "book"}, // missing property
		},
	}, nil)
	if !result.IsError {
		t.Fatal("expected tool error for missing property")
	}
}

func TestQueryObjectsHandler_FiltersRequired(t *testing.T) {
	v, _ := setupTestVault(t)

	result := callHandler(t, queryObjectsHandler(v), map[string]any{}, nil)
	if !result.IsError {
		t.Fatal("expected tool error when filters is absent")
	}
}

// --- list_backlinks ---

func TestListBacklinksHandler_WikiBacklinks(t *testing.T) {
	v, _ := setupTestVault(t)
	target, err := v.NewObject("book", "target-book", "")
	if err != nil {
		t.Fatalf("NewObject target: %v", err)
	}

	// Create three objects that wiki-link to the target.
	for i := range 3 {
		src, err := v.NewObject("book", "source-"+string(rune('a'+i)), "")
		if err != nil {
			t.Fatalf("NewObject src: %v", err)
		}
		src.Body = "See [[" + target.ID + "]]"
		if err := v.SaveObject(src); err != nil {
			t.Fatalf("SaveObject src: %v", err)
		}
	}
	events, _, err := v.Reconcile()
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if err := v.Project(events); err != nil {
		t.Fatalf("Project: %v", err)
	}

	var resp listBacklinksResponse
	callHandler(t, listBacklinksHandler(v), map[string]any{"id": target.ID}, &resp)

	if len(resp.WikiBacklinks) != 3 {
		t.Errorf("expected 3 wiki backlinks, got %d", len(resp.WikiBacklinks))
	}
	if len(resp.RelationBacklinks) != 0 {
		t.Errorf("expected 0 relation backlinks, got %d", len(resp.RelationBacklinks))
	}
}

func TestListBacklinksHandler_RelationBacklinks(t *testing.T) {
	v, _, _ := setupTestVaultWithRelation(t)
	target, err := v.NewObject("book", "target-book", "")
	if err != nil {
		t.Fatalf("NewObject target: %v", err)
	}
	src1, err := v.NewObject("book", "src-one", "")
	if err != nil {
		t.Fatalf("NewObject src1: %v", err)
	}
	src2, err := v.NewObject("book", "src-two", "")
	if err != nil {
		t.Fatalf("NewObject src2: %v", err)
	}
	if err := v.LinkObjects(src1.ID, "author", target.ID); err != nil {
		t.Fatalf("link src1: %v", err)
	}
	if err := v.LinkObjects(src2.ID, "author", target.ID); err != nil {
		t.Fatalf("link src2: %v", err)
	}

	var resp listBacklinksResponse
	callHandler(t, listBacklinksHandler(v), map[string]any{"id": target.ID}, &resp)

	if len(resp.RelationBacklinks) != 2 {
		t.Fatalf("expected 2 relation backlinks, got %d", len(resp.RelationBacklinks))
	}
	for _, rb := range resp.RelationBacklinks {
		if rb.Relation != "author" {
			t.Errorf("expected relation=author, got %s", rb.Relation)
		}
	}
}

func TestListBacklinksHandler_Empty(t *testing.T) {
	v, sampleID := setupTestVault(t)

	var resp listBacklinksResponse
	callHandler(t, listBacklinksHandler(v), map[string]any{"id": sampleID}, &resp)

	if len(resp.WikiBacklinks) != 0 || len(resp.RelationBacklinks) != 0 {
		t.Errorf("expected empty backlinks, got wiki=%d relation=%d", len(resp.WikiBacklinks), len(resp.RelationBacklinks))
	}
}

func TestListBacklinksHandler_AbbreviatedID(t *testing.T) {
	v, sampleID := setupTestVault(t)
	src, err := v.NewObject("book", "referrer", "")
	if err != nil {
		t.Fatalf("NewObject: %v", err)
	}
	src.Body = "See [[" + sampleID + "]]"
	if err := v.SaveObject(src); err != nil {
		t.Fatalf("SaveObject: %v", err)
	}
	events, _, err := v.Reconcile()
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}
	if err := v.Project(events); err != nil {
		t.Fatalf("Project: %v", err)
	}

	var resp listBacklinksResponse
	callHandler(t, listBacklinksHandler(v), map[string]any{"id": "book/clean-code"}, &resp)

	if len(resp.WikiBacklinks) != 1 {
		t.Errorf("expected 1 backlink via prefix resolution, got %d", len(resp.WikiBacklinks))
	}
}

// --- vault_stats ---

func TestVaultStatsHandler_PartialFill(t *testing.T) {
	v := newTestVault(t, "name: book\nproperties:\n  - name: rating\n    type: string\n")
	for i := range 10 {
		obj, err := v.NewObject("book", "book-"+string(rune('a'+i)), "")
		if err != nil {
			t.Fatalf("NewObject: %v", err)
		}
		if i < 6 {
			obj.Properties["rating"] = "5"
		}
		if err := v.SaveObject(obj); err != nil {
			t.Fatalf("SaveObject: %v", err)
		}
	}

	var resp vaultStatsResponse
	callHandler(t, vaultStatsHandler(v), map[string]any{"type": "book"}, &resp)

	if resp.Count != 10 {
		t.Errorf("expected count=10, got %d", resp.Count)
	}
	var rating *vaultStatsProperty
	for i := range resp.Properties {
		if resp.Properties[i].Name == "rating" {
			rating = &resp.Properties[i]
			break
		}
	}
	if rating == nil {
		t.Fatal("expected rating property in stats")
	}
	if rating.Filled != 6 {
		t.Errorf("expected filled=6, got %d", rating.Filled)
	}
	if rating.FillRate < 0.59 || rating.FillRate > 0.61 {
		t.Errorf("expected fill_rate≈0.6, got %v", rating.FillRate)
	}
}

func TestVaultStatsHandler_UnknownType(t *testing.T) {
	v, _ := setupTestVault(t)

	result := callHandler(t, vaultStatsHandler(v), map[string]any{"type": "nonexistent"}, nil)
	if !result.IsError {
		t.Fatal("expected tool error for unknown type")
	}
}

func TestVaultStatsHandler_NoObjects(t *testing.T) {
	v := newTestVault(t, "name: book\nproperties:\n  - name: status\n    type: string\n")

	var resp vaultStatsResponse
	callHandler(t, vaultStatsHandler(v), map[string]any{"type": "book"}, &resp)

	if resp.Count != 0 {
		t.Errorf("expected count=0, got %d", resp.Count)
	}
	for _, p := range resp.Properties {
		if p.Filled != 0 || p.FillRate != 0 {
			t.Errorf("%s: expected filled=0 fill_rate=0, got filled=%d fill_rate=%v", p.Name, p.Filled, p.FillRate)
		}
	}
}

// --- parseFilters / parseSort unit tests ---

func TestParseFilters(t *testing.T) {
	rules, err := parseFilters([]any{
		map[string]any{"property": "status", "operator": "is", "value": "reading"},
	})
	if err != nil {
		t.Fatalf("parseFilters: %v", err)
	}
	if len(rules) != 1 || rules[0].Property != "status" || rules[0].Value != "reading" {
		t.Errorf("unexpected result: %+v", rules)
	}
}

func TestParseFiltersErrors(t *testing.T) {
	cases := []struct {
		name string
		in   any
	}{
		{"nil", nil},
		{"not array", "filters"},
		{"entry not map", []any{"invalid"}},
		{"missing property", []any{map[string]any{"operator": "is"}}},
		{"missing operator", []any{map[string]any{"property": "status"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := parseFilters(c.in); err == nil {
				t.Errorf("expected error for %s", c.name)
			}
		})
	}
}

func TestParseSort(t *testing.T) {
	rules, err := parseSort([]any{
		map[string]any{"property": "updated_at", "direction": "desc"},
	})
	if err != nil {
		t.Fatalf("parseSort: %v", err)
	}
	if len(rules) != 1 || rules[0].Direction != "desc" {
		t.Errorf("unexpected sort rules: %+v", rules)
	}
}

func TestParseSortDefaultsAndErrors(t *testing.T) {
	// default direction = asc
	rules, err := parseSort([]any{map[string]any{"property": "name"}})
	if err != nil || rules[0].Direction != "asc" {
		t.Errorf("expected default direction asc, got %+v err=%v", rules, err)
	}

	// invalid direction
	if _, err := parseSort([]any{map[string]any{"property": "name", "direction": "sideways"}}); err == nil {
		t.Error("expected error for invalid direction")
	}

	// missing property
	if _, err := parseSort([]any{map[string]any{"direction": "asc"}}); err == nil {
		t.Error("expected error for missing property")
	}
}
