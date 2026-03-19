package core

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestViewConfig_YAML_EmptyOptionalFields(t *testing.T) {
	vc := ViewConfig{
		Name:   "default",
		Layout: ViewLayoutList,
	}

	data, err := yaml.Marshal(&vc)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	s := string(data)
	if strings.Contains(s, "filter:") {
		t.Error("empty filter should be omitted from YAML")
	}
	if strings.Contains(s, "sort:") {
		t.Error("empty sort should be omitted from YAML")
	}
	if strings.Contains(s, "group_by:") {
		t.Error("empty group_by should be omitted from YAML")
	}
}

func TestViewConfig_YAML_RoundTrip(t *testing.T) {
	original := ViewConfig{
		Name:   "by-rating",
		Layout: ViewLayoutList,
		Filter: []FilterRule{
			{Property: "status", Operator: "is", Value: "reading"},
		},
		Sort: []SortRule{
			{Property: "rating", Direction: "desc"},
		},
		GroupBy: []GroupRule{{Property: "genre"}},
	}

	data, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored ViewConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if restored.Name != original.Name {
		t.Errorf("Name = %q, want %q", restored.Name, original.Name)
	}
	if restored.Layout != original.Layout {
		t.Errorf("Layout = %q, want %q", restored.Layout, original.Layout)
	}
	if len(restored.Filter) != 1 {
		t.Fatalf("len(Filter) = %d, want 1", len(restored.Filter))
	}
	if restored.Filter[0].Operator != "is" {
		t.Errorf("Filter[0].Operator = %q, want %q", restored.Filter[0].Operator, "is")
	}
	if len(restored.Sort) != 1 {
		t.Fatalf("len(Sort) = %d, want 1", len(restored.Sort))
	}
	if restored.Sort[0].Direction != "desc" {
		t.Errorf("Sort[0].Direction = %q, want %q", restored.Sort[0].Direction, "desc")
	}
	if len(restored.GroupBy) != 1 || restored.GroupBy[0].Property != "genre" {
		t.Errorf("GroupBy = %v, want [{Property: genre}]", restored.GroupBy)
	}
}

func TestViewConfig_YAML_FilterWithoutValue(t *testing.T) {
	vc := ViewConfig{
		Name:   "empty-check",
		Layout: ViewLayoutList,
		Filter: []FilterRule{
			{Property: "author", Operator: "is_empty"},
		},
	}

	data, err := yaml.Marshal(&vc)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var restored ViewConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if restored.Filter[0].Value != "" {
		t.Errorf("Value = %q, want empty", restored.Filter[0].Value)
	}
}
