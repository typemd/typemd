package core

// ListTypes returns the names of all available types.
// It merges custom types from types/*/schema.yaml with built-in defaults.
func (v *Vault) ListTypes() []string {
	names, _ := v.repo.ListSchemas()
	return names
}
