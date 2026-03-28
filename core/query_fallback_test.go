package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// setupFallbackQueryService creates a QueryService with an error-returning index
// and the given objects on disk. Returns the QueryService and a cleanup function.
func setupFallbackQueryService(t *testing.T, objects map[string]string) *QueryService {
	t.Helper()
	dir := filepath.Join(os.TempDir(), "typemd-fallback-unit-"+mustULID())
	os.MkdirAll(dir, 0755)
	t.Cleanup(func() { os.RemoveAll(dir) })

	for relPath, content := range objects {
		fullPath := filepath.Join(dir, "objects", relPath)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte(content), 0644)
	}

	repo := NewLocalObjectRepository(dir)
	return NewQueryService(repo, &errorIndex{})
}

func TestQueryFallback_EmptyVault(t *testing.T) {
	qs := setupFallbackQueryService(t, nil)
	results, err := qs.Query(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestQueryFallback_WithFilter(t *testing.T) {
	qs := setupFallbackQueryService(t, map[string]string{
		"book/a-01jtest00000000000000001.md": "---\nname: A\nstatus: reading\n---\n",
		"book/b-01jtest00000000000000002.md": "---\nname: B\nstatus: done\n---\n",
	})
	results, err := qs.Query([]FilterRule{{Property: "status", Operator: "is", Value: "reading"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchFallback_EmptyKeyword(t *testing.T) {
	qs := setupFallbackQueryService(t, map[string]string{
		"book/a-01jtest00000000000000001.md": "---\nname: A\n---\nBody text\n",
	})
	results, err := qs.Search("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty keyword, got %d", len(results))
	}
}

func TestSearchFallback_NoMatch(t *testing.T) {
	qs := setupFallbackQueryService(t, map[string]string{
		"book/a-01jtest00000000000000001.md": "---\nname: Alpha\n---\nSome body\n",
	})
	results, err := qs.Search("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchFallback_MatchesDescription(t *testing.T) {
	qs := setupFallbackQueryService(t, map[string]string{
		"book/a-01jtest00000000000000001.md": "---\nname: Alpha\ndescription: A great book about golang\n---\n",
	})
	results, err := qs.Search("golang")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchFallback_SpecialCharacters(t *testing.T) {
	qs := setupFallbackQueryService(t, map[string]string{
		"book/a-01jtest00000000000000001.md": "---\nname: C++ Guide\n---\n",
	})
	results, err := qs.Search("C++")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

// errorWalkRepo is an ObjectRepository that returns an error from Walk().
type errorWalkRepo struct {
	ObjectRepository
}

func (r *errorWalkRepo) Walk() ([]*Object, error) {
	return nil, fmt.Errorf("filesystem error")
}

func TestQueryFallback_WalkError(t *testing.T) {
	qs := NewQueryService(&errorWalkRepo{}, &errorIndex{})
	_, err := qs.Query(nil)
	if err == nil {
		t.Fatal("expected error from walk, got nil")
	}
}

func TestSearchFallback_WalkError(t *testing.T) {
	qs := NewQueryService(&errorWalkRepo{}, &errorIndex{})
	_, err := qs.Search("test")
	if err == nil {
		t.Fatal("expected error from walk, got nil")
	}
}
