## Context

typemd provides object listing, querying, and searching, but no aggregate statistics. Users cannot quickly understand the shape of their vault — how many objects exist per type, or the distribution of property values within a type. The `QueryService` handles all read-side operations following CQRS, and `CountObjectsByType()` already exists on `Vault` as a simple count method.

## Goals / Non-Goals

**Goals:**
- Provide vault-wide summary (per-type counts, emoji, last updated)
- Provide single-type property aggregation for all aggregable property types
- Place aggregation logic in `QueryService` for reuse by CLI, MCP, and Web
- Support both human-readable table and `--json` output

**Non-Goals:**
- SQL-level aggregation (keep it simple with Go iteration for now)
- Cross-type aggregation (e.g., "average across all types")
- Histogram or chart rendering in the terminal
- Aggregation for `string` properties (e.g., distinct count, average length)

## Decisions

### 1. StatsResult types in core

Introduce `VaultStats` (vault-wide) and `TypeStats` (single-type) structs.

```go
type VaultStats struct {
    Types []TypeSummary
    Total int
}

type TypeSummary struct {
    Name        string
    Plural      string
    Emoji       string
    Count       int
    LastUpdated time.Time  // most recent updated_at across objects of this type
}

type TypeStats struct {
    TypeName   string
    Count      int
    Properties []PropertyStats
}

type PropertyStats struct {
    Name     string
    Type     string  // "number", "select", etc.
    Filled   int     // count of non-nil values
    Total    int     // total object count
    Stats    any     // type-specific: *NumberStats, *SelectStats, etc.
}
```

Type-specific stat structs:
- `NumberStats`: Sum, Avg, Min, Max float64
- `SelectStats`: Distribution map[string]int (option → count)
- `CheckboxStats`: TrueCount, FalseCount int
- `DateStats`: Earliest, Latest time.Time
- `RelationStats`: Count int (total links)

**Rationale**: Separate structs for vault-wide vs single-type allows different query patterns. `PropertyStats.Stats` uses `any` to accommodate type-specific data without a complex interface hierarchy.

### 2. QueryService methods

```go
func (qs *QueryService) VaultStats() (*VaultStats, error)
func (qs *QueryService) TypeStats(typeName string) (*TypeStats, error)
```

`VaultStats()` iterates all types via `ListTypes()` + `CountObjectsByType()` + queries for latest `updated_at`.

`TypeStats()` loads the type schema for property definitions, queries all objects of that type, then iterates to compute aggregations per property.

**Rationale**: QueryService is the CQRS query-side entry point. Go iteration over results is simple and sufficient for vaults under 10K objects.

### 3. CLI command structure

```
tmd stats              → vault-wide summary
tmd stats --type book  → single-type property stats
tmd stats --json       → JSON output (both modes)
```

Single `stats` command with optional `--type` flag, not subcommands.

**Rationale**: Mirrors the simplicity of `tmd doctor` — a single command with optional flags. The `--type` flag follows the pattern established in the issue description.

### 4. Output formatting

Vault-wide:
```
📚 books      12   last updated 2025-03-15
💡 ideas       8   last updated 2025-03-18
📝 notes      25   last updated 2025-03-19
────────────────────────────────────────────
Total         45
```

Single-type: section per property with type-appropriate stats.

**Rationale**: Compact, scannable. Emoji and plural name from TypeSchema. Right-aligned counts for readability.

## Risks / Trade-offs

- **[Performance]** Go iteration over all objects of a type loads them all into memory. → Acceptable for typical vaults (< 1K objects). If needed later, SQL aggregation can be added behind the same QueryService interface.
- **[Property value parsing]** Property values in frontmatter may be malformed (wrong type, missing). → Skip non-parseable values and count them as unfilled. Don't error on bad data.
- **[Built-in types]** `tag` and `page` have no YAML schema files by default. → Use built-in TypeSchema definitions from `BuiltinTypes()` for emoji/plural.
