package core

// ListTypes returns the names of all available types.
// It merges custom types from .typemd/types/*.yaml with built-in defaults.
func (v *Vault) ListTypes() []string {
	names, _ := v.repo.ListSchemas()
	return names
}
