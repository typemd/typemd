## Context

typemd currently supports 4 property types: string, number, enum, relation. Properties are stored as a JSON blob in the `objects.properties` column. The Property struct uses a `Values []string` field for enum options and a `Type string` field validated against a hardcoded allowlist.

The v0.2.0 milestone focuses on solidifying the type system. This change expands to 9 property types with proper validation. SQLite storage remains as JSON blob — typed storage for range queries and sorting will be added in a future issue alongside query syntax changes.

## Goals / Non-Goals

**Goals:**
- Define 9 property types with clear validation rules
- Replace `enum` with `select` and `values` with `options` (object array with value/label)
- Support `tmd migrate` for enum → select schema migration

**Non-Goals:**
- Typed SQLite storage (`object_properties` table) — separate issue, paired with query syntax
- Query syntax changes (`tmd query "rating>4"`) — separate issue
- `computed` / `formula` types — separate issue
- Property-level sorting in TUI — separate issue

## Decisions

### 1. Separate `date` and `datetime` types

**Decision:** Two distinct types rather than one auto-detecting type.

**Rationale:** Explicit typing avoids ambiguity. A `date` property always stores `YYYY-MM-DD`, a `datetime` always stores `YYYY-MM-DDTHH:MM:SS`. Users declare intent in the schema. ISO 8601 strings are naturally sortable in SQLite.

**Alternative:** Single `date` type that accepts both formats. Rejected because it makes validation and display inconsistent — sometimes showing time, sometimes not, depending on the value.

### 2. Options object array for select/multi_select

**Decision:** Replace `values: [a, b]` with `options: [{value: a, label: A}]`.

**Rationale:** The `label` field enables display names that differ from stored values (e.g., `value: in-progress`, `label: In Progress`). This is a common pattern in form builders and knowledge tools. `label` is optional and defaults to `value`.

**Alternative:** Keep `values` array and add a separate `labels` map. Rejected because it couples two parallel arrays, which is error-prone.

### 3. No SQLite schema changes (deferred)

**Decision:** Properties continue to be stored as JSON blob in `objects.properties`. No `object_properties` typed table in this iteration.

**Rationale:** A typed table is only valuable when paired with query syntax that can exploit it (e.g., `rating>4`, date ranges). Since query syntax changes are out of scope, the typed table would have zero consumers — pure overhead with no benefit. Building it now also forces premature decisions about multi_select storage (JSON array vs normalized rows) that are better made alongside the query design.

**Alternative:** Add typed table now for "future readiness". Rejected per YAGNI — dual-write complexity, cleanup logic, and additional tests for infrastructure nobody uses yet.

## Risks / Trade-offs

- **[Breaking change: enum → select]** → Mitigated by `tmd migrate` support and clear error messages when `enum` is encountered in validation.
- **[Breaking change: values → options]** → Mitigated by `tmd migrate` auto-converting `values: [a, b]` to `options: [{value: a}]`.
- **[YAML auto-parsing]** → YAML parsers auto-convert dates and booleans. `2026-01-01` becomes a `time.Time`, `true` becomes a Go `bool`. The sync/validation layer must handle Go native types, not just strings. This needs careful type coercion in `ValidateObject` and `SyncIndex`.
