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

func TestStripULID(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "strips ULID suffix",
			filename: "clean-code-01kk39c30x27ck7ahyc7ct4nyn",
			want:     "clean-code",
		},
		{
			name:     "strips ULID from multi-word name",
			filename: "the-phoenix-project-01kk39c30yx20m5n2wqys1dk02",
			want:     "the-phoenix-project",
		},
		{
			name:     "no ULID suffix returns original",
			filename: "golang-in-action",
			want:     "golang-in-action",
		},
		{
			name:     "short suffix not matching ULID pattern",
			filename: "test-123",
			want:     "test-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripULID(tt.filename)
			if got != tt.want {
				t.Errorf("StripULID(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestObject_DisplayName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "strips ULID from filename",
			filename: "clean-code-01kk39c30x27ck7ahyc7ct4nyn",
			want:     "clean-code",
		},
		{
			name:     "no ULID returns filename",
			filename: "golang-in-action",
			want:     "golang-in-action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &Object{Filename: tt.filename}
			got := obj.DisplayName()
			if got != tt.want {
				t.Errorf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestObject_DisplayID(t *testing.T) {
	tests := []struct {
		name     string
		objType  string
		filename string
		want     string
	}{
		{
			name:     "returns type/display-name",
			objType:  "book",
			filename: "clean-code-01kk39c30x27ck7ahyc7ct4nyn",
			want:     "book/clean-code",
		},
		{
			name:     "no ULID uses full filename",
			objType:  "book",
			filename: "golang-in-action",
			want:     "book/golang-in-action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &Object{Type: tt.objType, Filename: tt.filename}
			got := obj.DisplayID()
			if got != tt.want {
				t.Errorf("DisplayID() = %q, want %q", got, tt.want)
			}
		})
	}
}
