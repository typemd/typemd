package tui

import (
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func TestFilterRelationCandidates_EmptySearch(t *testing.T) {
	pe := &propEditor{
		mode: propModeRelationPick,
		relCandidates: []relationCandidate{
			{id: "person/alice-01abc", displayName: "alice"},
			{id: "person/bob-01def", displayName: "bob"},
		},
		relSearch: "",
	}
	pe.relFiltered = pe.relCandidates

	pe.filterRelationCandidates()

	if len(pe.relFiltered) != 2 {
		t.Errorf("expected 2 candidates, got %d", len(pe.relFiltered))
	}
}

func TestFilterRelationCandidates_SubstringMatch(t *testing.T) {
	pe := &propEditor{
		mode: propModeRelationPick,
		relCandidates: []relationCandidate{
			{id: "person/alice-01abc", displayName: "alice", lowerName: "alice"},
			{id: "person/bob-01def", displayName: "bob", lowerName: "bob"},
			{id: "person/ali-01ghi", displayName: "ali", lowerName: "ali"},
		},
		relSearch: "ali",
	}
	pe.relFiltered = pe.relCandidates

	pe.filterRelationCandidates()

	if len(pe.relFiltered) != 2 {
		t.Errorf("expected 2 matches for 'ali', got %d", len(pe.relFiltered))
	}
	for _, c := range pe.relFiltered {
		if c.displayName != "alice" && c.displayName != "ali" {
			t.Errorf("unexpected candidate: %s", c.displayName)
		}
	}
}

func TestFilterRelationCandidates_CaseInsensitive(t *testing.T) {
	pe := &propEditor{
		mode: propModeRelationPick,
		relCandidates: []relationCandidate{
			{id: "person/alice-01abc", displayName: "Alice", lowerName: "alice"},
			{id: "person/bob-01def", displayName: "Bob", lowerName: "bob"},
		},
		relSearch: "alice",
	}
	pe.relFiltered = pe.relCandidates

	pe.filterRelationCandidates()

	if len(pe.relFiltered) != 1 {
		t.Errorf("expected 1 match for 'alice' (case-insensitive), got %d", len(pe.relFiltered))
	}
}

func TestFilterRelationCandidates_NoMatch(t *testing.T) {
	pe := &propEditor{
		mode: propModeRelationPick,
		relCandidates: []relationCandidate{
			{id: "person/alice-01abc", displayName: "alice"},
			{id: "person/bob-01def", displayName: "bob"},
		},
		relSearch: "xyz",
	}
	pe.relFiltered = pe.relCandidates

	pe.filterRelationCandidates()

	if len(pe.relFiltered) != 0 {
		t.Errorf("expected 0 matches for 'xyz', got %d", len(pe.relFiltered))
	}
}

func TestFilteredToFullIndex(t *testing.T) {
	pe := &propEditor{
		relCandidates: []relationCandidate{
			{id: "person/alice-01abc", displayName: "alice", fullIndex: 0},
			{id: "person/bob-01def", displayName: "bob", fullIndex: 1},
			{id: "person/charlie-01ghi", displayName: "charlie", fullIndex: 2},
		},
		relFiltered: []relationCandidate{
			{id: "person/bob-01def", displayName: "bob", fullIndex: 1},
			{id: "person/charlie-01ghi", displayName: "charlie", fullIndex: 2},
		},
	}

	if idx := pe.filteredToFullIndex(0); idx != 1 {
		t.Errorf("filtered index 0 should map to full index 1, got %d", idx)
	}
	if idx := pe.filteredToFullIndex(1); idx != 2 {
		t.Errorf("filtered index 1 should map to full index 2, got %d", idx)
	}
	if idx := pe.filteredToFullIndex(-1); idx != -1 {
		t.Errorf("filtered index -1 should return -1, got %d", idx)
	}
	if idx := pe.filteredToFullIndex(5); idx != -1 {
		t.Errorf("out of bounds index should return -1, got %d", idx)
	}
}

func TestCurrentRelationIDs(t *testing.T) {
	tests := []struct {
		name     string
		dp       core.DisplayProperty
		expected []string
	}{
		{
			name:     "nil value",
			dp:       core.DisplayProperty{Value: nil},
			expected: nil,
		},
		{
			name:     "single string",
			dp:       core.DisplayProperty{Value: "person/alice-01abc"},
			expected: []string{"person/alice-01abc"},
		},
		{
			name:     "[]any of strings",
			dp:       core.DisplayProperty{Value: []any{"tag/go-01abc", "tag/rust-01def"}},
			expected: []string{"tag/go-01abc", "tag/rust-01def"},
		},
		{
			name:     "[]string",
			dp:       core.DisplayProperty{Value: []string{"tag/go-01abc"}},
			expected: []string{"tag/go-01abc"},
		},
		{
			name:     "empty string",
			dp:       core.DisplayProperty{Value: ""},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := currentRelationIDs(tt.dp)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d IDs, got %d", len(tt.expected), len(got))
				return
			}
			for i, id := range got {
				if id != tt.expected[i] {
					t.Errorf("ID[%d] = %q, want %q", i, id, tt.expected[i])
				}
			}
		})
	}
}

