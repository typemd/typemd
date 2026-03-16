package cmd

import (
	"testing"
)

func TestInitCmd_NoStartersFlag(t *testing.T) {
	f := initCmd.Flags().Lookup("no-starters")
	if f == nil {
		t.Fatal("expected --no-starters flag to exist")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %q", f.DefValue)
	}
}

func TestResolveDefaultType(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		expected string
	}{
		{"idea selected", []string{"idea", "book"}, "idea"},
		{"note only", []string{"note", "book"}, "note"},
		{"idea and note — idea wins", []string{"note", "idea"}, "idea"},
		{"book only — no default", []string{"book"}, ""},
		{"empty — no default", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveDefaultType(tt.names)
			if got != tt.expected {
				t.Errorf("resolveDefaultType(%v) = %q, want %q", tt.names, got, tt.expected)
			}
		})
	}
}
