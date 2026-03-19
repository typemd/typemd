package core

import (
	"testing"
)

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{"float64", float64(3.14), 3.14, false},
		{"float32", float32(2.5), 2.5, false},
		{"int", 42, 42.0, false},
		{"int64", int64(100), 100.0, false},
		{"string number", "3.5", 3.5, false},
		{"string int", "7", 7.0, false},
		{"invalid string", "abc", 0, true},
		{"empty string", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toFloat64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("toFloat64(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("toFloat64(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		{"bool true", true, true, false},
		{"bool false", false, false, false},
		{"string true", "true", true, false},
		{"string false", "false", false, false},
		{"string 1", "1", true, false},
		{"string 0", "0", false, false},
		{"invalid", "maybe", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("toBool(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("toBool(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestComputeNumberStats_MixedValues(t *testing.T) {
	objects := []*Object{
		{Properties: map[string]any{"rating": float64(5)}},
		{Properties: map[string]any{"rating": nil}},
		{Properties: map[string]any{"rating": "invalid"}},
		{Properties: map[string]any{"rating": float64(3)}},
		{Properties: map[string]any{}},
	}
	var filled int
	stats := computeNumberStats("rating", objects, &filled)
	if filled != 2 {
		t.Errorf("filled = %d, want 2", filled)
	}
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
	if stats.Avg != 4.0 {
		t.Errorf("avg = %v, want 4.0", stats.Avg)
	}
	if stats.Min != 3.0 {
		t.Errorf("min = %v, want 3.0", stats.Min)
	}
	if stats.Max != 5.0 {
		t.Errorf("max = %v, want 5.0", stats.Max)
	}
	if stats.Sum != 8.0 {
		t.Errorf("sum = %v, want 8.0", stats.Sum)
	}
}

func TestComputeNumberStats_AllNil(t *testing.T) {
	objects := []*Object{
		{Properties: map[string]any{"rating": nil}},
		{Properties: map[string]any{}},
	}
	var filled int
	stats := computeNumberStats("rating", objects, &filled)
	if filled != 0 {
		t.Errorf("filled = %d, want 0", filled)
	}
	if stats != nil {
		t.Errorf("expected nil stats, got %v", stats)
	}
}

func TestComputeSelectStats_EmptyValues(t *testing.T) {
	objects := []*Object{
		{Properties: map[string]any{"status": "active"}},
		{Properties: map[string]any{"status": nil}},
		{Properties: map[string]any{"status": ""}},
		{Properties: map[string]any{}},
	}
	var filled int
	stats := computeSelectStats("status", objects, &filled)
	if filled != 1 {
		t.Errorf("filled = %d, want 1", filled)
	}
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
	if stats.Distribution["active"] != 1 {
		t.Errorf("active count = %d, want 1", stats.Distribution["active"])
	}
}

func TestComputeCheckboxStats_MixedTypes(t *testing.T) {
	objects := []*Object{
		{Properties: map[string]any{"done": true}},
		{Properties: map[string]any{"done": false}},
		{Properties: map[string]any{"done": nil}},
		{Properties: map[string]any{}},
	}
	var filled int
	stats := computeCheckboxStats("done", objects, &filled)
	if filled != 2 {
		t.Errorf("filled = %d, want 2", filled)
	}
	if stats.TrueCount != 1 {
		t.Errorf("true = %d, want 1", stats.TrueCount)
	}
	if stats.FalseCount != 1 {
		t.Errorf("false = %d, want 1", stats.FalseCount)
	}
}

func TestComputeMultiSelectStats(t *testing.T) {
	objects := []*Object{
		{Properties: map[string]any{"genres": []any{"fiction", "sci-fi"}}},
		{Properties: map[string]any{"genres": []any{"fiction"}}},
		{Properties: map[string]any{"genres": nil}},
		{Properties: map[string]any{}},
	}
	var filled int
	stats := computeMultiSelectStats("genres", objects, &filled)
	if filled != 2 {
		t.Errorf("filled = %d, want 2", filled)
	}
	if stats.Distribution["fiction"] != 2 {
		t.Errorf("fiction = %d, want 2", stats.Distribution["fiction"])
	}
	if stats.Distribution["sci-fi"] != 1 {
		t.Errorf("sci-fi = %d, want 1", stats.Distribution["sci-fi"])
	}
}
