package core

import (
	"os"
	"sort"
	"strings"
)

// ListTypes returns the names of all available types.
// It merges custom types from .typemd/types/*.yaml with built-in defaults.
func (v *Vault) ListTypes() []string {
	seen := make(map[string]bool)

	// Custom types from YAML files
	entries, err := os.ReadDir(v.TypesDir())
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
				name := strings.TrimSuffix(e.Name(), ".yaml")
				seen[name] = true
			}
		}
	}

	// Built-in defaults
	for name := range defaultTypes {
		seen[name] = true
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
