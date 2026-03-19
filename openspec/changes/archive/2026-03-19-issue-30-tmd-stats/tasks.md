## 1. Core: Stats result types

- [x] 1.1 Write BDD scenarios for VaultStats (multiple types, empty vault, built-in types)
- [x] 1.2 Write BDD scenarios for TypeStats (number, select, multi_select, checkbox, date, relation aggregations, empty properties, non-existent type)
- [x] 1.3 Implement step definitions for VaultStats scenarios
- [x] 1.4 Implement step definitions for TypeStats scenarios
- [x] 1.5 Define StatsResult structs (VaultStats, TypeSummary, TypeStats, PropertyStats, NumberStats, SelectStats, CheckboxStats, DateStats, RelationStats)
- [x] 1.6 Add QueryService.VaultStats() method (iterate types, count objects, find latest updated_at)
- [x] 1.7 Add QueryService.TypeStats(typeName) method (load schema, query objects, compute per-property aggregations)
- [x] 1.8 Add unit tests for edge cases (malformed property values, zero objects, mixed filled/unfilled)

## 2. CLI: Stats command

- [x] 2.1 Create cmd/stats.go with `tmd stats` command, `--type` and `--json` flags
- [x] 2.2 Implement vault-wide summary output (emoji, plural name, count, last updated, total)
- [x] 2.3 Implement single-type property stats output (per-property section with type-appropriate formatting)
- [x] 2.4 Implement JSON output for both modes
- [x] 2.5 Add unit tests for output formatting and error handling (non-existent type)
