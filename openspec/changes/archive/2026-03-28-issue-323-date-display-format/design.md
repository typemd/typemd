## Context

`DisplayProperty.FormatValue()` in `core/display.go` hard-codes Go format strings for date (`2006-01-02`) and datetime (`2006-01-02T15:04:05`). System properties (`created_at`, `updated_at`) are `time.Time` values that flow through the same pipeline. The `dateFormatReplacer` in `core/name_template.go` already converts user-friendly tokens (`YYYY`, `MM`, `DD`, `HH`, `mm`, `ss`) to Go reference time — this can be reused.

All display consumers (TUI properties panel, TUI table view, CLI `tmd object show`) call `FormatValue()` / `Format()` on `DisplayProperty`, so changing the formatting at this single point propagates everywhere.

## Goals / Non-Goals

**Goals:**
- Add `date_format` and `datetime_format` config keys to `.typemd/config.yaml`
- Apply configured formats in `FormatValue()` for all date/datetime display
- Reuse existing `dateFormatReplacer` for token conversion
- Provide sensible defaults: `YYYY-MM-DD` for date, `YYYY-MM-DD HH:mm:ss` for datetime

**Non-Goals:**
- Changing storage format (always ISO 8601 / RFC3339)
- Changing input validation (still accepts ISO 8601 only)
- CLI output formatting beyond `tmd object show` (e.g., `tmd object list`)
- Web UI formatting (future)
- Timezone conversion or locale-aware formatting

## Decisions

### 1. Format strings stored on DisplayProperty

**Decision:** Add `DateFormat` and `DatetimeFormat` string fields to `DisplayProperty`. `BuildDisplayProperties()` in `QueryService` populates them from vault config. `FormatValue()` uses them when formatting date/datetime values.

**Rationale:** This keeps `FormatValue()` self-contained — it doesn't need access to vault config. The format is injected at construction time by the service layer, consistent with how other display metadata (emoji, pin, type) is already injected.

**Alternative considered:** Pass vault config to `FormatValue()` — rejected because it would change the method signature and require all callers to have config access.

### 2. Reuse `convertDateFormat()` from name_template.go

**Decision:** Export `ConvertDateFormat()` (capitalize) and call it from `FormatValue()`.

**Rationale:** The same token set (`YYYY`, `MM`, `DD`, `HH`, `mm`, `ss`) is already defined and tested in `dateFormatReplacer`. No need to duplicate.

### 3. Config placement at top level

**Decision:** Place `date_format` and `datetime_format` as top-level fields in `VaultConfig`, not namespaced under `tui` or `cli`.

**Rationale:** Date display format is cross-cutting — it applies to TUI, CLI, and future web UI. Top-level placement matches the cross-cutting nature. Config keys: `date_format`, `datetime_format`.

### 4. Fallback for empty/invalid format strings

**Decision:** If the config value is empty, use the default. If the format produces an invalid result (e.g., contains unrecognized tokens), the unrecognized tokens pass through as literal text — this is the natural behavior of `strings.NewReplacer`.

**Rationale:** No explicit validation needed. Users see immediate feedback in the TUI if their format is wrong.

## Risks / Trade-offs

- **Risk:** Adding fields to `DisplayProperty` increases struct size slightly → Acceptable, the struct is short-lived (display only).
- **Risk:** User sets an unusual format that truncates information → Mitigation: document format tokens in config key defaults/help text. The default is safe.
