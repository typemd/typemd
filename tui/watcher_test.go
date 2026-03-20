package tui

import (
	"testing"
)

func TestDeduplicatePaths(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{
			name:  "no duplicates",
			input: []string{"/a/b.md", "/c/d.md"},
			want:  2,
		},
		{
			name:  "with duplicates",
			input: []string{"/a/b.md", "/c/d.md", "/a/b.md"},
			want:  2,
		},
		{
			name:  "all same",
			input: []string{"/a/b.md", "/a/b.md", "/a/b.md"},
			want:  1,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  0,
		},
		{
			name:  "nil input",
			input: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deduplicatePaths(tt.input)
			if len(got) != tt.want {
				t.Errorf("deduplicatePaths() = %d paths, want %d", len(got), tt.want)
			}
		})
	}
}

func TestDeduplicatePaths_PreservesOrder(t *testing.T) {
	input := []string{"/c.md", "/a.md", "/b.md", "/a.md"}
	got := deduplicatePaths(input)
	expected := []string{"/c.md", "/a.md", "/b.md"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d paths, got %d", len(expected), len(got))
	}
	for i, p := range expected {
		if got[i] != p {
			t.Errorf("path[%d] = %q, want %q", i, got[i], p)
		}
	}
}
