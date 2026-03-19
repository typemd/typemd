## Context

TUI session state is persisted to `.typemd/tui-state.yaml` via the `SessionState` struct in `tui/state.go`. The struct captures sidebar navigation context (selected object/type, expanded groups, scroll, panel dimensions) but has no awareness of view mode (`panelView`).

The `viewMode` sub-model in `tui/view_mode.go` maintains its own state (typeName, viewName, cursor, scroll, expanded groups) during a session, but `captureState()` ignores it entirely.

## Goals / Non-Goals

**Goals:**
- Persist view mode context so the TUI resumes in the correct view on restart
- Persist cursor position, scroll offset, and expanded groups within the view
- Graceful fallback when saved type/view no longer exists
- Backward compatible: old state files without view fields work unchanged

**Non-Goals:**
- Persisting detail object state (which object was open within view mode)
- Persisting search/filter state within views
- Changing view mode's internal data model

## Decisions

### 1. Flat fields on SessionState (not a nested struct)

Add `ViewTypeName`, `ViewName`, `ViewCursor`, `ViewScroll`, and `ViewExpandedGroups` as top-level optional fields on `SessionState`, matching the existing pattern (all fields are flat with `omitempty`).

**Alternative considered:** A nested `ViewState` struct. Rejected because the existing `SessionState` is entirely flat, and YAML serialization with `omitempty` works cleanly with flat fields. A nested struct would be inconsistent with the current pattern.

### 2. View mode presence implies panelView restoration

If `ViewTypeName` and `ViewName` are both non-empty in the loaded state, the TUI enters view mode on startup. No separate "was in view mode" boolean is needed — the presence of these fields is sufficient.

### 3. Fallback chain for missing type/view

On startup, if view state fields are present:
1. Check if `ViewTypeName` exists as a type → if not, fall back to sidebar mode
2. Check if `ViewName` exists as a view for that type → if not, try "default" view → if no views exist, fall back to sidebar mode
3. If view loads successfully, apply `ViewCursor`, `ViewScroll`, `ViewExpandedGroups` with bounds clamping

This mirrors the existing fallback pattern for deleted objects (same-type fallback → first-object fallback).

### 4. captureState() checks rightPanel mode

When `rightPanel == panelView` and `viewMode` is non-nil, `captureState()` populates the view fields. When not in view mode, these fields remain zero-valued and `omitempty` excludes them from YAML.

### 5. View mode restoration happens after sidebar state

The startup flow in `Start()` applies sidebar state first (expanded groups, cursor), then checks for view state and enters view mode if applicable. This ensures the sidebar is in a valid state even if view restoration fails and falls back.

## Risks / Trade-offs

- **[Risk] View data changes between sessions (objects added/deleted)** → View mode's `load()` re-queries objects from the vault, so data is always fresh. Cursor/scroll clamping prevents out-of-bounds positions.
- **[Risk] Group labels depend on object data** → `ViewExpandedGroups` stores label strings. If grouping changes (e.g., group_by property renamed), stale labels are silently ignored (same pattern as sidebar `ExpandedGroups`).
- **[Trade-off] No detail object restoration** → Simpler implementation. Users land on the view list, not deep in an object. This matches the "resume context, not exact position" philosophy.
