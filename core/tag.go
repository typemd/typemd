package core

import (
	"database/sql"
	"fmt"
	"strings"
)

// checkTagNameUnique returns an error if a tag object with the given name already exists.
func (v *Vault) checkTagNameUnique(name string) error {
	var id string
	err := v.db.QueryRow(
		"SELECT id FROM objects WHERE type = ? AND json_extract(properties, '$.name') = ? LIMIT 1",
		TagTypeName, name,
	).Scan(&id)
	if err == nil {
		return fmt.Errorf("tag name %q already exists: %s", name, id)
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return fmt.Errorf("check tag name uniqueness: %w", err)
}

// resolveTagReference resolves a tag reference string to an object ID.
// If the value (after stripping "tag/" prefix) ends with a ULID suffix, it is
// treated as a full object ID. Otherwise, it is looked up by tag name using tagNameIndex.
// Returns the resolved ID and true, or empty string and false if not resolved.
func (v *Vault) resolveTagReference(ref string, diskTags map[string]*Object, tagNameIndex map[string]string) (string, bool) {
	if !strings.HasPrefix(ref, "tag/") {
		return "", false
	}
	slug := strings.TrimPrefix(ref, "tag/")

	// Check if it ends with a ULID suffix → treat as full ID
	if ulidSuffixPattern.MatchString(slug) {
		// Look up by ID in diskTags
		if _, ok := diskTags[ref]; ok {
			return ref, true
		}
		// Full ID reference to non-existent object → broken reference
		return "", false
	}

	// Look up by name via O(1) index
	if id, ok := tagNameIndex[slug]; ok {
		return id, true
	}

	return "", false
}
