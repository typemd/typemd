package core

import "testing"

func TestMatchFilter_NilPropertyValue(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{"rating": nil}}

	tests := []struct {
		operator string
		value    string
		want     bool
	}{
		{"is", "5", false},        // nil formats as "<nil>", not "5"
		{"is_not", "5", true},     // nil != "5"
		{"contains", "foo", false},
		{"does_not_contain", "foo", true},
		{"starts_with", "foo", false},
		{"ends_with", "foo", false},
		{"gt", "5", false},
		{"gte", "5", false},
		{"lt", "5", false},
		{"lte", "5", false},
		{"before", "2025-01-01", false},
		{"after", "2025-01-01", false},
		{"is_empty", "", true},
		{"is_not_empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.operator, func(t *testing.T) {
			rule := FilterRule{Property: "rating", Operator: tt.operator, Value: tt.value}
			got := MatchFilter(obj, rule)
			if got != tt.want {
				t.Errorf("MatchFilter(nil, %q, %q) = %v, want %v", tt.operator, tt.value, got, tt.want)
			}
		})
	}
}

func TestMatchFilter_NumericStringComparison(t *testing.T) {
	// "10" gt "5" should work numerically (10 > 5), not lexicographically ("10" < "5")
	obj := &Object{Type: "book", Properties: map[string]any{"pages": "10"}}

	tests := []struct {
		operator string
		value    string
		want     bool
	}{
		{"gt", "5", true},
		{"gt", "10", false},
		{"lt", "20", true},
		{"gte", "10", true},
		{"lte", "10", true},
	}

	for _, tt := range tests {
		t.Run(tt.operator+"_"+tt.value, func(t *testing.T) {
			rule := FilterRule{Property: "pages", Operator: tt.operator, Value: tt.value}
			got := MatchFilter(obj, rule)
			if got != tt.want {
				t.Errorf("MatchFilter(\"10\", %q, %q) = %v, want %v", tt.operator, tt.value, got, tt.want)
			}
		})
	}
}

func TestMatchFilter_EmptyVariants(t *testing.T) {
	tests := []struct {
		name  string
		val   any
		op    string
		want  bool
	}{
		{"nil is_empty", nil, "is_empty", true},
		{"empty string is_empty", "", "is_empty", true},
		{"null string is_empty", "null", "is_empty", true},
		{"non-empty is_empty", "hello", "is_empty", false},
		{"nil is_not_empty", nil, "is_not_empty", false},
		{"empty string is_not_empty", "", "is_not_empty", false},
		{"null string is_not_empty", "null", "is_not_empty", false},
		{"non-empty is_not_empty", "hello", "is_not_empty", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &Object{Type: "book", Properties: map[string]any{"field": tt.val}}
			rule := FilterRule{Property: "field", Operator: tt.op}
			got := MatchFilter(obj, rule)
			if got != tt.want {
				t.Errorf("MatchFilter(%v, %q) = %v, want %v", tt.val, tt.op, got, tt.want)
			}
		})
	}
}

func TestMatchFilter_PropertyNotInMap(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{}}

	tests := []struct {
		operator string
		value    string
		want     bool
	}{
		{"is", "reading", false},
		{"is_not", "reading", true},
		{"contains", "read", false},
		{"does_not_contain", "read", true},
		{"is_empty", "", true},
		{"is_not_empty", "", false},
		{"gt", "5", false},
	}

	for _, tt := range tests {
		t.Run(tt.operator, func(t *testing.T) {
			rule := FilterRule{Property: "status", Operator: tt.operator, Value: tt.value}
			got := MatchFilter(obj, rule)
			if got != tt.want {
				t.Errorf("MatchFilter(missing prop, %q, %q) = %v, want %v", tt.operator, tt.value, got, tt.want)
			}
		})
	}
}

func TestMatchFilter_BooleanPropertyWithIs(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{"locked": true}}

	rule := FilterRule{Property: "locked", Operator: "is", Value: "true"}
	if !MatchFilter(obj, rule) {
		t.Error("expected bool true to match is \"true\"")
	}

	rule = FilterRule{Property: "locked", Operator: "is", Value: "false"}
	if MatchFilter(obj, rule) {
		t.Error("expected bool true not to match is \"false\"")
	}
}

func TestMatchFilter_ArrayPropertyWithContains(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{
		"tags": []any{"fiction", "fantasy", "adventure"},
	}}

	rule := FilterRule{Property: "tags", Operator: "contains", Value: "fantasy"}
	if !MatchFilter(obj, rule) {
		t.Error("expected array containing 'fantasy' to match contains \"fantasy\"")
	}

	rule = FilterRule{Property: "tags", Operator: "contains", Value: "science"}
	if MatchFilter(obj, rule) {
		t.Error("expected array without 'science' not to match contains \"science\"")
	}

	// Test case-insensitivity
	rule = FilterRule{Property: "tags", Operator: "contains", Value: "Fantasy"}
	if !MatchFilter(obj, rule) {
		t.Error("expected case-insensitive match for 'Fantasy' in array")
	}
}

func TestMatchFilter_StringArrayPropertyWithContains(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{
		"authors": []string{"Tolkien", "Lewis"},
	}}

	rule := FilterRule{Property: "authors", Operator: "contains", Value: "tolkien"}
	if !MatchFilter(obj, rule) {
		t.Error("expected []string containing 'Tolkien' to match contains \"tolkien\"")
	}

	rule = FilterRule{Property: "authors", Operator: "does_not_contain", Value: "rowling"}
	if !MatchFilter(obj, rule) {
		t.Error("expected []string without 'rowling' to match does_not_contain \"rowling\"")
	}
}

func TestMatchFilter_UnknownOperator(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{"status": "reading"}}

	rule := FilterRule{Property: "status", Operator: "banana", Value: "reading"}
	if MatchFilter(obj, rule) {
		t.Error("expected unknown operator to return false")
	}
}

func TestMatchFilters_EmptyRules(t *testing.T) {
	obj := &Object{Type: "book", Properties: map[string]any{}}

	if !MatchFilters(obj, nil) {
		t.Error("expected empty rule set to match everything")
	}
	if !MatchFilters(obj, []FilterRule{}) {
		t.Error("expected empty rule slice to match everything")
	}
}

func TestMatchFilter_NumericPropertyTypes(t *testing.T) {
	// Test with actual numeric types (int, float64) not just strings
	tests := []struct {
		name string
		val  any
		op   string
		fv   string
		want bool
	}{
		{"int gt", 10, "gt", "5", true},
		{"float64 lt", 3.14, "lt", "4", true},
		{"int eq string", 42, "eq", "42", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &Object{Type: "book", Properties: map[string]any{"val": tt.val}}
			rule := FilterRule{Property: "val", Operator: tt.op, Value: tt.fv}
			got := MatchFilter(obj, rule)
			if got != tt.want {
				t.Errorf("MatchFilter(%v, %q, %q) = %v, want %v", tt.val, tt.op, tt.fv, got, tt.want)
			}
		})
	}
}
