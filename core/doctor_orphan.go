package core

import (
	"os"
	"path/filepath"
)

// OrphanDir represents a directory that doesn't match any known type schema.
type OrphanDir struct {
	Path string // relative path like "objects/note" or "templates/note"
	Kind string // "object" or "template"
}

// ScanOrphanDirs scans objects/ and templates/ for directories that don't
// correspond to any known type schema.
func ScanOrphanDirs(v *Vault) []OrphanDir {
	knownTypes := v.ListTypes()
	known := make(map[string]struct{}, len(knownTypes))
	for _, t := range knownTypes {
		known[t] = struct{}{}
	}

	var orphans []OrphanDir

	// Scan objects/ subdirectories.
	orphans = append(orphans, scanDir(v.ObjectsDir(), "objects", "object", known)...)

	// Scan templates/ subdirectories.
	orphans = append(orphans, scanDir(v.TemplatesDir(), "templates", "template", known)...)

	return orphans
}

// scanDir reads immediate subdirectories of dir and returns OrphanDir entries
// for any that don't appear in the known set. If dir doesn't exist, it returns
// nil without error.
func scanDir(dir, prefix, kind string, known map[string]struct{}) []OrphanDir {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var orphans []OrphanDir
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if _, ok := known[name]; !ok {
			orphans = append(orphans, OrphanDir{
				Path: filepath.Join(prefix, name),
				Kind: kind,
			})
		}
	}
	return orphans
}
