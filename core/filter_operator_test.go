package core

import (
	"strings"
	"testing"
)

func TestValidateFilterOperator_AllStringOps(t *testing.T) {
	ops := []string{"is", "is_not", "contains", "does_not_contain", "starts_with", "ends_with", "is_empty", "is_not_empty"}
	for _, op := range ops {
		if err := ValidateFilterOperator("string", op); err != nil {
			t.Errorf("string/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_AllNumberOps(t *testing.T) {
	ops := []string{"eq", "neq", "gt", "gte", "lt", "lte", "is_empty", "is_not_empty"}
	for _, op := range ops {
		if err := ValidateFilterOperator("number", op); err != nil {
			t.Errorf("number/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_AllDateOps(t *testing.T) {
	ops := []string{"eq", "before", "after", "on_or_before", "on_or_after", "is_empty", "is_not_empty"}
	for _, op := range ops {
		if err := ValidateFilterOperator("date", op); err != nil {
			t.Errorf("date/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_AllCheckboxOps(t *testing.T) {
	ops := []string{"is", "is_not"}
	for _, op := range ops {
		if err := ValidateFilterOperator("checkbox", op); err != nil {
			t.Errorf("checkbox/%s should be valid: %v", op, err)
		}
	}
}

func TestValidateFilterOperator_InvalidCrossType(t *testing.T) {
	cases := []struct {
		propType string
		operator string
	}{
		{"string", "gt"},
		{"number", "contains"},
		{"checkbox", "gt"},
		{"select", "contains"},
		{"date", "contains"},
	}
	for _, tc := range cases {
		if err := ValidateFilterOperator(tc.propType, tc.operator); err == nil {
			t.Errorf("%s/%s should be invalid", tc.propType, tc.operator)
		}
	}
}

func TestValidateFilterOperator_UnknownType(t *testing.T) {
	if err := ValidateFilterOperator("unknown", "is"); err == nil {
		t.Error("unknown type should fail")
	}
}

func TestOperatorsForType_String(t *testing.T) {
	ops := OperatorsForType("string")
	if len(ops) == 0 {
		t.Fatal("expected operators for string type")
	}
	// Verify canonical order: is before contains before is_empty
	isIdx, containsIdx, emptyIdx := -1, -1, -1
	for i, op := range ops {
		switch op {
		case "is":
			isIdx = i
		case "contains":
			containsIdx = i
		case "is_empty":
			emptyIdx = i
		}
	}
	if isIdx >= containsIdx || containsIdx >= emptyIdx {
		t.Errorf("unexpected order: is=%d contains=%d is_empty=%d", isIdx, containsIdx, emptyIdx)
	}
}

func TestOperatorsForType_Unknown(t *testing.T) {
	if ops := OperatorsForType("unknown"); ops != nil {
		t.Errorf("expected nil for unknown type, got %v", ops)
	}
}

func TestFilterRuleToSQL(t *testing.T) {
	tests := []struct {
		name       string
		rule       FilterRule
		wantErr    bool
		wantClause []string // substrings the clause must contain
		wantArgs   int
		checkArg0  any // if non-nil, assert args[0] equals this
	}{
		{
			name:    "unknown operator",
			rule:    FilterRule{Property: "x", Operator: "nope", Value: "y"},
			wantErr: true,
		},
		{
			name:    "unsafe property name",
			rule:    FilterRule{Property: "x; DROP TABLE", Operator: "is", Value: "y"},
			wantErr: true,
		},
		{
			name:       "is generates equality",
			rule:       FilterRule{Property: "status", Operator: "is", Value: "active"},
			wantClause: []string{"= ?"},
			wantArgs:   1,
			checkArg0:  "active",
		},
		{
			name:       "is_not generates inequality with NULL check",
			rule:       FilterRule{Property: "status", Operator: "is_not", Value: "done"},
			wantClause: []string{"IS NULL", "!= ?"},
			wantArgs:   1,
			checkArg0:  "done",
		},
		{
			name:       "contains wraps value with percent",
			rule:       FilterRule{Property: "title", Operator: "contains", Value: "Go"},
			wantClause: []string{"LIKE ?"},
			wantArgs:   1,
			checkArg0:  "%Go%",
		},
		{
			name:       "does_not_contain with NULL check",
			rule:       FilterRule{Property: "title", Operator: "does_not_contain", Value: "Go"},
			wantClause: []string{"IS NULL", "NOT LIKE ?"},
			wantArgs:   1,
			checkArg0:  "%Go%",
		},
		{
			name:       "starts_with appends percent",
			rule:       FilterRule{Property: "name", Operator: "starts_with", Value: "Go"},
			wantClause: []string{"LIKE ?"},
			wantArgs:   1,
			checkArg0:  "Go%",
		},
		{
			name:       "ends_with prepends percent",
			rule:       FilterRule{Property: "name", Operator: "ends_with", Value: "lang"},
			wantClause: []string{"LIKE ?"},
			wantArgs:   1,
			checkArg0:  "%lang",
		},
		{
			name:       "eq generates equality",
			rule:       FilterRule{Property: "rating", Operator: "eq", Value: "5"},
			wantClause: []string{"= ?"},
			wantArgs:   1,
		},
		{
			name:       "neq generates inequality with NULL check",
			rule:       FilterRule{Property: "rating", Operator: "neq", Value: "5"},
			wantClause: []string{"IS NULL", "!= ?"},
			wantArgs:   1,
		},
		{
			name:       "gt generates CAST comparison",
			rule:       FilterRule{Property: "rating", Operator: "gt", Value: "4"},
			wantClause: []string{"CAST(", "> ?"},
			wantArgs:   1,
		},
		{
			name:       "gte generates CAST comparison",
			rule:       FilterRule{Property: "rating", Operator: "gte", Value: "4"},
			wantClause: []string{"CAST(", ">= ?"},
			wantArgs:   1,
		},
		{
			name:       "lt generates CAST comparison",
			rule:       FilterRule{Property: "rating", Operator: "lt", Value: "4"},
			wantClause: []string{"CAST(", "< ?"},
			wantArgs:   1,
		},
		{
			name:       "lte generates CAST comparison",
			rule:       FilterRule{Property: "rating", Operator: "lte", Value: "4"},
			wantClause: []string{"CAST(", "<= ?"},
			wantArgs:   1,
		},
		{
			name:       "before generates date comparison",
			rule:       FilterRule{Property: "published", Operator: "before", Value: "2025-01-01"},
			wantClause: []string{"< ?"},
			wantArgs:   1,
			checkArg0:  "2025-01-01",
		},
		{
			name:       "after generates date comparison",
			rule:       FilterRule{Property: "published", Operator: "after", Value: "2025-01-01"},
			wantClause: []string{"> ?"},
			wantArgs:   1,
			checkArg0:  "2025-01-01",
		},
		{
			name:       "on_or_before generates date comparison",
			rule:       FilterRule{Property: "published", Operator: "on_or_before", Value: "2025-12-31"},
			wantClause: []string{"<= ?"},
			wantArgs:   1,
		},
		{
			name:       "on_or_after generates date comparison",
			rule:       FilterRule{Property: "published", Operator: "on_or_after", Value: "2025-01-01"},
			wantClause: []string{">= ?"},
			wantArgs:   1,
		},
		{
			name:       "is_empty generates null check",
			rule:       FilterRule{Property: "author", Operator: "is_empty"},
			wantClause: []string{"IS NULL"},
			wantArgs:   0,
		},
		{
			name:       "is_not_empty generates not-null check",
			rule:       FilterRule{Property: "author", Operator: "is_not_empty"},
			wantClause: []string{"IS NOT NULL"},
			wantArgs:   0,
		},
		{
			name:       "is with true converts to integer",
			rule:       FilterRule{Property: "done", Operator: "is", Value: "true"},
			wantClause: []string{"= ?"},
			wantArgs:   1,
			checkArg0:  1,
		},
		{
			name:       "is with false converts to integer",
			rule:       FilterRule{Property: "done", Operator: "is", Value: "false"},
			wantClause: []string{"= ?"},
			wantArgs:   1,
			checkArg0:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clause, args, err := FilterRuleToSQL(tc.rule)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(args) != tc.wantArgs {
				t.Fatalf("len(args) = %d, want %d", len(args), tc.wantArgs)
			}
			for _, substr := range tc.wantClause {
				if !strings.Contains(clause, substr) {
					t.Errorf("clause %q should contain %q", clause, substr)
				}
			}
			if tc.checkArg0 != nil && len(args) > 0 && args[0] != tc.checkArg0 {
				t.Errorf("args[0] = %v, want %v", args[0], tc.checkArg0)
			}
		})
	}
}
