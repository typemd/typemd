package cmd

import (
	"os"
	"testing"
)

func TestParseTemplateArg(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		wantType string
		wantName string
		wantErr  bool
	}{
		{"valid type/name", "book/review", "book", "review", false},
		{"valid with hyphens", "my-type/my-template", "my-type", "my-template", false},
		{"missing name", "book/", "", "", true},
		{"missing type", "/review", "", "", true},
		{"no slash", "review", "", "", true},
		{"empty string", "", "", "", true},
		{"multiple slashes takes first", "a/b/c", "a", "b/c", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeName, name, err := parseTemplateArg(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTemplateArg(%q) error = %v, wantErr %v", tt.arg, err, tt.wantErr)
				return
			}
			if typeName != tt.wantType {
				t.Errorf("parseTemplateArg(%q) typeName = %q, want %q", tt.arg, typeName, tt.wantType)
			}
			if name != tt.wantName {
				t.Errorf("parseTemplateArg(%q) name = %q, want %q", tt.arg, name, tt.wantName)
			}
		})
	}
}

func TestResolveEditor(t *testing.T) {
	// Save and restore environment
	origEditor := os.Getenv("EDITOR")
	origVisual := os.Getenv("VISUAL")
	defer func() {
		os.Setenv("EDITOR", origEditor)
		os.Setenv("VISUAL", origVisual)
	}()

	t.Run("EDITOR is set", func(t *testing.T) {
		os.Setenv("EDITOR", "nano")
		os.Setenv("VISUAL", "code")
		editor, err := resolveEditor()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if editor != "nano" {
			t.Errorf("got %q, want %q", editor, "nano")
		}
	})

	t.Run("EDITOR unset, VISUAL is set", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		os.Setenv("VISUAL", "code")
		editor, err := resolveEditor()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if editor != "code" {
			t.Errorf("got %q, want %q", editor, "code")
		}
	})

	t.Run("both unset, vi available", func(t *testing.T) {
		os.Unsetenv("EDITOR")
		os.Unsetenv("VISUAL")
		editor, err := resolveEditor()
		if err != nil {
			t.Skipf("vi not available on this system: %v", err)
		}
		if editor != "vi" {
			t.Errorf("got %q, want %q", editor, "vi")
		}
	})
}
