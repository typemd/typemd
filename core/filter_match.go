package core

import (
	"fmt"
	"strings"
)

// MatchFilter evaluates a single FilterRule against an Object in memory.
// Returns true if the object satisfies the filter condition.
func MatchFilter(obj *Object, rule FilterRule) bool {
	val := getObjectPropertyValue(obj, rule.Property)

	// Handle nil upfront: nil doesn't satisfy positive conditions,
	// but satisfies negation operators and is_empty.
	if val == nil {
		switch rule.Operator {
		case "is_empty":
			return true
		case "is_not", "neq", "does_not_contain":
			return true
		default:
			return false
		}
	}

	switch rule.Operator {
	// String / select equality
	case "is":
		return fmt.Sprintf("%v", val) == rule.Value
	case "is_not":
		return fmt.Sprintf("%v", val) != rule.Value

	// String containment (case-insensitive to match SQL LIKE behavior)
	case "contains":
		return containsCI(val, rule.Value)
	case "does_not_contain":
		return !containsCI(val, rule.Value)
	case "starts_with":
		s := fmt.Sprintf("%v", val)
		return strings.HasPrefix(strings.ToLower(s), strings.ToLower(rule.Value))
	case "ends_with":
		s := fmt.Sprintf("%v", val)
		return strings.HasSuffix(strings.ToLower(s), strings.ToLower(rule.Value))

	// Numeric comparison (eq/neq use numeric comparison unlike is/is_not)
	case "eq":
		c := compareNumeric(val, rule.Value)
		return c == 0
	case "neq":
		c := compareNumeric(val, rule.Value)
		return c != 0
	case "gt":
		c := compareNumeric(val, rule.Value)
		return c != compareError && c > 0
	case "gte":
		c := compareNumeric(val, rule.Value)
		return c != compareError && c >= 0
	case "lt":
		c := compareNumeric(val, rule.Value)
		return c != compareError && c < 0
	case "lte":
		c := compareNumeric(val, rule.Value)
		return c != compareError && c <= 0

	// Date comparison (ISO dates sort lexicographically)
	case "before":
		return fmt.Sprintf("%v", val) < rule.Value
	case "after":
		return fmt.Sprintf("%v", val) > rule.Value
	case "on_or_before":
		return fmt.Sprintf("%v", val) <= rule.Value
	case "on_or_after":
		return fmt.Sprintf("%v", val) >= rule.Value

	// Empty checks
	case "is_empty":
		return isEmpty(val)
	case "is_not_empty":
		return !isEmpty(val)

	default:
		return false
	}
}

// MatchFilters evaluates multiple FilterRules (AND logic) against an Object.
// Returns true only if all rules match. An empty rule set matches everything.
func MatchFilters(obj *Object, rules []FilterRule) bool {
	for _, rule := range rules {
		if !MatchFilter(obj, rule) {
			return false
		}
	}
	return true
}

// isEmpty returns true if val is nil, empty string, or the literal "null".
// Extends isNilOrEmpty with the "null" check for YAML null literal handling.
func isEmpty(val any) bool {
	if isNilOrEmpty(val) {
		return true
	}
	s, ok := val.(string)
	return ok && s == "null"
}

// containsCI checks if the string representation of val contains substr (case-insensitive).
// For slice values (multi_select, relation), checks if any element contains substr.
func containsCI(val any, substr string) bool {
	if val == nil {
		return false
	}

	lowerSubstr := strings.ToLower(substr)

	// Handle slice types for multi_select / relation properties
	switch v := val.(type) {
	case []any:
		for _, item := range v {
			if strings.Contains(strings.ToLower(fmt.Sprintf("%v", item)), lowerSubstr) {
				return true
			}
		}
		return false
	case []string:
		for _, item := range v {
			if strings.Contains(strings.ToLower(item), lowerSubstr) {
				return true
			}
		}
		return false
	}

	s := fmt.Sprintf("%v", val)
	return strings.Contains(strings.ToLower(s), lowerSubstr)
}

// compareError is the sentinel value returned by compareNumeric on error.
const compareError = -2

// compareNumeric compares val against target numerically.
// Returns -1, 0, or 1 like strings.Compare. Returns compareError on error.
func compareNumeric(val any, target string) int {
	if val == nil {
		return compareError
	}
	valF, err := toFloat64(val)
	if err != nil {
		return compareError
	}
	targetF, err := toFloat64(target)
	if err != nil {
		return compareError
	}
	switch {
	case valF < targetF:
		return -1
	case valF > targetF:
		return 1
	default:
		return 0
	}
}
