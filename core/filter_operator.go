package core

import "fmt"

// validOperators maps property types to their allowed filter operators.
var validOperators = map[string]map[string]bool{
	"string": {
		"is": true, "is_not": true,
		"contains": true, "does_not_contain": true,
		"starts_with": true, "ends_with": true,
		"is_empty": true, "is_not_empty": true,
	},
	"number": {
		"eq": true, "neq": true,
		"gt": true, "gte": true,
		"lt": true, "lte": true,
		"is_empty": true, "is_not_empty": true,
	},
	"date": {
		"eq": true, "before": true, "after": true,
		"on_or_before": true, "on_or_after": true,
		"is_empty": true, "is_not_empty": true,
	},
	"datetime": {
		"eq": true, "before": true, "after": true,
		"on_or_before": true, "on_or_after": true,
		"is_empty": true, "is_not_empty": true,
	},
	"select": {
		"is": true, "is_not": true,
		"is_empty": true, "is_not_empty": true,
	},
	"multi_select": {
		"contains": true, "does_not_contain": true,
		"is_empty": true, "is_not_empty": true,
	},
	"relation": {
		"contains": true, "does_not_contain": true,
		"is_empty": true, "is_not_empty": true,
	},
	"checkbox": {
		"is": true, "is_not": true,
	},
	"url": {
		"is": true, "is_not": true,
		"contains": true, "does_not_contain": true,
		"is_empty": true, "is_not_empty": true,
	},
}

// OperatorsForType returns the ordered list of valid filter operators for a property type.
// Returns nil if the property type is unknown.
func OperatorsForType(propertyType string) []string {
	ops, ok := validOperators[propertyType]
	if !ok {
		return nil
	}
	// Return operators in a stable, logical order.
	var ordered []string
	// Define the canonical order across all types.
	allOps := []string{
		"is", "is_not",
		"contains", "does_not_contain",
		"starts_with", "ends_with",
		"eq", "neq", "gt", "gte", "lt", "lte",
		"before", "after", "on_or_before", "on_or_after",
		"is_empty", "is_not_empty",
	}
	for _, op := range allOps {
		if ops[op] {
			ordered = append(ordered, op)
		}
	}
	return ordered
}

// ValidateFilterOperator checks that an operator is valid for a given property type.
func ValidateFilterOperator(propertyType, operator string) error {
	ops, ok := validOperators[propertyType]
	if !ok {
		return fmt.Errorf("unknown property type %q", propertyType)
	}
	if !ops[operator] {
		return fmt.Errorf("operator %q is not valid for property type %q", operator, propertyType)
	}
	return nil
}

// FilterRuleToSQL converts a FilterRule into a SQL WHERE clause fragment and args.
func FilterRuleToSQL(rule FilterRule) (clause string, args []any, err error) {
	if !isSafePropertyName(rule.Property) {
		return "", nil, fmt.Errorf("unsafe property name %q", rule.Property)
	}
	prop := fmt.Sprintf("json_extract(properties, '$.%s')", rule.Property)

	switch rule.Operator {
	// String / select equality
	case "is":
		return prop + " = ?", []any{rule.Value}, nil
	case "is_not":
		return "(" + prop + " IS NULL OR " + prop + " != ?)", []any{rule.Value}, nil

	// String containment
	case "contains":
		return prop + " LIKE ?", []any{"%" + rule.Value + "%"}, nil
	case "does_not_contain":
		return "(" + prop + " IS NULL OR " + prop + " NOT LIKE ?)", []any{"%" + rule.Value + "%"}, nil
	case "starts_with":
		return prop + " LIKE ?", []any{rule.Value + "%"}, nil
	case "ends_with":
		return prop + " LIKE ?", []any{"%" + rule.Value}, nil

	// Numeric comparison
	case "eq":
		return prop + " = ?", []any{rule.Value}, nil
	case "neq":
		return "(" + prop + " IS NULL OR " + prop + " != ?)", []any{rule.Value}, nil
	case "gt":
		return "CAST(" + prop + " AS REAL) > ?", []any{rule.Value}, nil
	case "gte":
		return "CAST(" + prop + " AS REAL) >= ?", []any{rule.Value}, nil
	case "lt":
		return "CAST(" + prop + " AS REAL) < ?", []any{rule.Value}, nil
	case "lte":
		return "CAST(" + prop + " AS REAL) <= ?", []any{rule.Value}, nil

	// Date comparison
	case "before":
		return prop + " < ?", []any{rule.Value}, nil
	case "after":
		return prop + " > ?", []any{rule.Value}, nil
	case "on_or_before":
		return prop + " <= ?", []any{rule.Value}, nil
	case "on_or_after":
		return prop + " >= ?", []any{rule.Value}, nil

	// Empty checks
	case "is_empty":
		return "(" + prop + " IS NULL OR " + prop + " = '' OR " + prop + " = 'null')", nil, nil
	case "is_not_empty":
		return "(" + prop + " IS NOT NULL AND " + prop + " != '' AND " + prop + " != 'null')", nil, nil

	default:
		return "", nil, fmt.Errorf("unknown filter operator %q", rule.Operator)
	}
}
