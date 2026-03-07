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

