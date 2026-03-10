## Context

The TUI (`tui/app.go`) currently hardcodes all initial state in the `Start()` function: first group expanded, cursor at position 0, default panel widths, focus on left panel. Every restart forces users to rebuild their navigation context.

The TUI already has a precedent for file-based configuration: `.typemd/tui.yaml` stores theme settings and is loaded at startup. The session state feature follows the same pattern but for ephemeral UI state.

## Goals / Non-Goals

**Goals:**
- Save TUI session state to `.typemd/tui-state.yaml` on exit
- Restore session state on next launch with graceful degradation
- Support two-tier storage (global + vault-level), vault overrides global
- Initial implementation: vault-level only

**Non-Goals:**
- Global-level state storage (`~/.config/tmd/tui-state.yaml`) — deferred to future iteration
- Search state persistence — search is ephemeral by design
- Undo/redo of state changes
- Real-time state sync across multiple TUI instances

## Decisions

### 1. YAML format for state file

**Decision**: Use YAML (`tui-state.yaml`) consistent with existing `.typemd/tui.yaml`.

**Rationale**: All `.typemd/` config files use YAML — keeping state in the same format maintains consistency. Users may want to hand-edit state (e.g., reset panel widths, change selected object). The project already uses `gopkg.in/yaml.v3` throughout.

**Alternative considered**: JSON for simpler marshal/unmarshal — rejected in favor of consistency with existing config convention.

### 2. Object ID for cursor position

**Decision**: Store `selectedObjectID` (e.g., `book/clean-code-01jqr...`) rather than cursor index.

**Rationale**: Object IDs are stable across sessions even when objects are added/deleted. A cursor index of 3 could point to a completely different object after changes.

### 3. Fallback strategy

**Decision**: When the saved object no longer exists, fall back to the first object in the same type group, then to the overall first object.

**Rationale**: Users likely want to stay in the same "neighborhood" of their vault. Falling back to the same type group preserves context better than jumping to the top.

### 4. Expanded groups stored as type names

**Decision**: Store expanded groups as an array of type name strings (e.g., `["book", "person"]`).

**Rationale**: Type names are stable identifiers. If a type is removed, the entry is simply ignored during restore.

### 5. Silent failure on load errors

**Decision**: If the state file is missing, corrupt, or contains invalid data, silently fall back to default behavior (current startup behavior).

**Rationale**: State persistence is a convenience feature, not critical. Users should never be blocked by a broken state file.

### 6. State saved on quit only

**Decision**: Write state file only when the user quits (q / ctrl+c), not continuously.

**Rationale**: Simplest approach, avoids filesystem overhead. A crash will lose state, which is acceptable — the feature is best-effort.

## Risks / Trade-offs

- **[Risk] State file becomes stale after external vault changes** → Mitigation: Fallback logic handles missing objects gracefully; unknown type groups in `expandedGroups` are silently ignored.
- **[Risk] Panel width values invalid for different terminal size** → Mitigation: Existing `clampPanelWidths()` logic already normalizes widths to terminal size; restored values go through the same path.
- **[Trade-off] No global state tier in v1** → Acceptable: vault-level covers the primary use case. Global tier can be added later with a merge strategy.
- **[Trade-off] Crash loses state** → Acceptable: save-on-quit is simpler and sufficient for a convenience feature.
