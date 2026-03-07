package core

import (
	"strings"
	"testing"
)

func TestGenerateULID(t *testing.T) {
	id, err := GenerateULID()
	if err != nil {
		t.Fatalf("GenerateULID() error = %v", err)
	}
	if len(id) != 26 {
		t.Errorf("ULID length = %d, want 26", len(id))
	}
	if id != strings.ToLower(id) {
		t.Errorf("ULID should be lowercase, got %q", id)
	}
	id2, err := GenerateULID()
	if err != nil {
		t.Fatalf("GenerateULID() error = %v", err)
	}
	if id == id2 {
		t.Error("two GenerateULID() calls returned same value")
	}
}
