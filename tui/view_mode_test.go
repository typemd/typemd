package tui

import (
	"testing"

	"github.com/typemd/typemd/core"
)

func TestBuildGroups_NoGroupBy(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{},
		objects: []*core.Object{
			{Properties: map[string]any{"name": "A"}},
			{Properties: map[string]any{"name": "B"}},
		},
	}
	vm.buildGroups()

	if len(vm.groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(vm.groups))
	}
	if vm.groups[0].Label != "" {
		t.Errorf("expected empty label, got %q", vm.groups[0].Label)
	}
	if len(vm.groups[0].Objects) != 2 {
		t.Errorf("expected 2 objects, got %d", len(vm.groups[0].Objects))
	}
}

func TestBuildGroups_SingleGroupRule(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{
			GroupBy: []core.GroupRule{{Property: "genre"}},
		},
		objects: []*core.Object{
			{Properties: map[string]any{"name": "A", "genre": "sci-fi"}},
			{Properties: map[string]any{"name": "B", "genre": "drama"}},
			{Properties: map[string]any{"name": "C", "genre": "sci-fi"}},
		},
	}
	vm.buildGroups()

	if len(vm.groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(vm.groups))
	}
	if vm.groups[0].Label != "sci-fi" {
		t.Errorf("expected first group label %q, got %q", "sci-fi", vm.groups[0].Label)
	}
	if len(vm.groups[0].Objects) != 2 {
		t.Errorf("expected 2 objects in sci-fi, got %d", len(vm.groups[0].Objects))
	}
	if vm.groups[1].Label != "drama" {
		t.Errorf("expected second group label %q, got %q", "drama", vm.groups[1].Label)
	}
}

func TestBuildGroups_MultiLevelGrouping(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{
			GroupBy: []core.GroupRule{
				{Property: "genre"},
				{Property: "status"},
			},
		},
		objects: []*core.Object{
			{Properties: map[string]any{"genre": "sci-fi", "status": "reading"}},
			{Properties: map[string]any{"genre": "sci-fi", "status": "done"}},
			{Properties: map[string]any{"genre": "drama", "status": "reading"}},
		},
	}
	vm.buildGroups()

	if len(vm.groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(vm.groups))
	}
	if vm.groups[0].Label != "sci-fi · reading" {
		t.Errorf("group 0 label = %q, want %q", vm.groups[0].Label, "sci-fi · reading")
	}
	if vm.groups[1].Label != "sci-fi · done" {
		t.Errorf("group 1 label = %q, want %q", vm.groups[1].Label, "sci-fi · done")
	}
	if vm.groups[2].Label != "drama · reading" {
		t.Errorf("group 2 label = %q, want %q", vm.groups[2].Label, "drama · reading")
	}
}

func TestBuildGroups_MissingPropertyValue(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{
			GroupBy: []core.GroupRule{
				{Property: "genre"},
				{Property: "status"},
			},
		},
		objects: []*core.Object{
			{Properties: map[string]any{"genre": "sci-fi"}},
			{Properties: map[string]any{"status": "reading"}},
		},
	}
	vm.buildGroups()

	if len(vm.groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(vm.groups))
	}
	if vm.groups[0].Label != "sci-fi · (none)" {
		t.Errorf("group 0 label = %q, want %q", vm.groups[0].Label, "sci-fi · (none)")
	}
	if vm.groups[1].Label != "(none) · reading" {
		t.Errorf("group 1 label = %q, want %q", vm.groups[1].Label, "(none) · reading")
	}
}

func TestBuildGroups_ThreeLevels(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{
			GroupBy: []core.GroupRule{
				{Property: "a"},
				{Property: "b"},
				{Property: "c"},
			},
		},
		objects: []*core.Object{
			{Properties: map[string]any{"a": "x", "b": "y", "c": "z"}},
		},
	}
	vm.buildGroups()

	if len(vm.groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(vm.groups))
	}
	if vm.groups[0].Label != "x · y · z" {
		t.Errorf("label = %q, want %q", vm.groups[0].Label, "x · y · z")
	}
}

func TestVisibleRows_WithGroupBy(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{
			GroupBy: []core.GroupRule{{Property: "genre"}},
		},
		objects: []*core.Object{
			{Properties: map[string]any{"genre": "sci-fi"}},
			{Properties: map[string]any{"genre": "drama"}},
		},
	}
	vm.buildGroups()
	rows := vm.visibleRows()

	// Should have: header + object + header + object = 4 rows
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}
	if !rows[0].isHeader {
		t.Error("row 0 should be header")
	}
	if rows[1].isHeader {
		t.Error("row 1 should not be header")
	}
}

func TestVisibleRows_WithoutGroupBy(t *testing.T) {
	vm := &viewMode{
		view: &core.ViewConfig{},
		objects: []*core.Object{
			{Properties: map[string]any{"name": "A"}},
			{Properties: map[string]any{"name": "B"}},
		},
	}
	vm.buildGroups()
	rows := vm.visibleRows()

	// No headers — just 2 object rows
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	for _, r := range rows {
		if r.isHeader {
			t.Error("should not have any header rows")
		}
	}
}
