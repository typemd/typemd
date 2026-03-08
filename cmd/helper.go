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
		if _, err := vault.SyncIndex(); err != nil {
			vault.Close()
			return nil, fmt.Errorf("reindex: %w", err)
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

