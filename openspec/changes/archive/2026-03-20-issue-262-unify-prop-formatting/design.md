## Context

Property values are formatted for display in multiple places:

1. **`DisplayProperty.Format()`** (`core/display.go`) — type-aware formatting (date, checkbox, multi_select, relation, backlink) that returns `"key: value"` strings. Used by detail panel (`tui/detail.go`) and CLI (`cmd/show.go`).

2. **`formatPropValue()`** (`tui/view_mode.go`) — simplified formatting (string, bool→✓, fallback %v) that returns value-only strings. Used by view mode table rows and preview panel.

The TUI function lacks type awareness — dates render as Go `time.Time` defaults, multi_select as raw slices, relations as raw object IDs with ULIDs.

## Goals / Non-Goals

**Goals:**
- Single formatting pipeline for property values across all TUI views and CLI
- View mode table gains proper date, multi_select, and relation formatting
- Consistent checkbox display (✓/empty) everywhere

**Non-Goals:**
- Changing `BuildDisplayProperties()` call patterns in view mode (avoid performance regression from querying backlinks/reverse relations per row)
- Adding new formatting types beyond what `Format()` already supports
- Changing the display format of non-checkbox types

## Decisions

### 1. Add `FormatValue()` to `DisplayProperty`, refactor `Format()` as wrapper

**Decision:** Add `FormatValue() string` that returns the formatted value without key prefix. Refactor `Format()` to `key + ": " + FormatValue()`.

**Why:** This is the minimal extraction. `FormatValue()` is exactly what view mode needs. `Format()` callers (detail panel, CLI) keep working unchanged.

**Alternative considered:** A standalone `FormatPropertyValue(value any, propType string) string` function. Rejected because it would lose access to `DisplayProperty` fields like `IsRelation`, `FromID`, and `IsBacklink` needed for relation/backlink formatting.

### 2. Construct `DisplayProperty` locally in view mode instead of calling `BuildDisplayProperties()`

**Decision:** In view mode, construct `DisplayProperty` structs directly from `vm.schema.Properties` and `obj.Properties[propName]`, then call `FormatValue()`.

**Why:** `BuildDisplayProperties()` queries the index for reverse relations and backlinks per object. View mode tables render many rows — calling it per row would be a performance regression. View mode columns are user-defined properties (not backlinks/reverse), so the lightweight construction suffices.

### 3. Unify checkbox format to ✓ (true) / empty (false)

**Decision:** Change `FormatValue()` checkbox output from `[x]`/`[ ]` to `✓`/empty string.

**Why:** The user prefers the `✓`/empty style from view mode. It's more visually clean in both table and detail contexts. The `[x]`/`[ ]` notation is Markdown-specific and less readable in a TUI.

**Impact:** This changes the output of `Format()` for checkbox properties from `"key: [x]"` to `"key: ✓"` and `"key: [ ]"` to `"key: "`. This affects detail panel and CLI display.

## Risks / Trade-offs

- **Checkbox format is a visible behavior change** → Acceptable since `✓`/empty is more natural for TUI/CLI display. No external API depends on this format.
- **View mode constructs DisplayProperty without full context** → By design. View columns don't include backlinks/reverse relations, so the lightweight path is correct.