func TestCurrentSingleRelationID(t *testing.T) {
	tests := []struct {
		name     string
		dp       core.DisplayProperty
		expected string
	}{
		{
			name:     "nil value",
			dp:       core.DisplayProperty{Value: nil},
			expected: "",
		},
		{
			name:     "string value",
			dp:       core.DisplayProperty{Value: "person/alice-01abc"},
			expected: "person/alice-01abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := currentSingleRelationID(tt.dp)
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRelationPickerRender_SingleSelect(t *testing.T) {
	pe := &propEditor{
		mode:      propModeRelationPick,
		editIndex: 0,
		relCandidates: []relationCandidate{
			{id: "person/alice-01abc", displayName: "alice"},
			{id: "person/bob-01def", displayName: "bob"},
		},
		relFiltered: []relationCandidate{
			{id: "person/alice-01abc", displayName: "alice"},
			{id: "person/bob-01def", displayName: "bob"},
		},
		relSearch:     "",
		pickerCursor:  0,
	}

	item := propItem{
		dp:       core.DisplayProperty{Key: "author", Value: nil, IsRelation: true},
		editable: true,
	}

	output := pe.renderRelationPicker(item, true)

	if output == "" {
		t.Error("render output should not be empty")
	}
	// Should contain (none) option and candidates
	if !strings.Contains(output, "(none)") {
		t.Error("single-select should show (none) option")
	}
	if !strings.Contains(output, "alice") {
		t.Error("should show alice candidate")
	}
	if !strings.Contains(output, "bob") {
		t.Error("should show bob candidate")
	}
}

func TestRelationPickerRender_MultiSelect(t *testing.T) {
	pe := &propEditor{
		mode:      propModeRelationMultiPick,
		editIndex: 0,
		relCandidates: []relationCandidate{
			{id: "tag/go-01abc", displayName: "go", fullIndex: 0},
			{id: "tag/rust-01def", displayName: "rust", fullIndex: 1},
		},
		relFiltered: []relationCandidate{
			{id: "tag/go-01abc", displayName: "go", fullIndex: 0},
			{id: "tag/rust-01def", displayName: "rust", fullIndex: 1},
		},
		relSearch:     "",
		relChecked:    []bool{true, false},
		pickerCursor:  0,
	}

	item := propItem{
		dp:       core.DisplayProperty{Key: "tags", Value: nil},
		editable: true,
	}

	output := pe.renderRelationPicker(item, true)

	if output == "" {
		t.Error("render output should not be empty")
	}
	// Should contain checkmarks
	if !strings.Contains(output, "☑") {
		t.Error("multi-select should show checked mark for selected item")
	}
	if !strings.Contains(output, "☐") {
		t.Error("multi-select should show unchecked mark for unselected item")
	}
}

func TestIsPicking(t *testing.T) {
	tests := []struct {
		mode     propEditMode
		expected bool
	}{
		{propModeNavigate, false},
		{propModeTextInput, false},
		{propModeSelectPick, true},
		{propModeMultiPick, true},
		{propModeRelationPick, true},
		{propModeRelationMultiPick, true},
	}

	for _, tt := range tests {
		pe := &propEditor{mode: tt.mode}
		if got := pe.isPicking(); got != tt.expected {
			t.Errorf("isPicking() for mode %d = %v, want %v", tt.mode, got, tt.expected)
		}
	}
}


