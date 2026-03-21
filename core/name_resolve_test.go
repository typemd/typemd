package core

import (
	"testing"
)

func TestBuildNameIndex(t *testing.T) {
	ctx := &syncContext{
		diskObjects: map[string]*Object{
			"person/john-doe-01abc": {
				ID:       "person/john-doe-01abc",
				Type:     "person",
				Filename: "john-doe-01abc",
				Properties: map[string]any{
					NameProperty: "John Doe",
				},
			},
			"book/clean-code-01xyz": {
				ID:       "book/clean-code-01xyz",
				Type:     "book",
				Filename: "clean-code-01xyz",
				Properties: map[string]any{
					NameProperty: "Clean Code",
				},
			},
		},
		nameIndex: make(map[string]map[string][]string),
	}

	buildNameIndex(ctx)

	// Should index by slug (filename without ULID)
	if ids := ctx.nameIndex["person"]["john-doe"]; len(ids) != 1 || ids[0] != "person/john-doe-01abc" {
		t.Errorf("expected person/john-doe → [person/john-doe-01abc], got %v", ids)
	}
	if ids := ctx.nameIndex["book"]["clean-code"]; len(ids) != 1 || ids[0] != "book/clean-code-01xyz" {
		t.Errorf("expected book/clean-code → [book/clean-code-01xyz], got %v", ids)
	}
}

func TestBuildNameIndex_DuplicateNames(t *testing.T) {
	ctx := &syncContext{
		diskObjects: map[string]*Object{
			"person/john-01aaa": {
				ID: "person/john-01aaa", Type: "person", Filename: "john-01aaa",
				Properties: map[string]any{NameProperty: "john"},
			},
			"person/john-01bbb": {
				ID: "person/john-01bbb", Type: "person", Filename: "john-01bbb",
				Properties: map[string]any{NameProperty: "john"},
			},
		},
		nameIndex: make(map[string]map[string][]string),
	}

	buildNameIndex(ctx)

	ids := ctx.nameIndex["person"]["john"]
	if len(ids) != 2 {
		t.Errorf("expected 2 entries for ambiguous name, got %d: %v", len(ids), ids)
	}
}

func TestResolveByName_UniqueMatch(t *testing.T) {
	nameIndex := map[string]map[string][]string{
		"person": {"john-doe": {"person/john-doe-01abc"}},
	}

	id, err := resolveByName(nameIndex, "person", "john-doe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "person/john-doe-01abc" {
		t.Errorf("expected person/john-doe-01abc, got %s", id)
	}
}

func TestResolveByName_NotFound(t *testing.T) {
	nameIndex := map[string]map[string][]string{
		"person": {"john-doe": {"person/john-doe-01abc"}},
	}

	_, err := resolveByName(nameIndex, "person", "nobody")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestResolveByName_Ambiguous(t *testing.T) {
	nameIndex := map[string]map[string][]string{
		"person": {"john": {"person/john-01aaa", "person/john-01bbb"}},
	}

	_, err := resolveByName(nameIndex, "person", "john")
	if err == nil {
		t.Fatal("expected error for ambiguous")
	}
	if _, ok := err.(*AmbiguousMatchError); !ok {
		t.Errorf("expected AmbiguousMatchError, got %T: %v", err, err)
	}
}

func TestResolveByName_UnknownType(t *testing.T) {
	nameIndex := map[string]map[string][]string{}

	_, err := resolveByName(nameIndex, "unknown", "foo")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestResolveRelationValue_FullID(t *testing.T) {
	nameIndex := map[string]map[string][]string{}

	resolved, changed, err := resolveRelationValue("person/john-doe-01jqr3k5mpbvn8e0f2g7h9txyz", nameIndex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false for full ID")
	}
	if resolved != "person/john-doe-01jqr3k5mpbvn8e0f2g7h9txyz" {
		t.Errorf("expected unchanged ID, got %s", resolved)
	}
}

func TestResolveRelationValue_Prefix(t *testing.T) {
	nameIndex := map[string]map[string][]string{
		"person": {"john-doe": {"person/john-doe-01abc"}},
	}

	resolved, changed, err := resolveRelationValue("person/john-doe", nameIndex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for prefix")
	}
	if resolved != "person/john-doe-01abc" {
		t.Errorf("expected person/john-doe-01abc, got %s", resolved)
	}
}

func TestResolveRelationValue_Empty(t *testing.T) {
	nameIndex := map[string]map[string][]string{}

	_, _, err := resolveRelationValue("", nameIndex)
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestResolveRelationValue_NoSlash(t *testing.T) {
	nameIndex := map[string]map[string][]string{}

	_, _, err := resolveRelationValue("invalid", nameIndex)
	if err == nil {
		t.Fatal("expected error for value without slash")
	}
}

func TestResolveRelationValue_TypeOnly(t *testing.T) {
	nameIndex := map[string]map[string][]string{}

	_, _, err := resolveRelationValue("person/", nameIndex)
	if err == nil {
		t.Fatal("expected error for type-only value")
	}
}
