package core

import "strings"

// resolveTagReference resolves a tag reference string to an object ID.
// If the value (after stripping "tag/" prefix) ends with a ULID suffix, it is
// treated as a full object ID. Otherwise, it is looked up by tag name using tagNameIndex.
// Returns the resolved ID and true, or empty string and false if not resolved.
func resolveTagReference(ref string, diskTags map[string]*Object, tagNameIndex map[string]string) (string, bool) {
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
