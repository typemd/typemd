## Context

The query pipeline currently accepts a `filter string` parameter using `key=value` format (equality only). This is parsed by `parseFilter()` into `[]filterCondition` structs, then translated to SQL WHERE clauses in `SQLiteObjectIndex.Query()`.

Meanwhile, a structured `FilterRule` type already exists in `core/view.go` with rich operator support (20+ operators via `FilterRuleToSQL()` in `core/filter_operator.go`). These two systems are disconnected — `ViewConfig.Filter` (which uses `[]FilterRule`) is never consulted during query execution.

After #263 integrated `FilterRuleToSQL` into the codebase, the legacy string format has no unique capability that `[]FilterRule` cannot express. The only external consumer (`tmd query` CLI) can be removed as a breaking change.

## Goals / Non-Goals

**Goals:**
- Unify on `[]FilterRule` as the single filter representation across the query pipeline
- Enable `ViewConfig.Filter` to flow directly into queries without conversion
- Remove dead code (`parseFilter()`, `filterCondition`, `cmd/query.go`)
- Maintain identical query behavior for all existing callers

**Non-Goals:**
- Adding new filter operators (already done in #263)
- Implementing `ViewConfig.Filter` application in view mode (that's #264)
- Building a new query CLI command to replace `tmd query` (can be done later if needed)
- Changing the SQLite schema or data model

## Decisions

### 1. Function signature: `[]FilterRule` parameter

**Decision:** Replace `filter string` with `filter []FilterRule` in all three layers.

```go
// ObjectIndex interface
Query(filter []FilterRule, sort ...SortRule) ([]*ObjectResult, error)

// QueryService
func (s *QueryService) Query(filter []FilterRule, sort ...SortRule) ([]*Object, error)

// Vault facade
func (v *Vault) QueryObjects(filter []FilterRule, sort ...SortRule) ([]*Object, error)
```

**Rationale:** `FilterRule` already exists and mirrors `SortRule` symmetrically. Both are value types defined in `core/`. A variadic `QueryOption` pattern would be overengineered for the current needs.

**Alternative considered:** `Query(opts ...QueryOption)` functional options pattern — rejected as unnecessary abstraction for two orthogonal parameters.

### 2. `type` column handling in `SQLiteObjectIndex.Query()`

**Decision:** Handle `type` property as a special case in `SQLiteObjectIndex.Query()`, mapping it to the `type` SQL column instead of routing through `FilterRuleToSQL()`.

```go
for _, rule := range filter {
    if rule.Property == "type" {
        // Special case: "type" is a SQL column, not in properties JSON
        whereClauses = append(whereClauses, "type = ?")
        args = append(args, rule.Value)
    } else {
        clause, ruleArgs, err := FilterRuleToSQL(rule)
        // ...
    }
}
```

**Rationale:** `FilterRuleToSQL()` always generates `json_extract(properties, '$.X')` expressions, but `type` is a top-level column in the `objects` table, not stored in the properties JSON blob. Keeping this mapping in `SQLiteObjectIndex.Query()` (the infrastructure layer) is correct — it's an implementation detail of the SQLite index, not a domain concern.

**Alternative considered:** Making `FilterRuleToSQL()` aware of column mappings — rejected because it would couple a domain utility to SQLite schema details.

### 3. Helper function for common filter patterns

**Decision:** Add a convenience function `TypeFilter(typeName string) []FilterRule` to reduce boilerplate at call sites.

```go
func TypeFilter(typeName string) []FilterRule {
    return []FilterRule{{Property: "type", Operator: "is", Value: typeName}}
}
```

Most internal callers use the pattern `"type=X"` which becomes `[]FilterRule{{Property: "type", Operator: "is", Value: X}}`. A helper keeps call sites clean without introducing unnecessary abstraction.

### 4. Remove `tmd query` without replacement

**Decision:** Remove `cmd/query.go` entirely. Do not add a new structured-filter CLI equivalent.

**Rationale:** `tmd query` was an early debugging tool. The TUI and MCP are the primary query interfaces. If a CLI query command is needed later, it can be designed with proper UX (flags like `--type`, `--filter`, etc.) as a separate issue.

## Risks / Trade-offs

- **[Breaking change]** `tmd query` removal breaks any user scripts relying on it → Acceptable: documented as breaking in issue, low usage expected for a CLI tool in early development.
- **[type special-case]** Hardcoding `type` in `SQLiteObjectIndex.Query()` adds a code smell → Mitigated: this is infrastructure-layer concern, clearly documented, and matches the existing pattern in `parseFilter()`.
- **[Operator support for type filter]** The `type` special case only handles equality (`is` operator). Advanced operators like `is_not` on type would need additional handling → Acceptable: no current use case for non-equality type filters. Can be extended when needed.
