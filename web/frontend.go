package web

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist
var frontendAssets embed.FS

// FrontendFS returns the embedded frontend assets rooted at frontend/dist.
// Returns nil if the dist directory is empty (dev mode).
func FrontendFS() fs.FS {
	sub, err := fs.Sub(frontendAssets, "frontend/dist")
	if err != nil {
		return nil
	}
	// Check if dist has content
	entries, err := fs.ReadDir(sub, ".")
	if err != nil || len(entries) == 0 {
		return nil
	}
	return sub
}
