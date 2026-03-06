package cmd

import "github.com/typemd/typemd/core"

// resolveVault creates a Vault with the given path, defaulting to "." if empty.
func resolveVault(path string) *core.Vault {
	if path == "" {
		return core.NewVault(".")
	}
	return core.NewVault(path)
}
