package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/typemd/typemd/core"
)

// resolveVault creates a Vault with the given path, defaulting to "." if empty.
func resolveVault(path string) *core.Vault {
	if path == "" {
		return core.NewVault(".")
	}
	return core.NewVault(path)
}

// openVault creates, opens, and optionally reindexes a vault.
// The caller must defer vault.Close().
func openVault(path string, reindex bool) (*core.Vault, error) {
	vault := resolveVault(path)
	if err := vault.Open(); err != nil {
		return nil, err
	}
	if reindex {
		result, err := vault.SyncIndex()
		if err != nil {
			vault.Close()
			return nil, fmt.Errorf("reindex: %w", err)
		}
		fmt.Println("Index synced successfully.")
		if len(result.Orphaned) > 0 {
			fmt.Printf("Warning: Found %d orphaned relation(s):\n", len(result.Orphaned))
			for _, o := range result.Orphaned {
				fmt.Printf("  %s -> %s (relation: %q)\n", o.FromID, o.ToID, o.Name)
			}
			fmt.Println("Orphaned relations have been removed from the index.")
		}
	}
	return vault, nil
}

// printObjects prints objects as JSON or one DisplayID per line.
func printObjects(objects []*core.Object, asJSON bool) error {
	if asJSON {
		data, err := json.MarshalIndent(objects, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal JSON: %w", err)
		}
		fmt.Println(string(data))
	} else {
		for _, obj := range objects {
			fmt.Println(obj.DisplayID())
		}
	}
	return nil
}

