## Why

Date and datetime values are displayed in raw ISO 8601 format (`2026-03-28T01:59:02+08:00` for `updated_at`, `2006-01-02T15:04:05` for datetime properties). This is hard to read. Users need vault-level config options `date_format` and `datetime_format` to control how dates appear in the TUI while keeping storage format unchanged.

## What Changes

- Add `date_format` and `datetime_format` config keys to `VaultConfig` with sensible defaults (`YYYY-MM-DD` and `YYYY-MM-DD HH:mm:ss`)
- Update `FormatValue()` in `core/display.go` to accept and apply format strings using the existing `dateFormatReplacer` from `name_template.go`
- All date/datetime display in TUI (properties panel, table view cells) and CLI (`tmd object show`) respects the configured formats
- System properties (`created_at`, `updated_at`) are also formatted through the same pipeline

## Capabilities

### New Capabilities
- `date-display-format`: Configurable date and datetime display formatting via vault config, applied through the existing `DisplayProperty.FormatValue()` pipeline

### Modified Capabilities

## Impact

- `core/config.go` — new `date_format` / `datetime_format` fields + config key registry entries
- `core/display.go` — `FormatValue()` updated to accept format strings, reusing `dateFormatReplacer` from `name_template.go`
- `core/name_template.go` — `convertDateFormat()` exported for reuse (or logic extracted)
- TUI — no changes needed (already calls `FormatValue()` / `Format()` via `DisplayProperty`)
- CLI `cmd/show.go` — no changes needed (already calls `p.Format()`)
