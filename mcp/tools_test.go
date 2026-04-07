package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	mcplib "github.com/mark3labs/mcp-go/mcp"
)

// newTestVault creates a temporary vault with the given type schema YAML.
func newTestVault(t *testing.T, schemaYAML string) *core.Vault {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)

	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { v.Close() })

	if schemaYAML != "" {
		typesDir := filepath.Join(dir, "types")
		if err := os.MkdirAll(filepath.Join(typesDir, "book"), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(typesDir, "book", "schema.yaml"), []byte(schemaYAML), 0644); err != nil {
			t.Fatalf("write schema: %v", err)
		}
	}

	return v
}

// setupTestVault creates a temporary vault with a type schema and sample objects.
// Returns the vault and the ID of the created sample object.
func setupTestVault(t *testing.T) (*core.Vault, string) {
	t.Helper()
	v := newTestVault(t, "name: book\nproperties:\n  - name: status\n    type: string\n")

	obj, err := v.NewObject("book", "clean-code", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	return v, obj.ID
}

func TestSearchHandler_HappyPath(t *testing.T) {
	vault, sampleID := setupTestVault(t)

	handler := searchHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "search"
	req.Params.Arguments = map[string]any{"keyword": "clean"}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	// Parse result text
	textContent := result.Content[0].(mcplib.TextContent)
	var summaries []objectSummary
	if err := json.Unmarshal([]byte(textContent.Text), &summaries); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if len(summaries) != 1 {
		t.Fatalf("expected 1 result, got %d", len(summaries))
	}
	if summaries[0].ID != sampleID {
		t.Errorf("expected ID %s, got %s", sampleID, summaries[0].ID)
	}
	if !strings.HasPrefix(summaries[0].ID, "book/clean-code-") {
		t.Errorf("expected ID to start with book/clean-code-, got %s", summaries[0].ID)
	}
}

func TestSearchHandler_EmptyKeyword(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := searchHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "search"
	req.Params.Arguments = map[string]any{"keyword": ""}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	// Empty keyword returns empty array, not error
	textContent := result.Content[0].(mcplib.TextContent)
	var summaries []objectSummary
	if err := json.Unmarshal([]byte(textContent.Text), &summaries); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if len(summaries) != 0 {
		t.Errorf("expected 0 results for empty keyword, got %d", len(summaries))
	}
}

func TestSearchHandler_NoResults(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := searchHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "search"
	req.Params.Arguments = map[string]any{"keyword": "nonexistent"}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	var summaries []objectSummary
	if err := json.Unmarshal([]byte(textContent.Text), &summaries); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if len(summaries) != 0 {
		t.Errorf("expected 0 results, got %d", len(summaries))
	}
}

func TestGetObjectHandler_HappyPath(t *testing.T) {
	vault, sampleID := setupTestVault(t)

	handler := getObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "get_object"
	req.Params.Arguments = map[string]any{"id": sampleID}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	var detail objectDetail
	if err := json.Unmarshal([]byte(textContent.Text), &detail); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if detail.ID != sampleID {
		t.Errorf("expected ID %s, got %s", sampleID, detail.ID)
	}
	if detail.Type != "book" {
		t.Errorf("expected type book, got %s", detail.Type)
	}
}

func TestGetObjectHandler_NotFound(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := getObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "get_object"
	req.Params.Arguments = map[string]any{"id": "book/nonexistent"}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for nonexistent object")
	}
}

func TestGetObjectHandler_InvalidID(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := getObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "get_object"
	req.Params.Arguments = map[string]any{"id": "invalid-id"}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for invalid ID format")
	}
}

// --- list_types ---

func TestListTypesHandler_WithCustomTypes(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := listTypesHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "list_types"
	req.Params.Arguments = map[string]any{}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	var types []typeSummary
	if err := json.Unmarshal([]byte(textContent.Text), &types); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	// Should include built-in types (tag, page) and custom type (book)
	if len(types) < 3 {
		t.Fatalf("expected at least 3 types (tag, page, book), got %d", len(types))
	}

	found := false
	for _, ts := range types {
		if ts.Name == "book" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find book type in results")
	}
}

func TestListTypesHandler_BuiltInTypes(t *testing.T) {
	v := newTestVault(t, "")

	handler := listTypesHandler(v)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "list_types"
	req.Params.Arguments = map[string]any{}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	var types []typeSummary
	if err := json.Unmarshal([]byte(textContent.Text), &types); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	// Should at least have built-in types
	if len(types) < 2 {
		t.Fatalf("expected at least 2 built-in types, got %d", len(types))
	}

	names := make(map[string]bool)
	for _, ts := range types {
		names[ts.Name] = true
	}
	if !names["tag"] {
		t.Error("expected built-in type 'tag'")
	}
	if !names["page"] {
		t.Error("expected built-in type 'page'")
	}
}

// --- create_object ---

