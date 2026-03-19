package core

import "testing"

func TestObjectPathToID(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		objectsDir string
		wantID     string
		wantOK     bool
	}{
		{
			name:       "valid object path",
			path:       "/vault/objects/book/clean-code-01abc.md",
			objectsDir: "/vault/objects",
			wantID:     "book/clean-code-01abc",
			wantOK:     true,
		},
		{
			name:       "non-md file",
			path:       "/vault/objects/book/readme.txt",
			objectsDir: "/vault/objects",
			wantID:     "",
			wantOK:     false,
		},
		{
			name:       "file directly in objects dir (no type subdir)",
			path:       "/vault/objects/stray-file.md",
			objectsDir: "/vault/objects",
			wantID:     "",
			wantOK:     false,
		},
		{
			name:       "too deep nesting",
			path:       "/vault/objects/book/sub/deep.md",
			objectsDir: "/vault/objects",
			wantID:     "",
			wantOK:     false,
		},
		{
			name:       "path outside objects dir",
			path:       "/other/path/book/test.md",
			objectsDir: "/vault/objects",
			wantID:     "",
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotOK := objectPathToID(tt.path, tt.objectsDir)
			if gotOK != tt.wantOK {
				t.Errorf("objectPathToID() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if gotID != tt.wantID {
				t.Errorf("objectPathToID() id = %q, want %q", gotID, tt.wantID)
			}
		})
	}
}
