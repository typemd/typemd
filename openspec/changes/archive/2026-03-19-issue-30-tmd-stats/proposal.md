## Why

There is no way to get aggregate statistics about objects in the vault. Users can only list or view individual objects. A `tmd stats` command would provide a dashboard-like overview — object counts per type, property distributions, and numeric aggregations — giving users a quick understanding of their vault's shape.

## What Changes

- Add `tmd stats` CLI command showing a per-type summary table (emoji, name, count, last updated)
- Add `tmd stats --type <name>` for single-type property aggregation (number avg/min/max, select frequency, checkbox ratio, date range, relation counts)
- Add `QueryService.Stats()` method in core for reusable aggregation logic (Go iteration over query results)
- Support `--json` flag for structured output

## Capabilities

### New Capabilities
- `vault-stats`: Aggregate statistics for the vault — per-type object counts and single-type property aggregation (number, select, multi_select, checkbox, date, relation)

### Modified Capabilities

(none)

## Impact

- **core/**: New `StatsResult` types + `QueryService.Stats()` method
- **cmd/**: New `stats.go` command file
- No breaking changes, no dependency additions