func TestCreateObjectHandler_HappyPath(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := createObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "create_object"
	req.Params.Arguments = map[string]any{
		"type": "book",
		"name": "test-book",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	var summary objectSummary
	if err := json.Unmarshal([]byte(textContent.Text), &summary); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if summary.Type != "book" {
		t.Errorf("expected type book, got %s", summary.Type)
	}
	if !strings.HasPrefix(summary.ID, "book/test-book-") {
		t.Errorf("expected ID to start with book/test-book-, got %s", summary.ID)
	}
}

func TestCreateObjectHandler_WithPropertiesAndBody(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := createObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "create_object"
	req.Params.Arguments = map[string]any{
		"type":       "book",
		"name":       "props-book",
		"properties": map[string]any{"status": "reading"},
		"body":       "A great book.",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	var summary objectSummary
	if err := json.Unmarshal([]byte(textContent.Text), &summary); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	// Verify the object was saved with properties and body
	obj, err := vault.ResolveObject(summary.ID)
	if err != nil {
		t.Fatalf("resolve object: %v", err)
	}
	if obj.Properties["status"] != "reading" {
		t.Errorf("expected status=reading, got %v", obj.Properties["status"])
	}
	if obj.Body != "A great book." {
		t.Errorf("expected body 'A great book.', got %q", obj.Body)
	}
}

func TestCreateObjectHandler_InvalidType(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := createObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "create_object"
	req.Params.Arguments = map[string]any{
		"type": "nonexistent",
		"name": "test",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for invalid type")
	}
}

// --- update_object ---

func TestUpdateObjectHandler_MergeProperties(t *testing.T) {
	vault, sampleID := setupTestVault(t)

	// Set initial property
	obj, _ := vault.ResolveObject(sampleID)
	obj.Properties["status"] = "reading"
	obj.Properties["rating"] = "5"
	if err := vault.SaveObject(obj); err != nil {
		t.Fatalf("save: %v", err)
	}

	handler := updateObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "update_object"
	req.Params.Arguments = map[string]any{
		"id":         sampleID,
		"properties": map[string]any{"status": "completed"},
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	// Verify merge: status updated, rating preserved
	updated, _ := vault.ResolveObject(sampleID)
	if updated.Properties["status"] != "completed" {
		t.Errorf("expected status=completed, got %v", updated.Properties["status"])
	}
	if updated.Properties["rating"] != "5" {
		t.Errorf("expected rating=5 (preserved), got %v", updated.Properties["rating"])
	}
}

func TestUpdateObjectHandler_ReplaceBody(t *testing.T) {
	vault, sampleID := setupTestVault(t)

	// Set initial body
	obj, _ := vault.ResolveObject(sampleID)
	obj.Body = "Old content"
	vault.SaveObject(obj)

	handler := updateObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "update_object"
	req.Params.Arguments = map[string]any{
		"id":   sampleID,
		"body": "New content",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	updated, _ := vault.ResolveObject(sampleID)
	if updated.Body != "New content" {
		t.Errorf("expected body 'New content', got %q", updated.Body)
	}
}

func TestUpdateObjectHandler_LockedObject(t *testing.T) {
	vault, sampleID := setupTestVault(t)

	if err := vault.SetLocked(sampleID, true); err != nil {
		t.Fatalf("set locked: %v", err)
	}

	handler := updateObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "update_object"
	req.Params.Arguments = map[string]any{
		"id":         sampleID,
		"properties": map[string]any{"status": "reading"},
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for locked object")
	}
}

func TestUpdateObjectHandler_NotFound(t *testing.T) {
	vault, _ := setupTestVault(t)

	handler := updateObjectHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "update_object"
	req.Params.Arguments = map[string]any{
		"id": "book/nonexistent",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for nonexistent object")
	}
}

// --- link_objects / unlink_objects ---

// setupTestVaultWithRelation creates a vault with a type that has a relation property.
func setupTestVaultWithRelation(t *testing.T) (*core.Vault, string, string) {
	t.Helper()
	v := newTestVault(t, "name: book\nproperties:\n  - name: author\n    type: relation\n    target: book\n")

	obj1, err := v.NewObject("book", "book-a", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	obj2, err := v.NewObject("book", "book-b", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	return v, obj1.ID, obj2.ID
}

func TestLinkObjectsHandler_HappyPath(t *testing.T) {
	vault, id1, id2 := setupTestVaultWithRelation(t)

	handler := linkObjectsHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "link_objects"
	req.Params.Arguments = map[string]any{
		"from_id":  id1,
		"relation": "author",
		"to_id":    id2,
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	if !strings.Contains(textContent.Text, "linked") {
		t.Errorf("expected linked status, got %s", textContent.Text)
	}
}

func TestLinkObjectsHandler_InvalidRelation(t *testing.T) {
	vault, id1, id2 := setupTestVaultWithRelation(t)

	handler := linkObjectsHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "link_objects"
	req.Params.Arguments = map[string]any{
		"from_id":  id1,
		"relation": "nonexistent",
		"to_id":    id2,
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected tool error for invalid relation")
	}
}

func TestUnlinkObjectsHandler_HappyPath(t *testing.T) {
	vault, id1, id2 := setupTestVaultWithRelation(t)

	// First link them
	if err := vault.LinkObjects(id1, "author", id2); err != nil {
		t.Fatalf("link: %v", err)
	}

	handler := unlinkObjectsHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "unlink_objects"
	req.Params.Arguments = map[string]any{
		"from_id":  id1,
		"relation": "author",
		"to_id":    id2,
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}

	textContent := result.Content[0].(mcplib.TextContent)
	if !strings.Contains(textContent.Text, "unlinked") {
		t.Errorf("expected unlinked status, got %s", textContent.Text)
	}
}

func TestUnlinkObjectsHandler_BothDirections(t *testing.T) {
	vault, id1, id2 := setupTestVaultWithRelation(t)

	// Link in both directions
	if err := vault.LinkObjects(id1, "author", id2); err != nil {
		t.Fatalf("link: %v", err)
	}

	handler := unlinkObjectsHandler(vault)
	req := mcplib.CallToolRequest{}
	req.Params.Name = "unlink_objects"
	req.Params.Arguments = map[string]any{
		"from_id":  id1,
		"relation": "author",
		"to_id":    id2,
		"both":     true,
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("unexpected tool error: %v", result.Content)
	}
}
