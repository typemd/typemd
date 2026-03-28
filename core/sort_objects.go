package core

import (
	"fmt"
	"sort"
)

// SortObjects sorts a slice of Objects in-place by the given sort rules.
// Uses sort.SliceStable for deterministic ordering of equal elements.
// Nil values are sorted last regardless of direction.
func SortObjects(objects []*Object, rules []SortRule) {
	if len(rules) == 0 || len(objects) <= 1 {
		return
	}
	sort.SliceStable(objects, func(i, j int) bool {
		for _, rule := range rules {
			vi := getObjectSortValue(objects[i], rule.Property)
			vj := getObjectSortValue(objects[j], rule.Property)

			// Nil values always sort last, regardless of direction.
			viNil := isNilOrEmpty(vi)
			vjNil := isNilOrEmpty(vj)
			if viNil && vjNil {
				continue
			}
			if viNil {
				return false // i is nil, should be after j
			}
			if vjNil {
				return true // j is nil, should be after i
			}

			cmp := compareSortValues(vi, vj)
			if cmp == 0 {
				continue
			}
			if rule.Direction == "desc" {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})
}

// getObjectSortValue extracts the value to sort on from an object.
// Returns nil if the property is missing or the object is nil.
func getObjectSortValue(obj *Object, property string) any {
	if obj == nil {
		return nil
	}
	return getObjectPropertyValue(obj, property)
}

// compareSortValues compares two values for sorting.
// Nil values are always sorted last (after any non-nil value).
// Attempts numeric comparison first, falls back to string comparison.
// Returns -1, 0, or 1.
func compareSortValues(a, b any) int {
	aNil := isNilOrEmpty(a)
	bNil := isNilOrEmpty(b)

	if aNil && bNil {
		return 0
	}
	if aNil {
		return 1 // nil sorts last
	}
	if bNil {
		return -1 // nil sorts last
	}

	// Try numeric comparison
	af, aErr := toFloat64(a)
	bf, bErr := toFloat64(b)
	if aErr == nil && bErr == nil {
		switch {
		case af < bf:
			return -1
		case af > bf:
			return 1
		default:
			return 0
		}
	}

	// Fall back to string comparison
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	switch {
	case as < bs:
		return -1
	case as > bs:
		return 1
	default:
		return 0
	}
}

// isNilOrEmpty returns true if the value is nil or an empty string.
func isNilOrEmpty(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok && s == "" {
		return true
	}
	return false
}

// getObjectPropertyValue extracts a property value from an Object,
// handling special properties "type" and "name" consistently.
func getObjectPropertyValue(obj *Object, property string) any {
	switch property {
	case "type":
		return obj.Type
	case "name":
		return obj.GetName()
	default:
		return obj.Properties[property]
	}
}
