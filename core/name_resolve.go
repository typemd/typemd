package core

import (
	"fmt"
	"strings"
)

// buildNameIndex populates ctx.nameIndex from walked objects.
// For each object, it indexes the slug (filename without ULID), the name
// property value, and each alias in the aliases system property.
// nameIndex[type][key] = []objectID — multiple IDs indicate ambiguity.
func buildNameIndex(ctx *syncContext) {
	for _, obj := range ctx.diskObjects {
		if ctx.nameIndex[obj.Type] == nil {
			ctx.nameIndex[obj.Type] = make(map[string][]string)
		}
		typeIdx := ctx.nameIndex[obj.Type]

		// Index by slug (filename without ULID)
		slug := StripULID(obj.Filename)
		typeIdx[slug] = append(typeIdx[slug], obj.ID)

		// Index by name property if different from slug
		if nameVal, ok := obj.Properties[NameProperty].(string); ok && nameVal != "" {
			nameSlug := Slugify(nameVal)
			if nameSlug != slug {
				typeIdx[nameSlug] = append(typeIdx[nameSlug], obj.ID)
			}
		}

		// Index by each alias
		for _, alias := range extractAliases(obj.Properties) {
			aliasSlug := Slugify(alias)
			if aliasSlug != "" {
				typeIdx[aliasSlug] = append(typeIdx[aliasSlug], obj.ID)
			}
		}
	}
}

// extractAliases returns the aliases from object properties as a []string.
// Handles both []string (in-memory) and []interface{} (from YAML parse).
func extractAliases(props map[string]any) []string {
	raw, ok := props[AliasesProperty]
	if !ok || raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case []string:
		return v
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// resolveByName looks up a name within a type using the name index.
// Returns the full object ID if exactly one match is found.
// Returns an AmbiguousMatchError if multiple matches exist.
// Returns a not-found error if no match exists.
func resolveByName(nameIndex map[string]map[string][]string, typeName, name string) (string, error) {
	typeIdx, ok := nameIndex[typeName]
	if !ok {
		return "", fmt.Errorf("no object found matching %q", typeName+"/"+name)
	}

	ids, ok := typeIdx[name]
	if !ok || len(ids) == 0 {
		return "", fmt.Errorf("no object found matching %q", typeName+"/"+name)
	}

	if len(ids) == 1 {
		return ids[0], nil
	}

	return "", &AmbiguousMatchError{Prefix: typeName + "/" + name, Matches: ids}
}

// isFullObjectID checks whether a relation value is a complete object ID
// (i.e., has the format "type/slug-ULID" with a valid ULID suffix).
func isFullObjectID(value string) bool {
	parts := strings.SplitN(value, "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		return false
	}
	return ulidSuffixPattern.MatchString(parts[1])
}

// resolveRelationValue resolves a single relation value to a full object ID.
// If the value already has a ULID suffix, it is returned as-is.
// Otherwise, it is treated as a type/name reference and resolved via the name index.
func resolveRelationValue(value string, nameIndex map[string]map[string][]string) (resolved string, changed bool, err error) {
	if value == "" {
		return "", false, fmt.Errorf("empty relation value")
	}

	parts := strings.SplitN(value, "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		return value, false, fmt.Errorf("invalid relation value format: %q", value)
	}

	// If already has ULID suffix, treat as full ID
	if ulidSuffixPattern.MatchString(parts[1]) {
		return value, false, nil
	}

	// Resolve by name
	fullID, err := resolveByName(nameIndex, parts[0], parts[1])
	if err != nil {
		return value, false, err
	}

	return fullID, true, nil
}
