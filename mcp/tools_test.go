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

// setupTestVault creates a temporary vault with a type schema and sample objects.
// Returns the vault and the ID of the created sample object.
func setupTestVault(t *testing.T) (*core.Vault, string) {
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

	// Create type schema
	typesDir := filepath.Join(dir, ".typemd", "types")
	schema := "name: book\nproperties:\n  - name: status\n    type: string\n"
	if err := os.WriteFile(filepath.Join(typesDir, "book.yaml"), []byte(schema), 0644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	// Create sample object
	obj, err := v.NewObject("book", "clean-code")
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
